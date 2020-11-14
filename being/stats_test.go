package being

import "testing"

func TestNewStatBlock(t * testing.T) {
	block := NewStatBlock()
	if len(block.StatArray) != 6 {
		t.Errorf("Expected StatBlock array length to be 6, got %d", len(block.StatArray))
	}
}

func TestGetStat(t * testing.T) {
	block := StatBlock{
		StatArray: []Stat{
			Stat{"STR", 10, 0},
			Stat{"DEX", 12, 1},
			Stat{"CON", 14, 2},
			Stat{"INT", 16, 3},
			Stat{"WIS", 18, 4},
			Stat{"CHA", 20, 5},
		},
	}
	str := block.GetStat("STR")
	if str.Name != "STR" {
		t.Errorf("Expected STR, got %s", str.Name)
	}
	dex := block.GetStat("DEX")
	if dex.Name != "DEX" {
		t.Errorf("Expected DEX, got %s", dex.Name)
	}
	con := block.GetStat("CON")
	if con.Name != "CON" {
		t.Errorf("Expected CON, got %s", con.Name)
	}
	int := block.GetStat("INT")
	if int.Name != "INT" {
		t.Errorf("Expected INT, got %s", int.Name)
	}
	wis := block.GetStat("WIS")
	if wis.Name != "WIS" {
		t.Errorf("Expected WIS, got %s", wis.Name)
	}
	cha := block.GetStat("CHA")
	if cha.Name != "CHA" {
		t.Errorf("Expected CHA, got %s", cha.Name)
	}
}