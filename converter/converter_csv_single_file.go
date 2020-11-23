package converter

import (
	"github.com/toudi/jpk_vat/common"
)

func (p *Parser) parseSAFTSections(line []string) {
	p.naglowki = line
	// prevSection oznacza poprzednią sekcję parsera.
	// chodzi o to, że sekcje mogą być ustawione w pliku w dowolnej kolejności
	// i aby uchronić się przed błędem indeksowania line[wieksze:mniejsze] ustawiamy
	// wskaźnik na poprzednią sekcję. tym sposobem kiedy napotkamy początek kolejnej
	// sekcji wiemy od razu w której kolumnie jest koniec poprzednio parsowanej.
	var prevSection *SekcjaParsera = nil
	// kolumnaStart := 0
	// najpierw iterujemy po kolumnach a później sprawdzamy do której sekcji można przypisać
	// daną kolumnę.
	for nrKolumny := 0; nrKolumny < len(line); nrKolumny++ {
		logger.Debugf("nrKolumny; naglowek; %d, %s\n", nrKolumny, line[nrKolumny])
		for _, sekcja := range p.sekcje {
			// jeśli nagłówki dla sekcji zostały już sparsowane to pomijamy ją
			// i próbujemy dopasować te które jeszcze nie zostały rozpoznane.
			if len(sekcja.kolejnoscPol) > 0 {
				continue
			}
			if line[nrKolumny] == sekcja.start {
				logger.Debugf("Znaleziono sekcję %s", sekcja.nazwa)
				sekcja.kolumnaStart = nrKolumny
				// domyślnie ustawmy koniec wiersza jako koniec sekcji - w późniejszej
				// iteracji nadpiszemy tą wartość
				sekcja.kolumnaKoniec = len(line)

				if prevSection != nil {
					prevSection.kolumnaKoniec = nrKolumny
					prevSection.kolejnoscPol = line[prevSection.kolumnaStart:prevSection.kolumnaKoniec]
					logger.Debugf("Ustawiam zakres sekcji %s na %d:%d\n", prevSection.nazwa, prevSection.kolumnaStart, prevSection.kolumnaKoniec)
				}

				prevSection = sekcja
				break
			}
		}
	}

	prevSection.kolejnoscPol = line[prevSection.kolumnaStart:len(line)]
	logger.Debugf("Ustawiam zakres sekcji %s na %d:%d\n", prevSection.nazwa, prevSection.kolumnaStart, prevSection.kolumnaKoniec)
}

var naglowek string
var pola map[string]string
var atrybuty map[string]string

func (p *Parser) parseLineSingleFile(line []string) {
	// iterujemy po sekcjach i staramy sie parsować elementy.
	for _, sekcja := range p.sekcje {
		// należy odnaleźć kolumnę ze startem sekcji.
		startSekcji := sekcja.kolumnaStart

		logger.Debugf("Próba parsowania sekcji %s (od kolumny %s/%d)", sekcja.nazwa, sekcja.start, startSekcji)
		if common.LineIsEmpty(line[sekcja.kolumnaStart:sekcja.kolumnaKoniec]) {
			// pusta sekcja, lecimy dalej.
			continue
		}

		// znaleźliśmy sekcję. zaczynamy parsowanie.
		logger.Debugf("Pole startowe znalezione. Rozpoczynam parsowanie")

		sekcja.pola = make(map[string]string)
		sekcja.atrybuty = make(map[string]string)

		for kol := startSekcji; kol < len(p.naglowki); kol++ {
			if kol >= sekcja.kolumnaKoniec || p.naglowki[kol] == "stop" {
				logger.Debugf("koniec sekcji")
				break
			}
			naglowek = p.naglowki[kol]
			sekcja.SetData(p.naglowki[kol], line[kol])
		}

		logger.Debugf("===> KONIEC SEKCJI <=== ")
		sekcja.finish(sekcja)
	}

}
