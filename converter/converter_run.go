package converter

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/toudi/jpk_vat/common"
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

	metadataTemplateVars.Metadata.SchemaVersion = jpk.naglowek.atrybuty["KodFormularza.wersjaSchemy"]
	metadataTemplateVars.Metadata.SystemCode = jpk.naglowek.atrybuty["KodFormularza.kodSystemowy"]

	err, output := jpk.zapiszDoPliku(statInfo, c.source, c.GeneratorOptions.UseCurrentDir)

	if err != nil {
		return fmt.Errorf("Nie udało się zapisać pliku JPK : %v", err)
	}

	c.SAFTFile = output

	if c.GeneratorOptions.GenerateMetadata {
		metadataTemplateVars.SourceMetadata.Filename = c.SAFTFileName()
		if metadataTemplateVars.SourceMetadata.Size, err = common.FileSize(output); err != nil {
			return fmt.Errorf("Nie udało się obliczyć rozmiaru pliku jpk: %v", err)
		}
		metadataTemplateVars.SourceMetadata.ContentHash = common.Sha256File(output)

		// pakowanie pliku JPK do archiwum
		if err = c.compressSAFTFile(); err != nil {
			return fmt.Errorf("Nie udało się spakować pliku JPK do archiwum: %v", err)
		}

		// pakowanie pomyślne, możemy dodać metadane.
		metadataTemplateVars.ArchiveMetadata.Filename = path.Base(c.compressedSAFTFile())
		if metadataTemplateVars.ArchiveMetadata.Size, err = common.FileSize(c.compressedSAFTFile()); err != nil {
			return fmt.Errorf("Nie udało się obliczyć rozmiaru archiwum")
		}
		metadataTemplateVars.ArchiveMetadata.ContentHash = common.Md5File(c.compressedSAFTFile())

		if err = c.encryptSAFTFile(); err != nil {
			return fmt.Errorf("Nie udało się zaszyfrować skompresowanego pliku JPK: %v", err)
		}
		metadataTemplateVars.EncryptedMetadata.ContentHash = common.Md5File(c.encryptedArchiveFile())

		if err = c.createSAFTMetadataFile(); err != nil {
			return fmt.Errorf("Nie udało się stworzyć pliku metadanych JPK: %v", err)
		}

		fmt.Printf("Zapis do pliku zakończony; Plik do podpisu: %s\n", c.saftMetadataFile())
		fmt.Printf("Aby podpisać plik:\n- użyj czytnika z podpisem kwalifikowanym\n- wejdź na stronę: https://www.gov.pl/web/gov/podpisz-jpkvat-profilem-zaufanym (link znajdziesz w katalogu z plikiem źródłowym)\n")
	} else {
		fmt.Printf("Konwersja zakończona.\nNie wybrano opcji generowania metadanych.\nProszę skorzystać z poniższego narzędzia aby dokończyć wysyłanie pliku:\nhttps://e-mikrofirma.mf.gov.pl/jpk-client\n")
	}

	return nil
}
