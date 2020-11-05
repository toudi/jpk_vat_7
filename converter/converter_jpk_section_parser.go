package converter

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
)

// typ definiuje parser struktury zapisanej w CSV.
type Parser struct {
	// potrzebne do inicjalizacji
	file      *os.File
	csvReader *csv.Reader
	// naglowki służą do sprawdzenia czy struktura się zgadza.
	naglowki []string
	// sekcje definiują nam co będziemy parsować.
	sekcje []*SekcjaParsera
	mode   int
}

func parser(filePath string, sekcje []*SekcjaParsera, delimiter string) error {
	var err error
	p, err := parserInit(filePath, delimiter)
	if err != nil {
		return err
	}
	defer p.Close()
	p.sekcje = sekcje
	return p.parsuj()
}

func parserInit(filePath string, delimiter string) (*Parser, error) {
	p := &Parser{mode: ParserModeSingleFile}
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	p.file = file
	reader := csv.NewReader(bufio.NewReader(p.file))
	reader.Comma = []rune(delimiter)[0]
	p.csvReader = reader

	return p, nil
}

func (p *Parser) parsuj() error {
	var nrLinii int = -1
	// var err error

	for {
		line, err := p.csvReader.Read()
		logger.Debugf("Odczytano rekord: %+v. Ilość pól: %d\n", line, len(line))
		if err == io.EOF {
			break
		} else if err != nil {
			return fmt.Errorf("Błąd odczytu CSV: %v; odczytany rekord: %v o długości: %d", err, line, len(line))
		}

		nrLinii++
		if nrLinii == 0 {
			// spróbujmy wykryć typ pliku.
			if line[0] == "SEKCJA" {
				logger.Debugf("Wykryto plik CSV z wieloma sekcjami w pojedynczym pliku; zmiana trybu parsera")
				p.mode = ParserModeSingleFileWithSections
				p.parseLineSingleFileWithSections(line)
			}

			if p.mode == ParserModeSingleFile {
				p.parseSAFTSections(line)
			}
		} else {
			if p.mode == ParserModeSingleFileWithSections {
				if err = p.parseLineSingleFileWithSections(line); err != nil {
					return fmt.Errorf("Błąd podczas parsowania pliku: %v", err)
				}
			} else {
				p.parseLineSingleFile(line)
			}
		}
	}
	return nil
}

func (p *Parser) Close() {
	p.file.Close()
}

func (j *JPK) parsujCSV(fileName string, delimiter string) error {
	return parser(fileName, []*SekcjaParsera{
		sekcjaNaglowek,
		sekcjaDeklaracjaNaglowek,
		sekcjaPodmiot,
		sekcjaSprzedaz,
		sekcjaSprzedazCtrl,
		sekcjaZakup,
		sekcjaZakupCtrl,
		sekcjaDeklaracjaPozycje,
		sekcjaDeklaracja,
	}, delimiter)
}
