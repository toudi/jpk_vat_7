package converter

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
)

func (c *Converter) compressSAFTFile() error {
	logger.Debugf("Kompresuję źródłowy plik JPK")
	var err error

	zipFile, err := os.Create(c.compressedSAFTFile())
	if err != nil {
		return fmt.Errorf("Nie udało się otworzyć pliku archiwum: %v", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	if err = addFileToZip(zipWriter, c.SAFTFile); err != nil {
		return fmt.Errorf("Nie udało się dodać pliku JPK do archiwum")
	}

	logger.Debugf("Pomyślnie skompresowano plik JPK: %s => %s", c.SAFTFileName(), c.compressedSAFTFile())

	return nil

}

func addFileToZip(zipWriter *zip.Writer, filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = path.Base(filename)

	// Change to deflate to gain better compression
	// see http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}
