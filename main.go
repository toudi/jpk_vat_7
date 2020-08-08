package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/toudi/jpk_vat/commands"
)

var err error
var availableCommands = []commands.Command{commands.GenerateCmd.Command, commands.UploadCmd.Command}

func usage() {
	fmt.Printf("Proszę użyć jednej z subkomend:\n\n")
	for _, cmd := range availableCommands {
		fmt.Printf("%-10s %s\n", cmd.FlagSet.Name(), cmd.Description)
	}
	fmt.Printf("\n")
	os.Exit(1)
}

func main() {
	log.Info("jpk_vat:: start programu")

	if len(os.Args) < 2 {
		usage()
	}

	for _, cmd := range availableCommands {
		if os.Args[1] == cmd.FlagSet.Name() {
			cmd.FlagSet.Parse(os.Args[2:])
			if err = cmd.Run(&cmd); err != nil {
				log.Errorf("Błąd wykonywania komendy %s: %v\n", os.Args[1], err)
				os.Exit(-1)
			}
			os.Exit(0)
		}
	}

	usage()
}
