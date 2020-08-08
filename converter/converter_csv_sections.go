package converter

// ten moduł zawiera informacje o sekcjach JPK. Sekcje pobierają informacje
// o kolumnach na podstawie konfiguracji w pliku toml.

type Sekcja struct {
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
	finish func(Sekcja)
}

var sekcjaNaglowek Sekcja
var sekcjaPodmiot Sekcja
var sekcjaSprzedaz Sekcja
var sekcjaSprzedazCtrl Sekcja
var sekcjaZakup Sekcja
var sekcjaZakupCtrl Sekcja
var sekcjaDeklaracja Sekcja
var sekcjaDeklaracjaNaglowek Sekcja
var sekcjaDeklaracjaPozycje Sekcja

func (j *JPK) inicjalizujSekcje() {
	sekcjaNaglowek = Sekcja{
		start: "KodFormularza",
		nazwa: "nagłówek",
		finish: func(s Sekcja) {
			j.naglowek.pola = s.pola
			j.naglowek.atrybuty = s.atrybuty
		},
	}

	sekcjaPodmiot = Sekcja{
		start: "typPodmiotu",
		nazwa: "podmiot",
		finish: func(s Sekcja) {
			j.podmiot.typPodmiotu = s.pola["typPodmiotu"]
			delete(s.pola, "typPodmiotu")
			if j.podmiot.typPodmiotu == "F" {
				j.podmiot.osobaFizyczna.pola = s.pola
			} else if j.podmiot.typPodmiotu == "NF" {
				j.podmiot.osobaNiefizyczna.pola = s.pola
			}
		},
	}

	sekcjaDeklaracjaNaglowek = Sekcja{
		nazwa: "nagłówek deklaracji VAT",
		start: "KodFormularzaDekl",
		finish: func(s Sekcja) {
			j.deklaracjaNaglowek.pola = s.pola
			j.deklaracjaNaglowek.atrybuty = s.atrybuty
		},
	}
	sekcjaDeklaracjaPozycje = Sekcja{
		nazwa: "pozycje szczegółowej deklaracji VAT",
		start: "P_10",
		finish: func(s Sekcja) {
			j.deklaracjaPozycjeSzczegolowe.pola = s.pola
			j.deklaracjaPozycjeSzczegolowe.atrybuty = s.atrybuty
		},
	}
	sekcjaDeklaracja = Sekcja{
		nazwa: "deklaracja",
		start: "Pouczenia",
		finish: func(s Sekcja) {
			j.deklaracja.pola = s.pola
		},
	}

	sekcjaSprzedaz = Sekcja{
		nazwa: "sprzedaż (wiersz)",
		start: "LpSprzedazy",
		finish: func(s Sekcja) {
			j.sprzedaz.wierszeSprzedazy = append(j.sprzedaz.wierszeSprzedazy, SekcjaJPK{
				pola: s.pola,
			})
		},
	}

	sekcjaSprzedazCtrl = Sekcja{
		nazwa: "sprzedaż (wiersz kontrolny)",
		start: "LiczbaWierszySprzedazy",
		finish: func(s Sekcja) {
			j.sprzedaz.sprzedazCtrl.pola = s.pola
		},
	}

	sekcjaZakup = Sekcja{
		nazwa: "zakup (wiersz)",
		start: "LpZakupu",
		finish: func(s Sekcja) {
			j.kupno.wierszeZakup = append(j.kupno.wierszeZakup, SekcjaJPK{
				pola: s.pola,
			})
		},
	}

	sekcjaZakupCtrl = Sekcja{
		nazwa: "zakup (wiersz kontrolny)",
		start: "LiczbaWierszyZakupow",
		finish: func(s Sekcja) {
			j.kupno.zakupCtrl.pola = s.pola
		},
	}
}
