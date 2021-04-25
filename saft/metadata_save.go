package saft

import (
	"crypto/aes"
	"encoding/base64"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/toudi/jpk_vat_7/common"
)

// nie wiem, być może dałoby radę jakoś sparsować to od razu do struktury saftFormCode
// ale nie byłem w stanie odnaleźć mapowania które trzeba byłoby wpakować po xml:""
// w skończonym czasie.
// przykład na którym bazowałem wziąłem stąd:
// https://stackoverflow.com/questions/42209427/unmarshal-namespaced-xml-tags-golang
type saftFormCode struct {
	Value         string `xml:",chardata"`
	SystemCode    string `xml:"kodSystemowy,attr"`
	SchemaVersion string `xml:"wersjaSchemy,attr"`
}

type saftHeader struct {
	FormCode saftFormCode `xml:"KodFormularza"`
}

type saftXml struct {
	XmlName xml.Name
	Header  saftHeader `xml:"Naglowek"`
}

func saftMetadataFileName(srcFile string) string {
	return strings.TrimSuffix(srcFile, ".xml") + "-metadata.xml"
}

func (m *SAFTMetadata) Save() error {
	var err error
	var saftHeader saftXml
	var saftXmlBytes []byte

	// sprawdźmy, czy pola FormCode są wypełnione.
	if m.TemplateVars.Metadata.SystemCode == "" {
		saftXmlBytes, err = ioutil.ReadFile(m.SaftFilePath)
		if err != nil {
			return fmt.Errorf("nie udało się odczytać zawartości pliku JPK: %v", err)
		}
		if err = xml.Unmarshal(saftXmlBytes, &saftHeader); err != nil {
			return fmt.Errorf("nie udało się sparsować nagłówka JPK do struktury: %v", err)
		}
		m.TemplateVars.Metadata.FormCode = saftHeader.Header.FormCode.Value
		m.TemplateVars.Metadata.SystemCode = saftHeader.Header.FormCode.SystemCode
		m.TemplateVars.Metadata.SchemaVersion = saftHeader.Header.FormCode.SchemaVersion
	}

	m.cipher, err = common.CipherInit(32)
	if err != nil {
		return fmt.Errorf("nie udało się zainicjować szyfru: %v", err)
	}

	m.TemplateVars.IV = make([]byte, aes.BlockSize)
	copy(m.TemplateVars.IV, m.cipher.IV)

	// uzupełniamy metadane pliku źródłowego
	m.TemplateVars.SourceMetadata.Read(m.SaftFilePath)
	// kompresujemy plik źródłowy
	compressedSAFTFileName, err := compressSAFTXml(m.SaftFilePath)
	if err != nil {
		return fmt.Errorf("nie udało się spakować pliku JPK do archiwum: %v", err)
	}
	m.TemplateVars.ArchiveMetadata.Read(compressedSAFTFileName)
	// szyfrujemy plik archiwum
	encryptedArchiveFileName, err := m.encryptSAFTArchive(compressedSAFTFileName)
	if err != nil {
		return fmt.Errorf("nie udało się zaszyfrować pliku archiwum: %v", err)
	}
	m.TemplateVars.EncryptedMetadata.Read(encryptedArchiveFileName)

	var funcMap = template.FuncMap{
		"base64":   base64.StdEncoding.EncodeToString,
		"filename": path.Base,
	}

	tmpl, err := template.New("jpk-metadata").Funcs(funcMap).Parse(saftMetaXmlTemplate)
	if err != nil {
		return fmt.Errorf("Nie udało się sparsować szablonu dla metainformacji JPK: %v", err)
	}

	metaFileName := saftMetadataFileName(m.SaftFilePath)

	metaFile, err := os.Create(metaFileName)
	if err != nil {
		return fmt.Errorf("Nie udało się otworzyć pliku metadanych do zapisu: %v", err)
	}

	if err = tmpl.Execute(metaFile, m.TemplateVars); err != nil {
		return fmt.Errorf("Nie udało się zapisać metadanych do pliku: %v", err)
	}

	return nil
}
