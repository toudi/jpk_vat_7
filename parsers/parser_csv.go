package parsers

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/toudi/jpk_vat_7/saft"
)

const (
	// Pojedynczy plik CSV, gdzie sekcje zdefiniowane są w
	// linii nr. 1 a dane do sekcji rozdzielane są nowymi wierszami
	SingleFile = iota
	// Pojedynczy plik CSV, gdzie jego struktura jest następująca:
	// SEKCJA;nazwa-sekcji
	// Kolumna;kolumna;kolumna
	// dane;dane;dane;dane;
	// SEKCJA;nazwa-sekcji-1
	// kolumna;kolumna; ...
	SingleFileWithSections = iota
)

type CSVParser struct {
	BaseParser

	fileHandler *os.File
	csvReader   *csv.Reader

	lineNo int
}

func (c *CSVParser) Parse(dst *saft.SAFT) error {
	var err error
	var line []string

	if c.fileHandler == nil {
		c.fileHandler, err = os.Open(c.Source)
		if err != nil {
			return err
		}

		defer c.fileHandler.Close()

		c.csvReader = csv.NewReader(c.fileHandler)
		c.csvReader.Comma = []rune(c.Options.CSVDelimiter)[0]
	}

	if err != nil {
		return err
	}

	for {
		line, err = c.csvReader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err = c.processLine(line, dst); err != nil {
			return fmt.Errorf("nie udało się przetworzyć linii %d z pliku %s: %v", c.lineNo, c.Source, err)
		}
	}
	return nil
}
