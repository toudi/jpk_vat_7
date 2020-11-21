package commands

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/toudi/jpk_vat/common"
)

type upoCommand struct {
	Command
}

type upoCommandArgs struct {
	File string
	// szablon dla nazywania UPO
	Template string
}

var UpoCommand *upoCommand
var UpoArgs = &upoCommandArgs{}
var docRefNo string

// var httpClient http.DefaultClient
var destFileName string
var destFileNameBuffer *bytes.Buffer

var destFileTemplateVars struct {
	File string
}

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
	UpoCommand.FlagSet.StringVar(&UpoArgs.Template, "t", "{{.File}}_UPO.pdf", "Szablon dla zapisywania UPO w PDF. {{.File}} zostanie zastąpione przez nazwę oryginalnego pliku")
	UpoCommand.FlagSet.SetOutput(os.Stdout)
}

func upoRun(c *Command) error {
	docRefNo = c.FlagSet.Arg(0)
	destUPOTEmplate, err := template.New("destUPOFilename").Parse(UpoArgs.Template)
	if err != nil {
		return fmt.Errorf("Nie udało się sparsować szablonu dla nazwy wynikowego UPO: %v", err)
	}
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
		destFileTemplateVars.File = strings.TrimSuffix(fileInfo.Name(), filepath.Ext(fileInfo.Name()))

	} else {
		// podano numer referencyjny jako parametr.
		destFileTemplateVars.File = docRefNo
	}
	// jeśli jesteśmy tutaj, to znaczy, że ustalono numer referencyjny.
	destFileNameBuffer = bytes.NewBufferString(destFileName)
	if err := destUPOTEmplate.Execute(destFileNameBuffer, destFileTemplateVars); err != nil {
		return fmt.Errorf("Nie udało się wygenerować nazwy pliku wynikowego UPO")
	}
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
	if err = ioutil.WriteFile(destFileNameBuffer.String(), upoContent, 0644); err != nil {
		return fmt.Errorf("Nie udało się zapisać UPO na dysk")
	}

	return nil
}
