package v7m_3

import (
	"strings"
	"time"

	"github.com/toudi/jpk_vat_7/saft/generators/abstract"
	"github.com/toudi/jpk_vat_7/saft/ksef"
	"github.com/toudi/jpk_vat_7/utils/xml"
)

type v7m_3 struct {
	root         *xml.Node
	ksefRegistry *ksef.KSeFRegistry
}

var defaults = map[string]string{
	"JPK.#xmlns":                                                   "http://crd.gov.pl/wzor/2025/12/19/14090/",
	"JPK.#xmlns:xsi":                                               "http://www.w3.org/2001/XMLSchema-instance",
	"JPK.Naglowek.KodFormularza":                                   "JPK_VAT",
	"JPK.Naglowek.KodFormularza#kodSystemowy":                      "JPK_V7M (3)",
	"JPK.Naglowek.KodFormularza#wersjaSchemy":                      "1-0E",
	"JPK.Naglowek.WariantFormularza":                               "3",
	"JPK.Naglowek.NazwaSystemu":                                    "WSI Pegasus",
	"JPK.Naglowek.CelZlozenia#poz":                                 "P_7",
	"JPK.Naglowek.CelZlozenia":                                     "1",
	"JPK.Deklaracja.Naglowek.KodFormularzaDekl":                    "VAT-7",
	"JPK.Deklaracja.Naglowek.KodFormularzaDekl#kodSystemowy":       "VAT-7 (23)",
	"JPK.Deklaracja.Naglowek.KodFormularzaDekl#kodPodatku":         "VAT",
	"JPK.Deklaracja.Naglowek.KodFormularzaDekl#rodzajZobowiazania": "Z",
	"JPK.Deklaracja.Naglowek.KodFormularzaDekl#wersjaSchemy":       "1-0E",
	"JPK.Deklaracja.Naglowek.WariantFormularzaDekl":                "23",
	"JPK.Deklaracja.Pouczenia":                                     "1",
	"JPK.Podmiot1#rola":                                            "Podatnik",
}

func Initialize() abstract.Generator {
	root := &xml.Node{Name: "JPK"}
	root.SetValuesFromMap(defaults)
	root.SetValue("JPK.Naglowek.DataWytworzeniaJPK", time.Now().Format(time.RFC3339))

	// now let's follow up by populating defaults from the parsed schema
	for keyName, defaultValue := range JPK_V7M_3RequiredDefaults {
		// but only if it's not a array element (we'll deal with these later)
		nodeParts := strings.Split(keyName, ".")
		nodeName := strings.Join(nodeParts[:len(nodeParts)-1], ".")

		if JPK_V7M_3ArrayElements[nodeName] {
			// this is an array element - skip it.
			continue
		}

		root.SetValue(keyName, defaultValue)
	}
	return &v7m_3{
		root: root,
	}
}

func (g *v7m_3) SetKSeFRegistryFile(registryFile string) error {
	var err error

	if registryFile != "" {
		g.ksefRegistry, err = ksef.LoadRegistry(registryFile)
		if err != nil {
			return err
		}
	}

	return nil
}
