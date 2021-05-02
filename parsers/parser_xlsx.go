package parsers

import (
	"fmt"

	"github.com/tealeg/xlsx/v3"
	"github.com/toudi/jpk_vat_7/saft"
)

type XLSXParser struct {
	BaseParser

	workbook *xlsx.File
}

type XLSXRow []string

func (r *XLSXRow) readCells(c *xlsx.Cell) error {
	value, err := c.FormattedValue()
	if err != nil {
		return err
	}

	*r = append(*r, value)

	return nil
}

func (x *XLSXParser) Parse(dst *saft.SAFT) error {
	var err error
	var row XLSXRow

	x.workbook, err = xlsx.OpenFile(x.Source)

	if err != nil {
		return fmt.Errorf("nie udało sie otworzyć pliku xlsx: %v", err)
	}

	if len(x.workbook.Sheets) == 0 {
		return fmt.Errorf("arkusz excela nie ma żadnego skoroszytu?: %v", err)
	}

	worksheet := x.workbook.Sheets[0]

	defer worksheet.Close()

	for rowNum := 0; rowNum < worksheet.MaxRow; rowNum++ {
		sheetRow, err := worksheet.Row(rowNum)
		if err != nil {
			return fmt.Errorf("nie udało się odczytać wiersza %d: %v", rowNum, err)
		}

		row = make(XLSXRow, 0)

		sheetRow.ForEachCell(row.readCells)

		if err = x.processLine(row, dst); err != nil {
			return fmt.Errorf("nie udało się przetworzyć linii %d: %v", rowNum, err)
		}
	}

	return nil
}
