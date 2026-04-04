package main

import (
	"log"
	"os"

	"github.com/potibm/kasseapparat/cmd"
	"github.com/potibm/kasseapparat/internal/app/exitcode"
)

var (
	version = "0.0.0"
)

func startup() int {
	if err := cmd.Execute(); err != nil {
		log.Fatalf("Fatal error while starting: %v", err)

		return int(exitcode.Software)
	}

	return int(exitcode.OK)
}

func main() {
	os.Exit(startup())
}
