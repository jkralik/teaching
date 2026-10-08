package game

import (
	"math"
	"math/rand/v2"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestGenerateMapHasSolidBorder(t *testing.T) {
	// Lekcia 26: Test je miesto, kam pridame ocakavanie pre nove pravidlo mapy.
	gameMap := GenerateMap("test", 12, 8)

	for x := 0; x < gameMap.Width; x++ {
		if gameMap.TileAt(x, 0) != TileRock || gameMap.TileAt(x, gameMap.Height-1) != TileRock {
			t.Fatalf("top and bottom borders must be rock")
		}
	}
	for y := 0; y < gameMap.Height; y++ {
		if gameMap.TileAt(0, y) != TileRock || gameMap.TileAt(gameMap.Width-1, y) != TileRock {
			t.Fatalf("left and right borders must be rock")
		}
	}
}

func TestTankDigsDirtBeforeMovingIntoTile(t *testing.T) {
	// Lekcia 16: Tento test ukazuje, ze kopanie meni stav hubu a tank sa este nepohne.
	gameMap := &Map{
		Name:   "test",
		Width:  5,
		Height: 5,
		Tiles: [][]Tile{
			[]Tile("XXXXX"),
			[]Tile("X.#.X"),
			[]Tile("X...X"),
			[]Tile("X...X"),
			[]Tile("XXXXX"),
		},
	}
	hub := NewHub(gameMap)
	tank := Tank{ID: "p1", X: 1.95, Y: 1.5, Direction: Right, Alive: true}
	for range digMoves() {
		hub.moveTank(&tank)
	}

	if tank.X != 1.95 || tank.Y != 1.5 {
		t.Fatalf("tank should dig first and stay in place, got %f,%f", tank.X, tank.Y)
	}
	if gameMap.TileAt(2, 1) != TileDirt {
		t.Fatalf("original map should not be mutated")
	}
	if hub.gameMap.TileAt(2, 1) != TileEmpty {
		t.Fatalf("hub map tile should be dug out")
	}
}

func TestTankMovesDiagonally(t *testing.T) {
	gameMap := &Map{
		Name:   "test",
		Width:  5,
		Height: 5,
		Tiles: [][]Tile{
			[]Tile("XXXXX"),
			[]Tile("X...X"),
			[]Tile("X...X"),
			[]Tile("X...X"),
			[]Tile("XXXXX"),
		},
	}
	hub := NewHub(gameMap)
	tank := Tank{ID: "p1", X: 1.5, Y: 2.5, Angle: math.Pi / 4, Alive: true}
	hub.moveTank(&tank)

	if math.Abs(tank.X-1.5-math.Sqrt2*tankMoveStep/2) > 1e-9 ||
		math.Abs(tank.Y-2.5+math.Sqrt2*tankMoveStep/2) > 1e-9 {
		t.Fatalf("tank should move diagonally, got %f,%f", tank.X, tank.Y)
	}
}

func TestBulletMovesAtItsAngle(t *testing.T) {
	gameMap := &Map{
		Name:   "test",
		Width:  5,
		Height: 5,
		Tiles: [][]Tile{
			[]Tile("XXXXX"),
			[]Tile("X...X"),
			[]Tile("X...X"),
			[]Tile("X...X"),
			[]Tile("XXXXX"),
		},
	}
	hub := NewHub(gameMap)
	bullet := Bullet{X: 1.5, Y: 2.5, Angle: math.Pi / 4}
	bulletStep := weaponSettings().BulletSpeed * gameTickInterval.Seconds()

	if !hub.moveBullet(&bullet) {
		t.Fatal("diagonal bullet should keep moving through empty tiles")
	}
	if math.Abs(bullet.X-1.5-math.Sqrt2*bulletStep/2) > 1e-9 ||
		math.Abs(bullet.Y-2.5+math.Sqrt2*bulletStep/2) > 1e-9 {
		t.Fatalf("bullet should travel diagonally, got %f,%f", bullet.X, bullet.Y)
	}
}

func TestWeaponSettingsControlCadenceAndReload(t *testing.T) {
	settings := weaponSettings()
	if settings.MagazineSize != 3 || settings.FireInterval <= 0 || settings.ReloadDuration <= 0 ||
		settings.BulletSpeed <= 0 || settings.ExplosionRadius <= 0 {
		t.Fatalf("weapon settings should contain positive gameplay values: %+v", settings)
	}

	now := time.Now()
	var state weaponState
	for shot := 0; shot < settings.MagazineSize; shot++ {
		shotTime := now.Add(time.Duration(shot) * settings.FireInterval)
		if !state.tryShoot(shotTime, settings) {
			t.Fatalf("shot %d should be allowed", shot+1)
		}
		if shot > 0 && state.tryShoot(shotTime.Add(settings.FireInterval-time.Nanosecond), settings) {
			t.Fatal("a shot before the cadence interval should be rejected")
		}
	}

	reloadEnds := now.Add(time.Duration(settings.MagazineSize-1)*settings.FireInterval + settings.ReloadDuration)
	if state.tryShoot(reloadEnds.Add(-time.Nanosecond), settings) {
		t.Fatal("shooting should be blocked until reload completes")
	}
	reloadingStatus := state.status(reloadEnds.Add(-500*time.Millisecond), settings)
	if !reloadingStatus.Reloading || reloadingStatus.ShotsRemaining != 0 ||
		reloadingStatus.ReloadRemainingMs != 500 {
		t.Fatalf("HUD should report reload progress: %+v", reloadingStatus)
	}
	readyStatus := state.status(reloadEnds, settings)
	if readyStatus.Reloading || readyStatus.ShotsRemaining != settings.MagazineSize {
		t.Fatalf("HUD should report a full magazine after reloading: %+v", readyStatus)
	}
	if !state.tryShoot(reloadEnds, settings) {
		t.Fatal("shooting should be allowed when reload completes")
	}
}

func TestExplosionRadiusDamagesNearbyTanks(t *testing.T) {
	gameMap := &Map{
		Name:   "test",
		Width:  6,
		Height: 5,
		Tiles: [][]Tile{
			[]Tile("XXXXXX"),
			[]Tile("X....X"),
			[]Tile("X....X"),
			[]Tile("X....X"),
			[]Tile("XXXXXX"),
		},
	}
	hub := NewHub(gameMap)
	close(hub.stop)
	hub.tanks["owner"] = Tank{ID: "owner", X: 1.5, Y: 2.5, Alive: true, Health: tankMaxHealth, MaxHealth: tankMaxHealth}
	hub.tanks["near"] = Tank{ID: "near", X: 3, Y: 2.5, Alive: true, Health: tankMaxHealth, MaxHealth: tankMaxHealth}
	hub.tanks["far"] = Tank{ID: "far", X: 4.5, Y: 2.5, Alive: true, Health: tankMaxHealth, MaxHealth: tankMaxHealth}
	settings := weaponSettings()
	settings.Damage = 60
	hub.explode(&Bullet{ID: "shot", OwnerID: "owner"}, 2.5, 2.5, settings)

	if hub.tanks["near"].Health != tankMaxHealth-60 || !hub.tanks["near"].Alive {
		t.Fatalf("first hit should only damage the tank, got health %d", hub.tanks["near"].Health)
	}
	if hub.tanks["far"].Health != tankMaxHealth || hub.tanks["owner"].Health != tankMaxHealth {
		t.Fatal("explosion should damage enemy tanks inside its radius only")
	}

	hub.explode(&Bullet{ID: "shot2", OwnerID: "owner"}, 2.5, 2.5, settings)
	if hub.tanks["near"].Alive || hub.tanks["near"].Health != 0 {
		t.Fatal("tank should be destroyed when its health reaches zero")
	}
	if hub.tanks["owner"].Score != 7 {
		t.Fatalf("owner should get 1 point per hit and 5 for destroying a tank, got %d", hub.tanks["owner"].Score)
	}
	if len(hub.snapshotLocked().Explosions) != 2 {
		t.Fatal("explosions should be included in the game snapshot")
	}
}

func TestDestroyedTankRespawnsAfterThreeSeconds(t *testing.T) {
	gameMap := &Map{
		Name:   "test",
		Width:  6,
		Height: 5,
		Tiles: [][]Tile{
			[]Tile("XXXXXX"),
			[]Tile("X....X"),
			[]Tile("X....X"),
			[]Tile("X....X"),
			[]Tile("XXXXXX"),
		},
	}
	hub := NewHub(gameMap)
	defer close(hub.stop)
	hub.tanks["owner"] = Tank{ID: "owner", X: 1.5, Y: 2.5, Alive: true, Health: tankMaxHealth, MaxHealth: tankMaxHealth}
	hub.tanks["target"] = Tank{ID: "target", X: 3.5, Y: 2.5, Alive: true, Health: tankMaxHealth, MaxHealth: tankMaxHealth}
	settings := weaponSettings()
	settings.Damage = tankMaxHealth

	hub.explode(&Bullet{ID: "lethal", OwnerID: "owner"}, 3.5, 2.5, settings)

	target := hub.tanks["target"]
	if target.Alive || target.Health != 0 {
		t.Fatal("lethal hit should destroy the tank")
	}
	if delay := time.Until(target.respawnAt); delay < tankRespawnDelay-time.Second/10 || delay > tankRespawnDelay {
		t.Fatalf("tank should be scheduled to respawn in 3 seconds, got %v", delay)
	}

	hub.tick()
	if hub.tanks["target"].Alive {
		t.Fatal("tank should remain destroyed until the respawn delay has elapsed")
	}

	target = hub.tanks["target"]
	target.respawnAt = time.Now().Add(-time.Millisecond)
	hub.tanks["target"] = target
	hub.tick()

	target = hub.tanks["target"]
	if !target.Alive || target.Health != tankMaxHealth {
		t.Fatalf("tank should respawn with full health, got alive=%t health=%d", target.Alive, target.Health)
	}
	if hub.gameMap.TileAt(int(math.Floor(target.X)), int(math.Floor(target.Y))) != TileEmpty {
		t.Fatalf("tank should respawn on an empty tile, got %v,%v", target.X, target.Y)
	}
	if target.X == 1.5 && target.Y == 2.5 {
		t.Fatal("tank should not respawn on the other player's tile")
	}
}

func TestFindSpawnPositionChoosesRandomFreeTile(t *testing.T) {
	gameMap := &Map{
		Name:   "test",
		Width:  6,
		Height: 5,
		Tiles: [][]Tile{
			[]Tile("XXXXXX"),
			[]Tile("X....X"),
			[]Tile("X....X"),
			[]Tile("X....X"),
			[]Tile("XXXXXX"),
		},
	}
	hub := NewHub(gameMap)
	defer close(hub.stop)
	hub.random = rand.New(rand.NewPCG(1, 2))
	hub.tanks["other"] = Tank{ID: "other", X: 1.5, Y: 2.5, Alive: true}

	positions := make(map[[2]float64]bool)
	for range 20 {
		x, y, ok := hub.findSpawnPositionLocked("respawning")
		if !ok {
			t.Fatal("expected to find a free spawn position")
		}
		if hub.gameMap.TileAt(int(math.Floor(x)), int(math.Floor(y))) != TileEmpty {
			t.Fatalf("spawn position %v,%v is not empty", x, y)
		}
		if x == 1.5 && y == 2.5 {
			t.Fatal("spawn position should not be occupied by another tank")
		}
		positions[[2]float64{x, y}] = true
	}
	if len(positions) < 2 {
		t.Fatalf("expected random selection to choose multiple positions, got %v", positions)
	}
}

func TestDisconnectedPlayerCanResumeTankAndExpiresAfterFiveMinutes(t *testing.T) {
	gameMap := &Map{
		Name:   "test",
		Width:  6,
		Height: 5,
		Tiles: [][]Tile{
			[]Tile("XXXXXX"),
			[]Tile("X....X"),
			[]Tile("X....X"),
			[]Tile("X....X"),
			[]Tile("XXXXXX"),
		},
	}
	hub := NewHub(gameMap)
	defer close(hub.stop)
	hub.arsenal = []WeaponSettings{weaponSettings()}

	oldConnection := hub.Join("old-id", "Alice")
	<-oldConnection
	tank := hub.tanks["old-id"]
	tank.Score = 42
	tank.Health = 65
	tank.Armor = 1
	hub.tanks["old-id"] = tank
	weapon := weaponState{shotsFired: 1}
	hub.weapons[weaponSlot{PlayerID: "old-id", Weapon: 0}] = weapon
	hub.Leave("old-id")

	disconnected := hub.tanks["old-id"]
	if disconnected.Connected || disconnected.disconnectedAt.IsZero() {
		t.Fatal("leaving should mark the existing tank as disconnected")
	}
	hub.Handle(Command{PlayerID: "old-id", Message: ClientMessage{Type: "shoot"}})
	if len(hub.bullets) != 0 || hub.weapons[weaponSlot{PlayerID: "old-id", Weapon: 0}] != weapon {
		t.Fatal("disconnected players should not be able to shoot")
	}

	newConnection := hub.Join("new-id", "aLiCe")
	<-newConnection
	resumed, ok := hub.tanks["new-id"]
	if !ok || resumed.ID != "new-id" || !resumed.Connected {
		t.Fatal("reconnecting with the same name should restore the existing tank")
	}
	if resumed.Score != 42 || resumed.Health != 65 || resumed.Armor != 1 ||
		resumed.X != disconnected.X || resumed.Y != disconnected.Y {
		t.Fatalf("reconnection should preserve tank state, got %+v", resumed)
	}
	if hub.weapons[weaponSlot{PlayerID: "new-id", Weapon: 0}] != weapon {
		t.Fatal("reconnection should preserve weapon state")
	}
	if _, ok := hub.tanks["old-id"]; ok {
		t.Fatal("reconnected tank should no longer use its old connection ID")
	}

	hub.Leave("new-id")
	expired := hub.tanks["new-id"]
	expired.disconnectedAt = time.Now().Add(-tankRetentionDuration - time.Second)
	hub.tanks["new-id"] = expired
	hub.tick()
	if _, ok := hub.tanks["new-id"]; ok {
		t.Fatal("disconnected tank should be deleted after five minutes")
	}
	if _, ok := hub.weapons[weaponSlot{PlayerID: "new-id", Weapon: 0}]; ok {
		t.Fatal("expired tank should have its weapon state deleted")
	}
}

func TestTanksCannotMoveThroughEachOther(t *testing.T) {
	gameMap := &Map{
		Name:   "test",
		Width:  6,
		Height: 5,
		Tiles: [][]Tile{
			[]Tile("XXXXXX"),
			[]Tile("X....X"),
			[]Tile("X....X"),
			[]Tile("X....X"),
			[]Tile("XXXXXX"),
		},
	}
	hub := NewHub(gameMap)
	defer close(hub.stop)
	hub.tanks["other"] = Tank{ID: "other", X: 2.3, Y: 2.5, Alive: true}
	tank := Tank{ID: "moving", X: 1.5, Y: 2.5, Angle: 0, Alive: true}

	hub.moveTank(&tank)
	if tank.X != 1.5 || tank.Y != 2.5 {
		t.Fatalf("tank should stop before overlapping another tank, got %v,%v", tank.X, tank.Y)
	}
	tank.Angle = math.Pi / 2
	hub.moveTank(&tank)
	if tank.Y >= 2.5 {
		t.Fatalf("tank should be able to move around another tank, got %v,%v", tank.X, tank.Y)
	}
}

func TestMapFilesExpandTilesToSmallerSquares(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.txt")
	if err := os.WriteFile(path, []byte("XXX\nX.X\nXXX\n"), 0600); err != nil {
		t.Fatal(err)
	}
	gameMap, err := LoadMap(path)
	if err != nil {
		t.Fatal(err)
	}
	if gameMap.Width != 6 || gameMap.Height != 6 {
		t.Fatalf("expected 6x6 expanded map, got %dx%d", gameMap.Width, gameMap.Height)
	}
	if gameMap.TileAt(0, 0) != TileRock || gameMap.TileAt(2, 2) != TileEmpty {
		t.Fatal("expanded map should preserve its tile pattern")
	}
	if gameMap.TileAt(2, 3) != TileEmpty {
		t.Fatal("each source tile should expand across two rows")
	}
}

// Lekcia 16: Zem sa neodkope hned, tank ju musi chvilu kopat.
func TestTankNeedsTimeToDigDirt(t *testing.T) {
	gameMap := &Map{
		Name:   "test",
		Width:  5,
		Height: 3,
		Tiles: [][]Tile{
			[]Tile("XXXXX"),
			[]Tile("X.##X"),
			[]Tile("XXXXX"),
		},
	}
	hub := NewHub(gameMap)
	close(hub.stop)
	tank := Tank{ID: "p1", X: 1.9, Y: 1.5, Angle: 0, Alive: true}

	for range digMoves() - 1 {
		hub.moveTank(&tank)
	}
	if hub.gameMap.TileAt(2, 1) != TileDirt || tank.X != 1.9 || tank.Score != 0 {
		t.Fatalf("dirt should still be there while digging, tile=%c x=%f score=%d", hub.gameMap.TileAt(2, 1), tank.X, tank.Score)
	}
	hub.moveTank(&tank)
	if hub.gameMap.TileAt(2, 1) != TileEmpty || tank.Score != 1 {
		t.Fatalf("dirt should be dug after %d moves", digMoves())
	}
	hub.moveTank(&tank)
	if tank.X <= 1.9 {
		t.Fatalf("tank should drive into the dug tunnel, got x=%f", tank.X)
	}
}

func digMoves() int {
	return int(math.Ceil(dirtHardness/tankMoveStep - 1e-9))
}

// Server obmedzi rychlost tanku, aj keby klient posielal "move" prilis casto.
func TestServerLimitsHowOftenTankCanMove(t *testing.T) {
	gameMap := &Map{
		Name:   "test",
		Width:  30,
		Height: 3,
		Tiles: [][]Tile{
			[]Tile("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"),
			[]Tile("X............................X"),
			[]Tile("XXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"),
		},
	}
	hub := NewHub(gameMap)
	close(hub.stop)
	hub.tanks["p1"] = Tank{ID: "p1", X: 1.5, Y: 1.5, Alive: true, Connected: true}

	angle := 0.0
	for range 50 {
		hub.Handle(Command{PlayerID: "p1", Message: ClientMessage{Type: "move", Direction: Right, Angle: &angle}})
	}
	maxSteps := int(tankMoveBurst/tankMoveInterval) + 1
	maxX := 1.5 + float64(maxSteps)*tankMoveStep + 1e-9
	if x := hub.tanks["p1"].X; x <= 1.5 || x > maxX {
		t.Fatalf("50 fast moves should move the tank at most %d steps (x <= %f), got x=%f", maxSteps, maxX, x)
	}

	tank := hub.tanks["p1"]
	tank.nextMoveAt = time.Now().Add(-time.Millisecond)
	hub.tanks["p1"] = tank
	before := tank.X
	hub.Handle(Command{PlayerID: "p1", Message: ClientMessage{Type: "move", Direction: Right, Angle: &angle}})
	if hub.tanks["p1"].X <= before {
		t.Fatal("tank should move again once its move interval passed")
	}
}
