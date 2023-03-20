package parsers

import (
	"fmt"
	"strings"

	"github.com/tealeg/xlsx/v3"
	"github.com/toudi/jpk_vat_7/saft"
)

type XLSXParser struct {
	BaseParser

	workbook *xlsx.File
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
			for key, _ := range x.workbook.Sheet {
				sheetNames = append(sheetNames, fmt.Sprintf("\"%s\"", key))
			}
			return fmt.Errorf("Podano nieistniejący arkusz: %s\nDostępne arkusze: %v", x.Options.XLSXSpreadsheetName, strings.Join(sheetNames, ", "))
		}
	}

	defer worksheet.Close()

	for rowNum := 0; rowNum < worksheet.MaxRow; rowNum++ {
		sheetRow, err := worksheet.Row(rowNum)
		if err != nil {
			return fmt.Errorf("nie udało się odczytać wiersza %d: %v", rowNum, err)
		}

		row = make(XLSXRow, 0)

		if err = sheetRow.ForEachCell(row.readCells); err != nil {
			return fmt.Errorf("nie udało się odczytać wiersza %d: %v", rowNum, err)
		}

		if err = x.processLine(row, dst); err != nil {
			return fmt.Errorf("nie udało się przetworzyć linii %d: %v", rowNum, err)
		}
	}

	return nil
}
