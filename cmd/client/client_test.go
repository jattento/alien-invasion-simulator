package client

import (
	"math/rand"
	"testing"

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
