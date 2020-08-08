package converter

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path"

	"github.com/toudi/jpk_vat/common"
)

const certsDir = "certyfikaty"

var certFile string

// jeśli funkcja zauważy, że plik certyfikatu istnieje to po prostu wcześniej wyjdzie.
func (c *Converter) downloadCertificate() error {
	gateway := ProductionGatewayURL
	if c.GatewayOptions.UseTestGateway {
		gateway = TestGatewayURL
	}
	url, _ := url.Parse(gateway)

	if !common.FileExists(certsDir) {
		if err = os.Mkdir(certsDir, 0775); err != nil {
			return err
		}
	}

	certFile = path.Join(certsDir, url.Host) + ".crt"
	certFileExists := common.FileExists(certFile)

	if certFileExists && !c.GatewayOptions.ForceSSLCertRead {
		return nil
	}

	conn, err := tls.Dial("tcp", url.Host+":443", &tls.Config{})
	if err != nil {
		return fmt.Errorf("Nie udało się nawiązać połączenia z bramką: %v", err)
	}
	defer conn.Close()
	var b bytes.Buffer

	if err = pem.Encode(&b, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: conn.ConnectionState().PeerCertificates[0].Raw,
	}); err != nil {
		return fmt.Errorf("Nie udało się zapisać certyfikatu do bufora: %v", err)
	}

	if !certFileExists {
		if _, err = os.Create(certFile); err != nil {
			return fmt.Errorf("Nie udało się stworzyć pliku z certyfikatem :%v", err)
		}
	}

	if err = ioutil.WriteFile(certFile, b.Bytes(), 0644); err != nil {
		return fmt.Errorf("Nie udało się zapisać certyfikatu do pliku: %v", err)
	}

	return nil
}

// funkcja szyfruje ciąg bajtów za pomocą klucza publicznego z certyfikatu ministerstwa
func (c *Converter) encryptKeyWithCertificate(key []byte) ([]byte, error) {
	var err error
	c.downloadCertificate()
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
	if err != nil {
		return nil, fmt.Errorf("Nie udało się odczytać klucza publicznego z certyfikatu: %v", err)
	}
	encryptedKeyBytes, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, key)
	if err != nil {
		return nil, fmt.Errorf("Nie udało się zaszyfrować klucza szyfrującego certyfikatem RSA: %v", err)
	}
	logger.Debugf("Klucz szyfrujący zaszyfrowany certyfikatem: %v", encryptedKeyBytes)
	logger.Debugf("Klucz szyfrujący zakodowany base64: %s", base64.StdEncoding.EncodeToString(encryptedKeyBytes))

	return encryptedKeyBytes, nil

}
