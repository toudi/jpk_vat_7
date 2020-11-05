package converter

import (
	"fmt"
	"strings"
)

const (
	// rozpoznawanie sekcji
	StateDetectSection = iota
	// znamy sekcje, musimy odnaleźć kolumny
	StateDetectHeaders = iota
	// parsowanie danych w obrębie sekcji.
	// w tym trybie parsowania dane zawsze zaczynają się od kolumny 0.
	StateParseData = iota
)

var sekcjaParsera *SekcjaParsera
var parserState int = StateDetectSection
var headers []string

func (p *Parser) parseLineSingleFileWithSections(line []string) error {
	var exists bool
	if line[0] == "SEKCJA" {
		fmt.Printf("parser state: %v (%v)\n", parserState, StateParseData)
		if parserState == StateParseData {
			// należy zakończyć parsowanie danych z poprzedniej sekcji.
			sekcjaParsera.finish(sekcjaParsera)
		}
		// próba odnalezienia sekcji
		sekcjaParsera, exists = sekcje[line[1]]
		if !exists {
			return fmt.Errorf("Błąd: Nieznana sekcja: %s", line[1])
		}
		parserState = StateDetectHeaders
	} else if line[0] == "" {
		// pusta linia, ignorujemy
		return nil
	} else if parserState == StateDetectHeaders {
		// parsujemy nagłówki
		headers = line
		sekcjaParsera.SetHeaders(line)
		parserState = StateParseData
	} else if parserState == StateParseData {
		// parsujemy dane
		for colIdx, colData := range line {
			// w tym trybie pracy każda sekcja ma inną ilość kolumn więc jeśli dana sekcja
			// jest krótsza niż najdłuższa w pliku nie ma sensu jej parsować.
			if headers[colIdx] != "" {
				sekcjaParsera.SetData(headers[colIdx], strings.TrimSpace(colData))
			}
		}
	}
	return nil
}

func (p *Parser) finish() {
	sekcjaParsera.finish(sekcjaParsera)
}
