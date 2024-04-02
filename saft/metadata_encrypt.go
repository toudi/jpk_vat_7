package saft

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/toudi/jpk_vat_7/common"
)

func encryptedArchiveFileName(srcFile string) string {
	return srcFile + ".aes"
}

func (m *SAFTMetadata) encryptSAFTArchive(srcFile string) (string, error) {
	var err error
	log.Debugf("Szyfruję plik %s", filepath.Base(srcFile))
	log.Debugf("Klucz szyfrujący: %v", m.cipher.Key)

	// odczytanie pliku .zip
	srcFileBytes, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return "", fmt.Errorf("Nie udało się odczytać pliku archiwum: %v", err)
	}
	encryptedBytes := m.cipher.Encrypt(srcFileBytes, true)
	if err != nil {
		return "", fmt.Errorf("Nie udało się zaszyfrować pliku archiwum: %v", err)
	}

	log.Debugf("wektor inicjalizujący (IV): %v", m.cipher.IV)

	dstFile, err := os.Create(encryptedArchiveFileName(srcFile))
	if err != nil {
		return "", fmt.Errorf("Nie udało się stworzyć zaszyfrowanego pliku: %v", err)
	}
	_, err = io.Copy(dstFile, bytes.NewReader(encryptedBytes))
	if err != nil {
		return "", fmt.Errorf("Nie udało się zapisać zaszyfrowanego pliku: %v", err)
	}

	log.Debugf(
		"Pomyślnie zaszyfrowano: %s => %s",
		filepath.Base(srcFile),
		filepath.Base(dstFile.Name()),
	)

	// zaszyfrowanie klucza kluczem publicznym z certyfikatu ministerstwa
	encryptedKey, err := m.encryptKeyWithCertificate(m.cipher.Key)
	if err != nil {
		return "", fmt.Errorf("Nie udało się zaszyfrować klucza certyfikatem ministerstwa: %v", err)
	}
	m.TemplateVars.EncryptionKey = make([]byte, len(encryptedKey))
	copy(m.TemplateVars.EncryptionKey, encryptedKey)
	return dstFile.Name(), nil

}

const certsDir = "certyfikaty"

// jeśli funkcja zauważy, że plik certyfikatu istnieje to po prostu wcześniej wyjdzie.
func (m *SAFTMetadata) locateCertFile() (string, error) {
	var err error

	gateway := common.ProductionGatewayURL
	if m.UseTestGateway {
		gateway = common.TestGatewayURL
	}
	url, _ := url.Parse(gateway)

	if !common.FileExists(certsDir) {
		if err = os.Mkdir(certsDir, 0775); err != nil {
			return "", err
		}
	}

	_certFile := filepath.Join(certsDir, url.Host) + ".crt"
	log.Infof("Plik certyfikatu: %s", _certFile)
	certFileExists := common.FileExists(_certFile)

	if !certFileExists {
		return "", fmt.Errorf(
			"Plik certyfikatu nie istnieje; proszę pobrać go ze strony ministerstwa lub repozytorium programu: %s",
			_certFile,
		)
	}
	return _certFile, nil
}

// funkcja szyfruje ciąg bajtów za pomocą klucza publicznego z certyfikatu ministerstwa
func (m *SAFTMetadata) encryptKeyWithCertificate(key []byte) ([]byte, error) {
	var err error
	var certFile string

	log.Debugf("Co będzie szyfrowane: %+v", key)
	certFile, err = m.locateCertFile()
	certFileBytes, err := ioutil.ReadFile(certFile)

	if err != nil {
		return nil, fmt.Errorf("Nie udało się odczytać pliku certyfikatu")
	}
	block, _ := pem.Decode(certFileBytes)
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("Nie udało się sparsować certyfikatu z pliku: %v", err)
	}
	publicKey := cert.PublicKey.(*rsa.PublicKey)
	log.Debugf("Klucz publiczny: %v; size=%d", publicKey, publicKey.Size())
	if err != nil {
		return nil, fmt.Errorf("Nie udało się odczytać klucza publicznego z certyfikatu: %v", err)
	}
	encryptedKeyBytes, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, key)
	if err != nil {
		return nil, fmt.Errorf(
			"Nie udało się zaszyfrować klucza szyfrującego certyfikatem RSA: %v",
			err,
		)
	}
	log.Debugf("Klucz szyfrujący zaszyfrowany certyfikatem: %v", encryptedKeyBytes)
	log.Debugf(
		"Klucz szyfrujący zakodowany base64: %s",
		base64.StdEncoding.EncodeToString(encryptedKeyBytes),
	)

	return encryptedKeyBytes, nil

}
