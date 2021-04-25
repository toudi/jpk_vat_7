package saft

import (
	"errors"
	"strings"
)

var ErrUnknownSAFTSection = errors.New("nierozpoznana sekcja JPK")

// struktury JPK
type SAFTData map[string]string

func (m *SAFTData) Attributes(field string) map[string]string {
	// atrybuty to wszystkie pola które zaczynają się od danej kolumny i mają kropkę w nazwie.
	out := make(map[string]string)
	var attrName string
	var prefix string = field + "."

	for attr, value := range *m {
		if strings.HasPrefix(attr, prefix) {
			attrName = strings.Replace(attr, prefix, "", 1)
			out[attrName] = value
		}
	}

	return out
}

type Podmiot struct {
	typPodmiotu      string
	osobaFizyczna    SAFTData
	osobaNiefizyczna SAFTData
}

func (p Podmiot) OsobaFizyczna() bool {
	return strings.ToUpper(p.typPodmiotu) == "F"
}

type Sprzedaz struct {
	wierszeSprzedazy []SAFTData
	sprzedazCtrl     SAFTData
}

type Kupno struct {
	wierszeZakup []SAFTData
	zakupCtrl    SAFTData
}

type SAFT struct {
	sectionFieldOrder            map[string][]string
	naglowek                     SAFTData
	deklaracjaNaglowek           SAFTData
	deklaracjaPozycjeSzczegolowe SAFTData
	deklaracjaPouczenia          SAFTData
	// deklaracja formularzVAT7
	podmiot  Podmiot
	sprzedaz Sprzedaz
	kupno    Kupno
}

func (s *SAFT) SetFieldOrder(section string, fieldOrder []string) error {
	if s.sectionFieldOrder == nil {
		s.sectionFieldOrder = make(map[string][]string)
	}
	s.sectionFieldOrder[section] = fieldOrder
	return nil
}

func (s *SAFT) AddData(section string, data SAFTData) error {
	if section == SectionNaglowek {
		data["CelZlozenia.poz"] = "P_7"
		s.naglowek = data
		return nil
	}
	if section == SectionDeklaracjaNaglowek {
		s.deklaracjaNaglowek = data
		return nil
	}
	if section == SectionDeklaracjaPozSzcz {
		s.deklaracjaPozycjeSzczegolowe = data
		return nil
	}
	if section == SectionDeklaracjaPouczenia {
		s.deklaracjaPouczenia = data
		return nil
	}
	if section == SectionZakup {
		if s.kupno.wierszeZakup == nil {
			s.kupno.wierszeZakup = make([]SAFTData, 0)
		}
		s.kupno.wierszeZakup = append(s.kupno.wierszeZakup, data)
		return nil
	}
	if section == SectionZakupCtrl {
		s.kupno.zakupCtrl = data
		return nil
	}
	if section == SectionPodmiot {
		s.podmiot.typPodmiotu = data["typPodmiotu"]
		delete(data, "typPodmiotu")
		if s.podmiot.OsobaFizyczna() {
			s.podmiot.osobaFizyczna = data
		} else {
			s.podmiot.osobaNiefizyczna = data
		}
		return nil
	}
	if section == SectionSprzedaz {
		if s.sprzedaz.wierszeSprzedazy == nil {
			s.sprzedaz.wierszeSprzedazy = make([]SAFTData, 0)
		}
		s.sprzedaz.wierszeSprzedazy = append(s.sprzedaz.wierszeSprzedazy, data)
		return nil
	}
	if section == SectionSprzedazCtrl {
		s.sprzedaz.sprzedazCtrl = data
		return nil
	}

	return ErrUnknownSAFTSection
}

func (s *SAFT) elementNamespace(section string, element string) string {
	if section == SectionPodmiot {
		if !(element == "Email" || element == "Telefon") && s.podmiot.OsobaFizyczna() {
			return "etd"
		}
	}

	return "tns"
}
