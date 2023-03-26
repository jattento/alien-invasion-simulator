package simulation

import (
	"errors"
	"io/ioutil"
	"os"
	"testing"

	"github.com/jattento/alien-invasion-simulator/internal/earth"
	"github.com/jattento/alien-invasion-simulator/internal/platform/datastructure"
	"github.com/jattento/alien-invasion-simulator/internal/platform/system"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type MockSystemManager struct{}

func (msm *MockSystemManager) LoadFile(path string) (system.LoadFileRecords, error) {
	if path == "invalid_file" {
		return nil, errors.New("invalid file")
	}

	return system.LoadFileRecords{
		"City1": {
			"north": "City2",
			"south": "City3",
			"east":  "City4",
			"west":  "City5",
		},
		"City2": {
			"south": "City1",
		},
		"City3": {
			"north": "City1",
		},
		"City4": {
			"west": "City1",
		},
		"City5": {
			"east": "City1",
		},
	}, nil
}

func TestNewInvasion(t *testing.T) {
	// Test loading file error
	msm := &MockSystemManager{}
	_, err := NewInvasion("invalid_file", 10, msm, 100)
	assert.Error(t, err)

	// Test valid city layout
	invasion, err := NewInvasion("some_file", 10, msm, 100)
	assert.NoError(t, err)
	assert.NotNil(t, invasion.planet)
	assert.Equal(t, 100, invasion.tickLimit)
	assert.NotEmpty(t, invasion.CityLayout)
}

func TestInvasion_alienPositions(t *testing.T) {
	invasion := &Invasion{
		planet: &earth.Planet{
			Aliens: map[string]*datastructure.Vertex{
				"A1": &datastructure.Vertex{Id: "City1"},
				"A2": &datastructure.Vertex{Id: "City2"},
				"A3": &datastructure.Vertex{Id: "City3"},
				"A4": &datastructure.Vertex{Id: "City3"},
				"A5": &datastructure.Vertex{Id: "City1"},
			},
		},
	}

	expected := map[string][]string{"City1": []string{"A5", "A1"}, "City2": []string{"A2"}, "City3": []string{"A3", "A4"}}

	actual := invasion.alienPositions()

	for city, aliens := range expected {
		assert.ElementsMatch(t, aliens, actual[city])
	}
}

func TestInvasion_Tick(t *testing.T) {
	// Create a temporary file with valid city layout for testing purposes
	file, err := ioutil.TempFile("", "city_layout")
	require.NoError(t, err)
	defer os.Remove(file.Name())

	_, err = file.WriteString("north=south_city\neast=east_city\nwest=west_city\n")
	require.NoError(t, err)

	invasion, err := NewInvasion(file.Name(), 2, &MockSystemManager{}, 2)
	require.NoError(t, err)

	// Call the Tick() method and check the return values
	ok, report := invasion.Tick()
	assert.True(t, ok)
	assert.Equal(t, 0, report.Tick)

	ok, report = invasion.Tick()
	assert.False(t, ok)
	assert.Equal(t, 1, report.Tick)
	assert.NotEmpty(t, report.AlienPositions)
}
