package v7m_3

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/samber/lo"
	"github.com/toudi/jpk_vat_7/saft/section"
	"github.com/toudi/jpk_vat_7/utils"
	"github.com/toudi/jpk_vat_7/utils/xml"
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

var edtNamespaceFields = []string{"NIP", "ImiePierwsze", "Nazwisko", "DataUrodzenia"}

var invoiceRefSourceFields = []string{
	"NrKSeF", "OFF", "BFK", "DI",
}

const ksefRefNoFieldName = "NrKSeF"

var (
	errUnknownSubjectType = errors.New("unknown subject type")
	errUnknownSection     = errors.New("unknown section")
)

func (g *v7m_3) SetData(sectionName string, data map[string]string) error {
	if slices.Contains([]string{section.Zakup, section.Sprzedaz}, sectionName) {
		// let's check if we can add KSeFRefNo from the registry
		if g.ksefRegistry != nil {
			// perfect. let's check if it is required
			ksefRefNo := data[ksefRefNoFieldName]
			if ksefRefNo == "" {
				var nip *string = nil
				invoiceRefNoFieldName := "DowodSprzedazy"
				if sectionName == section.Zakup {
					nip = lo.ToPtr(data["NrDostawcy"])
					invoiceRefNoFieldName = "DowodZakupu"
				}
				data[ksefRefNoFieldName] = g.ksefRegistry.GetKSeFNoByRefNo(
					nip, data[invoiceRefNoFieldName],
				)
			}
		}
	}
	nodeName := sectionToNode[sectionName]
	if sectionName == section.Podmiot {
		subjectType, exists := data["typPodmiotu"]
		if !exists {
			return errUnknownSubjectType
		}
		subjectType = strings.ToLower(subjectType)
		delete(data, "typPodmiotu")
		if subjectNodeName, exists := typPodmiotuToNode[subjectType]; !exists {
			return errUnknownSubjectType
		} else {
			nodeName += "." + subjectNodeName
		}
		if subjectType == "f" {
			g.root.SetValue("JPK.#xmlns:edt", edtNamespace)
			// need to rewrite the namespace in certain subelements
			dataAdjustedNS := make(map[string]string)
			for fieldName, value := range data {
				if slices.Contains(edtNamespaceFields, fieldName) {
					dataAdjustedNS["edt:"+fieldName] = value
				} else {
					dataAdjustedNS[fieldName] = value
				}
			}
			data = dataAdjustedNS
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
			// check if this is a date field
			fullNodeName := nodeName + "." + field
			if JPK_V7M_3DateElements[fullNodeName] {
				dateValue, err := utils.ReparseDateField(trimmedValue)
				if err != nil {
					return fmt.Errorf("Błąd parsowania daty: Komórka %s ma wartość %s i nie udało się jej dopasować do znanych formatów daty", field, trimmedValue)
				} else {
					trimmedValue = dateValue
				}
			}
			if JPK_V7M_3ArrayElements[nodeName] {
				node.SetValue(field, trimmedValue)
			} else {
				node.SetValue(fullNodeName, trimmedValue)
			}
		}
	}

	if JPK_V7M_3ArrayElements[nodeName] && !invoiceRefSourceDefined {
		node.SetValue("BFK", "1")
	}
	// if the invoice number was found then there's no point checking for the other fields to be
	// populated as they are mutually exclusive
	if !invoiceRefSourceDefined {
		// let's validate the required choice1 fields (i.e. a set of fields where at least one field needs to be
		// populated with a value of "1".)
		if err := g.checkRequiredChoice1Fields(node, nodeName, sectionName); err != nil {
			return err
		}
	}

	return nil
}

func (g *v7m_3) checkRequiredChoice1Fields(node *xml.Node, nodeName string, sectionName string) error {
	if err := g.checkMutualExclusion(node, nodeName, sectionName); err != nil {
		return err
	}
	return g.checkConditionalChoice(node, nodeName, sectionName)
}

func (g *v7m_3) checkMutualExclusion(node *xml.Node, nodeName, sectionName string) error {
	fieldsToCheck, exists := mutualExclusion[nodeName]
	if !exists {
		return nil
	}

	fieldPopulatedValue := node.ValueOfOrDefault(nodeName+"."+fieldsToCheck[0][0], "")
	populated, err := g.countPopulated(node, nodeName, fieldsToCheck[1])
	if err != nil {
		return err
	}

	if strings.TrimSpace(fieldPopulatedValue) != "" {
		if len(populated) > 0 {
			return fmt.Errorf("Błąd walidacji sekcji %s. Grupy pól %v oraz %v są wzajemnie rozłączne (można wypełnić wartość tylko w jednej z tych grup). Wypełnione pola: %v", sectionName, fieldsToCheck[0], fieldsToCheck[1], populated)
		}
	} else if len(populated) != 1 {
		if len(populated) > 1 {
			return fmt.Errorf("Błąd walidacji sekcji %s. Tylko jedno z pól (%v) musi mieć wartość równą 1. Wykryto pola: %v", sectionName, fieldsToCheck[1], populated)
		}
		return fmt.Errorf("Błąd walidacji sekcji %s. Przynajmniej jedno z pól (%v) lub (%v) w sekcji musi być wypełnione", sectionName, fieldsToCheck[0], fieldsToCheck[1])
	}
	return nil
}

func (g *v7m_3) checkConditionalChoice(node *xml.Node, nodeName, sectionName string) error {
	cc, exists := conditionalChoiceRules[nodeName]
	if !exists {
		return nil
	}

	triggerValue, err := strconv.ParseFloat(node.ValueOfOrDefault(nodeName+"."+cc.Trigger, "0"), 64)
	if err != nil {
		return err
	}

	populated, err := g.countPopulated(node, nodeName, cc.Choices)
	if err != nil {
		return err
	}

	if triggerValue > 0 {
		if len(populated) != 1 {
			if len(populated) > 1 {
				return fmt.Errorf("Błąd walidacji sekcji %s. Pole %s jest wypełnione więc tylko jedno z pól (%v) musi mieć wartość równą 1. Wykryto pola: %v", sectionName, cc.Trigger, cc.Choices, populated)
			}
			return fmt.Errorf("Błąd walidacji sekcji %s. Pole %s jest wypełnione, więc przynajmniej jedno z pól (%v) musi mieć wartość 1", sectionName, cc.Trigger, cc.Choices)
		}
	} else if len(populated) > 0 {
		return fmt.Errorf("Błąd walidacji sekcji %s. Pole %s jest puste, więc żadne z pól (%v) nie powinno być wypełnione. Wykryto pola: %v", sectionName, cc.Trigger, cc.Choices, populated)
	}
	return nil
}

func (g *v7m_3) countPopulated(node *xml.Node, nodeName string, fields []string) ([]string, error) {
	populated := []string{}
	for _, f := range fields {
		v := node.ValueOfOrDefault(nodeName+"."+f, "")
		if v != "" && v != "1" {
			return nil, fmt.Errorf("Błąd walidacji; Pole %s ma wartość %s (Oczekiwana wartość to 1 lub puste pole)", f, v)
		}
		if v == "1" {
			populated = append(populated, f)
		}
	}
	return populated, nil
}
