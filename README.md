jpk vat z deklaracją
====================
Klient JPK do przygotowywania plików JPK VAT z deklaracją (JPK_V7M) na podstawie plików CSV oraz ich opcjonalnej wysyłki.

Więcej informacji w wiki: https://github.com/toudi/jpk_vat_7/wiki

wersje binarne dla linuksa i windows znajdziesz tutaj: https://github.com/toudi/jpk_vat_7/releases

kompilacja ze źródeł:
---------------------

```
go get github.com/toudi/jpk_vat_7
```

w tym momencie w katalogu $GOPATH/bin otrzymasz binarkę jpk_vat_7 którą możesz przenieść gdziekolwiek

program można też skompilować poza GOPATH, w tym celu pobierz archiwum ze źródłami, rozpakuj w dowolne miejsce, otwórz terminal w katalogu gdzie rozpakowałeś źródła i stwórz pliki go.mod:

```
go mod init github.com/toudi/jpk_vat_7
cd common
go mod init github.com/toudi/jpk_vat_7/common
cd ../commands
go mod init github.com/toudi/jpk_vat_7/commands
cd ../converter
go mod init github.com/toudi/jpk_vat_7/converter
```

następnie wyedytuj go.mod w głównym katalogu projektu aby był w postaci następującej:

```
module github.com/toudi/jpk_vat_7

go 1.15

require (
	github.com/toudi/jpk_vat_7/commands v0.0.0
	github.com/toudi/jpk_vat_7/common v0.0.0
	github.com/toudi/jpk_vat_7/converter v0.0.0
)

replace github.com/toudi/jpk_vat_7/commands => ./commands

replace github.com/toudi/jpk_vat_7/common => ./common

replace github.com/toudi/jpk_vat_7/converter => ./converter

```

pozostałe pliki możesz pozostawić bez zmian.

wówczas pozostaje tylko kompilacja:

```
go build
```