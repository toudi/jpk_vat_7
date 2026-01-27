package saft

import "github.com/toudi/jpk_vat_7/saft/section"

// ten moduł zawiera informacje o sekcjach JPK. Sekcje pobierają informacje
// o kolumnach na podstawie konfiguracji w pliku toml.

type SAFTSection struct {
	Id       string
	StartCol string
}

var SAFTSections = []SAFTSection{
	{Id: section.Naglowek, StartCol: "NazwaSystemu"},
	{Id: section.Podmiot, StartCol: "typPodmiotu"},
	{Id: section.DeklaracjaNaglowek, StartCol: "KodFormularzaDekl"},
	{Id: section.DeklaracjaPozSzcz, StartCol: "P_10"},
	{Id: section.DeklaracjaPouczenia, StartCol: "Pouczenia"},
	{Id: section.Sprzedaz, StartCol: "LpSprzedazy"},
	{Id: section.SprzedazCtrl, StartCol: "LiczbaWierszySprzedazy"},
	{Id: section.Zakup, StartCol: "LpZakupu"},
	{Id: section.ZakupCtrl, StartCol: "LiczbaWierszyZakupow"},
}
