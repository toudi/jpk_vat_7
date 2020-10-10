package converter

import (
	"bytes"
	"crypto/aes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/toudi/jpk_vat/common"
)

func (c *Converter) encryptSAFTFile() error {
	logger.Debugf("Szyfruję plik %s", path.Base(c.compressedSAFTFile()))
	var err error

	logger.Debugf("Generuję klucz do szyfrowania pliku %s\n", c.compressedSAFTFile())
	c.cipher, err = common.CipherInit(32)
	if err != nil {
		return fmt.Errorf("Nie udało się zainicjować szyfru AES: %v", err)
	}

	metadataTemplateVars.IV = make([]byte, aes.BlockSize)
	copy(metadataTemplateVars.IV, c.cipher.IV)

	logger.Debugf("Klucz szyfrujący: %v", c.cipher.Key)

	// odczytanie pliku .zip
	archiveFileBytes, err := ioutil.ReadFile(c.compressedSAFTFile())
	if err != nil {
		return fmt.Errorf("Nie udało się odczytać pliku archiwum: %v", err)
	}
	encryptedArchiveFileBytes := c.cipher.Encrypt(archiveFileBytes, true)
	if err != nil {
		return fmt.Errorf("Nie udało się zaszyfrować pliku archiwum: %v", err)
	}

	logger.Debugf("wektor inicjalizujący (IV): %v", c.cipher.IV)

	encryptedFile, err := os.Create(c.encryptedArchiveFile())
	if err != nil {
		return fmt.Errorf("Nie udało się stworzyć zaszyfrowanego pliku")
	}
	_, err = io.Copy(encryptedFile, bytes.NewReader(encryptedArchiveFileBytes))
	if err != nil {
		return fmt.Errorf("Nie udało się zapisać zaszyfrowanego pliku")
	}

	logger.Debugf("Pomyślnie zaszyfrnowano: %s => %s", path.Base(c.compressedSAFTFile()), path.Base(c.encryptedArchiveFile()))

	// zaszyfrowanie klucza kluczem publicznym z certyfikatu ministerstwa
	encryptedKey, err := c.encryptKeyWithCertificate(c.cipher.Key)
	if err != nil {
		return fmt.Errorf("Nie udało się zaszyfrować klucza certyfikatem ministerstwa: %v", err)
	}
	metadataTemplateVars.EncryptionKey = make([]byte, len(encryptedKey))
	copy(metadataTemplateVars.EncryptionKey, encryptedKey)
	logger.Debugf("dane szablonu: %+v", metadataTemplateVars)
	return nil

}
