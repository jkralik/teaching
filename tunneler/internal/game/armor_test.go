package game

import (
	"math"
	"testing"
)

func TestArmorCatalogIsPlayable(t *testing.T) {
	// Lekcia 36: Dopln kontrolu, ktoru musi splnit kazdy novy pancier.
	catalog := armorCatalog()
	if len(catalog) == 0 {
		t.Fatal("armor catalog must contain at least one armor")
	}
	names := map[string]bool{}
	for index, armor := range catalog {
		if armor.Name == "" {
			t.Fatalf("armor %d must have a name", index)
		}
		if names[armor.Name] {
			t.Fatalf("armor name %q is used twice", armor.Name)
		}
		names[armor.Name] = true
		if armor.Divisor < 0 {
			t.Fatalf("armor %q must not have negative divisor", armor.Name)
		}
		if armor.SpeedFactor <= 0 || armor.SpeedFactor > 1 {
			t.Fatalf("armor %q speed factor must be in (0, 1], got %v", armor.Name, armor.SpeedFactor)
		}
	}
}

func TestArmorDividesDamage(t *testing.T) {
	cases := []struct {
		damage, divisor, want int
	}{
		{damage: 35, divisor: 0, want: 35},
		{damage: 35, divisor: 1, want: 35},
		{damage: 35, divisor: 2, want: 17},
		{damage: 10, divisor: 50, want: 1},
	}
	for _, tc := range cases {
		if got := armoredDamage(tc.damage, tc.divisor); got != tc.want {
			t.Fatalf("armoredDamage(%d, %d) = %d, want %d", tc.damage, tc.divisor, got, tc.want)
		}
	}
}

func testArmorHub(t *testing.T) *Hub {
	t.Helper()
	hub := testArsenalHub(t)
	hub.armors = []ArmorSettings{
		{Name: "Bez panciera", SpeedFactor: 1},
		{Name: "Bronzovy", Color: "#cd7f32", Divisor: 2, SpeedFactor: 0.5},
	}
	return hub
}

func TestPlayerSwitchesArmor(t *testing.T) {
	hub := testArmorHub(t)

	hub.Handle(Command{PlayerID: "p1", Message: ClientMessage{Type: "armor"}})
	if hub.tanks["p1"].Armor != 1 {
		t.Fatalf("armor command should select next armor, got %d", hub.tanks["p1"].Armor)
	}
	tank := hub.snapshotLocked().Tanks["p1"]
	if tank.ArmorName != "Bronzovy" || tank.ArmorColor != "#cd7f32" || tank.ArmorDivisor != 2 || tank.ArmorSpeed != 0.5 {
		t.Fatalf("snapshot should describe armor for the browser, got %+v", tank)
	}

	hub.Handle(Command{PlayerID: "p1", Message: ClientMessage{Type: "armor"}})
	if hub.tanks["p1"].Armor != 0 {
		t.Fatalf("armor should cycle back to the first one, got %d", hub.tanks["p1"].Armor)
	}
	invalid := 9
	hub.Handle(Command{PlayerID: "p1", Message: ClientMessage{Type: "armor", Armor: &invalid}})
	if hub.tanks["p1"].Armor != 0 {
		t.Fatalf("invalid armor should be ignored, got %d", hub.tanks["p1"].Armor)
	}
}

func TestArmoredTankTakesLessDamage(t *testing.T) {
	hub := testArmorHub(t)
	target := hub.tanks["p2"]
	target.Armor = 1
	hub.tanks["p2"] = target
	cannon := hub.arsenal[0]

	bullet := Bullet{ID: "b1", OwnerID: "p1"}
	hub.explode(&bullet, target.X, target.Y, cannon)

	want := tankMaxHealth - cannon.Damage/2
	if hub.tanks["p2"].Health != want {
		t.Fatalf("armored tank should lose half damage, got health %d, want %d", hub.tanks["p2"].Health, want)
	}
}

func TestArmoredTankMovesSlower(t *testing.T) {
	hub := testArmorHub(t)
	tank := hub.tanks["p1"]
	tank.X, tank.Y, tank.Armor = 1.5, 1.5, 1
	hub.tanks["p1"] = tank

	angle := 0.0
	hub.Handle(Command{PlayerID: "p1", Message: ClientMessage{Type: "move", Direction: Right, Angle: &angle}})

	moved := hub.tanks["p1"].X - 1.5
	if math.Abs(moved-tankMoveStep*0.5) > 1e-9 {
		t.Fatalf("armored tank should move half step, moved %v", moved)
	}
}
