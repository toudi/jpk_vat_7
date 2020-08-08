package commands

import (
	"flag"
	"fmt"

	log "github.com/sirupsen/logrus"
)

type uploadArgsType struct {
	UseTestGateway bool
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
}

func uploadRun(c *Command) error {
	args := c.Args.(*uploadArgsType)
	fmt.Printf("test=%v\n", args.UseTestGateway)
	return nil
}
