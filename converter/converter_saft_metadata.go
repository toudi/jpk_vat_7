package converter

import (
	"bytes"
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/toudi/jpk_vat_7/common"
)

var jpkAuthDataTemplate string = `<?xml version="1.0" encoding="UTF-8"?>
<podp:DaneAutoryzujace xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns:podp="http://e-deklaracje.mf.gov.pl/Repozytorium/Definicje/Podpis/">
	<podp:NIP>{{ .NIP }}</podp:NIP>
	<podp:ImiePierwsze>{{ .ImiePierwsze }}</podp:ImiePierwsze>
	<podp:Nazwisko>{{ .Nazwisko }}</podp:Nazwisko>
	<podp:DataUrodzenia>{{ .DataUrodzenia }}</podp:DataUrodzenia>
	<podp:Kwota>{{ printf "%.2f" .Income }}</podp:Kwota>
</podp:DaneAutoryzujace>`

var jpkMetaXmlTemplate string = `<?xml version="1.0" encoding="UTF-8"?>
<InitUpload xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:xsd="http://www.w3.org/2001/XMLSchema" xmlns="http://e-dokumenty.mf.gov.pl">
	<DocumentType>JPK</DocumentType>
	<Version>01.02.01.20160617</Version>
	<EncryptionKey algorithm="RSA" mode="ECB" padding="PKCS#1" encoding="Base64">{{ base64 .EncryptionKey }}</EncryptionKey>
	<DocumentList>
		<Document>
			<FormCode
				systemCode="{{.Metadata.SystemCode }}"
				schemaVersion="{{.Metadata.SchemaVersion }}">
				JPK_VAT
			</FormCode>
			<FileName>{{ .SourceMetadata.Filename }}</FileName>
			<ContentLength>{{.SourceMetadata.Size }}</ContentLength>
			<HashValue algorithm="SHA-256" encoding="Base64">{{ base64 .SourceMetadata.ContentHash }}</HashValue>
			<FileSignatureList filesNumber="1">
				<Packaging>
					<SplitZip type="split" mode="zip"/>
				</Packaging>
				<Encryption>
					<AES size="256" block="16" mode="CBC" padding="PKCS#7">
						<IV bytes="16" encoding="Base64">{{ base64 .IV }}</IV>
					</AES>
				</Encryption>
				<FileSignature>
					<OrdinalNumber>1</OrdinalNumber>
					<FileName>{{ .ArchiveMetadata.Filename }}</FileName>
					<ContentLength>{{ .ArchiveMetadata.Size }}</ContentLength>
					<HashValue algorithm="MD5" encoding="Base64">{{ base64 .EncryptedMetadata.ContentHash }}</HashValue>
				</FileSignature>
			</FileSignatureList>
		</Document>
	</DocumentList>
	{{ if .AuthDataXML }}
	<AuthData>
	    {{ base64 .AuthDataXML }}
	</AuthData>
	{{ end }}
</InitUpload>`

type FileMetadata struct {
	Filename    string
	Size        int64
	ContentHash []byte
}

type SAFTMetadataTemplateVars struct {
	// w teorii powinniśmy generować losowy IV dla każdego z plików
	// ale jako że i tak wysyłamy tylko jeden plik to nie ma sensu
	// bardziej komplikować programu.
	IV []byte
	// klucz szyfrujący archiwum ZIP, zaszyfrowany za pomocą algorytmu
	// RSA i użyciu klucza publicznego ministerstwa.
	EncryptionKey []byte

	// dane z pliku JPK
	Metadata struct {
		SchemaVersion string
		SystemCode    string
		FormCode      string
	}

	// dane poszczególnych plików, potrzebne do wygenerowania pliku metadanych
	SourceMetadata    FileMetadata // dane pliku źródłowego     (.xml)
	ArchiveMetadata   FileMetadata // dane pliku archiwum       (.zip)
	EncryptedMetadata FileMetadata // dane pliku zaszyfrowanego (.aes)
	// xml AuthData który będzie użyty tylko jeśli użyjemy autoryzacji za
	// pomocą kwoty przychodu.
	AuthDataXML []byte
}

type metadataGeneratorStateType struct {
	IV             []byte
	UseTestGateway bool
	SaftFilePath   string
	AuthData       common.AuthData
	TemplateVars   SAFTMetadataTemplateVars
}

