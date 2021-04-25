package commands

import (
	"flag"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/toudi/jpk_vat_7/common"
	"github.com/toudi/jpk_vat_7/parsers"
	"github.com/toudi/jpk_vat_7/saft"
)

type generateCommand struct {
	Command
}

var GenerateCmd *generateCommand
var generateArgs = &common.GeneratorOptions{}
var logPath string

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
	GenerateCmd.FlagSet.BoolVar(&generateArgs.UseCurrentDir, "cd", false, "użycie bieżącego katalogu do wygenerowania pliku wynikowego")
	GenerateCmd.FlagSet.BoolVar(&generateArgs.GenerateMetadata, "m", false, "generuj plik metadanych (wymagane jeśli nie zostanie użyty klient JPK Web)")
	GenerateCmd.FlagSet.StringVar(&generateArgs.EncodingConversionFile, "e", "", "użyj pliku z mapą konwersji znaków")
	GenerateCmd.FlagSet.StringVar(&logPath, "log", "", "Plik do zapisu logów; Jeśli wartość flagi będzie pusta logi zostaną przekierowane na wyjście standardowe")

	handleMetadataArgs(GenerateCmd.FlagSet)

	GenerateCmd.FlagSet.SetOutput(os.Stdout)
	GenerateCmd.FlagSet.Usage = func() {
		fmt.Printf("Użycie komendy: jpk_vat generuj [flagi] plik-lub-katalog%[1]s%[1]sflagi:%[1]s", common.LineBreak)
		GenerateCmd.FlagSet.PrintDefaults()
	}
}

func generateRun(c *Command) error {
	if len(c.FlagSet.Args()) == 0 {
		c.FlagSet.Usage()
	} else {
		if generateArgs.Verbose {
			log.SetLevel(log.DebugLevel)
		}
		parser, err := parsers.InitParser(c.FlagSet.Arg(0), generateArgs)
		if err != nil {
			return fmt.Errorf("Nie udało się zainicjować parsera wejścia: %v\n", err)
		}

		saftDoc := &saft.SAFT{}

		if err := parser.Parse(saftDoc); err != nil {
			log.Errorf("błąd parsowania: %v", err)
			return err
		}

		log.Debugf("Koniec parsowania")
		log.Debugf("Zapis do pliku: %s", parser.SAFTFileName())

		saftDoc.Save(parser.SAFTFileName())

		if generateArgs.GenerateMetadata {
			saftMeta := saft.Metadata
			saftMeta.SaftFilePath = parser.SAFTFileName()

			log.Debugf("Generowanie metadanych do podpisu")
			if err := saftMeta.Save(); err != nil {
				return fmt.Errorf("nie udało się wygenerować pliku metadanych: %v", err)
			}
		}
	}
	return nil
}
