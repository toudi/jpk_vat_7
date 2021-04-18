package commands

import (
	"flag"
	"fmt"
	"os"

	"github.com/toudi/jpk_vat_7/common"

	"github.com/toudi/jpk_vat_7/converter"
)

type generateArgsType struct {
	Verbose                bool
	TestGateway            bool
	GenerateAuthData       bool
	AuthData               common.AuthData
	EncodingConversionFile string
	UseCurrentDir          bool
	GenerateMetadata       bool
	CSVDelimiter           string
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

	GenerateCmd.FlagSet.StringVar(&generateArgs.CSVDelimiter, "d", ",", "separator pól CSV")
	GenerateCmd.FlagSet.BoolVar(&generateArgs.Verbose, "v", false, "tryb verbose (zwiększa poziom komunikatów wyjściowych)")
	GenerateCmd.FlagSet.BoolVar(&generateArgs.TestGateway, "t", false, "użycie bramki testowej do generowania metadanych")
	GenerateCmd.FlagSet.BoolVar(&generateArgs.UseCurrentDir, "cd", false, "użycie bieżącego katalogu do wygenerowania pliku wynikowego")
	GenerateCmd.FlagSet.BoolVar(&generateArgs.GenerateMetadata, "m", false, "generuj plik metadanych (wymagane jeśli nie zostanie użyty klient JPK Web)")
	GenerateCmd.FlagSet.StringVar(&generateArgs.EncodingConversionFile, "e", "", "użyj pliku z mapą konwersji znaków")

	handleMetadataArgs(GenerateCmd.FlagSet)

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
		converter.GeneratorOptions.UseCurrentDir = args.UseCurrentDir
		converter.GeneratorOptions.GenerateMetadata = args.GenerateMetadata
		converter.Delimiter = args.CSVDelimiter

		if args.EncodingConversionFile != "" {
			converter.PrepareEncodingConversionTable(args.EncodingConversionFile)
		}
		return converter.Run()
	}
	return nil
}
