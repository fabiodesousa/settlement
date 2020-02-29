package being

import (
	"fmt"

	"github.com/Pallinder/go-randomdata"
	"github.com/fabiodesousa/settlement/dice"
	"github.com/google/uuid"
)

// A Being has a name, sex and a slice of six stats
type Being struct {
	ID         string
	Name       string
	Sex        string
	Team       *Team
	Initiative float64
	HitPoints
	StatBlock
}

// rolls 1d100, <50 = male, else female
func determineSex() string {
	if dice.Roll(100) < 50 {
		return "male"
	}
	return "female"
}

// generates a name given a sex (we'll take gender into account later)
func makeName(sex string) string {
	if sex == "male" {
		return randomdata.FullName(randomdata.Male)
	}
	return randomdata.FullName(randomdata.Female)
}

// PrintStats prints each of the Beings stats
func PrintStats(b Being) {
	fmt.Printf("%s (%d/%d) [%s]\n", b.Name, b.HitPoints.Current, b.HitPoints.Max, b.MaxStat().Name)
	fmt.Printf("Team: %s\n", b.Team.Name)
	fmt.Printf("STR:%d (%d)\t", b.STR(), b.GetStat("STR").Mod)
	fmt.Printf("INT:%d (%d)\n", b.INT(), b.GetStat("INT").Mod)
	fmt.Printf("DEX:%d (%d)\t", b.DEX(), b.GetStat("DEX").Mod)
	fmt.Printf("WIS:%d (%d)\n", b.WIS(), b.GetStat("WIS").Mod)
	fmt.Printf("CON:%d (%d)\t", b.CON(), b.GetStat("CON").Mod)
	fmt.Printf("CHA:%d (%d)\n", b.CHA(), b.GetStat("CHA").Mod)
	return
}

// NewRandomBeing generates a random Being
func NewRandomBeing() Being {
	sex := determineSex()
	stats := NewStatBlock()
	c := Being{
		ID:        uuid.Must(uuid.NewRandom()).String(),
		Sex:       sex,
		Name:      makeName(sex),
		StatBlock: stats,
		HitPoints: stats.NewHitPoints(),
	}
	return c
}
