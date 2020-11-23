package commands

import (
	"flag"

	"github.com/toudi/jpk_vat_7/uploader"

	log "github.com/sirupsen/logrus"
)

type uploadArgsType struct {
	UseTestGateway bool
	Verbose        bool
}

type uploadCommand struct {
	Command
	logger *log.Logger
}

var uploadArgs = &uploadArgsType{}
var UploadCmd *uploadCommand

func init() {
	UploadCmd = &uploadCommand{
		Command: Command{
			FlagSet:     flag.NewFlagSet("wyslij", flag.ExitOnError),
			Description: "Wysyła podpisany plik JPK na serwer ministerstwa",
			Run:         uploadRun,
			Args:        uploadArgs,
		},
		logger: log.New(),
	}
	UploadCmd.FlagSet.BoolVar(&uploadArgs.UseTestGateway, "t", false, "użyj testowego środowiska")
	UploadCmd.FlagSet.BoolVar(&uploadArgs.Verbose, "v", false, "tryb verbose (zwiększa ilość informacji na wyjściu)")
}

func uploadRun(c *Command) error {
	fileName := c.FlagSet.Arg(0)
	uploader := uploader.UploaderInit(fileName, uploadArgs.Verbose)
	uploader.UseTestGateway = uploadArgs.UseTestGateway

	return uploader.UploadSAFTFile()
}
