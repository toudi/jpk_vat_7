package main

import (
	"fmt"
	"os"
	"strings"
)

const jpkXMLNS = "http://crd.gov.pl/xml/schematy/dziedzinowe/mf/2018/08/24/eD/DefinicjeTypy/"
const tnsXMLNS = "http://jpk.mf.gov.pl/wzor/2019/10/16/10167/"
const xsiXMLNS = "http://www.w3.org/2001/XMLSchema-instance"
const etdXMLNS = "http://crd.gov.pl/xml/schematy/dziedzinowe/mf/2018/08/24/eD/DefinicjeTypy/"

// kod powstal doslownie w jeden wieczor. nie mialem czasu na to aby
// bawic sie w eleganckie wywolywanie funkcji xml.marshall
// tym bardziej, ze xml w golang jakos nie bardzo mozna naklonic do wypluwania
// namespace'ow.
func (j *JPK) zapiszDoPliku(fileInfo os.FileInfo, fileName string) error {
	if fileInfo.IsDir() {
		fileName = fileInfo.Name() + ".jpk"
	} else {
		fileName += ".jpk"
	}

	xml, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Błąd tworzenia pliku wyjściowego: %v", err)
	}
	defer xml.Close()
	xml.Truncate(0)

	fmt.Fprintf(xml, "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n")

	fmt.Fprintf(xml, "<tns:JPK xmlns=\"%s\" xmlns:tns=\"%s\" xmlns:xsi=\"%s\">\n xmlns:edt=\"%s\">", jpkXMLNS, tnsXMLNS, xsiXMLNS, etdXMLNS)

	n := j.naglowek

	fmt.Fprintf(xml, " <tns:Naglowek>\n")
	fmt.Fprintf(xml, "   <tns:KodFormularza kodSystemowy=\"%s\" wersjaSchemy=\"%s\">%s</tns:KodFormularza>\n", n.kodSystemowy, n.wersjaSchemy, n.kodFormularza)
	fmt.Fprintf(xml, "   <tns:WariantFormularza>%s</tns:WariantFormularza>\n", n.wariantFormularza)
	fmt.Fprintf(xml, "   <tns:DataWytworzeniaJPK>%s</tns:DataWytworzeniaJPK>\n", j.dataWytworzenia.Format("2006-01-02T15:04:05.99"))
	fmt.Fprintf(xml, "   <tns:NazwaSystemu>%s</tns:NazwaSystemu>\n", n.nazwaSystemu)
	fmt.Fprintf(xml, "   <tns:CelZlozenia poz=\"%s\">%s</tns:CelZlozenia>\n", n.celZlozeniaPozycja, n.celZlozenia)
	fmt.Fprintf(xml, "   <tns:KodUrzedu>%s</tns:KodUrzedu>\n", n.kodUrzedu)
	fmt.Fprintf(xml, "   <tns:Rok>%s</tns:Rok>\n", n.rok)
	fmt.Fprintf(xml, "   <tns:Miesiac>%s</tns:Miesiac>\n", n.miesiac)

	d := j.deklaracja

	fmt.Fprintf(xml, "   <tns:KodFormularzaDekl kodSystemowy=\"%s\" kodPodatku=\"%s\" rodzajZobowiazania=\"%s\" wersjaSchemy=\"%s\">%s</tns:KodFormularzaDekl>\n", d.kodSystemowy, d.kodPodatku, d.rodzajZobowiazania, d.wersjaSchemy, d.kod)
	fmt.Fprintf(xml, "   <tns:WariantFormularzaDekl>%s</tns:WariantFormularzaDekl>\n", d.wariantFormularza)
	fmt.Fprintf(xml, " </tns:Naglowek>\n")
	fmt.Fprintf(xml, " <tns:Podmiot1 rola=\"Podatnik\">\n")
	if j.podmiot.osobaFizyczna {
		fmt.Fprintf(xml, "  <tns:OsobaFizyczna>\n")
		fmt.Fprintf(xml, "    <tns:NIP>%s</tns:NIP>\n", j.podmiot.NIP)
		fmt.Fprintf(xml, "    <tns:ImiePierwsze>%s</tns:ImiePierwsze>\n", j.podmiot.imie)
		fmt.Fprintf(xml, "    <tns:Nazwisko>%s</tns:Nazwisko>\n", j.podmiot.nazwisko)
		fmt.Fprintf(xml, "    <tns:DataUrodzenia>%s</tns:DataUrodzenia>\n", j.podmiot.dataUrodzenia)
		fmt.Fprintf(xml, "    <tns:Email>%s</tns:Email>\n", j.podmiot.email)
		fmt.Fprintf(xml, "  </tns:OsobaFizyczna>\n")
	} else {
		fmt.Fprintf(xml, "  <tns:OsobaNiefizyczna>\n")
		fmt.Fprintf(xml, "    <tns:NIP>%s</tns:NIP>\n", j.podmiot.NIP)
		fmt.Fprintf(xml, "    <tns:PelnaNazwa>%s</tns:PelnaNazwa>\n", j.podmiot.nazwa)
		fmt.Fprintf(xml, "    <tns:Email>%s</tns:Email>\n", j.podmiot.email)
		fmt.Fprintf(xml, "  </tns:OsobaNiefizyczna>\n")
	}
	fmt.Fprintf(xml, " </tns:Podmiot1>\n")
	// sekcja deklaracja (VAT-7)
	fmt.Fprintf(xml, " <tns:Deklaracja>\n")
	fmt.Fprintf(xml, "  <tns:PozycjeSzczegolowe>\n")

	// wypisujemy tylko te pozycje które są wymagane.
	for pozycja, wartosc := range d.pozycjeSzczegolowe {
		fmt.Fprintf(xml, "   <tns:%s>%s</tns:%s>\n", strings.ToUpper(pozycja), wartosc, strings.ToUpper(pozycja))
	}
	fmt.Fprintf(xml, "  </tns:PozycjeSzczegolowe>\n")
	fmt.Fprintf(xml, "  <tns:Pouczenia>%s</tns:Pouczenia>\n", d.pouczenia)
	fmt.Fprintf(xml, " </tns:Deklaracja>\n")

	// pozycje sprzedaży
	for _, sprzedaz := range j.sprzedaz {
		fmt.Fprintf(xml, " <tns:SprzedazWiersz>\n")
		fmt.Fprintf(xml, "  <tns:LPSprzedazy>%s</tns:LPSprzedazy>\n", sprzedaz.lpSprzedazy)
		fmt.Fprintf(xml, " </tns:SprzedazWiersz>\n")
	}
	fmt.Fprintf(xml, "</tns:JPK>")

	return nil
}
