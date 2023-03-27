package earth

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

func TestRandomAlienNames(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	randomizer := rand.New(rand.NewSource(rand.Int63()))

	names := randomAlienNames(1000, randomizer)

	if len(names) != 1000 {
		t.Errorf("Expected 1000 names, but got %d", len(names))
	}

	nameCounts := make(map[string]int)
	for _, name := range names {
		nameCounts[name]++
	}
	for name, count := range nameCounts {
		if count > 1 {
			t.Errorf("Name '%s' appears %d times", name, count)
		}
	}

	for i := 0; i < len(names); i++ {
		for j := i + 1; j < len(names); j++ {
			if names[i] == names[j] {
				if !strings.Contains(names[i], "I") {
					t.Errorf("Duplicate name '%s' does not have a roman numeral counter", names[i])
				}
				break
			}
		}
	}
}
