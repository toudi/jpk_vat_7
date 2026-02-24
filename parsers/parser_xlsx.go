package parsers

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/tealeg/xlsx/v3"
	"github.com/toudi/jpk_vat_7/saft"
)

const maxConsequitiveEmptyRows = 200

type XLSXParser struct {
	BaseParser

	workbook *xlsx.File
	// for some really weird reason excell can report max rows way greater than these that
	// are actually populated.
	consequitiveEmptyRows int
}

type XLSXRow []string

func (r *XLSXRow) readCells(c *xlsx.Cell) error {
	var value string
	var cellName string
	var err error

	value, err = c.FormattedValue()
	if err != nil {
		// komórki z datą są tak naprawdę liczbami - spróbujmy sprawdzić czy chodzi o
		// źle sformatowany styl komórki:
		if c.Type() == xlsx.CellTypeNumeric {
			numberFormat := c.GetNumberFormat()
			col, row := c.GetCoordinates()
			cellName = xlsx.GetCellIDStringFromCoords(col, row)
			if strings.Contains(numberFormat, ";@") && strings.Contains(numberFormat, "mm") {
				fmt.Printf("uwaga - wykryto nieprawidłowy format zapisu daty: %v w komórce %s; podmiana na yyyy-mm-dd\n", numberFormat, cellName)
				c.SetFormat("yyyy-mm-dd")
				value, err = c.FormattedValue()
			}
		}
		if err != nil {
			col, row := c.GetCoordinates()
			return fmt.Errorf("nie udało się odczytać zawartości z komórki (wiersz %d, kolumna %d): %v; kod formatu: %v", row+1, col, c.GetNumberFormat(), err)
		}
	}

	*r = append(*r, value)

	return nil
}

func (x *XLSXParser) Parse(dst *saft.SAFT) error {
	x.BaseParser.canIgnoreSectionStartingColumn = true
	var err error
	var row XLSXRow
	var exists bool

	x.workbook, err = xlsx.OpenFile(x.Source)
	if err != nil {
		return fmt.Errorf("nie udało sie otworzyć pliku xlsx: %v", err)
	}

	if len(x.workbook.Sheets) == 0 {
		return fmt.Errorf("arkusz excela nie ma żadnego skoroszytu?: %v", err)
	}

	worksheet := x.workbook.Sheets[0]

	if x.Options.XLSXSpreadsheetName != "" {
		if worksheet, exists = x.workbook.Sheet[x.Options.XLSXSpreadsheetName]; !exists {
			var sheetNames []string
			for key := range x.workbook.Sheet {
				sheetNames = append(sheetNames, fmt.Sprintf("\"%s\"", key))
			}
			return fmt.Errorf("Podano nieistniejący arkusz: %s\nDostępne arkusze: %v", x.Options.XLSXSpreadsheetName, strings.Join(sheetNames, ", "))
		}
	}

	defer worksheet.Close()

	for rowNum := 0; rowNum < worksheet.MaxRow; rowNum++ {
		// fmt.Printf("processing row %d\n", rowNum)
		sheetRow, err := worksheet.Row(rowNum)
		if err != nil {
			return fmt.Errorf("nie udało się odczytać wiersza %d: %v", rowNum, err)
		}

		row = make(XLSXRow, 0)

		if err = sheetRow.ForEachCell(row.readCells); err != nil {
			return fmt.Errorf("nie udało się odczytać wiersza %d: %v", rowNum, err)
		}

		// check if row is empty. for some really weird reason,
		// libreoffice created a spreadsheet that had max rows set to huge number
		var isEmpty bool = true
		for _, cell := range row {
			if cell != "" {
				isEmpty = false
				x.consequitiveEmptyRows = 0
				break
			}
		}

		if isEmpty {
			x.consequitiveEmptyRows += 1
			if x.consequitiveEmptyRows > maxConsequitiveEmptyRows {
				log.Warnf("Przekroczono maksymalną liczbę pustych wierszy (%d). Koniec parsowania", maxConsequitiveEmptyRows)
				return nil
			}
		}

		if err = x.processLine(row, dst); err != nil {
			return fmt.Errorf("nie udało się przetworzyć linii %d: %v", rowNum, err)
		}
	}

	return nil
}
