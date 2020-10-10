package commands

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/toudi/jpk_vat/common"
)

type statusCommand struct {
	Command
}

var StatusCmd *statusCommand

type statusResponseType struct {
	Code        int
	Description string
	Details     string
	UPO         string `json:"Upo"`
}

var statusResponse statusResponseType
var httpClient = http.DefaultClient

func init() {
	StatusCmd = &statusCommand{
		Command: Command{
			FlagSet:     flag.NewFlagSet("status", flag.ExitOnError),
			Description: "Sprawdza status weryfikacji pliku JPK oraz pobiera UPO w przypadku sukcesu.",
			Run:         statusRun,
			Args:        nil,
		},
	}
	StatusCmd.FlagSet.Usage = func() {
		fmt.Printf("Użycie komendy: jpk_vat status plik-z-rozszerzeniem.ref\n")
	}
}

func statusRun(c *Command) error {
	var err error
	// var workdir string

	if c.FlagSet.NArg() == 0 {
		GenerateCmd.FlagSet.Usage()
	} else {
		refFileName := c.FlagSet.Arg(0)
		// workdir = path.Dir(refFileName)
		URLStatusBytes, err := ioutil.ReadFile(refFileName)
		if err != nil {
			return fmt.Errorf("Nie udało się odczytać pliku z numerem referencyjnym")
		}
		request, err := http.NewRequest("GET", string(URLStatusBytes), nil)
		if err != nil {
			return fmt.Errorf("Nie udało się zainicjować sprawdzania statusu")
		}
		response, err := httpClient.Do(request)
		if err != nil {
			return fmt.Errorf("Nie udało się uzyskać statusu: %v", err)
		}
		defer response.Body.Close()
		if err = json.NewDecoder(response.Body).Decode(&statusResponse); err != nil {
			return fmt.Errorf("Nie udało się zdekodować odpowiedzi status: %v", err)
		}
		fmt.Printf("Status przetwarzania:\n")
		fmt.Printf(
			"Kod odpowiedzi: %d\nOpis: %s\nInformacje szczegółowe: %s\n",
			statusResponse.Code, statusResponse.Description, statusResponse.Details,
		)

		if response.StatusCode == 200 {
			UPOFileName := strings.Replace(refFileName, ".ref", ".upo", 1)
			if !common.FileExists(UPOFileName) {
				fmt.Printf("Status przesyłania JPK poprawny; pobieram UPO\n")

				if err = ioutil.WriteFile(UPOFileName, []byte(statusResponse.UPO), 0644); err != nil {
					return fmt.Errorf("Nie udało się zapisać UPO: %v", err)
				}
			}

		}
	}

	return err
}
