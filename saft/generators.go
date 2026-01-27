package saft

import (
	"github.com/toudi/jpk_vat_7/saft/generators/abstract"
	"github.com/toudi/jpk_vat_7/saft/generators/v7m_3"
)

type Namespaces struct {
	ETD string
	XSI string
	TNS string
}

type Generator struct {
	Namespaces Namespaces
	Defaults   map[string]map[string]string
}

const (
	GeneratorV7M_3 string = "v7m_3"
)

var Generators = map[string]abstract.GeneratorFactory{
	GeneratorV7M_3: v7m_3.Initialize,
}
