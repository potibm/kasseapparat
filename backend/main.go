package main

import (
	"log"
	"os"

	"github.com/potibm/kasseapparat/cmd"
)

func startup() int {
	if err := cmd.Execute(); err != nil {
		log.Printf("Fatal error while starting: %v", err)

		return int(cmd.ExitSoftware)
	}

	return int(cmd.ExitOK)
}

func main() {
	os.Exit(startup())
}
