package being

/*
	COMBAT is fairly simple:
		0. Initiative is rolled (1dn where n is num of combatants, plus DEX/10^floor(log(10)n))
		1. An attacker selects a target randomly
		2. An attacker attacks their target using their highest stat
		3. Higher stat wins and deals damage

		Depending on the highest stat for the attacker, there are modifiers:
		- STR: Do STRMOD damage instead of 1
		- DEX: Add DEXMOD to initiative roll
		- CON: Add CONMOD to hitpoints (HP initially calculated as 1d(6+CONMOD) and adds CONMOD to result if CONMOD is highest)
		- INT: You attack INTMOD targets (AOE)
		- WIS: Remove WISMOD strongest opponents from list of targets (higher chance of success)
		- CHA: On successful attack, you have CHAMOD/20 chance of converting them to your team (dealing no damage)
*/

import (
	"fmt"
	"math"
	"sort"

	"github.com/fabiodesousa/settlement/dice"
)

// HitPoints have a max value, current value, and can be modified
type HitPoints struct {
	Max     int
	Current int
}

// NewHitPoints rolls 1d(6+CONMOD) and adds CONMOD to result if CONMOD is highest
func (stats *StatBlock) NewHitPoints() HitPoints {
	var mod int
	if stats.MaxStat().Name == "CON" {
		mod = stats.MaxStat().Mod
	} else {
		mod = 0
	}
	hp := dice.Roll(6+stats.GetStat("CON").Mod) + mod
	if hp < 1 {
		hp = 1
	}
	return HitPoints{
		Max:     hp,
		Current: hp,
	}
}

// RollInitiative rolls initiative for a being given the number of combatants
func (b *Being) RollInitiative(n int) float64 {
	roll := float64(dice.Roll(n))
	base := math.Pow(10, math.Floor(math.Log10(float64(n))))
	roll += float64(b.DEX()) / base
	return roll
}

// SelectTarget takes a slice of enemy beings, and selects at least one target
func (b *Being) SelectTarget(enemies []Being) []Being {
	fmt.Println("selecting target")
	max := b.MaxStat()
	fmt.Println("Max stat is " + max.Name)
	fmt.Printf("%d enemies\n", len(enemies))
	var targets []Being
	switch stat := max.Name; stat {
	// if the attacker is using WIS, they get to remove WISMOD strongest opponents
	case "WIS":
		// sort the enemies by the weakest slice
		sort.Slice(enemies[:], func(i, j int) bool {
			return enemies[i].GetStat(max.Name).Value < enemies[j].GetStat(max.Name).Value
		})
		// if there are more enemies than WISMOD, return the weakest slice
		// otherwise return a slice of 1
		fmt.Println(enemies)
		if len(targets) > max.Mod {
			narrowed := enemies[:max.Mod]
			targets = append(targets, narrowed[dice.Roll(len(narrowed))])
		} else {
			targets = append(targets, enemies[0])
		}
	// if the attacker is using INT, they get to select INTMOD targets at random
	case "INT":
		fmt.Printf("Selecting %d out of %d enemies", max.Mod, len(enemies))
		for i := 0; i < max.Mod; i++ {
			fmt.Printf("%d,", len(enemies))
			if len(enemies) > 0 {
				x := dice.Roll(len(enemies)) - 1
				targets = append(targets, enemies[x])
				enemies[x] = enemies[len(enemies)-1]
				enemies[len(enemies)-1] = Being{}
				enemies = enemies[:len(enemies)-1]
			}
		}
	default:
		x := dice.Roll(len(enemies) - 1)
		targets = append(targets, enemies[x])
	}
	return targets
}
