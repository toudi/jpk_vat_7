package saft

import "github.com/toudi/jpk_vat_7/common"

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
	SourceMetadata    common.FileMetadata // dane pliku źródłowego     (.xml)
	ArchiveMetadata   common.FileMetadata // dane pliku archiwum       (.zip)
	EncryptedMetadata common.FileMetadata // dane pliku zaszyfrowanego (.aes)
	// xml AuthData który będzie użyty tylko jeśli użyjemy autoryzacji za
	// pomocą kwoty przychodu.
	AuthDataXML []byte
}

type SAFTMetadata struct {
	cipher         *common.Cipher
	UseTestGateway bool
	SaftFilePath   string
	AuthData       common.AuthData
	TemplateVars   SAFTMetadataTemplateVars
}

var Metadata = &SAFTMetadata{}
