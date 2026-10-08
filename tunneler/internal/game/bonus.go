package game

import (
	"fmt"
	"math"
	"time"
)

// BonusSettings popisuje jeden bonus, ktory moze tank zdvihnut z mapy.
// Nevyplnene polia (0 alebo false) nic nemenia, takze bonus moze mat aj viac efektov naraz.
type BonusSettings struct {
	Name          string
	Symbol        string        // pismeno na policku v prehliadaci
	Color         string        // farba policka s bonusom
	Duration      time.Duration // ako dlho efekt trva; 0 = jednorazovy efekt
	InstantReload bool          // hned nabije vsetky zbrane
	SpeedFactor   float64       // nasobok rychlosti tanku, napr. 1.5 alebo 0.5
	DamageFactor  float64       // nasobok poskodenia striel, napr. 2
	Invulnerable  bool          // tank nedostava poskodenie
	MysteryOnly   bool          // objavi sa iba ako prekvapenie na policku ?
}

// Lekcia 37: Nastavenia toho, ako sa bonusy objavuju na mape.
const (
	bonusSpawnInterval = 4 * time.Second  // ako casto sa objavi novy bonus
	bonusLifetime      = 15 * time.Second // ako dlho bonus lezi na mape, potom zmizne
	bonusMaxOnMap      = 3                // kolko bonusov moze byt na mape naraz
	mysteryChance      = 0.3              // pravdepodobnost, ze sa objavi policko ?
	bonusPickupRadius  = 0.7              // ako blizko musi tank prist k stredu policka
)

// mysteryBonus je cislo, ktore v BonusItem znamena policko ? s nahodnym efektom.
const mysteryBonus = -1

// bonusCatalog vrati vsetky bonusy.
// Lekcia 37: Pridaj dalsie bonusy, napriklad Turbo, Silnu strelu alebo Nesmrtelnost.
// Lekcia 38: Pridaj dalsie negativne efekty s MysteryOnly: true.
func bonusCatalog() []BonusSettings {
	return []BonusSettings{
		{Name: "Okamzity reload", Symbol: "R", Color: "#5ec8f2", InstantReload: true},
		{Name: "Blato", Symbol: "B", Color: "#8a6b4a", Duration: 5 * time.Second, SpeedFactor: 0.5, MysteryOnly: true},
		{Name: "Turbo", Symbol: "T", Color: "#8a6b4a", Duration: 5 * time.Second, SpeedFactor: 1.5, MysteryOnly: true},
	}
}

// BonusItem je bonus lezaci na policku mapy.
type BonusItem struct {
	ID        string    `json:"id"`
	X         int       `json:"x"`
	Y         int       `json:"y"`
	Bonus     int       `json:"bonus"` // poradie v bonusCatalog() alebo mysteryBonus
	Name      string    `json:"name"`
	Symbol    string    `json:"symbol"`
	Color     string    `json:"color"`
	ExpiresAt time.Time `json:"-"`
}

// ActiveBonus je bonus, ktory tank prave ma; snapshot ho posle do prehliadaca.
type ActiveBonus struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Color       string `json:"color"`
	RemainingMs int64  `json:"remainingMs"`
}

func (hub *Hub) bonusSettingsFor(index int) (BonusSettings, bool) {
	if index < 0 || index >= len(hub.bonusCatalog) {
		return BonusSettings{}, false
	}
	return hub.bonusCatalog[index], true
}

// hasBonus povie, ci tank uz ma tento bonus. Ten isty bonus moze mat tank iba raz.
func (tank Tank) hasBonus(index int, now time.Time) bool {
	expiresAt, ok := tank.bonuses[index]
	return ok && now.Before(expiresAt)
}

// tankSpeedFactor vynasobi vsetky aktivne bonusy, ktore menia rychlost.
func (hub *Hub) tankSpeedFactor(tank Tank, now time.Time) float64 {
	factor := 1.0
	for index := range tank.bonuses {
		settings, ok := hub.bonusSettingsFor(index)
		if ok && tank.hasBonus(index, now) && settings.SpeedFactor > 0 {
			factor *= settings.SpeedFactor
		}
	}
	return factor
}

// tankDamageFactor vynasobi vsetky aktivne bonusy, ktore menia poskodenie.
func (hub *Hub) tankDamageFactor(tank Tank, now time.Time) float64 {
	factor := 1.0
	for index := range tank.bonuses {
		settings, ok := hub.bonusSettingsFor(index)
		if ok && tank.hasBonus(index, now) && settings.DamageFactor > 0 {
			factor *= settings.DamageFactor
		}
	}
	return factor
}

// tankInvulnerable povie, ci ma tank aktivnu nesmrtelnost.
func (hub *Hub) tankInvulnerable(tank Tank, now time.Time) bool {
	for index := range tank.bonuses {
		settings, ok := hub.bonusSettingsFor(index)
		if ok && tank.hasBonus(index, now) && settings.Invulnerable {
			return true
		}
	}
	return false
}

