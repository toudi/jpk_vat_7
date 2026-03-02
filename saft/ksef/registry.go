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

type InvoiceIssuer struct {
	Nip  string `yaml:"nip"`
	Name string `yaml:"name"`
}

type invoice struct {
	RefNo     string         `yaml:"ref-no"`
	KSeFRefNo string         `yaml:"ksef-ref-no"`
	Type      uint8          `yaml:"type,omitempty"`
	Issuer    *InvoiceIssuer `yaml:"issuer,omitempty"`
}

type KSeFRegistry struct {
	Invoices []*invoice `yaml:"invoices"`

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

	if err = yaml.NewDecoder(registryFile).Decode(registry); err != nil {
		return nil, err
	}

	for index, invoice := range registry.Invoices {
		hash := invoiceHash{
			RefNo: invoice.RefNo,
		}
		if invoice.Type > 0 {
			hash.Nip = invoice.Issuer.Nip
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
