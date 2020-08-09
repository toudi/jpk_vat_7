package converter

import (
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Converter struct {
	source   string
	SAFTFile string

	GatewayOptions struct {
		UseTestGateway   bool
		ForceSSLCertRead bool
	}
}

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
	}

	// dane poszczególnych plików, potrzebne do wygenerowania pliku metadanych
	SourceMetadata    FileMetadata // dane pliku źródłowego     (.xml)
	ArchiveMetadata   FileMetadata // dane pliku archiwum       (.zip)
	EncryptedMetadata FileMetadata // dane pliku zaszyfrowanego (.aes)
}

var metadataTemplateVars SAFTMetadataTemplateVars

func (c *Converter) SAFTFileName() string {
	return path.Base(c.SAFTFile)
}
func (c *Converter) compressedSAFTFile() string {
	return strings.TrimSuffix(c.SAFTFile, ".xml") + ".zip"
}

func (c *Converter) encryptedArchiveFile() string {
	return c.compressedSAFTFile() + ".aes"
}

func ConverterInit(source string, verbose bool) *Converter {
	converter := &Converter{source: source}
	logger = log.New()
	if verbose {
		logger.SetLevel(log.DebugLevel)
	}

	return converter
}

var logger *log.Logger
