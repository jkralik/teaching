package game

type Tile byte

const (
	TileEmpty Tile = '.'
	TileDirt  Tile = '#'
	TileRock  Tile = 'X'
)

type Direction string

const (
	Up    Direction = "up"
	Down  Direction = "down"
	Left  Direction = "left"
	Right Direction = "right"
)

type Tank struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	X         int       `json:"x"`
	Y         int       `json:"y"`
	Direction Direction `json:"direction"`
	Alive     bool      `json:"alive"`
	Score     int       `json:"score"`
}

type Bullet struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"ownerId"`
	X         int       `json:"x"`
	Y         int       `json:"y"`
	Direction Direction `json:"direction"`
}

type Snapshot struct {
	MapName string            `json:"mapName"`
	Width   int               `json:"width"`
	Height  int               `json:"height"`
	Tiles   []string          `json:"tiles"`
	Tanks   map[string]Tank   `json:"tanks"`
	Bullets map[string]Bullet `json:"bullets"`
}

type ClientMessage struct {
	Type      string    `json:"type"`
	Direction Direction `json:"direction,omitempty"`
}

type ServerMessage struct {
	Type  string   `json:"type"`
	State Snapshot `json:"state,omitempty"`
	Error string   `json:"error,omitempty"`
}
