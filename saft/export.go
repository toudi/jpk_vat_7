package saft

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/toudi/jpk_vat_7/common"
)

// eksport struktury JPK do pliku XML

func (s *SAFT) Save(fileName string) error {
	var err error

	// sprawdźmy, czy katalog do zapisu istnieje.
	dirName := filepath.Dir(fileName)
	if !common.FileExists(dirName) {
		if err = os.MkdirAll(dirName, 0775); err != nil {
			return errors.Join(err, fmt.Errorf("błąd tworzenia katalogu wyjścia"))
		}
	}

	if err = os.WriteFile(filepath.Join(dirName, "podpisz-profilem-zaufanym.url"), []byte("[InternetShortcut]\nURL=https://www.gov.pl/web/gov/podpisz-jpkvat-z-deklaracja-profilem-zaufanym"), 0644); err != nil {
		return fmt.Errorf("nie udało się stworzyć pliku z linkiem do podpisu")
	}

	saftFile, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer saftFile.Close()

	return s.generator.Save(saftFile)
}
