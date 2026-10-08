package game

import (
	"testing"
	"time"
)

func TestWeaponCatalogHasPlayableWeapons(t *testing.T) {
	// Lekcia 35: Dopln dalsiu kontrolu, ktoru musi splnit kazda nova zbran.
	catalog := weaponCatalog()
	if len(catalog) == 0 {
		t.Fatal("weapon catalog must contain at least one weapon")
	}
	if len(catalog) > 9 {
		t.Fatalf("keys 1-9 can select at most 9 weapons, got %d", len(catalog))
	}

	names := make(map[string]bool)
	for index, weapon := range catalog {
		if weapon.Name == "" {
			t.Fatalf("weapon %d needs a name", index)
		}
		if names[weapon.Name] {
			t.Fatalf("weapon name %q is used twice", weapon.Name)
		}
		names[weapon.Name] = true
		if weapon.MagazineSize <= 0 || weapon.BulletSpeed <= 0 || weapon.ExplosionRadius <= 0 || weapon.Damage <= 0 {
			t.Fatalf("weapon %q needs positive magazine, speed, radius and damage", weapon.Name)
		}
		if weapon.FireInterval < 0 || weapon.ReloadDuration < 0 {
			t.Fatalf("weapon %q cannot have negative times", weapon.Name)
		}
	}
}

func testArsenalHub(t *testing.T) *Hub {
	t.Helper()
	gameMap := &Map{
		Name:   "test",
		Width:  12,
		Height: 5,
		Tiles: [][]Tile{
			[]Tile("XXXXXXXXXXXX"),
			[]Tile("X..........X"),
			[]Tile("X..........X"),
			[]Tile("X..........X"),
			[]Tile("XXXXXXXXXXXX"),
		},
	}
	hub := NewHub(gameMap)
	close(hub.stop)
	cannon := weaponSettings()
	machineGun := WeaponSettings{
		Name:            "Gulomet",
		Color:           "#7fd3ff",
		FireInterval:    0,
		MagazineSize:    10,
		ReloadDuration:  time.Second,
		BulletSpeed:     8,
		ExplosionRadius: 0.5,
		Damage:          10,
	}
	hub.arsenal = []WeaponSettings{cannon, machineGun}
	hub.Join("p1", "Ana")
	hub.Join("p2", "Boris")
	return hub
}

func TestPlayerSwitchesWeapons(t *testing.T) {
	hub := testArsenalHub(t)

	hub.Handle(Command{PlayerID: "p1", Message: ClientMessage{Type: "weapon"}})
	if hub.tanks["p1"].Weapon != 1 {
		t.Fatalf("weapon without number should switch to the next weapon, got %d", hub.tanks["p1"].Weapon)
	}
	hub.Handle(Command{PlayerID: "p1", Message: ClientMessage{Type: "weapon"}})
	if hub.tanks["p1"].Weapon != 0 {
		t.Fatal("switching after the last weapon should go back to the first")
	}

	second := 1
	hub.Handle(Command{PlayerID: "p1", Message: ClientMessage{Type: "weapon", Weapon: &second}})
	invalid := 7
	hub.Handle(Command{PlayerID: "p1", Message: ClientMessage{Type: "weapon", Weapon: &invalid}})
	if hub.tanks["p1"].Weapon != 1 {
		t.Fatal("weapon number should select a weapon and invalid numbers should be ignored")
	}

	status := hub.snapshotLocked().Weapons["p1"]
	if status.Name != "Gulomet" || status.Index != 1 || status.Count != 2 || status.MagazineSize != 10 || status.Damage != 10 {
		t.Fatalf("weapon status should describe the selected weapon, got %+v", status)
	}
}

func TestEachWeaponHasOwnMagazineAndDamage(t *testing.T) {
	hub := testArsenalHub(t)
	hub.tanks["p1"] = Tank{ID: "p1", X: 1.5, Y: 2.5, Alive: true, Connected: true, Health: tankMaxHealth, MaxHealth: tankMaxHealth}
	cannon := hub.arsenal[0]

	hub.weapons[weaponSlot{PlayerID: "p1", Weapon: 0}] = weaponState{
		shotsFired:     cannon.MagazineSize,
		lastShotAt:     time.Now(),
		reloadingUntil: time.Now().Add(time.Hour),
	}
	second := 1
	hub.Handle(Command{PlayerID: "p1", Message: ClientMessage{Type: "weapon", Weapon: &second}})
	hub.Handle(Command{PlayerID: "p1", Message: ClientMessage{Type: "shoot"}})

	if len(hub.bullets) != 1 {
		t.Fatalf("second weapon should shoot while the first one reloads, got %d bullets", len(hub.bullets))
	}
	for _, bullet := range hub.bullets {
		if bullet.Weapon != 1 || bullet.Color != "#7fd3ff" {
			t.Fatalf("bullet should remember its weapon, got %+v", bullet)
		}
		hub.explode(&bullet, 3, 2.5, hub.weaponSettingsFor(bullet.Weapon))
	}
	if hub.tanks["p2"].Health != tankMaxHealth-10 {
		t.Fatalf("explosion should use damage of the selected weapon, got health %d", hub.tanks["p2"].Health)
	}
}
