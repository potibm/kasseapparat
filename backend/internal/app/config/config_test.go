package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadVersionFromFileWithValidFile(t *testing.T) {
	// Arrange
	filename := "./VERSION"
	expected := "1.2.3\n"
	err := os.WriteFile(filename, []byte(expected), 0644)
	assert.NoError(t, err)

	defer os.Remove(filename) // Clean up

	// Act
	version := readVersionFromFile()

	// Assert
	assert.Equal(t, "1.2.3", version)
}

func TestReadVersionFromFileWithFileMissing(t *testing.T) {
	// Ensure file is absent
	_ = os.Remove("./VERSION")

	version := readVersionFromFile()

	assert.Equal(t, "0.0.0", version)
}
