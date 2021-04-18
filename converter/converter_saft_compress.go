package converter

import (
	"archive/zip"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

func compressedSAFTFileName(saftFilePath string) string {
	return strings.TrimSuffix(saftFilePath, ".xml") + ".zip"
}

func compressSAFTFile(saftFilePath string) (string, error) {
	var destFileName = compressedSAFTFileName(saftFilePath)
	logger.Debugf("Kompresuję źródłowy plik JPK")
	var err error

	zipFile, err := os.Create(destFileName)
	if err != nil {
		return "", fmt.Errorf("Nie udało się otworzyć pliku archiwum: %v", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	if err = addFileToZip(zipWriter, saftFilePath); err != nil {
		return "", fmt.Errorf("Nie udało się dodać pliku JPK do archiwum")
	}

	logger.Debugf("Pomyślnie skompresowano plik JPK: %s => %s", saftFilePath, destFileName)

	return zipFile.Name(), nil

}

func addFileToZip(zipWriter *zip.Writer, filename string) error {
	var err error
	fileEntry, err := zipWriter.Create(path.Base(filename))
	if err != nil {
		return err
	}
	fileContents, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	fileEntry.Write(fileContents)

	return nil
}
