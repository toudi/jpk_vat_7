package parsers

import (
	"fmt"
	"path"

	log "github.com/sirupsen/logrus"
	"github.com/toudi/jpk_vat_7/common"
	"github.com/toudi/jpk_vat_7/saft"
)

// co do zasady działania nie ma wymogu, żeby pliki miały takie konkretnie
// nazwy bo parser mógłby teoretycznie zrobić pętlę po wszystkich plikach
// w katalogu i próbować je parsować. Ale zrobiłem tak, żeby uprościć program
// i zmniejszyć ilość bałaganu
const plikNaglowek string = "naglowek.csv"
const plikDeklaracja string = "deklaracja.csv"
const plikSprzedaz string = "sprzedaz.csv"
const plikPodmiot string = "podmiot.csv"
const plikKupno string = "zakup.csv"

type CSVDirParser struct {
	BaseParser
}

func (d *CSVDirParser) Parse(dst *saft.SAFT) error {
	var fullPath string
	var csvParser *CSVParser
	var err error

	for _, file := range []string{plikNaglowek, plikDeklaracja, plikSprzedaz, plikPodmiot, plikKupno} {
		fullPath = path.Join(d.Source, file)
		if !common.FileExists(fullPath) {
			return fmt.Errorf("brak pliku %s", fullPath)
		}
		log.Debugf("przetwarzam plik %s", fullPath)
		csvParser = &CSVParser{BaseParser: BaseParser{Source: fullPath, Options: d.Options}}
		if err = csvParser.Parse(dst); err != nil {
			return fmt.Errorf("nie udało się przetworzyć pliku %s: %v", fullPath, err)
		}
	}

	return nil
}
