package earth

import (
	"math/rand"
	"testing"
)

func TestPlanet_NextDay(t *testing.T) {
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

	planet.randomizer = rand.New(rand.NewSource(0))

	if len(planet.Aliens) != 2 {
		t.Errorf("the planet should have 2 aliens, but has %d", len(planet.Aliens))
	}

	reports := planet.NextDay()

	if len(planet.Aliens) != 2 {
		t.Errorf("there should still be 2 aliens, but there are %d", len(planet.Aliens))
	}

	reports = planet.NextDay()
	for _, report := range reports {
		for _, alien := range report.InvolvedAliens {
			if planet.Aliens[alien] != nil {
				t.Errorf("alien %s should be in a destroyed city, but it is in %s", alien, planet.Aliens[alien].Id)
			}
		}
	}
}
