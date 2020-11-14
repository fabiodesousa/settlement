package main

import (
	"fmt"

	"github.com/fabiodesousa/settlement/being"
)

func main() {
	villagers := being.Team{Name: "villagers", Roster: make([]*being.Being, 0)}
	enemies := being.Team{Name: "bandits", Roster: make([]*being.Being, 0)}
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
	for (len(being.GetSurvivors(encounter.Attackers.Roster)) > 0 && len(being.GetSurvivors(encounter.Defenders.Roster)) > 0) {
		encounter.TurnCount++
		for _, b := range encounter.Initiative {
			if(b.IsAlive() != true) {
				fmt.Printf("%s is dead. Continuing.\n", b.Name)
				continue
			}
			var target []*being.Being
			if b.Team.Name == "villagers" {
				target = b.SelectTarget(encounter.Attackers.Roster)
			} else {
				target = b.SelectTarget(encounter.Defenders.Roster)
			}
			fmt.Printf("%s (%s %d, hp %d/%d) from team %s is attacking: ", b.Name, b.MaxStat().Name, b.MaxStat().Value, b.HitPoints.Current, b.HitPoints.Max, b.Team.Name)
			for i, t := range target {
				if(i > 0) {
					fmt.Printf(", ")
				}
				fmt.Printf("%s (%s %d, hp %d/%d)",t.Name, b.MaxStat().Name, t.GetStat(b.MaxStat().Name).Value, t.HitPoints.Current, t.HitPoints.Max)
			}
			fmt.Printf("\n")
			
			b.Attack(target)
			if (len(being.GetSurvivors(encounter.Attackers.Roster)) == 0 || len(being.GetSurvivors(encounter.Defenders.Roster)) == 0) {
				break
			}
		}
	}
	fmt.Printf("Encounter finished in %d turns", encounter.TurnCount)
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
