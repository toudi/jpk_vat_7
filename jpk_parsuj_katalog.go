package main

import (
	"fmt"
	"path/filepath"
)

const plikNaglowek string = "naglowek.csv"
const plikVat string = "vat.csv"
const plikSprzedaz string = "sprzedaz.csv"

func (j *JPK) parsujKatalog(fileName string) error {
	// stwórz obiekt parsera
	var err error

	parserNaglowek, err := parserInit(filepath.Join(fileName, plikNaglowek))
	if err != nil {
		return fmt.Errorf("Błąd tworzenia instancji parsera nagłówka: %v", err)
	}
	defer parserNaglowek.Close()

	parserNaglowek.sekcje = []Sekcja{
		Sekcja{
			nazwa: "nagłówek",
			pola: []Pole{
				Pole{naglowek: "kodFormularza", p: &j.naglowek.kodFormularza},
				Pole{naglowek: "kodSystemowy", p: &j.naglowek.kodSystemowy},
				Pole{naglowek: "wersjaSchemy", p: &j.naglowek.wersjaSchemy},
				Pole{naglowek: "wariantFormularza", p: &j.naglowek.wariantFormularza},
				Pole{naglowek: "rok", p: &j.naglowek.rok},
				Pole{naglowek: "miesiac", p: &j.naglowek.miesiac},
				Pole{naglowek: "kodUrzedu", p: &j.naglowek.kodUrzedu},
			},
		},
	}
	if err = parserNaglowek.parsuj(); err != nil {
		return fmt.Errorf("Błąd parsowania pliku nagłówka: %v", err)
	}

	parserVAT, err := parserInit(filepath.Join(fileName, plikVat))
	if err != nil {
		return fmt.Errorf("Błąd tworzenia instancji parsera VAT: %v", err)
	}
	defer parserVAT.Close()

	parserVAT.sekcje = []Sekcja{
		Sekcja{
			nazwa: "deklaracja VAT-7",
			pola: []Pole{
				Pole{naglowek: "dekl.kodFormularza", p: &j.deklaracja.kod},
				Pole{naglowek: "dekl.kodSystemowy", p: &j.deklaracja.kodSystemowy},
				Pole{naglowek: "dekl.kodPodatku", p: &j.deklaracja.kodPodatku},
				Pole{naglowek: "dekl.rodzajZobowiazania", p: &j.deklaracja.rodzajZobowiazania},
				Pole{naglowek: "dekl.wersjaSchemy", p: &j.deklaracja.wersjaSchemy},
				Pole{naglowek: "dekl.wariant", p: &j.deklaracja.wariantFormularza},
				Pole{naglowek: "dekl.pouczenia", p: &j.deklaracja.pouczenia},
				Pole{naglowek: "dekl.p_ordzu", p: &j.deklaracja.p_ordzu},
			},
			slownik: j.deklaracja.pozycjeSzczegolowe,
		},
	}

	if err = parserVAT.parsuj(); err != nil {
		return fmt.Errorf("Błąd parsowania pliku VAT: %v", err)
	}

	parserSprzedaz, err := parserInit(filepath.Join(fileName, plikSprzedaz))
	if err != nil {
		return fmt.Errorf("Błąd tworzenia instancji parsera sprzedaży: %v", err)
	}
	defer parserSprzedaz.Close()

	parserSprzedaz.sekcje = []Sekcja{
		Sekcja{
			nazwa: "wiersze sprzedaży",
			pobierzPola: func() []Pole {
				// dopisanie wiersza do tablicy.
				wierszSprzedaz := Sprzedaz{}
				j.sprzedaz = append(j.sprzedaz, wierszSprzedaz)

				len := len(j.sprzedaz)

				return []Pole{
					Pole{naglowek: "sprzedaz.lp", p: &j.sprzedaz[len-1].lpSprzedazy},
					Pole{naglowek: "sprzedaz.numer", p: &j.sprzedaz[len-1].nrKontrahenta},
					Pole{naglowek: "sprzedaz.nazwa", p: &j.sprzedaz[len-1].nazwaKontrahenta},
					Pole{naglowek: "sprzedaz.dowod", p: &j.sprzedaz[len-1].dowodSprzedazy},
					Pole{naglowek: "sprzedaz.dataWystawienia", p: &j.sprzedaz[len-1].dataWystawienia},
					Pole{naglowek: "sprzedaz.dataSprzedazy", p: &j.sprzedaz[len-1].dataSprzedazy},
					Pole{naglowek: "sprzedaz.typDokumentu", p: &j.sprzedaz[len-1].typDokumentu},
				}
			},
		},
	}

	if err = parserSprzedaz.parsuj(); err != nil {
		return fmt.Errorf("Błąd parsowania pliku sprzedaży: %v", err)
	}

	return nil
}
