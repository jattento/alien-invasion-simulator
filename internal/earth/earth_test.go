package earth

import (
	"math/rand"
	"testing"
)

func TestPlanet_NextDay(t *testing.T) {
	// Create a new Planet
	planet, err := New(map[string]map[Direction]string{
		"New York": {
			South: "Washington",
			West:  "Chicago",
		},
		"Washington": {
			North: "New York",
			West:  "San Francisco",
		},
		"Chicago": {
			East: "New York",
		},
		"San Francisco": {
			East: "Washington",
		},
	}, 2)
	if err != nil {
		t.Errorf("error while creating the planet: %v", err)
	}

	// Mock the randomizer so that it always returns the same alien movement direction
	planet.randomizer = rand.New(rand.NewSource(0))

	// Make sure that the planet was created correctly
	if len(planet.Aliens) != 2 {
		t.Errorf("the planet should have 2 aliens, but has %d", len(planet.Aliens))
	}

	// Run the NextDay method and get the battle reports
	reports := planet.NextDay()

	// Check the state of the aliens after the first day
	if len(planet.Aliens) != 2 {
		t.Errorf("there should still be 2 aliens, but there are %d", len(planet.Aliens))
	}

	// Run the NextDay method again and make sure that all aliens are in destroyed cities
	reports = planet.NextDay()
	for _, report := range reports {
		for _, alien := range report.InvolvedAliens {
			if planet.Aliens[alien] != nil {
				t.Errorf("alien %s should be in a destroyed city, but it is in %s", alien, planet.Aliens[alien].Id)
			}
		}
	}
}
