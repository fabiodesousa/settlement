package main

import (
	"fmt"

	"github.com/fabiodesousa/settlement/being"
)

func main() {
	villagers := being.Team{Name: "villagers", Roster: make([]being.Being, 0)}
	enemies := being.Team{Name: "bandits", Roster: make([]being.Being, 0)}
	encounter := being.Encounter{TurnCount: 0}
	for i := 0; i < 5; i++ {
		newVillager := being.NewRandomBeing()
		//being.PrintStats(newVillager)
		villagers.AddTeamMember(newVillager)
		newEnemy := being.NewRandomBeing()
		enemies.AddTeamMember(newEnemy)
	}
	encounter.Attackers = enemies
	encounter.Defenders = villagers
	encounter.RollInitiative()
	fmt.Printf("Villagers size: %v\n", len(villagers.Roster))

	/*
		targets := villagers[0].SelectTarget(enemies)
		fmt.Print("Attacker: ")
		being.PrintStats(villagers[0])
		fmt.Println("Defender(s):")
		for _, target := range targets {
			being.PrintStats(target)
		}
		villagers[0].Attack(targets)

		villagers[0].RollInitiative()
		fmt.Printf("Initiative: %.2f", villagers[0].Initiative)*/
}
