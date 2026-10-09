package game

import (
	"testing"
)

func testTeamHub(t *testing.T) *Hub {
	t.Helper()
	gameMap := &Map{
		Name:   "test",
		Width:  10,
		Height: 5,
		Tiles: [][]Tile{
			[]Tile("XXXXXXXXXX"),
			[]Tile("X........X"),
			[]Tile("X........X"),
			[]Tile("X........X"),
			[]Tile("XXXXXXXXXX"),
		},
	}
	hub := NewHub(gameMap)
	close(hub.stop)
	hub.teams = []TeamSettings{{Name: "Cerveni", Color: "#f00"}, {Name: "Modri", Color: "#00f"}}
	return hub
}

func TestTeamCatalogIsPlayable(t *testing.T) {
	teams := teamCatalog()
	if len(teams) < 2 {
		t.Fatal("there should be at least two teams")
	}
	for _, team := range teams {
		if team.Name == "" || team.Color == "" {
			t.Fatalf("team needs a name and a color, got %+v", team)
		}
	}
}

func TestPlayersAreBalancedIntoTeams(t *testing.T) {
	hub := testTeamHub(t)
	hub.Join("p1", "Ana")
	hub.Join("p2", "Boris")
	hub.Join("p3", "Cyril")
	hub.JoinTeam("p4", "Dana", 1)
	hub.JoinTeam("p5", "Emil", TeamNone)
	hub.JoinTeam("p6", "Fero", 99)

	if hub.tanks["p1"].Team != 1 || hub.tanks["p2"].Team != 2 || hub.tanks["p3"].Team != 1 {
		t.Fatalf("automatic choice should fill the smallest team, got %d %d %d",
			hub.tanks["p1"].Team, hub.tanks["p2"].Team, hub.tanks["p3"].Team)
	}
	if hub.tanks["p4"].Team != 1 {
		t.Fatal("player should join the chosen team")
	}
	if hub.tanks["p5"].Team != 3 {
		t.Fatalf("player without a selected team should get a personal team, got %d", hub.tanks["p5"].Team)
	}
	if hub.tanks["p6"].Team != 2 {
		t.Fatal("unknown team should fall back to automatic choice among selectable teams")
	}

	snapshot := hub.snapshotLocked()
	if snapshot.Tanks["p2"].TeamName != "Modri" || snapshot.Tanks["p2"].TeamColor != "#00f" {
		t.Fatalf("snapshot should contain team name and color, got %+v", snapshot.Tanks["p2"])
	}
	if len(snapshot.Teams) != 3 || snapshot.Teams[0].Players != 3 ||
		snapshot.Teams[1].Players != 2 || snapshot.Teams[2].Name != "Emil" ||
		snapshot.Teams[2].Players != 1 {
		t.Fatalf("snapshot should count players in teams, got %+v", snapshot.Teams)
	}
}

func TestTeamScoreIsSumOfPlayers(t *testing.T) {
	hub := testTeamHub(t)
	hub.tanks["a"] = Tank{ID: "a", Team: 1, Score: 3}
	hub.tanks["b"] = Tank{ID: "b", Team: 1, Score: 4}
	hub.tanks["c"] = Tank{ID: "c", Team: 2, Score: 5}
	hub.tanks["d"] = Tank{ID: "d", Team: TeamNone, Score: 100}

	teams := hub.teamStatusesLocked()
	if teams[0].Score != 7 || teams[1].Score != 5 {
		t.Fatalf("team score should be the sum of its players, got %+v", teams)
	}
}

func TestExplosionDoesNotHurtTeammates(t *testing.T) {
	hub := testTeamHub(t)
	hub.tanks["owner"] = Tank{ID: "owner", Team: 1, X: 1.5, Y: 2.5, Alive: true, Health: tankMaxHealth}
	hub.tanks["friend"] = Tank{ID: "friend", Team: 1, X: 3, Y: 2.5, Alive: true, Health: tankMaxHealth}
	hub.tanks["enemy"] = Tank{ID: "enemy", Team: 2, X: 3, Y: 3, Alive: true, Health: tankMaxHealth}
	settings := weaponSettings()
	settings.Damage = 10

	hub.explode(&Bullet{ID: "shot", OwnerID: "owner", Team: 1}, 3, 2.75, settings)

	if hub.tanks["friend"].Health != tankMaxHealth {
		t.Fatal("explosion should not hurt a teammate")
	}
	if hub.tanks["enemy"].Health != tankMaxHealth-10 {
		t.Fatal("explosion should hurt an enemy")
	}
	if hub.tanks["owner"].Score != 1 {
		t.Fatalf("owner should get points only for the enemy, got %d", hub.tanks["owner"].Score)
	}
}

