package parsers

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"strings"

	"github.com/toudi/jpk_vat_7/common"
	"github.com/toudi/jpk_vat_7/saft"
	"github.com/toudi/jpk_vat_7/saft/section"
)

type CSVWithSectionsParser struct {
	BaseParser

	// private fields
	currentSection        string
	currentSectionHeaders []string
}

func (p *CSVWithSectionsParser) Parse(dst *saft.SAFT) error {
	var lineBuffer bytes.Buffer
	srcFile, err := os.Open(p.Source)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	var csvReader *csv.Reader

	scanner := bufio.NewScanner(srcFile)
	lineNo := 0
	for scanner.Scan() {
		lineNo += 1
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
			// skip empty lines
			continue
		}
		lineBuffer.Reset()
		lineBuffer.WriteString(p.convertEncoding(line))
		csvReader = csv.NewReader(&lineBuffer)
		csvReader.Comma = rune(p.Options.CSVDelimiter[0])
		var fields []string
		if fields, err = csvReader.Read(); err != nil {
			return err
		}
		if strings.EqualFold(fields[0], section.SectionID) {
			// we have a new section. let's mark it in the state
			p.currentSection = fields[1]
			p.currentSectionHeaders = nil
			continue
		}
		if common.LineIsEmpty(fields) {
			continue
		}
		// we have the data. we just need to understand if they are header data or actual data.
		// the check is relatively simple:
		var rowData map[string]string = make(map[string]string)
		if p.currentSectionHeaders == nil {
			p.currentSectionHeaders = fields
		} else {
			if len(fields) > len(p.currentSectionHeaders) {
				return fmt.Errorf("błędna liczba pól w linii %d: ilość nagłówków: %d; ilość pól: %d", lineNo, len(p.currentSectionHeaders), len(fields))
			}
			// by iterating over the fields rather than the headers we can parse incomplete rows
			for fieldIdx, fieldValue := range fields {
				if trimmedValue := strings.TrimSpace(fieldValue); trimmedValue != "" {
					rowData[p.currentSectionHeaders[fieldIdx]] = trimmedValue
				}
			}
			if err = dst.AddData(p.currentSection, rowData); err != nil {
				return err
			}
		}
	}
	// that wasn't so hard now, was it ?
	return nil
}
