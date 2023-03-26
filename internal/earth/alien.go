package earth

import (
	"fmt"
	"github.com/jattento/alien-invasion-simulator/internal/platform/numeric"
	"math/rand"
)

// Generates a slice with `amount` random alien names.
// If a name is repeated, it appends a roman numeral counter to the end.
func randomAlienNames(amount int, randomizer *rand.Rand) []string {
	prefixes := []string{"Zorg", "Vort", "Gork", "Gorbl", "Borg", "Krel", "Mort", "Snag", "Thrag", "Zug"}
	suffixes := []string{"on", "ax", "ik", "ar", "or", "ul", "ith", "ol", "arx", "ath"}

	// A map to keep track of the number of occurrences of each name
	nameCounts := make(map[string]int)

	// A slice to store the generated names
	names := make([]string, amount)

	for i := 0; i < amount; i++ {
		name := fmt.Sprintf("%s %s",
			prefixes[randomizer.Intn(len(prefixes))],
			suffixes[randomizer.Intn(len(suffixes))],
		)

		count, alreadyTaken := nameCounts[name]
		if alreadyTaken {
			// If the name has already been generated, append a roman numeral counter
			count++
			nameCounts[name] = count
			names[i] = fmt.Sprintf("%s %s", name, numeric.ToRomanSystem(count))
			continue
		}

		nameCounts[name] = 1
		names[i] = name
	}

	return names
}
