package earth

import (
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/jattento/alien-invasion-simulator/internal/platform/datastructure"
)

type BattleReport struct {
	InvolvedAliens []string
	City           string
}

type Planet struct {
	graph *datastructure.Graph

	// Name:Position
	Aliens map[string]*datastructure.Vertex

	// This cache saves the state of the cities and Aliens at the moment the Planet is created
	// to be used at day zero without the need to process it again
	dayZeroCacheData map[*datastructure.Vertex][]string

	randomizer *rand.Rand
}

type Direction = int

const (
	North Direction = iota
	East
	South
	West
)

// New input looks like: <Bar:1:Foo>
// This function isn't designed to be performant but to give a nice interface,
// and since this function is going to be called once at the beginning it is acceptable.
func New(citiesAndAdjacent map[string]map[Direction]string, aliensAmount int) (*Planet, error) {
	randomizer := rand.New(rand.NewSource(time.Now().UnixNano()))

	p := Planet{
		graph:            new(datastructure.Graph),
		Aliens:           make(map[string]*datastructure.Vertex),
		randomizer:       randomizer,
		dayZeroCacheData: make(map[*datastructure.Vertex][]string),
	}

	// Since we need to make first all the AddVertex calls, we are going to save the AddEdge calls
	addEdgeCalls := make([]func() error, 0)

	// This slice is going to be used to generate the random alien positions.
	cities := make([]*datastructure.Vertex, 0)

	// Execute AddVertex calls and prepare AddEdge ones
	for city, adjacent := range citiesAndAdjacent {
		cityRef, err := p.graph.AddVertex(city)
		if err != nil {
			return nil, err
		}

		cities = append(cities, cityRef)
		p.dayZeroCacheData[cityRef] = make([]string, 0)

		cityAddEdgeCalls, err := p.buildAddEdgeCalls(city, adjacent)
		if err != nil {
			return nil, err
		}

		addEdgeCalls = append(addEdgeCalls, cityAddEdgeCalls...)
	}

	// Execute AddEdge calls
	for _, addEdgeCall := range addEdgeCalls {
		if err := addEdgeCall(); err != nil {
			return nil, err
		}
	}

	// This function priority to cities which were not selected already
	randomCitySelectorFunc := newRandomSelector(randomizer, cities)

	alienNames := randomAlienNames(aliensAmount, randomizer)

	for _, alienName := range alienNames {
		city := randomCitySelectorFunc()
		p.Aliens[alienName] = city
		p.dayZeroCacheData[city] = append(p.dayZeroCacheData[city], alienName)
	}

	return &p, nil
}

func (planet *Planet) buildAddEdgeCalls(city string, adjacentCities map[Direction]string) ([]func() error, error) {
	calls := make([]func() error, 0)

	for direction, adjacentCity := range adjacentCities {
		if direction > 3 {
			return nil, fmt.Errorf("invalid direction: %q -> %q -> %q", city, direction, adjacentCity)
		}

		calls = append(calls, func() error {
			err := planet.graph.AddEdge(direction, city, adjacentCity)
			if err != nil && errors.Is(err, datastructure.ErrEdgeDuplicated) {
				return nil
			}

			return err
		})
	}

	return calls, nil
}

func newRandomSelector[T comparable](randomizer *rand.Rand, items []T) func() T {
	selected := make(map[T]bool)
	selectedCount := 0

	return func() T {
		if selectedCount >= len(items) {
			return items[randomizer.Intn(len(items))]
		}

		var item T

		for {
			item = items[randomizer.Intn(len(items))]
			if !selected[item] {
				selected[item] = true
				selectedCount++
				break
			}
		}

		return item
	}
}

func (planet *Planet) NextDay() []BattleReport {
	if planet.dayZeroCacheData != nil {
		report := planet.processDay(planet.dayZeroCacheData)
		planet.dayZeroCacheData = nil

		return report
	}

	updatedData := make(map[*datastructure.Vertex][]string)

	// Alien movements...
	for alienId, alienLocation := range planet.Aliens {
		edges := append(alienLocation.AllEdges(), 4)

		destinationEdge := edges[planet.randomizer.Intn(len(edges))]

		// Stay at the same place
		var newDestination *datastructure.Vertex
		if destinationEdge == 4 {
			newDestination = alienLocation
		} else {
			newDestination = alienLocation.GetAdjacent(destinationEdge)
		}

		planet.Aliens[alienId] = newDestination

		updatedDataAliens, updatedDataExistForCity := updatedData[newDestination]
		if !updatedDataExistForCity {
			updatedData[newDestination] = make([]string, 0)
		}

		updatedData[newDestination] = append(updatedDataAliens, alienId)
	}

	return planet.processDay(updatedData)
}

// citiesCache: city:[alien1Id,alien2Id]
func (planet *Planet) processDay(citiesCache map[*datastructure.Vertex][]string) []BattleReport {
	destroyedCities := make(map[*datastructure.Vertex]struct{})
	destroyedAliens := make([]string, 0)
	reports := make([]BattleReport, 0)

	for city, aliens := range citiesCache {
		if len(aliens) > 1 {
			destroyedAliens = append(destroyedAliens, aliens...)
			destroyedCities[city] = struct{}{}
			reports = append(reports, BattleReport{City: city.Id, InvolvedAliens: aliens})
		}
	}

	for _, alien := range destroyedAliens {
		delete(planet.Aliens, alien)
	}

	for city := range destroyedCities {
		city.Disable()
	}

	return reports
}
