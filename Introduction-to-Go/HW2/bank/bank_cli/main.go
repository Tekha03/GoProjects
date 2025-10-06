package main

import "bank_cli/cli"

func main() {
	command := &cli.Command{}
	command.Init()
	command.Execute()
}
