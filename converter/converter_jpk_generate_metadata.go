package converter

import (
	"encoding/base64"
	"fmt"
	"os"
	"path"
	"strings"
	"text/template"
)

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
</InitUpload>`

func (c *Converter) saftMetadataFile() string {
	return strings.TrimSuffix(c.SAFTFile, ".xml") + "-metadata.xml"
}
func (c *Converter) createSAFTMetadataFile() error {
	var funcMap = template.FuncMap{
		"base64":   base64.StdEncoding.EncodeToString,
		"filename": path.Base,
	}

	var err error

	tmpl, err := template.New("jpk-metadata").Funcs(funcMap).Parse(jpkMetaXmlTemplate)
	if err != nil {
		return fmt.Errorf("Nie udało się sparsować szablonu dla metainformacji JPK: %v", err)
	}

	metaFileName := c.saftMetadataFile()

	metaFile, err := os.Create(metaFileName)
	if err != nil {
		return fmt.Errorf("Nie udało się otworzyć pliku metadanych do zapisu: %v", err)
	}

	if err = tmpl.Execute(metaFile, metadataTemplateVars); err != nil {
		return fmt.Errorf("Nie udało się zapisać metadanych do pliku: %v", err)
	}

	defer metaFile.Close()

	return err
}
