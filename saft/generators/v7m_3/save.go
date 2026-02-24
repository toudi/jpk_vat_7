package v7m_3

import (
	"io"
	"slices"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/toudi/jpk_vat_7/saft/section"
)

type incomeCtrlRules struct {
	aggregationFields          []string
	aggregationDecrementFields []string
	exclusionMarks             []string
}

var purchaseCtrlAggregationFields = []string{
	"K_41", "K_43", "K_44", "K_45", "K_46", "K_47",
}

var incomeCtrlProcessor = incomeCtrlRules{
	aggregationFields: []string{
		"K_16", "K_18", "K_20", "K_24", "K_26", "K_28", "K_30", "K_32", "K_33", "K_34",
		"K_35", "K_36", "K_360",
	},
	aggregationDecrementFields: []string{
		"K_35", "K_36", "K_360",
	},
	exclusionMarks: []string{"FP"},
}

type controlRow struct {
	numRows        int
	accumulatedVAT float64
}

func (g *v7m_3) Save(writer io.Writer) error {
	if err := g.checkRequiredChoice1Fields(
		g.root,
		sectionToNode[section.DeklaracjaPozSzcz],
		section.DeklaracjaPozSzcz,
	); err != nil {
		return err
	}

	_, err := g.root.LocateNode(sectionToNode[section.SprzedazCtrl])

	var incomeCtrlRow controlRow

	if err != nil {
		log.Warn("Nie znaleziono sekcji SprzedazCtrl. Odtwarzam ją na podstawie wierszy Ewidencja.Sprzedaz")
		incomeRows, err := g.root.ChildrenIterator(sectionToNode[section.Sprzedaz])
		if err != nil {
			return err
		}

		for incomeRow := range incomeRows {
			incomeCtrlRow.numRows += 1
			if slices.Contains(incomeCtrlProcessor.exclusionMarks, incomeRow.GetChildValue("TypDokumentu", "")) {
				continue
			}
			var accumulatedVAT float64 = 0.0
			for _, child := range incomeRow.Children {
				if slices.Contains(incomeCtrlProcessor.aggregationFields, child.Name) {
					vatAmount, err := strconv.ParseFloat(child.Value, 64)
					if err != nil {
						return err
					}
					if slices.Contains(incomeCtrlProcessor.aggregationDecrementFields, child.Name) {
						vatAmount *= -1
					}
					accumulatedVAT += vatAmount
				}
			}
			incomeCtrlRow.accumulatedVAT = accumulatedVAT
		}

		sprzedazCtrl, _ := g.root.CreateChild(strings.TrimPrefix(sectionToNode[section.SprzedazCtrl], "JPK."), false)

		sprzedazCtrl.SetValuesFromMap(
			map[string]string{
				"LiczbaWierszySprzedazy": strconv.Itoa(incomeCtrlRow.numRows),
				"PodatekNalezny":         strconv.FormatFloat(incomeCtrlRow.accumulatedVAT, 'f', 2, 64),
			},
		)
	}

	_, err = g.root.LocateNode(sectionToNode[section.ZakupCtrl])

	var purchaseCtrlRow controlRow

	if err != nil {
		log.Warn("Nie znaleziono sekcji ZakupCtrl. Odtwarzam ją na podstawie wierszy Ewidencja.Zakup")
		purchaseRows, err := g.root.ChildrenIterator(sectionToNode[section.Zakup])
		if err != nil {
			return err
		}

		for purchaseRow := range purchaseRows {
			purchaseCtrlRow.numRows += 1
			var accumulatedVAT float64 = 0.0
			for _, child := range purchaseRow.Children {
				if slices.Contains(purchaseCtrlAggregationFields, child.Name) {
					vatAmount, err := strconv.ParseFloat(child.Value, 64)
					if err != nil {
						return err
					}
					accumulatedVAT += vatAmount
				}
			}
			purchaseCtrlRow.accumulatedVAT += accumulatedVAT
		}

		zakupCtrl, _ := g.root.CreateChild(strings.TrimPrefix(sectionToNode[section.ZakupCtrl], "JPK."), false)
		zakupCtrl.SetValuesFromMap(
			map[string]string{
				"LiczbaWierszyZakupow": strconv.Itoa(purchaseCtrlRow.numRows),
				"PodatekNaliczony":     strconv.FormatFloat(purchaseCtrlRow.accumulatedVAT, 'f', 2, 64),
			},
		)
	}

	if err := g.root.ApplyOrdering(JPK_V7M_3ChildrenOrder); err != nil {
		return err
	}

	return g.root.DumpToWriter(writer, 0)
}
