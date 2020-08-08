package converter

import (
	"fmt"
	"path/filepath"
)

const plikNaglowek string = "naglowek.csv"
const plikDeklaracja string = "deklaracja.csv"
const plikSprzedaz string = "sprzedaz.csv"
const plikPodmiot string = "podmiot.csv"
const plikKupno string = "zakup.csv"

func (c *Converter) convertDirectory() error {
	var err error

	logger.Debugf("Tryb konwersji katalogu")

	if err = parser(filepath.Join(c.source, plikNaglowek), []Sekcja{sekcjaNaglowek}); err != nil {
		return fmt.Errorf("Błąd tworzenia instancji parsera nagłówka: %v", err)
	}

	if err = parser(filepath.Join(c.source, plikPodmiot), []Sekcja{sekcjaPodmiot}); err != nil {
		return fmt.Errorf("Błąd tworzenia instancji parsera podmiotu: %v", err)
	}

	if err = parser(filepath.Join(c.source, plikDeklaracja), []Sekcja{sekcjaDeklaracja}); err != nil {
		return fmt.Errorf("Błąd tworzenia instancji parsera deklaracji: %v", err)
	}

	if err = parser(filepath.Join(c.source, plikSprzedaz), []Sekcja{sekcjaSprzedaz, sekcjaSprzedazCtrl}); err != nil {
		return fmt.Errorf("Błąd tworzenia instancji parsera sprzedaży: %v", err)
	}

	if err = parser(filepath.Join(c.source, plikKupno), []Sekcja{sekcjaZakup, sekcjaZakupCtrl}); err != nil {
		return fmt.Errorf("Błąd tworzenia instancji parsera zakupu: %v", err)
	}

	return nil
}
