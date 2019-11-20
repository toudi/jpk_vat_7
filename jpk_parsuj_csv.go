package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"
)

type Pole struct {
	nrKol    int
	naglowek string
	// wskaźnik na miejsce gdzie zapisać sparsowaną wartość
	p *string
}

type Sekcja struct {
	start       int
	nazwa       string
	pobierzPola func() []Pole
	pola        []Pole
	slownik     map[string]string
}

// typ definiuje parser struktury zapisanej w CSV.
type Parser struct {
	// potrzebne do inicjalizacji
	file      *os.File
	csvReader *csv.Reader
	// naglowki służą do sprawdzenia czy struktura się zgadza.
	naglowki []string
	// sekcje definiują nam co będziemy parsować.
	sekcje []Sekcja
}

func parserInit(filePath string) (*Parser, error) {
	p := &Parser{}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	p.file = file
	reader := csv.NewReader(bufio.NewReader(p.file))
	reader.Comma = ';'
	p.csvReader = reader

	return p, nil
}

func (p *Parser) parsuj() error {
	var nrLinii int = -1
	var nrKolNaglowka int
	var pola []Pole

	for {
		line, err := p.csvReader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Błąd odczytu CSV: %v", err)
		}

		nrLinii++
		if nrLinii == 0 {
			p.naglowki = line
		} else {
			// iterujemy po sekcjach i staramy sie parsować elementy.
			for _, sekcja := range p.sekcje {
				log.Debugf("Próba parsowania sekcji %s (od kolumny %d)", sekcja.nazwa, sekcja.start)
				if line[sekcja.start] != "" {
					// znaleźliśmy sekcję. zaczynamy parsowanie.
					log.Debugf("Pole startowe znalezione. Rozpoczynam parsowanie")
					pola = sekcja.pola

					if sekcja.pobierzPola != nil {
						log.Debugf("Sekcja jest typu tablicowego.")
						pola = sekcja.pobierzPola()
					}
					//else {
					for kol, pole := range pola {
						nrKolNaglowka = sekcja.start + kol
						if pole.nrKol > 0 {
							nrKolNaglowka = pole.nrKol
						}
						if strings.ToUpper(p.naglowki[nrKolNaglowka]) != strings.ToUpper(pole.naglowek) {
							return fmt.Errorf("Nieprawidłowy nagłowek w linii %d; kolumna %d: %s (Oczekiwano %s)", nrLinii, nrKolNaglowka, p.naglowki[nrKolNaglowka], pole.naglowek)
						}
						if *pole.p != line[nrKolNaglowka] {
							*pole.p = line[nrKolNaglowka]
							log.Debugf("%s=%s", p.naglowki[nrKolNaglowka], *pole.p)
						}
					}
					// koniec parsowania pól na podstawie mapowania.
					// spróbujemy sparsować pola dynamiczne
					if len(p.naglowki) > nrKolNaglowka {
						log.Debugf("Pozostały pola do sparsowania. Wczytywanie ich do słownika.")
						if len(line) < len(p.naglowki) {
							return fmt.Errorf("Linia %d nie ma wymaganej ilości kolumn (%d)", nrLinii, len(p.naglowki))
						}
						nrKolNaglowka++
						for i := nrKolNaglowka; i < len(p.naglowki); i++ {
							if line[i] != "" {
								sekcja.slownik[p.naglowki[i]] = line[i]
								log.Debugf("%s => %s", p.naglowki[i], line[i])
							} else {
								log.Debugf("Pomijanie kolumny %d - jest pusta", i)
							}
						}
					}
					//}
				}
				log.Debugf("===> KONIEC SEKCJI <=== ")
			}
		}
	}
	return nil
}

func (p *Parser) Close() {
	p.file.Close()
}

func (j *JPK) parsujCSV(fileName string) error {
	parser, err := parserInit(fileName)
	if parser == nil {
		return fmt.Errorf("Błąd tworzenia instancji parsera: %v", err)
	}
	// tworzymy sekcje.
	sekcjaNaglowek := Sekcja{
		nazwa: "nagłówek",
		start: 0,
		pola: []Pole{
			Pole{naglowek: "kodFormularza", p: &j.naglowek.kodFormularza},
			Pole{naglowek: "kodSystemowy", p: &j.naglowek.kodSystemowy},
			Pole{naglowek: "wersjaSchemy", p: &j.naglowek.wersjaSchemy},
			Pole{naglowek: "wariantFormularza", p: &j.naglowek.wariantFormularza},
			Pole{naglowek: "kodUrzedu", p: &j.naglowek.kodUrzedu, nrKol: 69},
		},
	}
	sekcjaDeklaracja := Sekcja{
		nazwa: "deklaracja VAT-7",
		start: 70,
		pola: []Pole{
			Pole{naglowek: "dekl.kodFormularza", p: &j.deklaracja.kod},
			Pole{naglowek: "dekl.kodSystemowy", p: &j.deklaracja.kodSystemowy},
			Pole{naglowek: "dekl.kodPodatku", p: &j.deklaracja.kodPodatku},
			Pole{naglowek: "dekl.rodzajZobowiazania", p: &j.deklaracja.rodzajZobowiazania},
			Pole{naglowek: "dekl.wersjaSchemy", p: &j.deklaracja.wersjaSchemy},
			Pole{naglowek: "dekl.wariant", p: &j.deklaracja.wariantFormularza},
		},
	}

	parser.sekcje = []Sekcja{sekcjaNaglowek, sekcjaDeklaracja}
	defer parser.Close()
	return parser.parsuj()
}
