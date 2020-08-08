package commands

import (
	"flag"
	"fmt"

	"github.com/toudi/jpk_vat/converter"
)

type generateArgsType struct {
	Verbose           bool
	TestGateway       bool
	RereadCertificate bool
}

type generateCommand struct {
	Command
}

var GenerateCmd *generateCommand
var generateArgs = &generateArgsType{}

func init() {
	GenerateCmd = &generateCommand{
		Command: Command{
			FlagSet:     flag.NewFlagSet("generuj", flag.ExitOnError),
			Description: "Konwertuje plik CSV (lub katalog z plikami CSV) do pliku JPK",
			Run:         generateRun,
			Args:        generateArgs,
		},
	}
	GenerateCmd.FlagSet.Usage = func() {
		fmt.Printf("Użycie komendy: jpk_vat generuj [flagi] plik-lub-katalog\n\nflagi:\n")
		GenerateCmd.FlagSet.PrintDefaults()
	}
	GenerateCmd.FlagSet.BoolVar(&generateArgs.Verbose, "v", false, "tryb verbose (zwiększa poziom komunikatów wyjściowych)")
	GenerateCmd.FlagSet.BoolVar(&generateArgs.TestGateway, "t", false, "użycie bramki testowej do generowania metadanych")
	GenerateCmd.FlagSet.BoolVar(&generateArgs.RereadCertificate, "r", false, "wymuś ponowne odczytanie certyfikatu SSL")
}

func generateRun(c *Command) error {
	args := generateArgs

	if len(c.FlagSet.Args()) == 0 {
		c.FlagSet.Usage()
	} else {
		converter := converter.ConverterInit(c.FlagSet.Args()[0], args.Verbose)
		converter.GatewayOptions.UseTestGateway = args.TestGateway
		converter.GatewayOptions.ForceSSLCertRead = args.RereadCertificate
		return converter.Run()
	}
	return nil
}
