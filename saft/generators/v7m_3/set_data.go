package v7m_3

import (
	"errors"
	"slices"
	"strings"

	"github.com/toudi/jpk_vat_7/saft/section"
)

var sectionToNode = map[string]string{
	section.Naglowek:          "JPK.Naglowek",
	section.Podmiot:           "JPK.Podmiot1",
	section.DeklaracjaPozSzcz: "JPK.Deklaracja.PozycjeSzczegolowe",
	section.Sprzedaz:          "JPK.Ewidencja.SprzedazWiersz",
	section.SprzedazCtrl:      "JPK.Ewidencja.SprzedazCtrl",
	section.Zakup:             "JPK.Ewidencja.ZakupWiersz",
	section.ZakupCtrl:         "JPK.Ewidencja.ZakupCtrl",
}

var typPodmiotuToNode = map[string]string{
	"f":  "OsobaFizyczna",
	"nf": "OsobaNiefizyczna",
}

var invoiceRefSourceFields = []string{
	"NrKSeF", "OFF", "BFK", "DI",
}

var (
	errUnknownSubjectType = errors.New("unknown subject type")
	errUnknownSection     = errors.New("unknown section")
)

func (g *v7m_3) SetData(sectionName string, data map[string]string) error {
	nodeName := sectionToNode[sectionName]
	if sectionName == section.Podmiot {
		subjectType, exists := data["typPodmiotu"]
		if !exists {
			return errUnknownSubjectType
		}
		delete(data, "typPodmiotu")
		if subjectNodeName, exists := typPodmiotuToNode[strings.ToLower(subjectType)]; !exists {
			return errUnknownSubjectType
		} else {
			nodeName += "." + subjectNodeName
		}
	}

	if nodeName == "" {
		return errors.Join(errUnknownSection, errors.New(sectionName))
	}

	node := g.root

	if JPK_V7M_3ArrayElements[nodeName] {
		ewidencjaNode, _ := g.root.GetOrCreateChild("Ewidencja", false)
		node, _ = ewidencjaNode.CreateChild(strings.TrimPrefix(nodeName, "JPK.Ewidencja."), true)
	}

	// let's assign the default values for required fields. they can be later overwritten
	// by user-provided values
	for fieldNameWithNode, defaultValue := range JPK_V7M_3RequiredDefaults {
		if !strings.HasPrefix(fieldNameWithNode, nodeName) {
			continue
		}
		fieldName := fieldNameWithNode
		if JPK_V7M_3ArrayElements[nodeName] {
			fieldName = strings.TrimPrefix(fieldName, nodeName+".")
		}
		node.SetValue(fieldName, defaultValue)
	}

	var invoiceRefSourceDefined bool

	for field, value := range data {
		if trimmedValue := strings.TrimSpace(value); trimmedValue != "" {
			if slices.Contains(invoiceRefSourceFields, field) {
				invoiceRefSourceDefined = true
			}
			if JPK_V7M_3ArrayElements[nodeName] {
				node.SetValue(field, trimmedValue)
			} else {
				node.SetValue(nodeName+"."+field, trimmedValue)
			}
		}
	}

	if JPK_V7M_3ArrayElements[nodeName] && !invoiceRefSourceDefined {
		node.SetValue("BFK", "1")
	}

	return nil
}
