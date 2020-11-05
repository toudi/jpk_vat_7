package converter

import (
	"strings"
)

// ten moduł zawiera informacje o sekcjach JPK. Sekcje pobierają informacje
// o kolumnach na podstawie konfiguracji w pliku toml.

type SekcjaParsera struct {
	// kolumna która oznacza start sekcji
	start string
	// numer kolumny która jest pierwszą w sekcji
	kolumnaStart int
	// numer kolumny która jest ostatnią w sekcji
	kolumnaKoniec int
	// nazwa sekcji (tylko do logów)
	nazwa    string
	pola     map[string]string
	atrybuty map[string]string
	// funkcja która zostanie wywołana po zakończeniu parsowania sekcji
	finish       func(*SekcjaParsera)
	kolejnoscPol []string
}

var sekcjaNaglowek *SekcjaParsera
var sekcjaPodmiot *SekcjaParsera
var sekcjaSprzedaz *SekcjaParsera
var sekcjaSprzedazCtrl *SekcjaParsera
var sekcjaZakup *SekcjaParsera
var sekcjaZakupCtrl *SekcjaParsera
var sekcjaDeklaracja *SekcjaParsera
var sekcjaDeklaracjaNaglowek *SekcjaParsera
var sekcjaDeklaracjaPozycje *SekcjaParsera

func (sp *SekcjaParsera) SetHeaders(headers []string) {
	sp.kolejnoscPol = headers
}

func (sp *SekcjaParsera) SetData(field string, value string) {
	if value != "" {
		value = strings.ReplaceAll(value, "&", "&amp;")
		if encodingConversion != nil {
			value = convertEncoding(value)
		}

		if strings.Contains(field, ".") {
			// to jest atrybut
			logger.Debugf("Znalazłem atrybut: %s o wartości %s", field, value)
			sp.atrybuty[field] = value
		} else {
			logger.Debugf("Znalazłem pole %s o wartości %s", field, value)
			sp.pola[field] = value
		}
	} else {
		logger.Debugf("Pole %s ma pustą wartość - pomijanie", field)
	}
}

var sekcje map[string]*SekcjaParsera

func (j *JPK) inicjalizujSekcje() {
	sekcje = make(map[string]*SekcjaParsera)

	sekcjaNaglowek = &SekcjaParsera{
		start: "KodFormularza",
		nazwa: "NAGLOWEK",
		finish: func(s *SekcjaParsera) {
			j.naglowek.pola = s.pola
			j.naglowek.atrybuty = s.atrybuty
			j.naglowek.atrybuty["CelZlozenia.poz"] = "P_7"
			j.naglowek.sekcjaParsera = s
		},
		pola:     make(map[string]string),
		atrybuty: make(map[string]string),
	}
	sekcje[sekcjaNaglowek.nazwa] = sekcjaNaglowek

	sekcjaPodmiot = &SekcjaParsera{
		start: "typPodmiotu",
		nazwa: "PODMIOT",
		finish: func(s *SekcjaParsera) {
			j.podmiot.typPodmiotu = s.pola["typPodmiotu"]
			delete(s.pola, "typPodmiotu")
			if j.podmiot.typPodmiotu == "F" {
				j.podmiot.osobaFizyczna.pola = s.pola
				j.podmiot.osobaFizyczna.sekcjaParsera = s
			} else if j.podmiot.typPodmiotu == "NF" {
				j.podmiot.osobaNiefizyczna.pola = s.pola
				j.podmiot.osobaNiefizyczna.sekcjaParsera = s
			}
		},
		pola:     make(map[string]string),
		atrybuty: make(map[string]string),
	}
	sekcje[sekcjaPodmiot.nazwa] = sekcjaPodmiot

	sekcjaDeklaracjaNaglowek = &SekcjaParsera{
		nazwa: "DEKLARACJA-NAGLOWEK",
		start: "KodFormularzaDekl",
		finish: func(s *SekcjaParsera) {
			j.deklaracjaNaglowek.sekcjaParsera = s
			j.deklaracjaNaglowek.pola = s.pola
			j.deklaracjaNaglowek.atrybuty = s.atrybuty
		},
		pola:     make(map[string]string),
		atrybuty: make(map[string]string),
	}
	sekcje[sekcjaDeklaracjaNaglowek.nazwa] = sekcjaDeklaracjaNaglowek

	sekcjaDeklaracjaPozycje = &SekcjaParsera{
		nazwa: "DEKLARACJA-POZ-SZCZ",
		start: "P_10",
		finish: func(s *SekcjaParsera) {
			j.deklaracjaPozycjeSzczegolowe.sekcjaParsera = s
			j.deklaracjaPozycjeSzczegolowe.pola = s.pola
			j.deklaracjaPozycjeSzczegolowe.atrybuty = s.atrybuty
		},
		pola:     make(map[string]string),
		atrybuty: make(map[string]string),
	}
	sekcje[sekcjaDeklaracjaPozycje.nazwa] = sekcjaDeklaracjaPozycje

	sekcjaDeklaracja = &SekcjaParsera{
		nazwa: "DEKLARACJA-POUCZENIA",
		start: "Pouczenia",
		finish: func(s *SekcjaParsera) {
			j.deklaracja.sekcjaParsera = s
			j.deklaracja.pola = s.pola
		},
		pola:     make(map[string]string),
		atrybuty: make(map[string]string),
	}
	sekcje[sekcjaDeklaracja.nazwa] = sekcjaDeklaracja

	sekcjaSprzedaz = &SekcjaParsera{
		nazwa: "SPRZEDAZ",
		start: "LpSprzedazy",
		finish: func(s *SekcjaParsera) {
			j.sprzedaz.wierszeSprzedazy = append(j.sprzedaz.wierszeSprzedazy, SekcjaJPK{
				pola:          s.pola,
				sekcjaParsera: s,
			})
		},
		pola:     make(map[string]string),
		atrybuty: make(map[string]string),
	}
	sekcje[sekcjaSprzedaz.nazwa] = sekcjaSprzedaz

	sekcjaSprzedazCtrl = &SekcjaParsera{
		nazwa: "SPRZEDAZ-CTRL",
		start: "LiczbaWierszySprzedazy",
		finish: func(s *SekcjaParsera) {
			j.sprzedaz.sprzedazCtrl.pola = s.pola
			j.sprzedaz.sprzedazCtrl.sekcjaParsera = s
		},
		pola:     make(map[string]string),
		atrybuty: make(map[string]string),
	}
	sekcje[sekcjaSprzedazCtrl.nazwa] = sekcjaSprzedazCtrl

	sekcjaZakup = &SekcjaParsera{
		nazwa: "ZAKUP",
		start: "LpZakupu",
		finish: func(s *SekcjaParsera) {
			j.kupno.wierszeZakup = append(j.kupno.wierszeZakup, SekcjaJPK{
				pola:          s.pola,
				sekcjaParsera: s,
			})
		},
		pola:     make(map[string]string),
		atrybuty: make(map[string]string),
	}
	sekcje[sekcjaZakup.nazwa] = sekcjaZakup

	sekcjaZakupCtrl = &SekcjaParsera{
		nazwa: "ZAKUP-CTRL",
		start: "LiczbaWierszyZakupow",
		finish: func(s *SekcjaParsera) {
			j.kupno.zakupCtrl.pola = s.pola
			j.kupno.zakupCtrl.sekcjaParsera = s
		},
		pola:     make(map[string]string),
		atrybuty: make(map[string]string),
	}
	sekcje[sekcjaZakupCtrl.nazwa] = sekcjaZakupCtrl
}
