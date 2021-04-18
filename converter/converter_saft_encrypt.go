package converter

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
)

func encryptedArchiveFileName(srcFile string) string {
	return srcFile + ".aes"
}

func (g *MetadataGenerator) encryptSAFTFile(srcFile string) error {
	var err error
	logger.Debugf("Szyfruję plik %s", path.Base(srcFile))
	logger.Debugf("Klucz szyfrujący: %v", g.cipher.Key)

	// odczytanie pliku .zip
	srcFileBytes, err := ioutil.ReadFile(srcFile)
	if err != nil {
		return fmt.Errorf("Nie udało się odczytać pliku archiwum: %v", err)
	}
	encryptedBytes := g.cipher.Encrypt(srcFileBytes, true)
	if err != nil {
		return fmt.Errorf("Nie udało się zaszyfrować pliku archiwum: %v", err)
	}

	logger.Debugf("wektor inicjalizujący (IV): %v", g.cipher.IV)

	dstFile, err := os.Create(encryptedArchiveFileName(srcFile))
	if err != nil {
		return fmt.Errorf("Nie udało się stworzyć zaszyfrowanego pliku: %v", err)
	}
	_, err = io.Copy(dstFile, bytes.NewReader(encryptedBytes))
	if err != nil {
		return fmt.Errorf("Nie udało się zapisać zaszyfrowanego pliku: %v", err)
	}

	logger.Debugf("Pomyślnie zaszyfrowano: %s => %s", path.Base(srcFile), path.Base(dstFile.Name()))

	// zaszyfrowanie klucza kluczem publicznym z certyfikatu ministerstwa
	encryptedKey, err := g.encryptKeyWithCertificate(g.cipher.Key)
	if err != nil {
		return fmt.Errorf("Nie udało się zaszyfrować klucza certyfikatem ministerstwa: %v", err)
	}
	g.state.TemplateVars.EncryptionKey = make([]byte, len(encryptedKey))
	copy(g.state.TemplateVars.EncryptionKey, encryptedKey)
	logger.Debugf("dane szablonu: %+v", g.state.TemplateVars)
	return nil

}
