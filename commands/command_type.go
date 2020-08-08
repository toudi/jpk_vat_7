package commands

import "flag"

type CommandCallable = func(c *Command) error

type Command struct {
	FlagSet     *flag.FlagSet
	Description string
	Args        interface{}
	Run         CommandCallable
}
