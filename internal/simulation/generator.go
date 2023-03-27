package simulation

import (
	"io"
	"math/rand"
	"os"
	"strings"

	"github.com/jattento/alien-invasion-simulator/internal/platform/numeric"
)

func generateFile(name string, matrix, size int) error {
	worldsSpecsFile, err := os.Create(name)
	if err != nil {
		return err
	}

	defer func() {
		_ = worldsSpecsFile.Close()
	}()

	generatedCitySpecs := generateCities(size, matrix)

	_, err = io.Copy(worldsSpecsFile, strings.NewReader(generatedCitySpecs))
	if err != nil {
		return err
	}

	return nil
}

func generateCities(numCities, size int) string {
	cities := make([][]string, size)
	for i := 0; i < size; i++ {
		cities[i] = make([]string, size)
		for j := 0; j < size; j++ {
			cities[i][j] = ""
		}
	}

	// Pool of city names
	cityPool := []string{"new_york", "los_angeles", "buenos_aires", "paris", "phoenix", "shanghai", "san_antonio", "new_delhi", "dallas", "san_jose", "austin", "jacksonville", "fort_worth", "columbus", "san_francisco", "brasilia", "south_africa", "denver", "washington", "medellin"}

	// Add cities to the matrix randomly
	for i := 0; i < numCities; i++ {
		x := rand.Intn(size)
		y := rand.Intn(size)
		for cities[x][y] != "" {
			x = rand.Intn(size)
			y = rand.Intn(size)
		}
		if i < len(cityPool) {
			cities[x][y] = cityPool[i]
		} else {
			cities[x][y] = cityPool[(i-len(cityPool))%len(cityPool)] + "_" + numeric.ToRomanSystem(i-len(cityPool)+1)
		}
	}

	directions := []string{"north", "south", "east", "west"}
	result := ""
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if cities[i][j] != "" {
				result += cities[i][j]
				for _, dir := range directions {
					x, y := i, j
					if dir == "north" {
						x--
					} else if dir == "south" {
						x++
					} else if dir == "west" {
						y--
					} else if dir == "east" {
						y++
					}
					if x >= 0 && x < size && y >= 0 && y < size && cities[x][y] != "" {
						result += " " + dir + "=" + cities[x][y]
					}
				}
				result += "\n"
			}
		}
	}
	return result
}
