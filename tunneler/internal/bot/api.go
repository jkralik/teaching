// Balik bot je AI tank. Je to samostatny program, ktory sa k serveru pripoji
// presne ako prehliadac: cez HTTP (/api/maps/{mapa}) a WebSocket (/ws).
// Server nevie, ze hra robot - bot nema ziadne vyhody ani specialne prikazy.
package bot

import "math"

// State je stav hry zo spravy "state", ktoru server posiela cez WebSocket.
// Obsahuje len polia, ktore bot pouziva; ostatne polia z JSON sa ignoruju.
// Lekcia 40: Potrebujes dalsiu informaciu zo servera? Pridaj pole s rovnakym JSON nazvom.
type State struct {
	MapName string            `json:"mapName"`
	Width   int               `json:"width"`
	Height  int               `json:"height"`
	Tiles   []string          `json:"tiles"`
	Tanks   map[string]Tank   `json:"tanks"`
	Bullets map[string]Bullet `json:"bullets"`
	Bonuses []Bonus           `json:"bonuses"`
	// Weapons: ID tanku -> stav jeho aktualnej zbrane.
	Weapons map[string]Weapon `json:"weapons"`
}

// Tank je jeden tank na mape (aj bot sam).
type Tank struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	X            float64 `json:"x"`
	Y            float64 `json:"y"`
	Angle        float64 `json:"angle"`
	Alive        bool    `json:"alive"`
	Connected    bool    `json:"connected"`
	Health       int     `json:"health"`
	MaxHealth    int     `json:"maxHealth"`
	Team         int     `json:"team"` // 0 = bez timu
	Score        int     `json:"score"`
	Weapon       int     `json:"weapon"`
	Invulnerable bool    `json:"invulnerable"`
}

// Bullet je letiaca strela.
type Bullet struct {
	ID      string  `json:"id"`
	OwnerID string  `json:"ownerId"`
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	Angle   float64 `json:"angle"`
	Team    int     `json:"team"`
}

// Bonus lezi na policku X, Y (stred policka je X+0.5, Y+0.5).
type Bonus struct {
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

// Weapon je stav zbrane jedneho tanku (naboje, prebijanie).
type Weapon struct {
	Name           string `json:"name"`
	Index          int    `json:"index"`
	Count          int    `json:"count"`
	ShotsRemaining int    `json:"shotsRemaining"`
	Reloading      bool   `json:"reloading"`
}

// Policka mapy, rovnake znaky ako v maps/*.txt.
const (
	TileEmpty byte = '.'
	TileDirt  byte = '#'
	TileRock  byte = 'X'
)

// TileAt vrati policko mapy. Mimo mapy je vzdy kamen.
func (state State) TileAt(x int, y int) byte {
	if y < 0 || y >= len(state.Tiles) || x < 0 || x >= len(state.Tiles[y]) {
		return TileRock
	}
	return state.Tiles[y][x]
}

// TileUnder vrati policko pod bodom (napr. pod tankom alebo strelou).
func (state State) TileUnder(x float64, y float64) byte {
	return state.TileAt(int(math.Floor(x)), int(math.Floor(y)))
}

// Action je to, co sa bot rozhodol urobit v tomto kroku.
// Lekcia 40: Sem mozes pridat dalsiu akciu, napriklad prepnutie zbrane.
type Action struct {
	Move  bool    // otocit sa na Angle a spravit krok
	Angle float64 // smer v radianoch: 0 = doprava, Pi/2 = hore
	Shoot bool    // vystrelit smerom Angle
}

// serverMessage je sprava zo servera (rovnaka, aku dostava prehliadac).
type serverMessage struct {
	Type     string `json:"type"`
	State    State  `json:"state"`
	PlayerID string `json:"playerId"`
	Error    string `json:"error"`
}

// clientMessage je prikaz pre server (rovnaky, aky posiela prehliadac).
type clientMessage struct {
	Type      string   `json:"type"`
	Direction string   `json:"direction,omitempty"`
	Angle     *float64 `json:"angle,omitempty"`
	Weapon    *int     `json:"weapon,omitempty"`
}
