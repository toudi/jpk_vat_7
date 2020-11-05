package converter

import "strings"

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
				//kolumnaStart = nrKolumny
				if i > 0 {
					p.sekcje[i-1].kolumnaKoniec = nrKolumny
					logger.Debugf("Ustawiam koniec sekcji %s na kolumne %d\n", p.sekcje[i-1].nazwa, nrKolumny)
				}
				break
			}
		}
	}
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
		if line[startSekcji] == "" {
			// pusta sekcja, lecimy dalej.
			continue
		}

		// znaleźliśmy sekcję. zaczynamy parsowanie.
		logger.Debugf("Pole startowe znalezione. Rozpoczynam parsowanie")
		pola = sekcja.pola
		atrybuty = sekcja.atrybuty

		if pola == nil {
			logger.Debugf("pusta mapa, tworze nową")
			pola = make(map[string]string)
			atrybuty = make(map[string]string)
		}

		for kol := startSekcji; kol < len(p.naglowki); kol++ {
			if sekcja.kolejnoscPol == nil {
				sekcja.kolejnoscPol = make([]string, 0)
			}
			if kol >= sekcja.kolumnaKoniec || p.naglowki[kol] == "stop" {
				logger.Debugf("koniec sekcji")
				break
			}
			naglowek = p.naglowki[kol]
			sekcja.kolejnoscPol = append(sekcja.kolejnoscPol, naglowek)
			if line[kol] != "" {
				line[kol] = strings.ReplaceAll(line[kol], "&", "&amp;")
				if encodingConversion != nil {
					line[kol] = convertEncoding(line[kol])
				}
				logger.Debugf("Znalazłem pole: %s (%s)", naglowek, line[kol])
				if strings.Contains(naglowek, ".") {
					// to jest atrybut.
					atrybuty[naglowek] = line[kol]
				} else {
					pola[naglowek] = strings.TrimRight(line[kol], " ")
				}
			} else {
				logger.Debugf("Pomijanie pola %s - pusta wartość", naglowek)
			}
		}

		sekcja.pola = pola
		sekcja.atrybuty = atrybuty

		logger.Debugf("===> KONIEC SEKCJI <=== ")
		sekcja.finish(sekcja)
	}

}
