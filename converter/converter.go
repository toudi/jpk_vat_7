package converter

import (
	"path"

	"github.com/toudi/jpk_vat_7/common"

	log "github.com/sirupsen/logrus"
)

type Converter struct {
	source    string
	SAFTFile  string
	cipher    *common.Cipher
	Delimiter string

	GatewayOptions struct {
		UseTestGateway bool
	}

	GeneratorOptions struct {
		GenerateAuthData bool
		AuthData         common.AuthData
		UseCurrentDir    bool
		GenerateMetadata bool
	}
}

var encodingConversion map[byte]string

func (c *Converter) SAFTFileName() string {
	return path.Base(c.SAFTFile)
}

func ConverterInit(source string, verbose bool) *Converter {
	converter := &Converter{source: source}
	logger = log.New()
	if verbose {
		logger.SetLevel(log.DebugLevel)
	}

	return converter
}

var logger *log.Logger
