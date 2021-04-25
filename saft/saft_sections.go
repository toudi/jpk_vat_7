package saft

// ten moduł zawiera informacje o sekcjach JPK. Sekcje pobierają informacje
// o kolumnach na podstawie konfiguracji w pliku toml.

const SectionNaglowek = "NAGLOWEK"
const SectionPodmiot = "PODMIOT"
const SectionDeklaracjaNaglowek = "DEKLARACJA-NAGLOWEK"
const SectionDeklaracjaPozSzcz = "DEKLARACJA-POZ-SZCZ"
const SectionDeklaracjaPouczenia = "DEKLARACJA-POUCZENIA"
const SectionSprzedaz = "SPRZEDAZ"
const SectionSprzedazCtrl = "SPRZEDAZ-CTRL"
const SectionZakup = "ZAKUP"
const SectionZakupCtrl = "ZAKUP-CTRL"

type SAFTSection struct {
	Id       string
	StartCol string
}

var SAFTSections = []SAFTSection{
	{Id: SectionNaglowek, StartCol: "KodFormularza"},
	{Id: SectionPodmiot, StartCol: "typPodmiotu"},
	{Id: SectionDeklaracjaNaglowek, StartCol: "KodFormularzaDekl"},
	{Id: SectionDeklaracjaPozSzcz, StartCol: "P_10"},
	{Id: SectionDeklaracjaPouczenia, StartCol: "Pouczenia"},
	{Id: SectionSprzedaz, StartCol: "LpSprzedazy"},
	{Id: SectionSprzedazCtrl, StartCol: "LiczbaWierszySprzedazy"},
	{Id: SectionZakup, StartCol: "LpZakupu"},
	{Id: SectionZakupCtrl, StartCol: "LiczbaWierszyZakupow"},
}
