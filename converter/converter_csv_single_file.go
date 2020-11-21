package converter

import (
	"github.com/toudi/jpk_vat/common"
)

func (p *Parser) parseSAFTSections(line []string) {
	p.naglowki = line
	kolumnaStart := 0
	for i, sekcja := range p.sekcje {
		logger.Debugf("Sprawdzam sekcje: %s (od kolumny %s)\n", sekcja.nazwa, sekcja.start)
		for nrKolumny := kolumnaStart; nrKolumny < len(line); nrKolumny++ {
			logger.Debugf("nrKolumny; naglowek; %d, %s\n", nrKolumny, line[nrKolumny])
			if line[nrKolumny] == sekcja.start {
				p.sekcje[i].kolumnaStart = nrKolumny
				p.sekcje[i].kolumnaKoniec = len(line)
				kolumnaStart = nrKolumny
				if i > 0 {
					p.sekcje[i-1].kolumnaKoniec = nrKolumny
					p.sekcje[i-1].kolejnoscPol = line[p.sekcje[i-1].kolumnaStart:p.sekcje[i-1].kolumnaKoniec]
					logger.Debugf("Ustawiam koniec sekcji %s na kolumne %d\n", p.sekcje[i-1].nazwa, nrKolumny)
				}
				break
			}
		}
	}
	p.sekcje[len(p.sekcje)-1].kolejnoscPol = line[p.sekcje[len(p.sekcje)-1].kolumnaStart:len(line)]
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
