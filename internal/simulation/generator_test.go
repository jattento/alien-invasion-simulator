package simulation

import (
	"os"
	"testing"
)

func TestGenerateFile(t *testing.T) {
	// Test case 1: Valid input
	err := generateFile("test1.txt", 5, 10)
	if err != nil {
		t.Errorf("generateFile() error = %v; want nil", err)
	}
	// Clean up
	_ = os.Remove("test1.txt")

	// Test case 2: Invalid file name
	err = generateFile("", 5, 10)
	if err == nil {
		t.Errorf("generateFile() error = nil; want non-nil")
	}
}

func TestGenerateCities(t *testing.T) {
	// Test case 1: Valid input
	cities := generateCities(10, 5)
	if cities == "" {
		t.Errorf("generateCities() = ''; want non-empty string")
	}

	// Test case 2: Invalid number of cities
	cities = generateCities(0, 5)
	if cities != "" {
		t.Errorf("generateCities() = %v; want ''", cities)
	}
}
