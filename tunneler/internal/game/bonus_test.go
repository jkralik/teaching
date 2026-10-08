package game

import (
	"math"
	"testing"
	"time"
)

// Lekcia 37: Tento test skontroluje, ze kazdy bonus v katalogu ma rozumne nastavenia.
func TestBonusCatalogIsPlayable(t *testing.T) {
	catalog := bonusCatalog()
	normal := 0
	for _, bonus := range catalog {
		if bonus.Name == "" || bonus.Symbol == "" || bonus.Color == "" {
			t.Fatalf("bonus needs name, symbol and color: %+v", bonus)
		}
		if bonus.SpeedFactor < 0 || bonus.DamageFactor < 0 || bonus.Duration < 0 {
			t.Fatalf("bonus %q has negative settings", bonus.Name)
		}
		if bonus.Duration > 20*time.Second {
			t.Fatalf("bonus %q lasts too long: %s", bonus.Name, bonus.Duration)
		}
		if !bonus.MysteryOnly {
			normal++
		}
	}
	if normal == 0 {
		t.Fatal("at least one bonus should appear on the map without the ? tile")
	}
}

func testBonusHub(t *testing.T) *Hub {
	t.Helper()
	hub := testArsenalHub(t)
	hub.bonusCatalog = []BonusSettings{
		{Name: "Okamzity reload", Symbol: "R", Color: "#5ec8f2", InstantReload: true},
		{Name: "Turbo", Symbol: "T", Color: "#7dff8a", Duration: 5 * time.Second, SpeedFactor: 2},
		{Name: "Silna strela", Symbol: "D", Color: "#ff7a59", Duration: 5 * time.Second, DamageFactor: 2},
		{Name: "Nesmrtelnost", Symbol: "N", Color: "#fff27a", Duration: 5 * time.Second, Invulnerable: true},
	}
	return hub
}

// placeBonus polozi bonus na policko (2,2), kde stoji tank p1.
func placeBonus(hub *Hub, id string, bonus int) {
	hub.bonuses[id] = BonusItem{ID: id, X: 2, Y: 2, Bonus: bonus, ExpiresAt: time.Now().Add(time.Minute)}
}

func moveRight(hub *Hub, playerID string) {
	angle := 0.0
	hub.Handle(Command{PlayerID: playerID, Message: ClientMessage{Type: "move", Direction: Right, Angle: &angle}})
}

func TestReloadBonusRefillsWeapon(t *testing.T) {
	hub := testBonusHub(t)
	hub.weapons[weaponSlot{PlayerID: "p1", Weapon: 0}] = weaponState{
		shotsFired:     3,
		reloadingUntil: time.Now().Add(time.Minute),
	}
	placeBonus(hub, "r", 0)

	moveRight(hub, "p1")

	if len(hub.bonuses) != 0 {
		t.Fatal("tank should pick up the bonus")
	}
	status := hub.snapshotLocked().Weapons["p1"]
	if status.Reloading || status.ShotsRemaining != status.MagazineSize {
		t.Fatalf("reload bonus should refill the weapon, got %+v", status)
	}
}

func TestTankHasSameBonusOnlyOnce(t *testing.T) {
	hub := testBonusHub(t)
	placeBonus(hub, "first", 1)
	moveRight(hub, "p1")
	if !hub.tanks["p1"].hasBonus(1, time.Now()) {
		t.Fatal("tank should get the Turbo bonus")
	}

	hub.bonuses["second"] = BonusItem{ID: "second", X: 3, Y: 2, Bonus: 1, ExpiresAt: time.Now().Add(time.Minute)}
	delete(hub.tanks, "p2")
	for range 5 {
		moveRight(hub, "p1")
	}
	if _, ok := hub.bonuses["second"]; !ok {
		t.Fatal("tank already has Turbo, the second one should stay on the map")
	}
	if active := hub.snapshotLocked().Tanks["p1"].Bonuses; len(active) != 1 || active[0].Name != "Turbo" {
		t.Fatalf("snapshot should list one active Turbo, got %+v", active)
	}
}

func TestSpeedBonusMakesTankFaster(t *testing.T) {
	hub := testBonusHub(t)
	delete(hub.tanks, "p2")
	tank := hub.tanks["p1"]
	tank.bonuses = map[int]time.Time{1: time.Now().Add(time.Minute)}
	hub.tanks["p1"] = tank
	startX := tank.X

	moveRight(hub, "p1")

	moved := hub.tanks["p1"].X - startX
	if math.Abs(moved-tankMoveStep*2) > 1e-9 {
		t.Fatalf("Turbo should double the step, moved %f", moved)
	}
}

func TestBonusesChangeDamage(t *testing.T) {
	hub := testBonusHub(t)
	shooter := hub.tanks["p2"]
	shooter.bonuses = map[int]time.Time{2: time.Now().Add(time.Minute)}
	hub.tanks["p2"] = shooter
	target := hub.tanks["p1"]
	settings := weaponSettings()

	bullet := Bullet{ID: "b1", OwnerID: "p2", DamageFactor: hub.tankDamageFactor(shooter, time.Now())}
	hub.explode(&bullet, target.X, target.Y, settings)
	if got := hub.tanks["p1"].Health; got != tankMaxHealth-2*settings.Damage {
		t.Fatalf("strong bullet should deal double damage, health %d", got)
	}

	target = hub.tanks["p1"]
	target.bonuses = map[int]time.Time{3: time.Now().Add(time.Minute)}
	hub.tanks["p1"] = target
	health := target.Health
	hub.explode(&Bullet{ID: "b2", OwnerID: "p2"}, target.X, target.Y, settings)
	if hub.tanks["p1"].Health != health {
		t.Fatal("invulnerable tank should not lose health")
	}
	if !hub.snapshotLocked().Tanks["p1"].Invulnerable {
		t.Fatal("snapshot should mark the invulnerable tank")
	}
}

func TestMysteryTileGivesRandomEffect(t *testing.T) {
	hub := testBonusHub(t)
	hub.bonusCatalog = []BonusSettings{
		{Name: "Blato", Symbol: "B", Color: "#8a6b4a", Duration: 5 * time.Second, SpeedFactor: 0.5, MysteryOnly: true},
	}
	placeBonus(hub, "mystery", mysteryBonus)

	moveRight(hub, "p1")

	if !hub.tanks["p1"].hasBonus(0, time.Now()) {
		t.Fatal("? tile should give a random bonus, here the only one is Blato")
	}
}

func TestBonusSpawnsOnEmptyTileAndExpires(t *testing.T) {
	hub := testBonusHub(t)
	now := hub.nextBonusAt

	if !hub.updateBonusesLocked(now) || len(hub.bonuses) != 1 {
		t.Fatalf("a bonus should appear after %s", bonusSpawnInterval)
	}
	for _, item := range hub.bonuses {
		if hub.gameMap.TileAt(item.X, item.Y) != TileEmpty {
			t.Fatalf("bonus must lie on an empty tile, got %c", hub.gameMap.TileAt(item.X, item.Y))
		}
	}

	hub.nextBonusAt = now.Add(time.Hour)
	hub.updateBonusesLocked(now.Add(bonusLifetime))
	if len(hub.bonuses) != 0 {
		t.Fatal("bonus should disappear after its lifetime")
	}
}
