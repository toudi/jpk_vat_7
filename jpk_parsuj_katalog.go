package main

import (
	"fmt"
	"path/filepath"
)

const plikNaglowek string = "naglowek.csv"
const plikDeklaracja string = "deklaracja.csv"
const plikSprzedaz string = "sprzedaz.csv"
const plikPodmiot string = "podmiot.csv"
const plikKupno string = "zakup.csv"

func (j *JPK) parsujKatalog(fileName string) error {
	// stwórz obiekt parsera
	var err error

	if err = parser(filepath.Join(fileName, plikNaglowek), []Sekcja{sekcjaNaglowek}); err != nil {
		return fmt.Errorf("Błąd tworzenia instancji parsera nagłówka: %v", err)
	}

	if err = parser(filepath.Join(fileName, plikPodmiot), []Sekcja{sekcjaPodmiot}); err != nil {
		return fmt.Errorf("Błąd tworzenia instancji parsera podmiotu: %v", err)
	}

	if err = parser(filepath.Join(fileName, plikDeklaracja), []Sekcja{sekcjaDeklaracja}); err != nil {
		return fmt.Errorf("Błąd tworzenia instancji parsera deklaracji: %v", err)
	}

	if err = parser(filepath.Join(fileName, plikSprzedaz), []Sekcja{sekcjaSprzedaz, sekcjaSprzedazCtrl}); err != nil {
		return fmt.Errorf("Błąd tworzenia instancji parsera sprzedaży: %v", err)
	}

	if err = parser(filepath.Join(fileName, plikKupno), []Sekcja{sekcjaZakup, sekcjaZakupCtrl}); err != nil {
		return fmt.Errorf("Błąd tworzenia instancji parsera zakupu: %v", err)
	}

	return nil
}
