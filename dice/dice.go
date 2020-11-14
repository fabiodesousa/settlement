package dice

import (
	"math/rand"
	"time"
)

type DiceRoller interface {
	Roll(int) int
}

type Dice struct {}

// Roll rolls a singe dice of the specified size
func (d Dice) Roll(size int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(size) + 1
}
