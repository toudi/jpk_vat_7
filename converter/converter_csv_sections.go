package converter

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
	finish       func(SekcjaParsera)
	kolejnoscPol []string
}

var sekcjaNaglowek SekcjaParsera
var sekcjaPodmiot SekcjaParsera
var sekcjaSprzedaz SekcjaParsera
var sekcjaSprzedazCtrl SekcjaParsera
var sekcjaZakup SekcjaParsera
var sekcjaZakupCtrl SekcjaParsera
var sekcjaDeklaracja SekcjaParsera
var sekcjaDeklaracjaNaglowek SekcjaParsera
var sekcjaDeklaracjaPozycje SekcjaParsera

func (j *JPK) inicjalizujSekcje() {
	sekcjaNaglowek = SekcjaParsera{
		start: "KodFormularza",
		nazwa: "nagłówek",
		finish: func(s SekcjaParsera) {
			j.naglowek.pola = s.pola
			j.naglowek.atrybuty = s.atrybuty
			j.naglowek.atrybuty["CelZlozenia.poz"] = "P_7"
			j.naglowek.sekcjaParsera = &s
		},
	}

	sekcjaPodmiot = SekcjaParsera{
		start: "typPodmiotu",
		nazwa: "podmiot",
		finish: func(s SekcjaParsera) {
			j.podmiot.typPodmiotu = s.pola["typPodmiotu"]
			delete(s.pola, "typPodmiotu")
			if j.podmiot.typPodmiotu == "F" {
				j.podmiot.osobaFizyczna.pola = s.pola
				j.podmiot.osobaFizyczna.sekcjaParsera = &s
			} else if j.podmiot.typPodmiotu == "NF" {
				j.podmiot.osobaNiefizyczna.pola = s.pola
				j.podmiot.osobaNiefizyczna.sekcjaParsera = &s
			}
		},
	}

	sekcjaDeklaracjaNaglowek = SekcjaParsera{
		nazwa: "nagłówek deklaracji VAT",
		start: "KodFormularzaDekl",
		finish: func(s SekcjaParsera) {
			j.deklaracjaNaglowek.sekcjaParsera = &s
			j.deklaracjaNaglowek.pola = s.pola
			j.deklaracjaNaglowek.atrybuty = s.atrybuty
		},
	}
	sekcjaDeklaracjaPozycje = SekcjaParsera{
		nazwa: "pozycje szczegółowej deklaracji VAT",
		start: "P_10",
		finish: func(s SekcjaParsera) {
			j.deklaracjaPozycjeSzczegolowe.sekcjaParsera = &s
			j.deklaracjaPozycjeSzczegolowe.pola = s.pola
			j.deklaracjaPozycjeSzczegolowe.atrybuty = s.atrybuty
		},
	}
	sekcjaDeklaracja = SekcjaParsera{
		nazwa: "deklaracja",
		start: "Pouczenia",
		finish: func(s SekcjaParsera) {
			j.deklaracja.sekcjaParsera = &s
			j.deklaracja.pola = s.pola
		},
	}

	sekcjaSprzedaz = SekcjaParsera{
		nazwa: "sprzedaż (wiersz)",
		start: "LpSprzedazy",
		finish: func(s SekcjaParsera) {
			j.sprzedaz.wierszeSprzedazy = append(j.sprzedaz.wierszeSprzedazy, SekcjaJPK{
				pola:          s.pola,
				sekcjaParsera: &s,
			})
		},
	}

	sekcjaSprzedazCtrl = SekcjaParsera{
		nazwa: "sprzedaż (wiersz kontrolny)",
		start: "LiczbaWierszySprzedazy",
		finish: func(s SekcjaParsera) {
			j.sprzedaz.sprzedazCtrl.pola = s.pola
			j.sprzedaz.sprzedazCtrl.sekcjaParsera = &s
		},
	}

	sekcjaZakup = SekcjaParsera{
		nazwa: "zakup (wiersz)",
		start: "LpZakupu",
		finish: func(s SekcjaParsera) {
			j.kupno.wierszeZakup = append(j.kupno.wierszeZakup, SekcjaJPK{
				pola:          s.pola,
				sekcjaParsera: &s,
			})
		},
	}

	sekcjaZakupCtrl = SekcjaParsera{
		nazwa: "zakup (wiersz kontrolny)",
		start: "LiczbaWierszyZakupow",
		finish: func(s SekcjaParsera) {
			j.kupno.zakupCtrl.pola = s.pola
			j.kupno.zakupCtrl.sekcjaParsera = &s
		},
	}
}
