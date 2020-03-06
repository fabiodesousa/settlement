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
	"sort"

	"github.com/fabiodesousa/settlement/dice"
)

// HitPoints have a max value, current value, and can be modified
type HitPoints struct {
	Max     int
	Current int
}

// AddTeamMember sets the Team for a Being
func (t *Team) AddTeamMember(being Being) {
	being.Team = t
	t.Roster = append(t.Roster, being)
}

// RemoveTeamMember removes a specified being from the Team
func (t *Team) RemoveTeamMember(being Being) {
	for i := range t.Roster {
		if t.Roster[i].ID == being.ID {
			t.Roster[i] = t.Roster[len(t.Roster)-1]
			t.Roster[len(t.Roster)-1] = Being{}
			t.Roster = t.Roster[:len(t.Roster)-1]
		}
	}

}

// DefectTo removes a being from their current team and adds them to the provided team
func (b Being) DefectTo(t *Team) {
	b.Team.RemoveTeamMember(b)
	t.AddTeamMember(b)
}

// Team has a Name and Roster
type Team struct {
	Name   string
	Roster []Being
}

// Encounter has an Initiative, TurnCount, Attackers, and Defenders
type Encounter struct {
	Initiative InitiativeOrder
	TurnCount  int
	Attackers  Team
	Defenders  Team
}

// InitiativeOrder has an array of being pointers sorted by initiative, plus current spot in order
type InitiativeOrder struct {
	Order       []Being
	CurrentTurn int
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

// PrintInitiative prints the current turn in the order and who is on deck
func (i InitiativeOrder) PrintInitiative() {
	var current = i.Order[i.CurrentTurn].Name
	var next string
	if i.CurrentTurn+1 < len(i.Order) {
		next = i.Order[i.CurrentTurn+1].Name
	} else {
		next = i.Order[0].Name
	}
	fmt.Printf("Current turn: %s\n", current)
	fmt.Printf("On deck: %s\n", next)

}

// PrintTotalInitiative prints the Initiative Order
func (i InitiativeOrder) PrintTotalInitiative() {
	for _, x := range i.Order {
		if x.Initiative >= 10 {
			fmt.Printf("%.2f | %s\n", x.Initiative, x.Name)
		} else {
			fmt.Printf(" %.2f | %s\n", x.Initiative, x.Name)
		}
	}
}

// RollInitiative rolls initiative for each being in an encounter and adds them to the resulting InitiativeOrder
func (e Encounter) RollInitiative() InitiativeOrder {
	//var size = len(e.Attackers.Roster) + len(e.Defenders.Roster) - 2
	var init = InitiativeOrder{CurrentTurn: 0, Order: make([]Being, 0)}
	for _, b := range e.Attackers.Roster {
		b.Initiative = RollInitiative(b)
		init.Order = append(init.Order, b)
	}
	for _, b := range e.Defenders.Roster {
		b.Initiative = RollInitiative(b)
		init.Order = append(init.Order, b)
	}
	// sort the initiative
	sort.Slice(init.Order, func(i, j int) bool {
		return init.Order[i].Initiative > init.Order[j].Initiative
	})
	return init
}

// RollInitiative rolls initiative for a being using a d20 + DEXMOD, and raw Dex for tie breakers as the decimal
func RollInitiative(b Being) float64 {
	roll := float64(dice.Roll(20) + b.GetStat("DEX").Mod)
	base := 100.00
	roll += float64(b.DEX()) / base
	return roll
}

// SelectTarget takes a slice of enemy being pointers, and selects at least one target
func (b Being) SelectTarget(enemies []Being) []Being {
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
		if len(targets) > max.Mod {
			narrowed := enemies[:max.Mod]
			targets = append(targets, narrowed[dice.Roll(len(narrowed))])
		} else {
			targets = append(targets, enemies[0])
		}
	// if the attacker is using INT, they get to select INTMOD targets at random
	case "INT":
		fmt.Printf("Selecting %d out of %d enemies", max.Mod+1, len(enemies))
		for i := 0; i <= max.Mod; i++ {
			if len(enemies) > 0 {
				x := dice.Roll(len(enemies)) - 1
				targets = append(targets, enemies[x])
				enemies[x] = enemies[len(enemies)-1]
				//enemies[len(enemies)-1] = Being{}
				enemies = enemies[:len(enemies)-1]
			}
		}
	default:
		x := dice.Roll(len(enemies) - 1)
		targets = append(targets, enemies[x])
	}
	return targets
}

// Attack takes an array of enemies and deals damage or converts them
func (b Being) Attack(enemies []Being) {
	for i := 0; i < len(enemies); i++ {
		max := b.MaxStat()
		fmt.Printf("%s (%s %d) attacks %s (%s %d)\n", b.Name, max.Name, max.Value, enemies[i].Name, max.Name, enemies[i].GetStat(max.Name).Value)
		if max.Value >= enemies[i].GetStat(max.Name).Value {
			fmt.Println("It's a hit!")
			// Charisma check
			if max.Name == "CHA" && dice.Roll(100) > (100-max.Mod*5) {
				enemies[i].DefectTo(b.Team)
				fmt.Printf("%s has been converted to Team %s\n", enemies[i].Name, enemies[i].Team.Name)
			} else if max.Name == "STR" {
				enemies[i].HitPoints.Current -= max.Mod
				fmt.Printf("%d damage! %s is at %d/%d\n", max.Mod, enemies[i].Name, enemies[i].HitPoints.Current, enemies[i].HitPoints.Max)
			} else {
				enemies[i].HitPoints.Current--
				fmt.Printf("%s is at %d/%d\n", enemies[i].Name, enemies[i].HitPoints.Current, enemies[i].HitPoints.Max)
			}
		}
	}
}
