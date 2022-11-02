module github.com/toudi/jpk_vat_7

go 1.15

require (
	github.com/sirupsen/logrus v1.7.0
	github.com/tealeg/xlsx/v3 v3.2.0
	github.com/toudi/jpk_vat_7/commands v0.0.0
	github.com/toudi/jpk_vat_7/common v0.0.0
	golang.org/x/sys v0.1.0 // indirect
)

replace github.com/toudi/jpk_vat_7/commands => ./commands

replace github.com/toudi/jpk_vat_7/common => ./common
