package saft

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/toudi/jpk_vat_7/common"
)

// eksport struktury JPK do pliku XML

const tnsXMLNS = "http://crd.gov.pl/wzor/2021/12/27/11148/"
const xsiXMLNS = "http://www.w3.org/2001/XMLSchema-instance"
const etdXMLNS = "http://crd.gov.pl/xml/schematy/dziedzinowe/mf/2021/06/08/eD/DefinicjeTypy/"

// kod powstal doslownie w jeden wieczor. nie mialem czasu na to aby
// bawic sie w eleganckie wywolywanie funkcji xml.marshall
// tym bardziej, ze xml w golang jakos nie bardzo mozna naklonic do wypluwania
// namespace'ow.
func (s *SAFT) exportSection(xml *os.File, sectionId string, section SAFTData) {
	// niestety, ponieważ plik JPK musi się walidować zgodnie ze schematem ministerstwa
	// oznacza to, iż kolejność pól w drzewie ma znaczenie (ministerstwo używa typu Sequence)
	// aby nie robić gigantycznej tablicy z kolejnością kolumn ani też żeby nie implementować
	// jakiejś wymyślnej metody do sortowania wymuszamy aby sekcje w pliku wejściowym zawierały
	// już posortowane kolumny. Wbrew pozorom nie jest to takie trudne, tym bardziej, że program
	// dostarcza przykładowy szablon.
	for _, fieldName := range s.sectionFieldOrder[sectionId] {
		// pomijamy pola które mają kropkę w nazwie bo są to atrybuty - zajmiemy się nimi
		// później. pomijamy także pola z pustą wartością.
		if strings.Contains(fieldName, ".") || section[fieldName] == "" {
			continue
		}
		// pobieramy namespace dla pola - w 99% przypadków jest to "tns" ale są od tego wyjątki
		// w sekcji OsobaFizyczna
		namespace := s.elementNamespace(sectionId, fieldName)
		fmt.Fprintf(xml, "    <%s:%s", namespace, fieldName)
		// jeśli istnieją jakieś atrybuty dla pola to wypisujemy je
		for attrName, attrValue := range section.Attributes(fieldName) {
			fmt.Fprintf(xml, " %s=\"%s\"", attrName, attrValue)
		}
		// i domykamy taga.
		fmt.Fprintf(xml, ">%s</%s:%s>\n", section[fieldName], namespace, fieldName)
	}
}

func (s *SAFT) Save(fileName string) error {
	var err error

	// sprawdźmy, czy katalog do zapisu istnieje.
	dirName := path.Dir(fileName)
	if !common.FileExists(dirName) {
		os.MkdirAll(dirName, 0775)
	}

	if err = ioutil.WriteFile(path.Join(dirName, "podpisz-profilem-zaufanym.url"), []byte("[InternetShortcut]\nURL=https://www.gov.pl/web/gov/podpisz-jpkvat-profilem-zaufanym"), 0644); err != nil {
		return fmt.Errorf("nie udało się stworzyć pliku z linkiem do podpisu")
	}

	xml, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("błąd tworzenia pliku wyjściowego: %v", err)
	}
	defer xml.Close()
	xml.Truncate(0)

	fmt.Fprintf(xml, "<?xml version=\"1.0\" encoding=\"utf-8\"?>\n")
	fmt.Fprintf(xml, "<tns:JPK xmlns:etd=\"%s\" xmlns:tns=\"%s\" xmlns:xsi=\"%s\">\n", etdXMLNS, tnsXMLNS, xsiXMLNS)
	fmt.Fprintf(xml, " <tns:Naglowek>\n")
	s.exportSection(xml, SectionNaglowek, s.naglowek)
	fmt.Fprintf(xml, " </tns:Naglowek>\n")

	fmt.Fprintf(xml, " <tns:Podmiot1 rola=\"Podatnik\">\n")
	if s.podmiot.OsobaFizyczna() {
		fmt.Fprintf(xml, "  <tns:OsobaFizyczna>\n")
		s.exportSection(xml, SectionPodmiot, s.podmiot.osobaFizyczna)
		fmt.Fprintf(xml, "  </tns:OsobaFizyczna>\n")
	} else {
		fmt.Fprintf(xml, "  <tns:OsobaNiefizyczna>\n")
		s.exportSection(xml, SectionPodmiot, s.podmiot.osobaNiefizyczna)
		fmt.Fprintf(xml, "  </tns:OsobaNiefizyczna>\n")

	}
	fmt.Fprintf(xml, " </tns:Podmiot1>\n")
	// sekcja deklaracja (VAT-7)
	fmt.Fprintf(xml, " <tns:Deklaracja>\n")
	fmt.Fprintf(xml, "  <tns:Naglowek>\n")

	s.exportSection(xml, SectionDeklaracjaNaglowek, s.deklaracjaNaglowek)

	fmt.Fprintf(xml, "  </tns:Naglowek>\n")
	fmt.Fprintf(xml, "  <tns:PozycjeSzczegolowe>\n")

	s.exportSection(xml, SectionDeklaracjaPozSzcz, s.deklaracjaPozycjeSzczegolowe)

	fmt.Fprintf(xml, "  </tns:PozycjeSzczegolowe>\n")

	s.exportSection(xml, SectionDeklaracjaPouczenia, s.deklaracjaPouczenia)

	fmt.Fprintf(xml, " </tns:Deklaracja>\n")

	fmt.Fprintf(xml, " <tns:Ewidencja>\n")
	// pozycje sprzedaży
	for _, sprzedaz := range s.sprzedaz.wierszeSprzedazy {
		fmt.Fprintf(xml, " <tns:SprzedazWiersz>\n")
		s.exportSection(xml, SectionSprzedaz, sprzedaz)
		fmt.Fprintf(xml, " </tns:SprzedazWiersz>\n")
	}

	fmt.Fprintf(xml, " <tns:SprzedazCtrl>\n")
	s.exportSection(xml, SectionSprzedazCtrl, s.sprzedaz.sprzedazCtrl)
	fmt.Fprintf(xml, " </tns:SprzedazCtrl>\n")

	// pozycje kupna
	for _, zakup := range s.kupno.wierszeZakup {
		fmt.Fprintf(xml, " <tns:ZakupWiersz>\n")
		s.exportSection(xml, SectionZakup, zakup)
		fmt.Fprintf(xml, " </tns:ZakupWiersz>\n")
	}

	fmt.Fprintf(xml, " <tns:ZakupCtrl>\n")
	s.exportSection(xml, SectionZakupCtrl, s.kupno.zakupCtrl)
	fmt.Fprintf(xml, " </tns:ZakupCtrl>\n")
	fmt.Fprintf(xml, " </tns:Ewidencja>\n")
	fmt.Fprintf(xml, "</tns:JPK>")

	return nil
}
