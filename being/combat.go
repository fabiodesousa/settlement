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
func (t *Team) AddTeamMember(being *Being) {
	being.Team = t
	t.Roster = append(t.Roster, being)
}

// RemoveTeamMember removes a specified being from the Team
func (t *Team) RemoveTeamMember(being *Being) {
	for i, b := range t.Roster {
		if b.ID == being.ID {
			t.Roster[i] = t.Roster[len(t.Roster)-1]
			t.Roster[len(t.Roster)-1] = &Being{}
			t.Roster = t.Roster[:len(t.Roster)-1]
		}
	}

}

// DefectTo removes a being from their current team and adds them to the provided team
func (b *Being) DefectTo(t *Team) {
	b.Team.RemoveTeamMember(b)
	t.AddTeamMember(b)
}

// Team has a Name and Roster
type Team struct {
	Name   string
	Roster []*Being
}

func GetSurvivors(t []*Being) []*Being {
	result := []*Being{}
	for _, b := range t {
		if b.IsAlive() {
			result = append(result, b)
		}
	}
	return result
}

// Encounter has an Initiative, TurnCount, Attackers, and Defenders
type Encounter struct {
	Initiative []*Being
	InitiativeTracker int
	TurnCount  int
	Attackers  Team
	Defenders  Team
	Winner string
	Log []*CombatLog
}

type CombatLog struct {
	Attacker *Being
	Defenders []*Being
	Hits []bool
	Kills []bool
	Conversions []bool
}

func (e *Encounter) SortInitiative() {
	sort.SliceStable(e.Initiative, func(i, j int) bool {
		return e.Initiative[i].Initiative > e.Initiative[j].Initiative
	})
}

func (e *Encounter) SetInitiative(init []*Being) {
	e.Initiative = init
}

func (e *Encounter) AddToInitiative(b *Being) {
	e.SetInitiative(append(e.Initiative, b))
}

// NewHitPoints rolls 1d(6+CONMOD) and adds CONMOD to result if CONMOD is highest
func (stats *StatBlock) NewHitPoints() HitPoints {
	d := dice.Dice{}
	var mod int
	if stats.MaxStat().Name == "CON" {
		mod = stats.MaxStat().Mod
	} else {
		mod = 0
	}
	hp := d.Roll(6+stats.GetStat("CON").Mod) + mod
	if hp < 1 {
		hp = 1
	}
	return HitPoints{
		Max:     hp,
		Current: hp,
	}
}

// RollInitiative rolls initiative for each being in an encounter and adds them in order to the InitiativeOrder
func (e *Encounter) RollInitiative() {
	for _, i := range e.Attackers.Roster {
		i.RollInitiative()
		e.AddToInitiative(i)
	}
	for _, i := range e.Defenders.Roster {
		i.RollInitiative()
		e.AddToInitiative(i)
	}
	e.SortInitiative()
}

// RollInitiative rolls initiative for a being using a d20
func (b *Being) RollInitiative() {
	d := dice.Dice{}
	roll := float64(d.Roll(20))
	base := 100.00
	roll += float64(b.DEX()) / base
	b.SetInitiative(roll)
}

// SelectTarget takes a slice of enemy being pointers, and selects at least one target
func (b Being) SelectTarget(enemyTeam []*Being) []*Being {
	d := dice.Dice{}
	max := b.MaxStat()
	enemies := GetSurvivors(enemyTeam)
	if len(enemies) == 0 {
		return []*Being{}
	}
	var targets []*Being
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
			targets = append(targets, narrowed[d.Roll(len(narrowed))])
		} else {
			targets = append(targets, enemies[0])
		}
	// if the attacker is using INT, they get to select INTMOD targets at random
	case "INT":
		for i := 0; i < max.Mod; i++ {
			if len(enemies) > 0 {
				x := d.Roll(len(enemies)) - 1
				targets = append(targets, enemies[x])
				enemies[x] = enemies[len(enemies)-1]
				//enemies[len(enemies)-1] = Being{}
				enemies = enemies[:len(enemies)-1]
			}
		}
	default:
		x := d.Roll(len(enemies)) -1
		targets = append(targets, enemies[x])
	}
	return targets
}

// Attack takes an array of enemies and deals damage or converts them
func (b *Being) Attack(e *Being) []bool {
	result := []bool{false, false, false}
	d := dice.Dice{}
	max := b.MaxStat()
	fmt.Printf("%s (%s %d) attacks %s (%s %d)\n", b.Name, max.Name, max.Value, e.Name, max.Name, e.GetStat(max.Name).Value)
	if max.Value >= e.GetStat(max.Name).Value {
		result[0] = true
		fmt.Println("It's a hit!")
		// Charisma check
		if max.Name == "CHA" && d.Roll(100) > (100-max.Mod*5) {
			e.DefectTo(b.Team)
			result[2] = true
			fmt.Printf("%s has been converted to Team %s\n", e.Name, e.Team.Name)
		} else if max.Name == "STR" {
			e.HitPoints.Current -= max.Mod
			fmt.Printf("%d damage! %s is at %d/%d\n", max.Mod, e.Name, e.HitPoints.Current, e.HitPoints.Max)
		} else {
			e.HitPoints.Current--
			fmt.Printf("%s is at %d/%d\n", e.Name, e.HitPoints.Current, e.HitPoints.Max)
		}
		if(e.IsAlive() != true) { 
			result[1] = true
			fmt.Printf("%s has been killed!\n", e.Name)
		}
	}
	return result
}
