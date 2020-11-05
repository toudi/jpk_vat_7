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

	if err = parser(filepath.Join(c.source, plikNaglowek), []*SekcjaParsera{sekcjaNaglowek}, c.Delimiter); err != nil {
		return fmt.Errorf("Błąd tworzenia instancji parsera nagłówka: %v", err)
	}

	if err = parser(filepath.Join(c.source, plikPodmiot), []*SekcjaParsera{sekcjaPodmiot}, c.Delimiter); err != nil {
		return fmt.Errorf("Błąd tworzenia instancji parsera podmiotu: %v", err)
	}

	if err = parser(filepath.Join(c.source, plikDeklaracja), []*SekcjaParsera{sekcjaDeklaracja}, c.Delimiter); err != nil {
		return fmt.Errorf("Błąd tworzenia instancji parsera deklaracji: %v", err)
	}

	if err = parser(filepath.Join(c.source, plikSprzedaz), []*SekcjaParsera{sekcjaSprzedaz, sekcjaSprzedazCtrl}, c.Delimiter); err != nil {
		return fmt.Errorf("Błąd tworzenia instancji parsera sprzedaży: %v", err)
	}

	if err = parser(filepath.Join(c.source, plikKupno), []*SekcjaParsera{sekcjaZakup, sekcjaZakupCtrl}, c.Delimiter); err != nil {
		return fmt.Errorf("Błąd tworzenia instancji parsera zakupu: %v", err)
	}

	return nil
}
