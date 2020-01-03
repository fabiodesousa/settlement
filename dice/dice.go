package dice

import (
	"math/rand"
	"time"
)

// Roll rolls a singe dice of the specified size
func Roll(size int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(size) + 1
}
