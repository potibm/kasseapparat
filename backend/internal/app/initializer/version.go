package initializer

import (
	"log"
	"os"
	"strings"
)

var version string = "0.0.0"

func InitializeVersion() {
	versionFilePath := "./VERSION"

	content, err := os.ReadFile(versionFilePath)
	if err != nil {
		log.Printf("Error reeading the version file: %v", err)

		return
	}

	version = strings.TrimSpace(string(content))
}

func GetVersion() string {
	return version
}

func OutputVersion() {
	log.Printf("Kasseapparat %s\n", version)
}
