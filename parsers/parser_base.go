package parsers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/toudi/jpk_vat_7/common"
	"github.com/toudi/jpk_vat_7/saft"
)

const SectionHeader = "SEKCJA"

type SAFTSection struct {
	Id       string
	ColStart int
	ColEnd   int
}

// typ Parser zawiera podstawową implementację parsera JPK. Chodzi o to, że
// zarówno CSV jak i xlsx możemy traktować jako system zapisu danych tabelarycznych
// a co za tym idzie - zawierających sekcje JPK. Tak więc znakomitą część logiki
// możemy zawrzeć w tej funkcji, żeby nie powtarzać kodu.
//
// właściwy parser CSV oraz XLSX będzie tylko nieco inaczej interpretować otwieranie
// zamykanie pliku oraz pobieranie wartośći komórek
type BaseParser struct {
	// w przypadku pliku CSV lub xlsx będzie to ścieżka do pliku.
	// w przypadku parsera katalogu będzie to nazwa katalogu do sparsowania.
	Source string
	// opcje generatora. Są one wspólne dla wszystkich parserów więc dany parser
	// musi sobie z nich wyciągnąć to co mmu potrzebne.
	Options *common.GeneratorOptions

	// mapowanie kolumny na jej nazwe. Przydatne później ponieważ chcemy przekazać
	// mapę kolumna => wartość do obiektu JPK.
	headerIndex map[int]string
	// sekcje które udało się rozpoznać w danym pliku.
	saftSections []*SAFTSection

	// konwersja znaków w danych wejściowych - w zasadzie używana tylko przy CSV
	// aczkolwiek gdyby kiedyś chciało mi się doimplementować obsługę plików
	// xls (a nie xlsx) to również będzie można zastosować :-)
	encodingConversion map[byte]string
}

// parseSAFTSections zachowuje się identycznie dla każdego sposobu zapisu
// - zarówno CSV jak i xlsx i próbuje odgadnąć jakie sekcje udało się
// rozpoznać w pliku wejściowym.
func (b *BaseParser) parseSAFTSections(line []string, dst *saft.SAFT) {
	b.headerIndex = make(map[int]string)

	var section *SAFTSection = nil
	var lastSection *SAFTSection = nil

	// ponieważ iterujemy po sekcjach kilkukrotnie musimy wiedzieć, czy
	// już jakąś odwiedziliśmy
	var processedSections map[string]bool = make(map[string]bool)
	b.saftSections = make([]*SAFTSection, 0)

	for colIdx, column := range line {
		b.headerIndex[colIdx] = column

		for _, _section := range saft.SAFTSections {
			if _, exists := processedSections[_section.Id]; exists {
				continue
			}
			if column == _section.StartCol {
				section = &SAFTSection{Id: _section.Id, ColStart: colIdx}
				if len(b.saftSections) > 0 {
					lastSection = b.saftSections[len(b.saftSections)-1]
					lastSection.ColEnd = colIdx
					dst.SetFieldOrder(lastSection.Id, line[lastSection.ColStart:lastSection.ColEnd])
				}
				b.saftSections = append(b.saftSections, section)
				processedSections[_section.Id] = true
			}
		}
	}

	if len(b.saftSections) > 0 {
		lastSection = b.saftSections[len(b.saftSections)-1]
		lastSection.ColEnd = len(line)
		dst.SetFieldOrder(lastSection.Id, line[lastSection.ColStart:lastSection.ColEnd])
	}
}

