package game

import "testing"

func TestGenerateMapHasSolidBorder(t *testing.T) {
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
	tank := Tank{ID: "p1", X: 1, Y: 1, Direction: Right, Alive: true}
	hub.moveTank(&tank)

	if tank.X != 1 || tank.Y != 1 {
		t.Fatalf("tank should dig first and stay in place, got %d,%d", tank.X, tank.Y)
	}
	if gameMap.TileAt(2, 1) != TileDirt {
		t.Fatalf("original map should not be mutated")
	}
	if hub.gameMap.TileAt(2, 1) != TileEmpty {
		t.Fatalf("hub map tile should be dug out")
	}
}
