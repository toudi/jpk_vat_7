package commands

import (
	"bytes"
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
var UpoDownloadURL = "https://e-mikrofirma.mf.gov.pl/jpk-client/api/status/%s/pdf"

type statusResponseType struct {
	Code        int
	Description string
	Details     string
	UPO         string `json:"Upo"`
}

var statusResponse statusResponseType
var httpClient = http.DefaultClient
var downloadPDF = false
var refNo string

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
		if bytes.HasPrefix(URLStatusBytes, []byte(common.ProductionGatewayURL)) {
			// środowisko testowe wysyła UPO jedynie w XML
			downloadPDF = true
			refNoParts := strings.Split(string(URLStatusBytes), "/")
			refNo = refNoParts[len(refNoParts)-1]
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
			UPOFileName := strings.Replace(refFileName, ".ref", "-upo.xml", 1)
			if downloadPDF {
				UPOFileName = strings.Replace(refFileName, ".ref", "-upo.pdf", 1)
			}
			if !common.FileExists(UPOFileName) {
				fmt.Printf("Status przesyłania JPK poprawny; pobieram UPO\n")

				if downloadPDF {
					upoDownloadReq, err := http.NewRequest("GET", fmt.Sprintf(UpoDownloadURL, refNo), nil)
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
					if err = ioutil.WriteFile(UPOFileName, upoContent, 0644); err != nil {
						return fmt.Errorf("Nie udało się zapisać UPO na dysk")
					}
				} else {
					if err = ioutil.WriteFile(UPOFileName, []byte(statusResponse.UPO), 0644); err != nil {
						return fmt.Errorf("Nie udało się zapisać UPO: %v", err)
					}
				}
			}

		}
	}

	return err
}
