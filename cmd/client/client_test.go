package client

import (
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/jattento/alien-invasion-simulator/internal/earth"
	"github.com/stretchr/testify/assert"
)

func TestKillLog(t *testing.T) {
	report := earth.BattleReport{
		City:           "New York",
		InvolvedAliens: []string{"Alien1", "Alien2"},
	}
	log := killLog(report)

	assert.Contains(t, log, "Alien1")
	assert.Contains(t, log, "Alien2")
	assert.Contains(t, log, "New York")
}

func TestRandomInt(t *testing.T) {
	rand.Seed(42)

	output := randomInt(1, 10)
	assert.Equal(t, 6, output, "randomInt() returned %d, expected 8", output)

	output = randomInt(5, 5)
	assert.Equal(t, 5, output, "randomInt() returned %d, expected 5", output)
}

func TestDeleteCity(t *testing.T) {
	cities := make(map[string]map[earth.Direction]string)
	cities["New York"] = map[earth.Direction]string{earth.North: "Toronto", earth.South: "Philadelphia", earth.East: "Boston"}
	cities["Toronto"] = map[earth.Direction]string{earth.South: "New York"}
	cities["Boston"] = map[earth.Direction]string{earth.West: "New York"}
	deleteCity(cities, "New York")

	assert.Equal(t, map[string]map[earth.Direction]string{
		"Boston":  map[int]string{},
		"Toronto": map[int]string{},
	}, cities)
}

func TestFinalLogs(t *testing.T) {
	logsCh := make(chan string)
	remainingCities := make(map[string]map[earth.Direction]string)
	remainingCities["New York"] = map[earth.Direction]string{earth.North: "Toronto", earth.South: "Philadelphia", earth.East: "Boston"}
	remainingCities["Toronto"] = map[earth.Direction]string{earth.South: "New York"}

	go finalLogs(logsCh, remainingCities)

	expectedOutput := "--------" +
		"New York north=Toronto south=Philadelphia east=Boston" +
		"Toronto south=New York" +
		"These are the remaining cities..." +
		"--------" +
		"Just remember, if any actual aliens come to visit, don't blame me if this isn't accurucate." +
		"Congratulations on completing the alien simulation!"

	timeout := time.After(5 * time.Second) // Wait for 5 seconds before timing out
	for i := 0; i < 7; i++ {
		select {
		case actualOutput := <-logsCh:
			assert.Contains(t, expectedOutput, strings.Split(actualOutput, " ")[0])
		case <-timeout:
			t.Error("Test timed out") // Test failed due to timeout
			return
		}
	}
}
