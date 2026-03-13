package ksef

import (
	"os"

	"github.com/goccy/go-yaml"
	"github.com/samber/lo"
)

type invoiceHash struct {
	Nip   string
	RefNo string
}

type invoice struct {
	RefNo     string `yaml:"ref-no"`
	KSeFRefNo string `yaml:"ksef-ref-no"`
	Issuer    string `yaml:"issuer,omitempty"`
}

type KSeFRegistry struct {
	Invoices []*invoice

	refNoIdex map[invoiceHash]int
}

func LoadRegistry(filename string) (*KSeFRegistry, error) {
	registryFile, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	registry := &KSeFRegistry{
		refNoIdex: make(map[invoiceHash]int),
	}

	if err = yaml.NewDecoder(registryFile).Decode(&registry.Invoices); err != nil {
		return nil, err
	}

	for index, invoice := range registry.Invoices {
		hash := invoiceHash{
			RefNo: invoice.RefNo,
			Nip:   invoice.Issuer,
		}
		registry.refNoIdex[hash] = index
	}

	return registry, nil
}

func (r *KSeFRegistry) GetKSeFNoByRefNo(nip *string, refNo string) string {
	if invoice, exists := r.refNoIdex[invoiceHash{
		Nip:   lo.FromPtrOr(nip, ""),
		RefNo: refNo,
	}]; exists {
		return r.Invoices[invoice].KSeFRefNo
	}

	return ""
}