// processLine przetwarza linię z pliku CSV lub xlsx i dodaje dane z rozpoznanych
// sekcji do obiektu JPK *dst.
//
// parser który rozszerza BaseParser zajmuje się przygotowaniem linii dzięki czemu
// wszystko co związane z kodowaniem itp zaimplementujemy w CSVParser
func (b *BaseParser) processLine(line []string, dst *saft.SAFT) error {
	var dataEmpty bool
	var headerField string
	var err error

	if common.LineIsEmpty(line) {
		return nil
	}
	// odczytana linia nie jest pusta.
	//
	// są 3 możliwości:
	// 1/ sprawdźmy, czy linia zaczyna się od prefiksu SEKCJA oraz w kolumnie 1 ma
	//    jedną z obsługiwanych sekcji. Jeśli tak to kolejna linia będzie mieć
	//    nagłówki
	if line[0] == SectionHeader {
		for _, section := range saft.SAFTSections {
			if line[1] == section.Id {
				// ok, z pewnością jest to rozgraniczenie pliku które informuje
				// parser o tym, że kolejna linia będzie zawierać nagłówek sekcji
				b.headerIndex = nil
				b.saftSections = nil
				// koniec obsługi
				return nil
			}
		}
	}
	// 2/ sprawdźmy, czy mamy zdefiniowane sekcje.
	//    Jeśli nie, to oznacza to, że bieżąca linia zawiera definicję nagłówków.
	if b.saftSections == nil {
		b.parseSAFTSections(line, dst)
		return nil
	}
	// 3/ przypadek ostatni czyli ani linia nie zawiera prefiksu ani nie jest to
	//    linia z nagłówkiem - czyli po prostu jest to linia danych.
	for _, section := range b.saftSections {
		data := make(saft.SAFTData)
		dataEmpty = true
		for idx := section.ColStart; idx < section.ColEnd; idx++ {
			headerField = b.headerIndex[idx]
			if headerField == "" {
				continue
			}
			data[headerField] = line[idx]
			if line[idx] != "" {
				data[headerField] = b.convertEncoding(data[headerField])
				data[headerField] = b.addCDataIfNecesary(data[headerField])
				dataEmpty = false
			}
		}
		if !dataEmpty {
			if err = dst.AddData(section.Id, data); err != nil {
				return fmt.Errorf("nie udało się dodać danych do sekcji %s", section.Id)
			}
		}
	}

	return nil
}

func (b *BaseParser) SAFTFileName() string {
	sourceBaseName := filepath.Base(b.Source)
	outputBaseName := sourceBaseName + "-jpk.xml"

	// opcja generowania pliku w katalogu bieżącym
	if b.Options.UseCurrentDir {
		return outputBaseName
	}

	today := time.Now()
	// opcja generowania w katalogu z sesjami.
	return filepath.Join(
		"sesje",
		strconv.Itoa(today.Year()),
		fmt.Sprintf("%02d", today.Month()),
		outputBaseName,
	)
}

func (b *BaseParser) convertEncoding(data string) string {
	if b.Options.EncodingConversionFile == "" {
		return data
	}

	if b.encodingConversion == nil {
		b.prepareEncodingConversionTable()
	}

	inputBytes := []byte(data)
	outputBytes := []byte{}
	var dstByte []byte

	for _, srcByte := range inputBytes {
		dstByte = []byte{srcByte}
		if dstChar, exists := b.encodingConversion[srcByte]; exists {
			dstByte = bytes.Replace(dstByte, []byte{srcByte}, []byte(dstChar), 1)
		}
		outputBytes = append(outputBytes, dstByte...)
	}

	return string(outputBytes)

}

func (b *BaseParser) prepareEncodingConversionTable() {
	fileBytes, err := ioutil.ReadFile(b.Options.EncodingConversionFile)

	if err != nil {
		log.Errorf("Nie udało się otworzyć pliku z konwersją znaków")
		return
	}

	b.encodingConversion = make(map[byte]string)

	for _, line := range strings.Split(string(fileBytes), common.LineBreak) {
		mapping := strings.Split(line, ":")

		if len(mapping) == 2 {
			srcByteHex := strings.Trim(mapping[0], " ")
			dstChar := strings.Trim(mapping[1], " \r")

			srcByte, err := strconv.ParseUint(srcByteHex, 0, 8)
			if err == nil {
				b.encodingConversion[byte(srcByte)] = dstChar
			}
		}
	}

	log.Debugf("Odczytano tablicę konwersji znaków: %v\n", b.encodingConversion)
}

func (b *BaseParser) addCDataIfNecesary(value string) string {
	if strings.Contains(value, "&") {
		return fmt.Sprintf("<![CDATA[%s]]>", value)
	}

	return value
}
