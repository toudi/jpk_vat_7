module github.com/toudi/jpk_vat_7

go 1.15

require (
	aqwari.net/xml v0.0.0-20200724195937-ae380bb65a55 // indirect
	github.com/labstack/gommon v0.3.0 // indirect
	github.com/sirupsen/logrus v1.6.0
	github.com/toudi/jpk_vat_7/commands v0.0.0
	github.com/toudi/jpk_vat_7/common v0.0.0
	github.com/toudi/jpk_vat_7/converter v0.0.0 // indirect
	golang.org/x/net v0.0.0-20200707034311-ab3426394381 // indirect
	golang.org/x/tools v0.0.0-20200808161706-5bf02b21f123 // indirect
)

replace github.com/toudi/jpk_vat_7/commands => ./commands

replace github.com/toudi/jpk_vat_7/common => ./common

replace github.com/toudi/jpk_vat_7/converter => ./converter
