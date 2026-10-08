package bot

import (
	"math"
	"math/rand"
	"strings"
	"time"
)

// ===== Nastavenia mozgu - toto mozes menit (lekcia 39) =====

const (
	// shootDistance: na kolko policok bot striela.
	shootDistance = 8.0
	// followDistance: copilot sa drzi takto blizko pri svojom hracovi.
	followDistance = 4.0
	// stuckTicks: kolko krokov sa tank nemusi pohnut, kym si mysli, ze sa zasekol.
	// Kopanie jedneho policka zeme trva asi 6 krokov, preto musi byt cislo vacsie.
	stuckTicks = 24
	// escapeTicks: kolko krokov ide zaseknuty tank nahodnym smerom.
	escapeTicks = 16
	// wanderTicks: ako dlho ide tank jednym smerom, ked nikoho nevidi.
	wanderTicks = 80
)

// Brain je mozog AI tanku. Decide sa vola kazdych 30 ms.
type Brain struct {
	// Follow je meno hraca, ktoremu bot pomaha (copilot). Prazdne = hra sam za seba.
	Follow string

	random      *rand.Rand
	lastX       float64
	lastY       float64
	stuck       int     // kolko krokov sa tank nepohol
	escape      int     // kolko krokov este utekat nahodnym smerom
	escapeAngle float64 // smer utekania
	wander      int     // kolko krokov este ist smerom wanderAngle
	wanderAngle float64
}

// NewBrain vytvori novy mozog. follow je meno hraca pre copilota alebo "".
func NewBrain(follow string) *Brain {
	return &Brain{
		Follow: follow,
		random: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

// Decide je hlavne rozhodnutie bota: kam ist a ci strielat.
// me je moj tank, state je cela mapa so vsetkymi tankami.
// Lekcia 40: Sem pridavaj nove pravidla, napr. ustup pri malom zdravi.
func (brain *Brain) Decide(me Tank, state State) Action {
	action := Action{}

	// 1. Strielanie: najblizsi nepriatel, ktoreho vidim, dostane gulku.
	enemy, found := nearestEnemy(me, state)
	if found && distance(me.X, me.Y, enemy.X, enemy.Y) <= shootDistance &&
		clearShot(state, me.X, me.Y, enemy.X, enemy.Y) {
		action.Angle = angleTo(me.X, me.Y, enemy.X, enemy.Y)
		action.Shoot = true
		return action // pri strielani stojime (otocenie spravi aj maly krok)
	}

	// 2. Pohyb: zaseknuty tank chvilu ide nahodnym smerom.
	if brain.isStuck(me) {
		action.Move = true
		action.Angle = brain.escapeAngle
		return action
	}

	// 3. Copilot ide za svojim hracom, ak je daleko.
	if owner, ok := brain.findOwner(me, state); ok &&
		distance(me.X, me.Y, owner.X, owner.Y) > followDistance {
		action.Move = true
		action.Angle = angleTo(me.X, me.Y, owner.X, owner.Y)
		return action
	}

	// 4. Inak ide za nepriatelom.
	if found {
		action.Move = true
		action.Angle = angleTo(me.X, me.Y, enemy.X, enemy.Y)
		return action
	}

	// 5. Nikoho nevidi - tula sa po mape.
	if brain.wander <= 0 {
		brain.wander = wanderTicks
		brain.wanderAngle = brain.randomAngle()
	}
	brain.wander--
	action.Move = true
	action.Angle = brain.wanderAngle
	return action
}

// isStuck zisti, ci sa tank dlho nepohol, a ak ano, vyberie nahodny smer uteku.
func (brain *Brain) isStuck(me Tank) bool {
	if brain.escape > 0 {
		brain.escape--
		return true
	}
	if distance(me.X, me.Y, brain.lastX, brain.lastY) < 0.01 {
		brain.stuck++
	} else {
		brain.stuck = 0
	}
	brain.lastX, brain.lastY = me.X, me.Y
	if brain.stuck < stuckTicks {
		return false
	}
	brain.stuck = 0
	brain.escape = escapeTicks
	brain.escapeAngle = brain.randomAngle()
	brain.wander = 0 // po uteku si vyberie novy smer tulania
	return true
}

// randomAngle vrati jeden z 8 smerov (vodorovne, zvisle, sikmo).
func (brain *Brain) randomAngle() float64 {
	return float64(brain.random.Intn(8)) * math.Pi / 4
}

// findOwner najde hraca, ktoremu copilot pomaha.
func (brain *Brain) findOwner(me Tank, state State) (Tank, bool) {
	if brain.Follow == "" {
		return Tank{}, false
	}
	for _, tank := range state.Tanks {
		if tank.ID != me.ID && tank.Alive && strings.EqualFold(tank.Name, brain.Follow) {
			return tank, true
		}
	}
	return Tank{}, false
}

// ===== Pomocne funkcie =====

// nearestEnemy najde najblizsi zivy tank, ktory nie je v mojom time.
func nearestEnemy(me Tank, state State) (Tank, bool) {
	best := Tank{}
	bestDistance := math.Inf(1)
	for _, tank := range state.Tanks {
		if tank.ID == me.ID || !tank.Alive || !tank.Connected || sameTeam(me, tank) {
			continue
		}
		if d := distance(me.X, me.Y, tank.X, tank.Y); d < bestDistance {
			best, bestDistance = tank, d
		}
	}
	return best, !math.IsInf(bestDistance, 1)
}

// sameTeam: tim 0 znamena "bez timu" - taky hrac nema spoluhracov.
func sameTeam(a Tank, b Tank) bool {
	return a.Team != 0 && a.Team == b.Team
}

// clearShot zisti, ci medzi dvoma bodmi nie je kamen.
// Zem strela prestrieli (vykope), kamen nie.
// Lekcia 40: Skus zakazat aj zem - bot potom striela len ked naozaj trafi.
func clearShot(state State, fromX, fromY, toX, toY float64) bool {
	steps := int(distance(fromX, fromY, toX, toY) / 0.2)
	for i := 1; i < steps; i++ {
		t := float64(i) / float64(steps)
		if state.TileUnder(fromX+(toX-fromX)*t, fromY+(toY-fromY)*t) == TileRock {
			return false
		}
	}
	return true
}

// angleTo vrati uhol z bodu (x, y) na bod (tx, ty).
// Pozor: na obrazovke ide y dole, preto je pred nim minus.
func angleTo(x, y, tx, ty float64) float64 {
	return math.Atan2(-(ty - y), tx-x)
}

// distance vrati vzdialenost dvoch bodov v polickach.
func distance(x, y, tx, ty float64) float64 {
	return math.Hypot(tx-x, ty-y)
}