// boostedDamage vynasobi poskodenie bonusom a zaokruhli ho. Zasah zoberie aspon 1 zivot.
func boostedDamage(damage int, factor float64) int {
	if factor <= 0 || factor == 1 {
		return damage
	}
	return max(1, int(math.Round(float64(damage)*factor)))
}

// spawnBonusLocked polozi novy bonus na nahodne prazdne policko.
func (hub *Hub) spawnBonusLocked(now time.Time) bool {
	if len(hub.bonuses) >= bonusMaxOnMap {
		return false
	}
	normal := make([]int, 0, len(hub.bonusCatalog))
	for index, settings := range hub.bonusCatalog {
		if !settings.MysteryOnly {
			normal = append(normal, index)
		}
	}
	bonus := mysteryBonus
	if len(normal) > 0 && hub.random.Float64() >= mysteryChance {
		bonus = normal[hub.random.IntN(len(normal))]
	}
	if len(hub.bonusCatalog) == 0 {
		return false
	}

	for range 30 {
		x := hub.random.IntN(max(1, hub.gameMap.Width))
		y := hub.random.IntN(max(1, hub.gameMap.Height))
		if !hub.freeForBonusLocked(x, y) {
			continue
		}
		item := BonusItem{
			ID:        fmt.Sprintf("bonus-%d-%d-%d", x, y, now.UnixNano()),
			X:         x,
			Y:         y,
			Bonus:     bonus,
			Name:      "Prekvapenie",
			Symbol:    "?",
			Color:     "#b07cff",
			ExpiresAt: now.Add(bonusLifetime),
		}
		if settings, ok := hub.bonusSettingsFor(bonus); ok {
			item.Name, item.Symbol, item.Color = settings.Name, settings.Symbol, settings.Color
		}
		hub.bonuses[item.ID] = item
		return true
	}
	return false
}

func (hub *Hub) freeForBonusLocked(x int, y int) bool {
	if hub.gameMap.TileAt(x, y) != TileEmpty {
		return false
	}
	for _, item := range hub.bonuses {
		if item.X == x && item.Y == y {
			return false
		}
	}
	for _, tank := range hub.tanks {
		if int(math.Floor(tank.X)) == x && int(math.Floor(tank.Y)) == y {
			return false
		}
	}
	return true
}

// pickUpBonusesLocked da tanku bonusy, na ktore prave nabehol.
func (hub *Hub) pickUpBonusesLocked(tank *Tank, now time.Time) {
	for id, item := range hub.bonuses {
		if math.Hypot(tank.X-(float64(item.X)+0.5), tank.Y-(float64(item.Y)+0.5)) > bonusPickupRadius {
			continue
		}
		index := item.Bonus
		if index == mysteryBonus {
			index = hub.randomMysteryBonusLocked(*tank, now)
		}
		if index < 0 || tank.hasBonus(index, now) {
			// Tank uz tento bonus ma, nechame ho na mape pre ostatnych.
			continue
		}
		hub.applyBonusLocked(tank, index, now)
		delete(hub.bonuses, id)
	}
}

// randomMysteryBonusLocked vyberie nahodny bonus (aj negativny), ktory tank este nema.
func (hub *Hub) randomMysteryBonusLocked(tank Tank, now time.Time) int {
	choices := make([]int, 0, len(hub.bonusCatalog))
	for index := range hub.bonusCatalog {
		if !tank.hasBonus(index, now) {
			choices = append(choices, index)
		}
	}
	if len(choices) == 0 {
		return -1
	}
	return choices[hub.random.IntN(len(choices))]
}

// applyBonusLocked spusti efekt bonusu.
// Lekcia 37: Ak pridas do BonusSettings novy efekt, naprogramuj ho tu alebo tam, kde sa pouziva.
func (hub *Hub) applyBonusLocked(tank *Tank, index int, now time.Time) {
	settings, ok := hub.bonusSettingsFor(index)
	if !ok {
		return
	}
	if settings.InstantReload {
		for slot := range hub.weapons {
			if slot.PlayerID == tank.ID {
				delete(hub.weapons, slot)
			}
		}
	}
	if settings.Duration > 0 {
		if tank.bonuses == nil {
			tank.bonuses = make(map[int]time.Time)
		}
		tank.bonuses[index] = now.Add(settings.Duration)
	}
}

// activeBonusesLocked pripravi zoznam aktivnych bonusov tanku pre snapshot.
func (hub *Hub) activeBonusesLocked(tank Tank, now time.Time) []ActiveBonus {
	active := make([]ActiveBonus, 0, len(tank.bonuses))
	for index := range hub.bonusCatalog {
		if !tank.hasBonus(index, now) {
			continue
		}
		settings := hub.bonusCatalog[index]
		active = append(active, ActiveBonus{
			Name:        settings.Name,
			Symbol:      settings.Symbol,
			Color:       settings.Color,
			RemainingMs: max(1, tank.bonuses[index].Sub(now).Milliseconds()),
		})
	}
	return active
}
