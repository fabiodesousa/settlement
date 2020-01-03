package being

import (
	"github.com/fabiodesousa/settlement/dice"
)

// A StatBlock has an array of stats and functions for fetching them by name
type StatBlock struct {
	StatArray []Stat
}

// MaxStat returns the highest Stat for a StatBlock
func (block *StatBlock) MaxStat() Stat {
	var maxStat Stat = block.StatArray[0]
	for _, e := range block.StatArray {
		if e.Value > maxStat.Value {
			maxStat = e
		}
	}
	return maxStat
}

// GetStat returns the specified stat given a name, defaults to maxstat
func (block *StatBlock) GetStat(name string) Stat {
	switch name {
	case "STR":
		return block.StatArray[0]
	case "DEX":
		return block.StatArray[1]
	case "CON":
		return block.StatArray[2]
	case "INT":
		return block.StatArray[3]
	case "WIS":
		return block.StatArray[4]
	case "CHA":
		return block.StatArray[5]
	default:
		return block.MaxStat()
	}
}

// STR returns the Strength stat for a StatBlock
func (block *StatBlock) STR() int {
	return block.StatArray[0].Value
}

// DEX returns the Dexterity stat for a StatBlock
func (block *StatBlock) DEX() int {
	return block.StatArray[1].Value
}

// CON returns the Constitution stat for a StatBlock
func (block *StatBlock) CON() int {
	return block.StatArray[2].Value
}

// INT returns the Intelligence stat for a StatBlock
func (block *StatBlock) INT() int {
	return block.StatArray[3].Value
}

// WIS returns the Wisdom stat for a StatBlock
func (block *StatBlock) WIS() int {
	return block.StatArray[4].Value
}

// CHA returns the Charisma stat for a StatBlock
func (block *StatBlock) CHA() int {
	return block.StatArray[5].Value
}

// NewStatBlock randomly generates a new array of stats
func NewStatBlock() StatBlock {
	array := StatBlock{
		StatArray: []Stat{
			NewStat("STR"),
			NewStat("DEX"),
			NewStat("CON"),
			NewStat("INT"),
			NewStat("WIS"),
			NewStat("CHA"),
		},
	}
	return array
}

// A Stat has a name and value
type Stat struct {
	Name  string
	Value int
	Mod   int
}

// NewStat rolls a stat given a name
func NewStat(name string) Stat {
	v := rollStat()
	var mod int
	if v < 10 {
		mod = (v - 11) / 2
	} else {
		mod = (v - 10) / 2
	}
	s := Stat{
		Name:  name,
		Value: v,
		Mod:   mod,
	}
	return s
}

// roll 4d6kh3
func rollStat() int {
	rolls := []int{dice.Roll(6), dice.Roll(6), dice.Roll(6), dice.Roll(6)}
	sum := 0
	var minValue int
	for i, e := range rolls {
		sum += e
		if i == 0 || e < minValue {
			minValue = e
		}
	}
	sum -= minValue
	return sum
}
