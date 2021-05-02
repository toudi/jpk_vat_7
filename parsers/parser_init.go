package parsers

import (
	"fmt"
	"os"
	"strings"

	"github.com/toudi/jpk_vat_7/common"
)

func InitParser(input string, options *common.GeneratorOptions) (Parser, error) {
	statInfo, err := os.Stat(input)

	if err != nil {
		return nil, fmt.Errorf("nie udało się odczytać informacji o wejściu: %v", err)
	}

	baseParser := BaseParser{Source: input, Options: options}

	if statInfo.IsDir() {
		return &CSVDirParser{BaseParser: baseParser}, nil
	}
	// nie jest to katalog więc sprawdźmy jaki to typ pliku
	if strings.HasSuffix(input, ".xlsx") {
		return &XLSXParser{BaseParser: baseParser}, nil
	}
	return &CSVParser{BaseParser: baseParser}, nil
}
