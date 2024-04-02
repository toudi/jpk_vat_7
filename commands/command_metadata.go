package commands

import (
	"flag"
	"fmt"

	"github.com/toudi/jpk_vat_7/saft"
)

type metadataArgsType struct {
	Verbose    bool
	SourceFile string
}

type metadataCommand struct {
	Command
}

var MetadataCmd *metadataCommand
var metadataArgs = &metadataArgsType{}

func handleMetadataArgs(flagSet *flag.FlagSet) {
	flagSet.BoolVar(
		&saft.Metadata.UseTestGateway,
		"t",
		false,
		"użycie bramki testowej do generowania metadanych",
	)
	flagSet.BoolVar(
		&saft.Metadata.AuthData.Enable,
		"a",
		false,
		"wygeneruj strukturę AuthData (alternatywa dla podpisu kwalifikowanego)",
	)
	flagSet.Float64Var(
		&saft.Metadata.AuthData.Income,
		"a:i",
		0.0,
		"użyj autoryzacji za pomocą kwoty przychodu",
	)
	flagSet.StringVar(&saft.Metadata.AuthData.NIP, "a:n", "", "numer NIP dla autoryzacji")
	flagSet.StringVar(&saft.Metadata.AuthData.PESEL, "a:p", "", "numer PESEL dla autoryzacji")
	flagSet.StringVar(
		&saft.Metadata.AuthData.ImiePierwsze,
		"a:fn",
		"",
		"pole ImiePierwsze dla autoryzacji",
	)
	flagSet.StringVar(&saft.Metadata.AuthData.Nazwisko, "a:ln", "", "pole Nazwisko dla autoryzacji")
	flagSet.StringVar(
		&saft.Metadata.AuthData.DataUrodzenia,
		"a:bd",
		"",
		"pole DataUrodzenia dla autoryzacji. Format: YYYY-MM-DD",
	)
	flagSet.StringVar(
		&saft.Metadata.TemplateVars.Metadata.SchemaVersion,
		"m:sv",
		"",
		"atrybut schemaVersion w nagłowku metadanych. Jeśli nie zostanie podany, wartość zostanie wyciągnięta z pliku źródłowego",
	)
	flagSet.StringVar(
		&saft.Metadata.TemplateVars.Metadata.SystemCode,
		"m:sc",
		"",
		"atrybut systemCode w nagłowku metadanych. Jeśli nie zostanie podany, wartość zostanie wyciągnięta z pliku źródłowego",
	)
	flagSet.StringVar(
		&saft.Metadata.TemplateVars.Metadata.FormCode,
		"m:fc",
		"",
		"wartość pola FormCode w nagłowku metadanych. Jeśli nie zostanie podana, wartość zostanie wyciągnięta z pliku źródłowego",
	)
}

func init() {
	MetadataCmd = &metadataCommand{
		Command: Command{
			FlagSet:     flag.NewFlagSet("metadane", flag.ExitOnError),
			Description: "Generuje metadane dla wcześniej wygenerowanego pliku JPK",
			Run:         metadataRun,
			Args:        metadataArgs,
		},
	}

	MetadataCmd.FlagSet.BoolVar(
		&metadataArgs.Verbose,
		"v",
		false,
		"tyb verbose (zwiększa poziom komunikatów)",
	)
	// MetadataCmd.FlagSet.BoolVar(&converter.MetadataGeneratorState.UseTestGateway, "t", false, "użycie bramki testowej do generowania metadanych")

	handleMetadataArgs(MetadataCmd.FlagSet)
}

func metadataRun(c *Command) error {
	// var err error

	if len(c.FlagSet.Args()) == 0 {
		c.FlagSet.Usage()
		return fmt.Errorf("błędne wywołanie komendy")
	}

	saftMeta := saft.Metadata
	saftMeta.SaftFilePath = c.FlagSet.Arg(0)

	return saftMeta.Save()
}
