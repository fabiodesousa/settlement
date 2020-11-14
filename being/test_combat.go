package being

import "testing"

func TestAddTeamMember(t *testing.T) {
	team := Team{}
	b := Being{}
	team.AddTeamMember(&b)
	if len(team.Roster) != 1 {
		t.Errorf("Expected team roster to be len 1, got %d", len(team.Roster))
	}
}

func TestRemoveTeamMember(t *testing.T) {
	b := Being{ID:"uuid"}
	team := Team{Roster: []*Being{&b}}
	team.RemoveTeamMember(&b)
	if len(team.Roster) != 0 {
		t.Errorf("Expected team roster to be len 0, got %d", len(team.Roster))
	}
}

func TestDefectTo(t *testing.T) {
	b := Being{ID:"uuid"}
	team1 := Team{Name:"team1", Roster: []*Being{&b}}
	team2 := Team{Name:"team2"}
	b.DefectTo(&team2)
	if len(team1.Roster) != 0 {
		t.Errorf("Expected team1 roster to be len 0, got %d", len(team1.Roster))
	}
	if len(team2.Roster) != 1 {
		t.Errorf("Expected team2 roster to be len 1, got %d", len(team2.Roster))
	}
}