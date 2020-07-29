package main

type SekcjaJPK struct {
	// kolejne pola, przekładane na tagi XML
	// np. KodFormularza
	pola map[string]string
	// atrybuty pól. atrybut musi zaczynać się nazwą pola, np.
	// KodFormularza.kodSystemowy
	atrybuty map[string]string
}

type Podmiot struct {
	typPodmiotu      string
	osobaFizyczna    SekcjaJPK
	osobaNiefizyczna SekcjaJPK
}

type Sprzedaz struct {
	wierszeSprzedazy []SekcjaJPK
	sprzedazCtrl     SekcjaJPK
}

type Kupno struct {
	wierszeZakup []SekcjaJPK
	zakupCtrl    SekcjaJPK
}

type JPK struct {
	// dataWytworzenia time.Time

	//
	naglowek                     SekcjaJPK
	deklaracja                   SekcjaJPK
	deklaracjaNaglowek           SekcjaJPK
	deklaracjaPozycjeSzczegolowe SekcjaJPK
	// deklaracja formularzVAT7
	podmiot  Podmiot
	sprzedaz Sprzedaz
	kupno    Kupno
}
