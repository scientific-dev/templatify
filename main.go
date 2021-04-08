package main

import (
	"fmt"
	"os"
)

var (
	command  Command
	Commands map[string]func() = map[string]func(){
		"save":        Save,
		"info":        Info,
		"all":         All,
		"list":        List,
		"remove":      Remove,
		"removeall":   RemoveAll,
		"use":         Use,
		"help":        Help,
		"pre-scripts": RunPreScripts,
		"get":         Download,
		"init":        MakeInit,
	}
)

func main() {
	args := os.Args[1:]

	if len(args) < 1 {
		fmt.Fprintln(stdout, "[31mERROR[0m No command has been provided to execute.")
		os.Exit(0)
	}

	call, exists := Commands[args[0]]
	command = CreateCommand()

	if exists {
		call()
	} else {
		UnknownCommand(args[0])
	}
}
