package utils

import (
	"errors"
	"time"

	log "github.com/sirupsen/logrus"
)

var dateFormats = []string{
	time.DateOnly,
	"02/01/2006",
	"2/01/2006",
	"2/1/2006",
	"02-01-2006",
	"2-1-2006",
}

var errUnableToParseDate = errors.New("błąd parsowania daty")

func ReparseDateField(input string) (result string, err error) {
	for _, format := range dateFormats {
		parsedValue, err := time.Parse(format, input)
		if err == nil {
			result = parsedValue.Format(time.DateOnly)
			if format != time.DateOnly {
				log.Warnf("Wykryto nieprawidłowy format daty. Konwersja wartości %s do %s", input, result)
			}
			return result, nil
		}
	}

	return "", errUnableToParseDate
}
