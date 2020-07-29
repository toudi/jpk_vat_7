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

func parser(filePath string, sekcje []Sekcja) error {
	var err error
	p, err := parserInit(filePath)
	if err != nil {
		return err
	}
	defer p.Close()
	p.sekcje = sekcje
	return p.parsuj()
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
	// var nrKolNaglowka int
	var naglowek string
	var pola map[string]string
	var atrybuty map[string]string

	for {
		line, err := p.csvReader.Read()
		log.Debugf("Odczytano rekord: %+v. Ilość pól: %d\n", line, len(line))
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Błąd odczytu CSV: %v; odczytany rekord: %v o długości: %d", err, line, len(line))
		}

		nrLinii++
		if nrLinii == 0 {
			p.naglowki = line
			kolumnaStart := 0
			for i, sekcja := range p.sekcje {
				log.Debugf("Sprawdzam sekcje: %s (od kolumny %s)\n", sekcja.nazwa, sekcja.start)
				for nrKolumny := kolumnaStart; nrKolumny < len(line); nrKolumny++ {
					log.Debugf("nrKolumny; naglowek; %d, %s\n", nrKolumny, line[nrKolumny])
					if line[nrKolumny] == sekcja.start {
						p.sekcje[i].kolumnaStart = nrKolumny
						p.sekcje[i].kolumnaKoniec = len(line)
						kolumnaStart = nrKolumny
						if i > 0 {
							p.sekcje[i-1].kolumnaKoniec = nrKolumny
							log.Debugf("Ustawiam koniec sekcji %s na kolumne %d\n", p.sekcje[i-1].nazwa, nrKolumny)
						}
						break
					}
				}
			}
		} else {
			// iterujemy po sekcjach i staramy sie parsować elementy.
			for _, sekcja := range p.sekcje {
				// należy odnaleźć kolumnę ze startem sekcji.
				startSekcji := sekcja.kolumnaStart

				log.Debugf("Próba parsowania sekcji %s (od kolumny %s/%d)", sekcja.nazwa, sekcja.start, startSekcji)
				if line[startSekcji] == "" {
					// pusta sekcja, lecimy dalej.
					continue
				}

				// znaleźliśmy sekcję. zaczynamy parsowanie.
				log.Debugf("Pole startowe znalezione. Rozpoczynam parsowanie")
				pola = sekcja.pola
				atrybuty = sekcja.atrybuty

				if pola == nil {
					log.Debugf("pusta mapa, tworze nową")
					pola = make(map[string]string)
					atrybuty = make(map[string]string)
				}

				for kol := startSekcji; kol < len(p.naglowki); kol++ {
					if kol >= sekcja.kolumnaKoniec || p.naglowki[kol] == "stop" {
						log.Debugf("koniec sekcji")
						break
					}
					naglowek = p.naglowki[kol]
					if line[kol] != "" {
						log.Debugf("Znalazłem pole: %s (%s)", naglowek, line[kol])
						if strings.Contains(naglowek, ".") {
							// to jest atrybut.
							atrybuty[naglowek] = line[kol]
						} else {
							pola[naglowek] = strings.TrimRight(line[kol], " ")
						}
					} else {
						log.Debugf("Pomijanie pola %s - pusta wartość", naglowek)
					}
				}

				sekcja.pola = pola
				sekcja.atrybuty = atrybuty

				log.Debugf("===> KONIEC SEKCJI <=== ")
				sekcja.finish(sekcja)
			}
		}
	}
	return nil
}

func (p *Parser) Close() {
	p.file.Close()
}

func (j *JPK) parsujCSV(fileName string) error {
	return parser(fileName, []Sekcja{
		sekcjaNaglowek,
		sekcjaDeklaracjaNaglowek,
		sekcjaPodmiot,
		sekcjaSprzedaz,
		sekcjaSprzedazCtrl,
		sekcjaZakup,
		sekcjaZakupCtrl,
		sekcjaDeklaracjaPozycje,
		sekcjaDeklaracja,
	})
}
