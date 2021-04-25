package parsers

import (
	"github.com/toudi/jpk_vat_7/saft"
)

type SAFTItem struct {
	Section string
	Field   string
	Value   string
}

type Parser interface {
	Parse(dst *saft.SAFT) error
	SAFTFileName() string
}
