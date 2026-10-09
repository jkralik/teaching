package game

import "time"

type Tile byte

const (
	TileEmpty Tile = '.'
	TileDirt  Tile = '#'
	TileRock  Tile = 'X'
)

type Direction string

const (
	Up        Direction = "up"
	Down      Direction = "down"
	Left      Direction = "left"
	Right     Direction = "right"
	UpLeft    Direction = "up-left"
	UpRight   Direction = "up-right"
	DownLeft  Direction = "down-left"
	DownRight Direction = "down-right"
)

type Tank struct {
	// Lekcia 11: Score je pripraveny ako nova vlastnost, ktoru klient dostane v JSON.
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	X         float64   `json:"x"`
	Y         float64   `json:"y"`
	Angle     float64   `json:"angle"`
	Direction Direction `json:"direction"`
	Alive     bool      `json:"alive"`
	Connected bool      `json:"connected"`
	Health    int       `json:"health"`
	MaxHealth int       `json:"maxHealth"`
	// Armor je poradie panciera v armorCatalog(), ostatne polia doplni snapshot pre prehliadac.
	Armor        int     `json:"armor"`
	ArmorName    string  `json:"armorName"`
	ArmorColor   string  `json:"armorColor"`
	ArmorDivisor int     `json:"armorDivisor"`
	ArmorSpeed   float64 `json:"armorSpeed"`
	Weapon       int     `json:"weapon"`
	Score        int     `json:"score"`
	// Team je cislo timu (0 = bez timu), TeamName a TeamColor doplni snapshot.
	Team      int    `json:"team"`
	TeamName  string `json:"teamName"`
	TeamColor string `json:"teamColor"`
	// Bonuses a Invulnerable doplni snapshot z aktivnych bonusov tanku.
	Bonuses      []ActiveBonus `json:"bonuses"`
	Invulnerable bool          `json:"invulnerable"`

	// Rozkopane policko zeme a kolko krokov uz tank kope (prehliadac to nepotrebuje).
	digging     bool
	digX, digY  int
	digProgress float64

	// bonuses: poradie bonusu v bonusCatalog() -> kedy efekt skonci.
	bonuses map[int]time.Time

	respawnAt      time.Time
	disconnectedAt time.Time
	// nextMoveAt: kedy smie tank spravit dalsi krok (limit rychlosti na serveri).
	nextMoveAt time.Time
}

type Bullet struct {
	ID        string    `json:"id"`
	OwnerID   string    `json:"ownerId"`
	X         float64   `json:"x"`
	Y         float64   `json:"y"`
	Angle     float64   `json:"angle"`
	Direction Direction `json:"direction"`
	Weapon    int       `json:"weapon"`
	Color     string    `json:"color"`
	// Team je tim strelca v chvili vystrelu, aby strela nezranila spoluhracov.
	Team int `json:"team"`
	// DamageFactor je nasobok poskodenia z bonusu v chvili vystrelu.
	DamageFactor float64 `json:"-"`
}

type Explosion struct {
	X         float64   `json:"x"`
	Y         float64   `json:"y"`
	Radius    float64   `json:"radius"`
	ExpiresAt time.Time `json:"-"`
}

type Snapshot struct {
	// Lekcia 13: Dopln sem napriklad PlayerCount a napln ho pri vytvarani snapshotu.
	MapName    string                  `json:"mapName"`
	Width      int                     `json:"width"`
	Height     int                     `json:"height"`
	Tiles      []string                `json:"tiles"`
	Tanks      map[string]Tank         `json:"tanks"`
	Bullets    map[string]Bullet       `json:"bullets"`
	Explosions []Explosion             `json:"explosions"`
	Weapons    map[string]WeaponStatus `json:"weapons"`
	Bonuses    []BonusItem             `json:"bonuses"`
	Teams      []TeamStatus            `json:"teams"`
}

type ClientMessage struct {
	Type      string    `json:"type"`
	Direction Direction `json:"direction,omitempty"`
	Angle     *float64  `json:"angle,omitempty"`
	Weapon    *int      `json:"weapon,omitempty"`
	Armor     *int      `json:"armor,omitempty"`
}

type ServerMessage struct {
	Type     string   `json:"type"`
	State    Snapshot `json:"state,omitempty"`
	PlayerID string   `json:"playerId,omitempty"`
	Error    string   `json:"error,omitempty"`
}