func TestBulletFliesThroughTeammate(t *testing.T) {
	hub := testTeamHub(t)
	hub.tanks["friend"] = Tank{ID: "friend", Team: 1, X: 3.5, Y: 2.5, Alive: true, Health: tankMaxHealth}
	bullet := Bullet{ID: "shot", OwnerID: "owner", Team: 1, X: 2.5, Y: 2.5}

	for step := 0; step < 50 && bullet.X <= 3.5; step++ {
		if !hub.moveBullet(&bullet) {
			t.Fatal("bullet should fly through a teammate")
		}
	}
	if bullet.X <= 3.5 {
		t.Fatalf("bullet should get past the teammate, got x=%.2f", bullet.X)
	}

	hub.tanks["enemy"] = Tank{ID: "enemy", Team: 2, X: bullet.X + 0.5, Y: 2.5, Alive: true, Health: tankMaxHealth}
	if hub.moveBullet(&bullet) {
		t.Fatal("bullet should explode at an enemy")
	}
}

func TestSoloPlayerGetsRandomColor(t *testing.T) {
	hub := testTeamHub(t)
	hub.JoinTeam("solo", "Ana", TeamNone)
	snapshot := hub.snapshotLocked()
	tank := snapshot.Tanks["solo"]
	color := tank.TeamColor
	if len(color) != 7 || color[0] != '#' {
		t.Fatalf("personal team should get a color like #rrggbb, got %q", color)
	}
	if tank.Team == TeamNone || tank.TeamName != "Ana" {
		t.Fatalf("player without a selected team should get a named personal team, got %+v", tank)
	}

	hub.Leave("solo")
	hub.JoinTeam("back", "Ana", TeamNone)
	if resumed := hub.snapshotLocked().Tanks["back"]; resumed.Team != tank.Team || resumed.TeamColor != color {
		t.Fatal("returning player should keep the same team and color")
	}
	hub.RemovePlayer("back")
	if len(hub.teams) != hub.playableTeamCount {
		t.Fatalf("explicitly disconnecting after reconnect should remove the personal team, got %+v", hub.teams)
	}

	hub.JoinTeam("other", "Boris", TeamNone)
	other := hub.tanks["other"]
	if team, ok := hub.teamSettingsFor(other.Team); !ok || team.OwnerID != "other" {
		t.Fatal("each player without a selected team should get a separate team")
	}

	hub.JoinTeam("red", "Boris", 1)
	if hub.tanks["red"].Team != 1 {
		t.Fatal("player who selects a catalog team should join that team")
	}
}

func TestHslToHex(t *testing.T) {
	cases := map[float64]string{0: "#ff0000", 120: "#00ff00", 240: "#0000ff"}
	for hue, want := range cases {
		if got := hslToHex(hue, 1, 0.5); got != want {
			t.Fatalf("hue %.0f: want %s, got %s", hue, want, got)
		}
	}
}

func TestReconnectedPlayerKeepsTeam(t *testing.T) {
	hub := testTeamHub(t)
	hub.JoinTeam("old", "Ana", 2)
	hub.Leave("old")
	hub.JoinTeam("new", "Ana", 1)

	if hub.tanks["new"].Team != 2 {
		t.Fatal("returning player should stay in the original team")
	}
}

func TestRemovePlayerClearsPlayerStateAndPersonalTeam(t *testing.T) {
	hub := testTeamHub(t)
	hub.JoinTeam("solo", "Ana", TeamNone)
	hub.JoinTeam("other", "Boris", TeamNone)
	hub.JoinTeam("red", "Cyril", 1)
	hub.weapons[weaponSlot{PlayerID: "solo", Weapon: 0}] = weaponState{shotsFired: 1}
	hub.bullets["solo-shot"] = Bullet{ID: "solo-shot", OwnerID: "solo", Team: hub.tanks["solo"].Team}

	hub.RemovePlayer("solo")

	if _, ok := hub.tanks["solo"]; ok {
		t.Fatal("explicitly disconnected player's tank should be removed")
	}
	if _, ok := hub.weapons[weaponSlot{PlayerID: "solo", Weapon: 0}]; ok {
		t.Fatal("explicitly disconnected player's weapon state should be removed")
	}
	if _, ok := hub.bullets["solo-shot"]; ok {
		t.Fatal("explicitly disconnected player's bullets should be removed")
	}
	if _, ok := hub.clients["solo"]; ok {
		t.Fatal("explicitly disconnected player's client should be removed")
	}
	if len(hub.teams) != hub.playableTeamCount+1 || hub.teams[hub.playableTeamCount].Name != "Boris" {
		t.Fatalf("only the remaining player's personal team should remain, got %+v", hub.teams)
	}
	if hub.tanks["other"].Team != hub.playableTeamCount+1 || hub.tanks["red"].Team != 1 {
		t.Fatal("removing a personal team should preserve other team assignments")
	}
}
