package main

import (
	"log"
	"os"

	"github.com/potibm/kasseapparat/cmd"
	"github.com/potibm/kasseapparat/internal/app/exitcode"
)

func startup() int {
	if err := cmd.Execute(); err != nil {
		log.Printf("Fatal error while starting: %v", err)

		return int(exitcode.Software)
	}

	return int(exitcode.OK)
}

func main() {
	os.Exit(startup())
}
