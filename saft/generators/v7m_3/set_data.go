package v7m_3

import (
	"errors"
	"fmt"
	"slices"
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

// this code is just utterly horrible but I do not have the time to fix it.
// the whole jpk generation code is due for a massive rewrite anyway
func (g *v7m_3) checkRequiredChoice1Fields(node *xml.Node, nodeName string, sectionName string) error {
	if fieldsToCheck, exists := requiredChoice1Fields[nodeName]; exists {
		var populatedFields int = 0

		fieldPopulated := fieldsToCheck[0][0]
		// these rules are mutually exclusive. meaning - if fieldPopulated is not empty then we would need to check that populatedFields === 0
		fieldPopulatedValue := node.ValueOfOrDefault(nodeName+"."+fieldPopulated, "")

		for _, fieldName := range fieldsToCheck[1] {
			fullFieldName := nodeName + "." + fieldName
			dataAtField := node.ValueOfOrDefault(fullFieldName, "")
			if dataAtField != "" {
				if dataAtField == "1" {
					populatedFields += 1
				} else {
					return fmt.Errorf("Błąd walidacji sekcji %s; Pole %s ma wartość %s (Oczekiwana wartość to 1 lub puste pole)", sectionName, fieldName, dataAtField)
				}
			}
		}

		// now for the final check:
		trimmedValue := strings.Trim(fieldPopulatedValue, " ")
		if trimmedValue != "" {
			// so if the field on the "left" side of the group is populated, we have to make sure that none of the values
			// on the "right" side of the mutually exclusive group are populated
			if populatedFields != 0 {
				return fmt.Errorf("Błąd walidacji sekcji %s. Grupy pól %v oraz %v są wzajemnie rozłączne (można wypełnić wartość tylko w jednej z tych grup)", sectionName, fieldsToCheck[0], fieldsToCheck[1])
			}
		} else {
			// the "left" side is not populated therefore let's check if any of the fields on the "right" side are populated.
			if populatedFields != 1 {
				if populatedFields > 1 {
					return fmt.Errorf("Błąd walidacji sekcji %s. Tylko jedno z pól (%v) musi mieć wartość równą 1. Wykryto ilość pól: %d", sectionName, fieldsToCheck, populatedFields)
				}
				return fmt.Errorf("Błąd walidacji sekcji %s. Przynajmniej jedno z pól (%v) lub (%v) w sekcji musi być wypełnione", sectionName, fieldsToCheck[0], fieldsToCheck[1])
			}
		}
	}

	return nil
}
