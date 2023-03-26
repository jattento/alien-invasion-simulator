package client

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/jattento/alien-invasion-simulator/internal/earth"
	"github.com/jattento/alien-invasion-simulator/internal/interface/terminal"
	"github.com/jattento/alien-invasion-simulator/internal/platform/numeric"
	"github.com/jattento/alien-invasion-simulator/internal/platform/system"
	"github.com/jattento/alien-invasion-simulator/internal/simulation"
)

const _defaultName = "world_specs.txt"

// Run blocks until the program is ended or an error happen
func Run(days, aliens, cities, matrixN int, configPath string) error {
	if configPath == "" {
		configPath = _defaultName

		if err := createMockFile(configPath, matrixN, cities); err != nil {
			return err
		}

		defer func() { _ = os.Remove(configPath) }()
	}

	sim, err := simulation.NewInvasion(configPath, aliens, system.NewManager(), days)
	if err != nil {
		return err
	}

	logsCh := make(chan string)
	DaysCh := make(chan string)
	citiesCh := make(chan []string)

	worldMatrix := worldMap{cities: make([]city, 0), citiesIndex: make(map[string]int), alive: aliens}
	for cityInfo := range sim.CityLayout {
		worldMatrix.save(city{name: cityInfo})
	}

	t := terminal.New(os.Stdout, logsCh, DaysCh, citiesCh)

	go func() {
		for keepTicking := true; keepTicking && worldMatrix.alive > 0; {
			now := time.Now()

			var report simulation.TickReport
			keepTicking, report = sim.Tick()

			worldMatrix.clear()

			for _, battleReport := range report.Battles {
				worldMatrix.save(city{name: battleReport.City, destroyed: true, aliens: battleReport.InvolvedAliens})
				logsCh <- killLog(battleReport)
			}
			for cityName := range report.AlienPositions {
				worldMatrix.save(city{name: cityName, aliens: report.AlienPositions[cityName]})
			}

			citiesCh <- worldMatrix.prettySlice()
			DaysCh <- fmt.Sprintf("🕒  :  %v   |   👽  :  %v   |   💀  :  %v   |   🏡  :  %v   |   🔥  :  %v",
				report.Tick, worldMatrix.alive, worldMatrix.dead, worldMatrix.notDestroyed, worldMatrix.destroyed)

			time.Sleep(time.Duration(atomic.LoadInt64(&t.WaitTime)) - (time.Now().Sub(now)))
		}

		logsCh <- "Congratulations on completing the alien simulation! "
		logsCh <- "Just remember, if any actual aliens come to visit, don't blame me if this isn't accurucate."

	}()

	if err = t.Run(); err != nil {
		return err
	}

	return nil

}

func killLog(report earth.BattleReport) string {
	skull := "💀"
	knife := "🗡️"
	gun := "🔫"
	bomb := "💣"
	wrench := "🔧"
	poison := "🧪"
	syringe := "💉"
	fire := "🔥"
	paperClip := "📎"

	weapons := []string{skull, knife, gun, bomb, wrench, poison, syringe, fire, paperClip}
	weapon := weapons[randomInt(0, len(weapons)-1)]

	fmtText := fmt.Sprintf("👽 %q", report.InvolvedAliens[0])
	for i := 1; i < len(report.InvolvedAliens); i++ {
		if i == len(report.InvolvedAliens)-1 {
			fmtText += fmt.Sprintf(" and 👽 %q", report.InvolvedAliens[i])
			continue
		}
		fmtText += fmt.Sprintf(", 👽 %q", report.InvolvedAliens[i])
	}

	return fmt.Sprintf("%s killed each other in %q in a %s  duel %s",
		fmtText, report.City, weapon, skull)
}

func randomInt(min int, max int) int {
	return min + rand.Intn(max-min+1)
}

func createMockFile(name string, matrix, size int) error {
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
