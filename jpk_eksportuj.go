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
func zapisSekcji(xml *os.File, sekcja SekcjaJPK, exclude []string) {
	for pole, wartosc := range sekcja.pola {
		pomin := false
		if exclude != nil {
			for _, e := range exclude {
				if e == pole {
					pomin = true
					break
				}
			}
		}

		if pomin {
			continue
		}

		fmt.Fprintf(xml, "   <tns:%s", pole)
		for atrybut, wartoscAtrybutu := range sekcja.atrybuty {
			if strings.HasPrefix(atrybut, pole) {
				fmt.Fprintf(xml, " %s=\"%s\"", strings.ReplaceAll(atrybut, pole+".", ""), wartoscAtrybutu)
			}
		}
		fmt.Fprintf(xml, ">%s</tns:%s>\n", wartosc, pole)
	}
}

func (j *JPK) zapiszDoPliku(fileInfo os.FileInfo, fileName string) error {
	if fileInfo.IsDir() {
		fileName = fileInfo.Name() + ".jpk"
	} else {
		fileName += ".jpk"
	}

	xml, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("Błąd tworzenia pliku wyjściowego: %v", err)
	}
	defer xml.Close()
	xml.Truncate(0)

	fmt.Fprintf(xml, "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n")

	fmt.Fprintf(xml, "<tns:JPK xmlns=\"%s\" xmlns:tns=\"%s\" xmlns:xsi=\"%s\">\n xmlns:edt=\"%s\">\n", jpkXMLNS, tnsXMLNS, xsiXMLNS, etdXMLNS)

	fmt.Fprintf(xml, " <tns:Naglowek>\n")
	// fmt.Fprintf(xml, "   <tns:DataWytworzeniaJPK>%s</tns:DataWytworzeniaJPK>\n", j.dataWytworzenia.Format("2006-01-02T15:04:05.99"))
	zapisSekcji(xml, j.naglowek, nil)

	// d := j.deklaracja

	// fmt.Fprintf(xml, "   <tns:KodFormularzaDekl kodSystemowy=\"%s\" kodPodatku=\"%s\" rodzajZobowiazania=\"%s\" wersjaSchemy=\"%s\">%s</tns:KodFormularzaDekl>\n", d.kodSystemowy, d.kodPodatku, d.rodzajZobowiazania, d.wersjaSchemy, d.kod)
	// fmt.Fprintf(xml, "   <tns:WariantFormularzaDekl>%s</tns:WariantFormularzaDekl>\n", d.wariantFormularza)
	fmt.Fprintf(xml, " </tns:Naglowek>\n")
	fmt.Fprintf(xml, " <tns:Podmiot1 rola=\"Podatnik\">\n")
	if j.podmiot.typPodmiotu == "F" {
		fmt.Fprintf(xml, "  <tns:OsobaFizyczna>\n")
		zapisSekcji(xml, j.podmiot.osobaFizyczna, nil)
		fmt.Fprintf(xml, "  </tns:OsobaFizyczna>\n")

	} else if j.podmiot.typPodmiotu == "NF" {
		fmt.Fprintf(xml, "  <tns:OsobaNiefizyczna>\n")
		zapisSekcji(xml, j.podmiot.osobaNiefizyczna, nil)
		fmt.Fprintf(xml, "  </tns:OsobaNiefizyczna>\n")
	}
	fmt.Fprintf(xml, " </tns:Podmiot1>\n")
	// sekcja deklaracja (VAT-7)
	fmt.Fprintf(xml, " <tns:Deklaracja>\n")
	fmt.Fprintf(xml, "  <tns:PozycjeSzczegolowe>\n")

	zapisSekcji(xml, j.deklaracja, []string{"Pouczenia"})

	fmt.Fprintf(xml, "  </tns:PozycjeSzczegolowe>\n")
	fmt.Fprintf(xml, "  <tns:Pouczenia>%s</tns:Pouczenia>\n", j.deklaracja.pola["Pouczenia"])
	fmt.Fprintf(xml, " </tns:Deklaracja>\n")

	// pozycje sprzedaży
	for _, sprzedaz := range j.sprzedaz.wierszeSprzedazy {
		fmt.Fprintf(xml, " <tns:SprzedazWiersz>\n")
		zapisSekcji(xml, sprzedaz, nil)
		fmt.Fprintf(xml, " </tns:SprzedazWiersz>\n")
	}

	fmt.Fprintf(xml, " <tns:SprzedazCtrl>\n")
	zapisSekcji(xml, j.sprzedaz.sprzedazCtrl, nil)
	fmt.Fprintf(xml, " </tns:SprzedazCtrl>\n")

	// pozycje kupna
	for _, zakup := range j.kupno.wierszeZakup {
		fmt.Fprintf(xml, " <tns:ZakupWiersz>\n")
		zapisSekcji(xml, zakup, nil)
		fmt.Fprintf(xml, " </tns:ZakupWiersz>\n")
	}

	fmt.Fprintf(xml, " <tns:ZakupCtrl>\n")
	zapisSekcji(xml, j.kupno.zakupCtrl, nil)
	fmt.Fprintf(xml, " </tns:ZakupCtrl>\n")

	fmt.Fprintf(xml, "</tns:JPK>")

	return nil
}
