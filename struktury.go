package main

import "time"

type formularzVAT7 struct {
	kod                string
	kodSystemowy       string
	kodPodatku         string
	rodzajZobowiazania string
	wersjaSchemy       string
	wariantFormularza  string
	pouczenia          string
	// jesli pole ma wartosc 1 użytkownik oznacza, że jest świadomy
	// tego, że wysyłany formularz jest podstawą do roszczeń prawnych
	// przez skarb państwa
	p_ordzu string
	// Uzasadnienie przyczyn złożenia korekty
	pozycjeSzczegolowe map[string]string
}

type Naglowek struct {
	kodSystemowy       string
	wersjaSchemy       string
	kodFormularza      string
	wariantFormularza  string
	nazwaSystemu       string
	celZlozenia        string
	celZlozeniaPozycja string
	kodUrzedu          string
	rok                string
	miesiac            string
}

type Podmiot struct {
	osobaFizyczna bool
	NIP           string
	imie          string
	nazwisko      string
	dataUrodzenia string
	email         string
	nazwa         string //tylko dla osoby niefizycznej.
}

type Sprzedaz struct {
	lpSprzedazy string
}

type JPK struct {
	dataWytworzenia time.Time

	//
	naglowek   Naglowek
	deklaracja formularzVAT7
	podmiot    Podmiot
	sprzedaz   []Sprzedaz
}
