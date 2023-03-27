package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrettySlice(t *testing.T) {
	world := worldMap{
		cities: []city{
			{
				name:      "New York",
				aliens:    []string{"Zog", "Krog", "Morg"},
				destroyed: false,
			},
			{
				name:      "San Francisco",
				aliens:    []string{"Gorg", "Lorg"},
				destroyed: true,
			},
		},
		citiesIndex:  map[string]int{"New York": 0, "San Francisco": 1},
		notDestroyed: 1,
		destroyed:    1,
		alive:        3,
		dead:         0,
	}

	expectedOutput := []string{"🏠🌳New York🌳🏠(👽Zog, 👽Krog, 👽Morg)", "🔥🔥San Francisco🔥🔥(💀️Gorg, 💀️Lorg)"}
	actualOutput := world.prettySlice()
	assert.Equal(t, expectedOutput, actualOutput)
}

func TestFmtAliens(t *testing.T) {
	city1 := city{
		name:      "New York",
		aliens:    []string{},
		destroyed: false,
	}

	city2 := city{
		name:      "San Francisco",
		aliens:    []string{"Gorg", "Lorg"},
		destroyed: true,
	}

	expectedOutput1 := ""
	expectedOutput2 := "💀️Gorg, 💀️Lorg"

	actualOutput1 := city1.FmtAliens()
	actualOutput2 := city2.FmtAliens()

	assert.Equal(t, expectedOutput1, actualOutput1)
	assert.Equal(t, expectedOutput2, actualOutput2)
}

func TestFmtName(t *testing.T) {
	city1 := city{
		name:      "New York",
		aliens:    []string{},
		destroyed: false,
	}

	city2 := city{
		name:      "San Francisco",
		aliens:    []string{},
		destroyed: true,
	}

	expectedOutput1 := "🏠🌳New York🌳🏠"
	expectedOutput2 := "🔥🔥San Francisco🔥🔥"

	actualOutput1 := city1.fmtName()
	actualOutput2 := city2.fmtName()

	assert.Equal(t, expectedOutput1, actualOutput1)
	assert.Equal(t, expectedOutput2, actualOutput2)
}

func TestClear(t *testing.T) {
	world := worldMap{
		cities: []city{
			{
				name:      "New York",
				aliens:    []string{"Zog", "Krog", "Morg"},
				destroyed: false,
			},
			{
				name:      "San Francisco",
				aliens:    []string{"Gorg", "Lorg"},
				destroyed: true,
			},
		},
		citiesIndex: map[string]int{"New York": 0, "San Francisco": 1},
	}

	world.clear()

	for _, c := range world.cities {
		if !c.destroyed && len(c.aliens) != 0 {
			t.Errorf("Unexpected aliens in city %s. Expected empty list, but got: %v", c.name, c.aliens)
		}
	}
}

func TestSave(t *testing.T) {
	world := &worldMap{
		cities:       make([]city, 0),
		citiesIndex:  make(map[string]int),
		notDestroyed: 0,
		destroyed:    0,
		alive:        0,
		dead:         0,
	}

	c1 := city{name: "Paris", aliens: []string{"Alien1", "Alien2"}, destroyed: false}
	world.save(c1)
	assert.Equal(t, 1, len(world.cities))
	assert.Equal(t, "Paris", world.cities[0].name)

	c2 := city{name: "Paris", aliens: []string{"Alien3"}, destroyed: false}
	world.save(c2)
	assert.Equal(t, 1, len(world.cities))
	assert.Equal(t, 1, len(world.cities[0].aliens))
	assert.Equal(t, "Alien3", world.cities[0].aliens[0])

	c3 := city{name: "Paris", aliens: nil, destroyed: true}
	world.save(c3)
	assert.True(t, world.cities[0].destroyed)
	assert.Equal(t, 0, world.notDestroyed)
	assert.Equal(t, 1, world.destroyed)
	assert.Equal(t, 0, world.alive)
}
