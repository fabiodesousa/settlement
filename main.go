package main

import (
	"fmt"

	"github.com/fabiodesousa/settlement/being"
	//"github.com/fabiodesousa/settlement/combat"
)

func main() {
	var villagers []being.Being
	var enemies []being.Being
	for i := 0; i < 5; i++ {
		newVillager := being.NewRandomBeing()
		newVillager.AssignTeam("villagers")
		newEnemy := being.NewRandomBeing()
		newEnemy.AssignTeam("bandits")
		villagers = append(villagers, newVillager)
		enemies = append(enemies, newEnemy)
	}
	targets := villagers[0].SelectTarget(enemies)
	fmt.Print("Attacker: ")
	being.PrintStats(villagers[0])
	fmt.Println("Defender(s):")
	for _, target := range targets {
		being.PrintStats(target)
	}
	villagers[0].Attack(targets)
}