type MetadataGenerator struct {
	cipher *common.Cipher
	state  *metadataGeneratorStateType
}

var MetadataGeneratorState = &metadataGeneratorStateType{
	AuthData:     common.AuthData{},
	TemplateVars: SAFTMetadataTemplateVars{},
}

var generator *MetadataGenerator

func saftMetadataFileName(srcFile string) string {
	return strings.TrimSuffix(srcFile, ".xml") + "-metadata.xml"
}

func MetadataGeneratorInit() (*MetadataGenerator, error) {
	var err error

	generator = &MetadataGenerator{state: MetadataGeneratorState}
	generator.cipher, err = common.CipherInit(32)

	if err != nil {
		return nil, fmt.Errorf("nie udało się zainicjować szyfru AES: %v", err)
	}

	generator.state.IV = make([]byte, aes.BlockSize)
	copy(generator.state.IV, generator.cipher.IV)

	logger.Debugf("Klucz szyfrujący: %v", generator.cipher.Key)

	return generator, nil
}

func (g *MetadataGenerator) GenerateMetadata(srcFile string) error {
	var err error
	var funcMap = template.FuncMap{
		"base64":   base64.StdEncoding.EncodeToString,
		"filename": path.Base,
	}

	// uzupełniamy metadane pliku źródłowego
	populateFileMetadata(srcFile, &g.state.TemplateVars.SourceMetadata)

	// najpierw należy spakować plik wejściowy
	compressedSAFTFile, err := compressSAFTFile(g.state.SaftFilePath)
	if err != nil {
		return fmt.Errorf("nie udało się skompresować pliku JPK: %v", err)
	}
	populateFileMetadata(compressedSAFTFile, &g.state.TemplateVars.ArchiveMetadata)

	// następnie szyfrujemy plik archiwum.
	if err = g.encryptSAFTFile(compressedSAFTFile); err != nil {
		return fmt.Errorf("nie udało się zaszyfrować pliku JPK: %v", err)
	}
	encryptedSAFTFile := encryptedArchiveFileName(compressedSAFTFile)
	populateFileMetadata(encryptedSAFTFile, &g.state.TemplateVars.EncryptedMetadata)

	if g.state.AuthData.Enable {
		var authDataXMLBuffer bytes.Buffer
		tmpl, err := template.New("jpk-authdata").Funcs(funcMap).Parse(jpkAuthDataTemplate)
		if err != nil {
			return fmt.Errorf("Nie udało się sparsować szablonu dla danych autoryzujących JPK: %v", err)
		}
		fmt.Printf("Dane autoryzujące: %+v", g.state.AuthData)
		if err = tmpl.Execute(&authDataXMLBuffer, g.state.AuthData); err != nil {
			return fmt.Errorf("Nie udało się wygenerować dokumentu AuthData: %v", err)
		}

		fmt.Printf("Dane autoryzujące: %s\n", authDataXMLBuffer.String())

		encryptedAuthDataXML := g.cipher.Encrypt(authDataXMLBuffer.Bytes(), true)
		g.state.TemplateVars.AuthDataXML = make([]byte, len(encryptedAuthDataXML))
		copy(g.state.TemplateVars.AuthDataXML, encryptedAuthDataXML)
	}

	tmpl, err := template.New("jpk-metadata").Funcs(funcMap).Parse(jpkMetaXmlTemplate)
	if err != nil {
		return fmt.Errorf("Nie udało się sparsować szablonu dla metainformacji JPK: %v", err)
	}

	metaFileName := saftMetadataFileName(srcFile)

	metaFile, err := os.Create(metaFileName)
	if err != nil {
		return fmt.Errorf("Nie udało się otworzyć pliku metadanych do zapisu: %v", err)
	}

	if err = tmpl.Execute(metaFile, g.state.TemplateVars); err != nil {
		return fmt.Errorf("Nie udało się zapisać metadanych do pliku: %v", err)
	}

	defer metaFile.Close()

	return nil
}

func populateFileMetadata(srcFile string, dest *FileMetadata) {
	dest.Filename = path.Base(srcFile)
	dest.ContentHash = common.Sha256File(srcFile)
	dest.Size, _ = common.FileSize(srcFile)
}

func (g *MetadataGenerator) Run() error {
	return nil
}
