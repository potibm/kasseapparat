package initializer

import (
	"log"
	"os"
)

var version string = "0.0.0"

func InitializeVersion() {
	versionFilePath := "./VERSION"

	content, err := os.ReadFile(versionFilePath)
	if err != nil {
		log.Printf("Fehler beim Lesen der Versionsdatei: %v", err)
		return
	}

	version = string(content)
}

func GetVersion() string {
	return version
}

func OutputVersion() {
	log.Printf("Kasseapparat %s\n", version)
}