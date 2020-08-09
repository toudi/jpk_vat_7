package converter

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

func (c *Converter) encryptSAFTFile() error {
	logger.Debugf("Szyfruję plik %s", path.Base(c.compressedSAFTFile()))
	var err error

	logger.Debugf("Generuję klucz do szyfrowania pliku %s\n", c.compressedSAFTFile())
	key := make([]byte, 32)
	iv := make([]byte, 16)

	_, err = rand.Read(key)
	if err != nil {
		return fmt.Errorf("Nie udało się wygenerować klucza szyfrującego: %v", err)
	}

	if _, err = rand.Read(iv); err != nil {
		return fmt.Errorf("Nie udało się odczytać wektora inicjalizującego: %v", err)
	}

	logger.Debugf("Klucz szyfrujący: %v", key)

	block, err := aes.NewCipher(key)
	if err != nil {
		return fmt.Errorf("Nie udało się stworzyć instancji szyfru")
	}

	// odczytanie pliku .zip
	archiveFileBytes, err := ioutil.ReadFile(c.compressedSAFTFile())
	if err != nil {
		return fmt.Errorf("Nie udało się odczytać pliku archiwum: %v", err)
	}
	plaintext, _ := pkcs7Pad(archiveFileBytes, block.BlockSize())
	ciphertext := make([]byte, len(plaintext))

	metadataTemplateVars.IV = make([]byte, aes.BlockSize)
	copy(metadataTemplateVars.IV, iv)
	logger.Debugf("wektor inicjalizujący (IV): %v", iv)

	bm := cipher.NewCBCEncrypter(block, iv)
	bm.CryptBlocks(ciphertext, plaintext)

	encryptedFile, err := os.Create(c.encryptedArchiveFile())
	if err != nil {
		return fmt.Errorf("Nie udało się stworzyć zaszyfrowanego pliku")
	}
	_, err = io.Copy(encryptedFile, bytes.NewReader(ciphertext))
	if err != nil {
		return fmt.Errorf("Nie udało się zapisać zaszyfrowanego pliku")
	}

	logger.Debugf("Pomyślnie zaszyfrnowano: %s => %s", path.Base(c.compressedSAFTFile()), path.Base(c.encryptedArchiveFile()))

	// zaszyfrowanie klucza kluczem publicznym z certyfikatu ministerstwa
	encryptedKey, err := c.encryptKeyWithCertificate(key)
	if err != nil {
		return fmt.Errorf("Nie udało się zaszyfrować klucza certyfikatem ministerstwa: %v", err)
	}
	metadataTemplateVars.EncryptionKey = make([]byte, len(encryptedKey))
	copy(metadataTemplateVars.EncryptionKey, encryptedKey)
	logger.Debugf("dane szablonu: %+v", metadataTemplateVars)
	return nil

}

// https://gist.github.com/huyinghuan/7bf174017bf54efb91ece04a48589b22
var (
	// ErrInvalidBlockSize indicates hash blocksize <= 0.
	ErrInvalidBlockSize = errors.New("invalid blocksize")

	// ErrInvalidPKCS7Data indicates bad input to PKCS7 pad or unpad.
	ErrInvalidPKCS7Data = errors.New("invalid PKCS7 data (empty or not padded)")

	// ErrInvalidPKCS7Padding indicates PKCS7 unpad fails to bad input.
	ErrInvalidPKCS7Padding = errors.New("invalid padding on input")
)

// pkcs7Pad right-pads the given byte slice with 1 to n bytes, where
// n is the block size. The size of the result is x times n, where x
// is at least 1.
func pkcs7Pad(data []byte, blockSize int) ([]byte, error) {
	if blockSize <= 0 {
		return nil, ErrInvalidBlockSize
	}
	if data == nil || len(data) == 0 {
		return nil, ErrInvalidPKCS7Data
	}
	padLen := blockSize - len(data)%blockSize
	padding := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(data, padding...), nil
}
