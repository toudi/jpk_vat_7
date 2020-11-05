package commands

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/toudi/jpk_vat/common"
)

type upoCommand struct {
	Command
}

type upoCommandArgs struct {
	File string
}

var UpoCommand *upoCommand
var UpoArgs = &upoCommandArgs{}
var docRefNo string

// var httpClient http.DefaultClient
var destFileName string

func init() {
	UpoCommand = &upoCommand{
		Command: Command{
			FlagSet:     flag.NewFlagSet("upo", flag.ExitOnError),
			Description: "Pobiera UPO dla podanego identyfikatora (Jedynie dla środowiska produkcyjnego)",
			Run:         upoRun,
			Args:        UpoArgs,
		},
	}

	UpoCommand.FlagSet.StringVar(&UpoArgs.File, "f", "", "Nazwa pliku w którym zapisany jest numer referencyjny")
	UpoCommand.FlagSet.SetOutput(os.Stdout)
}

func upoRun(c *Command) error {
	docRefNo = c.FlagSet.Arg(0)
	if docRefNo == "" {
		// jeśli nie podano numeru jako parametr programu, sprawdź czy podano nazwę pliku z numerem
		if UpoArgs.File == "" {
			// jeśli nie, oznacza to błąd parametrów wejściowych
			c.FlagSet.Usage()
			return nil
		}
		docRefNoBytes, err := ioutil.ReadFile(UpoArgs.File)
		if err != nil {
			fmt.Printf("Błąd odczytu pliku z numerem dokumentu: %v%s", err, common.LineBreak)
			return nil
		}
		docRefNo = strings.TrimSpace(string(docRefNoBytes))
		fileInfo, err := os.Stat(UpoArgs.File)
		if err != nil {
			fmt.Printf("Błąd pobierania informacji o pliku z numerem referencyjnym: %v%s", err, common.LineBreak)
			return nil
		}
		destFileName = "upo_" + strings.TrimSuffix(fileInfo.Name(), filepath.Ext(fileInfo.Name())) + ".pdf"

	} else {
		// podano numer referencyjny jako parametr.
		destFileName = "upo_" + docRefNo + ".pdf"
	}
	// jeśli jesteśmy tutaj, to znaczy, że ustalono numer referencyjny.

	upoDownloadReq, err := http.NewRequest("GET", fmt.Sprintf(UpoDownloadURL, docRefNo), nil)
	if err != nil {
		return fmt.Errorf("Nie udało się zainicjować pobierania UPO")
	}
	upoDownloadResponse, err := httpClient.Do(upoDownloadReq)
	if err != nil {
		return fmt.Errorf("Nie udało się pobrać UPO")
	}
	defer upoDownloadResponse.Body.Close()
	upoContent, err := ioutil.ReadAll(upoDownloadResponse.Body)
	if err != nil {
		return fmt.Errorf("Nie udało się odczytać bajtów UPO z odpowiedzi HTTP")
	}
	if err = ioutil.WriteFile(destFileName, upoContent, 0644); err != nil {
		return fmt.Errorf("Nie udało się zapisać UPO na dysk")
	}

	return nil
}
