package being

import "testing"

func TestIsAlive(t *testing.T) {
	b := Being{
		HitPoints: HitPoints{1, 1},
	}
	if b.IsAlive() != true {
		t.Errorf("Expected IsAlive to return true, returned %t", b.IsAlive())
	}
	b.HitPoints.Current--
	if b.IsAlive() {
		t.Errorf("Expected IsAlive to return false, returned %t", b.IsAlive())
	}
}

func TestSetInitiative(t *testing.T) {
	b := Being{}
	b.SetInitiative(1.1)
	if b.Initiative != 1.1 {
		t.Errorf("Expected initiative to be 1.1, got %f", b.Initiative)
	}
}

