package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/toudi/jpk_vat/common"

	"github.com/toudi/jpk_vat/converter"
)

type generateArgsType struct {
	Verbose          bool
	TestGateway      bool
	GenerateAuthData bool
	AuthData         common.AuthData
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

	GenerateCmd.FlagSet.BoolVar(&generateArgs.Verbose, "v", false, "tryb verbose (zwiększa poziom komunikatów wyjściowych)")
	GenerateCmd.FlagSet.BoolVar(&generateArgs.TestGateway, "t", false, "użycie bramki testowej do generowania metadanych")
	GenerateCmd.FlagSet.BoolVar(&generateArgs.GenerateAuthData, "a", false, "wygeneruj strukturę AuthData (alternatywa dla podpisu kwalifikowanego)")
	GenerateCmd.FlagSet.Float64Var(&generateArgs.AuthData.Income, "a:i", 0.0, "użyj autoryzacji za pomocą kwoty przychodu")
	GenerateCmd.FlagSet.StringVar(&generateArgs.AuthData.NIP, "a:n", "", "numer NIP dla autoryzacji")
	GenerateCmd.FlagSet.StringVar(&generateArgs.AuthData.ImiePierwsze, "a:fn", "", "pole ImiePierwsze dla autoryzacji")
	GenerateCmd.FlagSet.StringVar(&generateArgs.AuthData.Nazwisko, "a:ln", "", "pole Nazwisko dla autoryzacji")
	GenerateCmd.FlagSet.StringVar(&generateArgs.AuthData.DataUrodzenia, "a:bd", "", "pole DataUrodzenia dla autoryzacji. Format: YYYY-MM-DD")

	GenerateCmd.FlagSet.SetOutput(os.Stdout)
	GenerateCmd.FlagSet.Usage = func() {
		fmt.Printf("Użycie komendy: jpk_vat generuj [flagi] plik-lub-katalog%[1]s%[1]sflagi:%[1]s", common.LineBreak)
		GenerateCmd.FlagSet.PrintDefaults()
	}
}

func generateRun(c *Command) error {
	args := generateArgs

	if len(c.FlagSet.Args()) == 0 {
		c.FlagSet.Usage()
	} else {
		converter := converter.ConverterInit(c.FlagSet.Args()[0], args.Verbose)
		converter.GatewayOptions.UseTestGateway = args.TestGateway
		converter.GeneratorOptions.GenerateAuthData = args.GenerateAuthData
		converter.GeneratorOptions.AuthData = args.AuthData
		return converter.Run()
	}
	return nil
}
