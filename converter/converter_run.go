package converter

import (
	"fmt"
	"os"
	"time"
)

var jpk *JPK
var err error

func (c *Converter) Run() error {

	logger.Debugf("Rozpoczynam konwersję")

	statInfo, err := os.Stat(c.source)
	if err != nil {
		return fmt.Errorf("Nie można odczytać informacji o pliku/katalogu: %v", err)
	}

	jpk = &JPK{
		created: time.Now(),
	}

	jpk.inicjalizujSekcje()
	jpk.podmiot.osobaFizyczna.namespace = "etd"
	jpk.podmiot.osobaFizyczna.namespacePol = make(map[string]string)
	jpk.podmiot.osobaFizyczna.namespacePol["Email"] = "tns"
	jpk.podmiot.osobaFizyczna.namespacePol["Telefon"] = "tns"

	if statInfo.IsDir() {
		err = c.convertDirectory()
	} else {
		err = c.convertFile()
	}

	if err != nil {
		return fmt.Errorf("Błąd parsowania: %v", err)
	}

	err, output := jpk.zapiszDoPliku(statInfo, c.source, c.GeneratorOptions.UseCurrentDir)

	if err != nil {
		return fmt.Errorf("Nie udało się zapisać pliku JPK : %v", err)
	}

	c.SAFTFile = output

	if c.GeneratorOptions.GenerateMetadata {
		MetadataGeneratorState.TemplateVars.Metadata.SchemaVersion = jpk.naglowek.atrybuty["KodFormularza.wersjaSchemy"]
		MetadataGeneratorState.TemplateVars.Metadata.SystemCode = jpk.naglowek.atrybuty["KodFormularza.kodSystemowy"]

		generator, err := MetadataGeneratorInit()
		if err != nil {
			return fmt.Errorf("Nie udało się zainicjalizować generatora metadanych JPK: %v", err)
		}
		generator.GenerateMetadata(output)

		fmt.Printf("Zapis do pliku zakończony; Plik do podpisu: %s\n", saftMetadataFileName(output))
		fmt.Printf("Aby podpisać plik:\n- użyj czytnika z podpisem kwalifikowanym\n- wejdź na stronę: https://www.gov.pl/web/gov/podpisz-jpkvat-profilem-zaufanym (link znajdziesz w katalogu z plikiem źródłowym)\n")
	} else {
		fmt.Printf("Konwersja zakończona.\nNie wybrano opcji generowania metadanych.\nProszę skorzystać z poniższego narzędzia aby dokończyć wysyłanie pliku:\nhttps://e-mikrofirma.mf.gov.pl/jpk-client\n")
	}

	return nil
}
