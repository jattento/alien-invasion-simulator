package client

import (
	"fmt"
	"strings"
)

type worldMap struct {
	cities       []city
	citiesIndex  map[string]int
	notDestroyed int
	destroyed    int
	alive        int
	dead         int
}

type city struct {
	name      string
	aliens    []string
	destroyed bool
}

func (world *worldMap) prettySlice() []string {
	var output []string

	for _, c := range world.cities {
		output = append(output, fmt.Sprintf("%s(%s)", c.fmtName(), c.FmtAliens()))
	}

	return output
}

func (c city) FmtAliens() string {
	if len(c.aliens) == 0 {
		return ""
	}

	if c.destroyed {
		return fmt.Sprint("💀️" + strings.Join(c.aliens, ", 💀️"))
	}

	return fmt.Sprint("👽" + strings.Join(c.aliens, ", 👽"))
}

func (c city) fmtName() string {
	if c.destroyed {
		return fmt.Sprint("🔥🔥" + c.name + "🔥🔥")
	}

	return fmt.Sprint("🏠🌳" + c.name + "🌳🏠")
}

// Clear aliens positions
func (world *worldMap) clear() {
	for i := 0; i < len(world.cities); i++ {
		if !world.cities[i].destroyed {
			world.cities[i].aliens = make([]string, 0)
		}
	}
}

func (world *worldMap) save(c city) {
	if i, exist := world.citiesIndex[c.name]; exist {
		existingCity := world.cities[i]
		if existingCity.destroyed {
			return
		}

		if c.destroyed {
			world.notDestroyed--
			world.destroyed++
			world.alive -= len(c.aliens)
			world.dead += len(c.aliens)
		}

		world.cities[i] = c
		return
	}

	world.notDestroyed++

	world.cities = append(world.cities, c)
	world.citiesIndex[c.name] = len(world.cities) - 1
}
