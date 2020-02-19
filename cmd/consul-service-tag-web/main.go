package main

import (
	"consul-service-tag-web/cmd/consul-service-tag-web/commands"
	"os"
)

func main() {

	if err := commands.Run(os.Args[1:]); err != nil {
		os.Exit(1)
	}

}
