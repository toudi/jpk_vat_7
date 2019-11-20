package main

import (
	"flag"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

var verbose bool = false

func main() {
	flag.BoolVar(&verbose, "v", false, "tryb verbose (zwiększa ilość informacji na wyjściu)")
	flag.Parse()

	log.SetLevel(log.InfoLevel)
	if verbose {
		log.SetLevel(log.DebugLevel)
	}
	log.Debugf("jpk-gen:: start programu")

	if len(flag.Args()) < 1 {
		log.Errorf("Nieprawidłowe wywołanie programu. jpk_gen plik-lub-katalog\n")
		os.Exit(-1)
	}

	jpk := &JPK{
		dataWytworzenia: time.Now(),
		naglowek: Naglowek{
			kodSystemowy:       "JPK_V7M (1)",
			wersjaSchemy:       "1-0",
			kodFormularza:      "JPK_VAT",
			wariantFormularza:  "1",
			nazwaSystemu:       "WSI PEGASUS",
			celZlozenia:        "1",
			celZlozeniaPozycja: "P_7",
			kodUrzedu:          "",
			rok:                "0",
			miesiac:            "0",
		},

		deklaracja: formularzVAT7{
			kod:                "VAT-7",
			kodSystemowy:       "VAT-7 (21)",
			kodPodatku:         "VAT",
			rodzajZobowiazania: "Z",
			wersjaSchemy:       "1-0E",
			wariantFormularza:  "21",
			pozycjeSzczegolowe: make(map[string]string),
		},

		podmiot: Podmiot{
			osobaFizyczna: true,
		},
	}

	fileName := flag.Args()[0]

	fileInfo, err := os.Stat(fileName)

	if fileInfo.IsDir() {
		err = jpk.parsujKatalog(fileName)
	} else {
		err = jpk.parsujCSV(fileName)
	}

	if err != nil {
		log.Errorf("Błąd parsowania: %v", err)
	} else {
		log.Infof("Parsowanie zakończone sukcesem")
		err = jpk.zapiszDoPliku(fileInfo, fileName)
	}
	if err == nil {
		log.Infof("Zapis do pliku pomyślny")
	} else {
		log.Errorf("Błąd zapisu do pliku: %v", err)
	}
}
