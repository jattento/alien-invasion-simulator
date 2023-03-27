package simulation

import (
	"errors"
	"os"

	"github.com/jattento/alien-invasion-simulator/internal/earth"
	"github.com/jattento/alien-invasion-simulator/internal/platform/system"
)

type Invasion struct {
	planet    *earth.Planet
	tickCount int
	tickLimit int

	CityLayout map[string]map[earth.Direction]string
}

type TickReport struct {
	Battles        []earth.BattleReport
	AlienPositions map[string][]string
	Tick           int
}

type SystemManager interface {
	LoadFile(path string) (system.LoadFileRecords, error)
}

const _defaultName = "world_specs.txt"

var _directionToEnum = map[string]earth.Direction{
	"north": earth.North,
	"south": earth.South,
	"east":  earth.East,
	"west":  earth.West,
}

// Cities returns a copy of the map layout
func (invasion Invasion) Cities() map[string]map[earth.Direction]string {
	mapCopy := make(map[string]map[earth.Direction]string)
	for k, m := range invasion.CityLayout {
		mapCopy[k] = make(map[earth.Direction]string)
		for direction, adjacent := range m {
			mapCopy[k][direction] = adjacent
		}
	}

	return mapCopy
}

func NewInvasion(planetSpecsFile string, aliensAmount int, systemManager SystemManager, tickLimit, cities, matrixN int) (*Invasion, error) {
	if planetSpecsFile == "" {
		planetSpecsFile = _defaultName

		if err := generateFile(planetSpecsFile, matrixN, cities); err != nil {
			return nil, err
		}

		defer func() { _ = os.Remove(planetSpecsFile) }()
	}

	recs, err := systemManager.LoadFile(planetSpecsFile)
	if err != nil {
		return nil, err
	}

	if valid := validateFileRecords(recs); !valid {
		return nil, errors.New("invalid city layout")
	}

	earthCityLayout := make(map[string]map[earth.Direction]string)
	for city, directions := range recs {
		earthCityLayout[city] = make(map[earth.Direction]string)
		for direction, adjacentCity := range directions {
			earthCityLayout[city][_directionToEnum[direction]] = adjacentCity
		}
	}

	planet, err := earth.New(earthCityLayout, aliensAmount)
	if err != nil {
		return nil, err
	}

	return &Invasion{planet: planet, tickLimit: tickLimit, CityLayout: earthCityLayout}, nil
}

func (invasion Invasion) alienPositions() map[string][]string {
	positions := make(map[string][]string)

	for alien, vertex := range invasion.planet.Aliens {
		positions[vertex.Id] = append(positions[vertex.Id], alien)
	}

	return positions
}

// Tick first return value indicates if the function should continue to be called
func (invasion *Invasion) Tick() (bool, TickReport) {
	invasion.tickCount++

	return invasion.tickCount < invasion.tickLimit, TickReport{
		Battles:        invasion.planet.NextDay(),
		Tick:           invasion.tickCount - 1,
		AlienPositions: invasion.alienPositions(),
	}
}

func validateFileRecords(fileRecords system.LoadFileRecords) bool {
	for _, cityRoads := range fileRecords {
		for key := range cityRoads {
			if key != "north" && key != "south" && key != "west" && key != "east" {
				return false
			}
		}

		_, hasNorth := cityRoads["north"]
		_, hasEast := cityRoads["east"]
		_, hasSouth := cityRoads["south"]
		_, hasWest := cityRoads["west"]

		if hasSouth {
			if _, exists := fileRecords[cityRoads["south"]]; !exists {
				return false
			}
			if _, hasSubNorth := fileRecords[cityRoads["south"]]["north"]; !hasSubNorth {
				return false
			}
		}
		if hasWest {
			if _, exists := fileRecords[cityRoads["west"]]; !exists {
				return false
			}
			if _, hasSubEast := fileRecords[cityRoads["west"]]["east"]; !hasSubEast {
				return false
			}
		}
		if hasEast {
			if _, exists := fileRecords[cityRoads["east"]]; !exists {
				return false
			}
			if _, hasSubWest := fileRecords[cityRoads["east"]]["west"]; !hasSubWest {
				return false
			}
		}
		if hasNorth {
			if _, exists := fileRecords[cityRoads["north"]]; !exists {
				return false
			}
			if _, hasSubSouth := fileRecords[cityRoads["north"]]["south"]; !hasSubSouth {
				return false
			}
		}
	}

	return true
}
