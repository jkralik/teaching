package game

import (
	"fmt"
	"math"
	"math/rand/v2"
	"strings"
	"sync"
	"time"
)

// Rychlost tanku = tankMoveStep / tankMoveInterval (0.1 policka za 30 ms = 3.3 policka za sekundu).
// Mensi krok castejsie = plynulejsi pohyb pri rovnakej rychlosti.
const tankMoveStep = 0.1

// Ochrana proti podvadzaniu: server dovoli tanku najviac jeden krok za tankMoveInterval,
// aj keby klient posielal spravy "move" castejsie. Prehliadac posiela krok kazdych 30 ms.
// tankMoveBurst je mala rezerva, aby pohyb nesekal, ked spravy pridu nerovnomerne cez siet.
const tankMoveInterval = 30 * time.Millisecond
const tankMoveBurst = 2 * tankMoveInterval

// Ako casto server posunie strely a posle novy stav hracom (30 ms = asi 33 obrazkov za sekundu).
const gameTickInterval = 30 * time.Millisecond

// Lekcia 16: Tvrdost zeme. Kazdy pohyb vykope tolko, aky ma tank krok (tankMoveStep).
// 1.0 = kopanie policka trva rovnako dlho ako prejazd cez prazdne policko, 2.0 = dvakrat dlhsie.
const dirtHardness = 1.2

const explosionDuration = 350 * time.Millisecond
const tankRespawnDelay = 3 * time.Second
const tankRetentionDuration = 5 * time.Minute
const tankCollisionRadius = 0.4

type Command struct {
	PlayerID   string
	PlayerName string
	Message    ClientMessage
}

// weaponSlot oddeli zasobnik kazdej zbrane kazdeho hraca.
type weaponSlot struct {
	PlayerID string
	Weapon   int
}

type Hub struct {
	mu         sync.Mutex
	gameMap    *Map
	tanks      map[string]Tank
	bullets    map[string]Bullet
	arsenal    []WeaponSettings
	armors     []ArmorSettings
	teams      []TeamSettings
	weapons    map[weaponSlot]weaponState
	explosions map[string]Explosion
	// Lekcia 37: Bonusy z katalogu a bonusy, ktore prave lezia na mape.
	bonusCatalog []BonusSettings
	bonuses      map[string]BonusItem
	nextBonusAt  time.Time
	random       *rand.Rand
	clients      map[string]chan ServerMessage
	stop         chan struct{}
}

func NewHub(gameMap *Map) *Hub {
	hub := &Hub{
		gameMap:      gameMap.Clone(),
		tanks:        make(map[string]Tank),
		bullets:      make(map[string]Bullet),
		arsenal:      weaponCatalog(),
		armors:       armorCatalog(),
		teams:        teamCatalog(),
		weapons:      make(map[weaponSlot]weaponState),
		explosions:   make(map[string]Explosion),
		bonusCatalog: bonusCatalog(),
		bonuses:      make(map[string]BonusItem),
		nextBonusAt:  time.Now().Add(bonusSpawnInterval),
		random:       rand.New(rand.NewPCG(uint64(time.Now().UnixNano()), 37)),
		clients:      make(map[string]chan ServerMessage),
		stop:         make(chan struct{}),
	}
	go hub.loop()
	return hub
}

// Join pripoji hraca a tim mu vyberie server automaticky.
func (hub *Hub) Join(playerID string, playerName string) chan ServerMessage {
	return hub.JoinTeam(playerID, playerName, TeamAuto)
}

// JoinTeam pripoji hraca do timu team (TeamAuto, TeamNone alebo 1, 2, ...).
// Hrac, ktory sa vracia k svojmu tanku, zostava vo svojom povodnom time.
func (hub *Hub) JoinTeam(playerID string, playerName string, team int) chan ServerMessage {
	// Lekcia 18: Join prida hraca do zdielanej mapy tanks a posle prvy snapshot.
	hub.mu.Lock()
	defer hub.mu.Unlock()

	out := make(chan ServerMessage, 8)
	hub.clients[playerID] = out
	resumed := false
	for previousID, tank := range hub.tanks {
		if tank.Connected || tank.disconnectedAt.IsZero() ||
			!strings.EqualFold(strings.TrimSpace(tank.Name), strings.TrimSpace(playerName)) {
			continue
		}
		if time.Since(tank.disconnectedAt) >= tankRetentionDuration {
			hub.removeTankLocked(previousID)
			continue
		}
		delete(hub.tanks, previousID)
		tank.ID = playerID
		tank.Connected = true
		tank.disconnectedAt = time.Time{}
		hub.tanks[playerID] = tank
		hub.reassignPlayerStateLocked(previousID, playerID)
		resumed = true
		break
	}
	if !resumed {
		tank := Tank{
			ID:        playerID,
			Name:      playerName,
			Angle:     0,
			Direction: Right,
			Alive:     true,
			Connected: true,
			Health:    tankMaxHealth,
			MaxHealth: tankMaxHealth,
			Team:      hub.chooseTeamLocked(team),
		}
		if tank.Team == TeamNone {
			tank.soloColor = hub.randomSoloColorLocked()
		}
		spawnX := float64(2+len(hub.tanks)%max(1, hub.gameMap.Width-4)) + 0.5
		spawnY := 2.5
		if hub.gameMap.TileAt(int(math.Floor(spawnX)), int(math.Floor(spawnY))) == TileEmpty &&
			!hub.tankPositionBlockedLocked(playerID, spawnX, spawnY) {
			tank.X, tank.Y = spawnX, spawnY
		} else if x, y, ok := hub.findSpawnPositionLocked(playerID); ok {
			tank.X, tank.Y = x, y
		} else {
			tank.Alive = false
			tank.respawnAt = time.Now()
		}
		hub.tanks[playerID] = tank
	}
	out <- ServerMessage{Type: "state", State: hub.snapshotLocked(), PlayerID: playerID}
	hub.broadcastExceptLocked(playerID, ServerMessage{Type: "state", State: hub.snapshotLocked()})
	return out
}

func (hub *Hub) Leave(playerID string) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	if out, ok := hub.clients[playerID]; ok {
		close(out)
		delete(hub.clients, playerID)
		if tank, ok := hub.tanks[playerID]; ok {
			tank.Connected = false
			tank.disconnectedAt = time.Now()
			hub.tanks[playerID] = tank
		}
		hub.broadcastLocked(ServerMessage{Type: "state", State: hub.snapshotLocked()})
	}
}

func (hub *Hub) Handle(command Command) {
	// Lekcia 08: Najdi hraca podla ID v mape tanks. Skus pomocne vyhladavanie oddelit do funkcie.
	// Lekcia 17: ClientMessage prichadza z klavesnice; tu sa meni na hernu akciu.
	// Lekcia 20: Mutex chrani kriticku sekciu, v ktorej citame a menime zdielany stav.
	hub.mu.Lock()
	defer hub.mu.Unlock()

	tank, ok := hub.tanks[command.PlayerID]
	if !ok || !tank.Alive || !tank.Connected {
		return
	}

	switch command.Message.Type {
	case "move":
		tank.Direction = command.Message.Direction
		if command.Message.Angle != nil {
			tank.Angle = *command.Message.Angle
		} else {
			tank.Angle = angleForDirection(tank.Direction)
		}
		// Otocit sa moze tank vzdy, posunut sa len ak neprekrocil limit krokov.
		if takeMoveTurn(&tank, time.Now()) {
			hub.moveTank(&tank)
		}
		hub.tanks[command.PlayerID] = tank
	case "weapon":
		// Lekcia 33: Bez cisla prepne na dalsiu zbran, s cislom vyberie konkretnu.
		count := len(hub.arsenal)
		if count == 0 {
			return
		}
		if command.Message.Weapon == nil {
			tank.Weapon = (tank.Weapon + 1) % count
		} else if *command.Message.Weapon >= 0 && *command.Message.Weapon < count {
			tank.Weapon = *command.Message.Weapon
		} else {
			return
		}
		hub.tanks[command.PlayerID] = tank
	case "armor":
		// Lekcia 36: Bez cisla prepne na dalsi pancier, s cislom vyberie konkretny.
		count := len(hub.armors)
		if count == 0 {
			return
		}
		if command.Message.Armor == nil {
			tank.Armor = (tank.Armor + 1) % count
		} else if *command.Message.Armor >= 0 && *command.Message.Armor < count {
			tank.Armor = *command.Message.Armor
		} else {
			return
		}
		hub.tanks[command.PlayerID] = tank
	case "shoot":
		// Lekcia 28: Sem mozes pridat nove herne pravidlo, napriklad specialny typ strely.
		now := time.Now()
		slot := weaponSlot{PlayerID: command.PlayerID, Weapon: tank.Weapon}
		weapon := hub.weapons[slot]
		settings := hub.weaponSettingsFor(tank.Weapon)
		if !weapon.tryShoot(now, settings) {
			return
		}
		hub.weapons[slot] = weapon
		bullet := Bullet{
			ID:        fmt.Sprintf("%s-%d", tank.ID, time.Now().UnixNano()),
			OwnerID:   tank.ID,
			X:         tank.X,
			Y:         tank.Y,
			Angle:     tank.Angle,
			Direction: tank.Direction,
			Weapon:    tank.Weapon,
			Color:     settings.Color,
			Team:      tank.Team,
			// Lekcia 37: Strela si zapamata silu bonusu z chvile vystrelu.
			DamageFactor: hub.tankDamageFactor(tank, now),
		}
		if hub.moveBullet(&bullet) {
			hub.bullets[bullet.ID] = bullet
		}
	}

	hub.broadcastLocked(ServerMessage{Type: "state", State: hub.snapshotLocked()})
}

func (hub *Hub) loop() {
	// Lekcia 19: Herna slucka bezi vo vlastnej gorutine a caka na ticker alebo stop kanal.
	ticker := time.NewTicker(gameTickInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			hub.tick()
		case <-hub.stop:
			return
		}
	}
}

func (hub *Hub) tick() {
	// Lekcia 21: Kazdy tick posunie strely a pripadne odvysiela novy stav.
	hub.mu.Lock()
	defer hub.mu.Unlock()

	changed := false
	now := time.Now()
	for id, tank := range hub.tanks {
		if tank.Connected || tank.disconnectedAt.IsZero() || now.Sub(tank.disconnectedAt) < tankRetentionDuration {
			continue
		}
		hub.removeTankLocked(id)
		changed = true
	}
	for id, explosion := range hub.explosions {
		if !now.Before(explosion.ExpiresAt) {
			delete(hub.explosions, id)
			changed = true
		}
	}
	for id, tank := range hub.tanks {
		if tank.Alive || tank.respawnAt.IsZero() || now.Before(tank.respawnAt) {
			continue
		}
		x, y, ok := hub.findSpawnPositionLocked(id)
		if !ok {
			continue
		}
		tank.X, tank.Y = x, y
		tank.Alive = true
		tank.Health = tankMaxHealth
		tank.Angle = 0
		tank.Direction = Right
		tank.digging = false
		tank.digProgress = 0
		tank.bonuses = nil
		tank.respawnAt = time.Time{}
		hub.tanks[id] = tank
		changed = true
	}
	for id, weapon := range hub.weapons {
		if weapon.reloadingUntil.IsZero() {
			continue
		}
		if now.Before(weapon.reloadingUntil) {
			changed = true
			continue
		}
		weapon.shotsFired = 0
		weapon.reloadingUntil = time.Time{}
		hub.weapons[id] = weapon
		changed = true
	}
	if hub.updateBonusesLocked(now) {
		changed = true
	}
	for id, bullet := range hub.bullets {
		if !hub.moveBullet(&bullet) {
			delete(hub.bullets, id)
			changed = true
			continue
		}
		hub.bullets[id] = bullet
		changed = true
	}

	if changed {
		hub.broadcastLocked(ServerMessage{Type: "state", State: hub.snapshotLocked()})
	}
}

// takeMoveTurn povie, ci tank smie teraz spravit krok. Kazdy krok posunie
// nextMoveAt o tankMoveInterval, takze dlhodobo ide tank najviac
// 1 krok za tankMoveInterval, nech klient posiela spravy akokolvek rychlo.
func takeMoveTurn(tank *Tank, now time.Time) bool {
	if tank.nextMoveAt.After(now.Add(tankMoveBurst)) {
		return false
	}
	if tank.nextMoveAt.Before(now) {
		tank.nextMoveAt = now
	}
	tank.nextMoveAt = tank.nextMoveAt.Add(tankMoveInterval)
	return true
}

func (hub *Hub) moveTank(tank *Tank) {
	// Lekcia 14: nextPosition vypocita ciel podla smeru. Pridaj alebo otestuj dalsi smer.
	// Lekcia 16: Pri TileDirt sa mapa zmeni na TileEmpty a tank ziska bod.
	// Lekcia 36: Tazsi pancier spomali tank.
	// Lekcia 37: Bonusy mozu tank zrychlit alebo spomalit.
	now := time.Now()
	step := armoredSpeed(tankMoveStep, hub.armorSettingsFor(tank.Armor)) * hub.tankSpeedFactor(*tank, now)
	nextX, nextY := nextPosition(tank.X, tank.Y, tank.Angle, step)
	// Lekcia 15: Toto je miesto pre pravidlo kolizie pred zmenou pozicie tanku.
	if hub.tankPositionBlockedLocked(tank.ID, nextX, nextY) {
		return
	}
	tile := hub.gameMap.TileAt(int(math.Floor(nextX)), int(math.Floor(nextY)))
	if tile == TileRock {
		return
	}
	if tile == TileDirt {
		// Lekcia 16: Tank kope to iste policko niekolko krokov, az potom zem zmizne.
		digX, digY := int(math.Floor(nextX)), int(math.Floor(nextY))
		if !tank.digging || tank.digX != digX || tank.digY != digY {
			tank.digging, tank.digX, tank.digY, tank.digProgress = true, digX, digY, 0
		}
		tank.digProgress += step
		if tank.digProgress >= dirtHardness-1e-9 {
			hub.gameMap.SetTile(digX, digY, TileEmpty)
			tank.Score++
			tank.digging = false
		}
		return
	}
	tank.digging = false
	tank.X = nextX
	tank.Y = nextY
	hub.pickUpBonusesLocked(tank, now)
}

func (hub *Hub) moveBullet(bullet *Bullet) bool {
	// Lekcia 22: Pri zhode suradnic moze strela vyradit cudzi tank a pridat body vlastnikovi.
	// Lekcia 34: Tu mozes pouzit novu vlastnost zbrane, napriklad ukoncit strelu po dostrele.
	settings := hub.weaponSettingsFor(bullet.Weapon)
	bulletStep := settings.BulletSpeed * gameTickInterval.Seconds()
	nextX, nextY := nextPosition(bullet.X, bullet.Y, bullet.Angle, bulletStep)
	tile := hub.gameMap.TileAt(int(math.Floor(nextX)), int(math.Floor(nextY)))
	if tile == TileRock {
		hub.explode(bullet, nextX, nextY, settings)
		return false
	}
	if tile == TileDirt {
		hub.gameMap.SetTile(int(math.Floor(nextX)), int(math.Floor(nextY)), TileEmpty)
		hub.explode(bullet, nextX, nextY, settings)
		return false
	}
	hit := false
	for id, tank := range hub.tanks {
		// Spoluhraca strela preleti, ak je friendlyFire vypnuty.
		if id == bullet.OwnerID || !tank.Alive || !canHurt(bullet.Team, tank.Team) {
			continue
		}
		if math.Hypot(tank.X-nextX, tank.Y-nextY) <= settings.ExplosionRadius {
			hit = true
			break
		}
	}
	if hit {
		hub.explode(bullet, nextX, nextY, settings)
		return false
	}
	bullet.X = nextX
	bullet.Y = nextY
	return true
}

func (hub *Hub) explode(bullet *Bullet, x float64, y float64, settings WeaponSettings) {
	hub.explosions[bullet.ID] = Explosion{
		X:         x,
		Y:         y,
		Radius:    settings.ExplosionRadius,
		ExpiresAt: time.Now().Add(explosionDuration),
	}
	for id, tank := range hub.tanks {
		if id == bullet.OwnerID || !tank.Alive || math.Hypot(tank.X-x, tank.Y-y) > settings.ExplosionRadius {
			continue
		}
		// Timy: vybuch nezrani spoluhraca (a strelec zan nedostane body).
		if !canHurt(bullet.Team, tank.Team) {
			continue
		}
		// Lekcia 37: Nesmrtelny tank zasah ignoruje.
		if hub.tankInvulnerable(tank, time.Now()) {
			continue
		}
		damage := boostedDamage(settings.Damage, bullet.DamageFactor)
		tank.Health = max(0, tank.Health-armoredDamage(damage, hub.armorSettingsFor(tank.Armor).Divisor))
		destroyed := tank.Health == 0
		if destroyed {
			tank.Alive = false
			tank.respawnAt = time.Now().Add(tankRespawnDelay)
		}
		hub.tanks[id] = tank
		if owner, ok := hub.tanks[bullet.OwnerID]; ok {
			owner.Score++
			if destroyed {
				owner.Score += 5
			}
			hub.tanks[bullet.OwnerID] = owner
		}
	}
}

func (hub *Hub) findSpawnPositionLocked(playerID string) (float64, float64, bool) {
	type spawnPosition struct {
		x, y int
	}
	available := make([]spawnPosition, 0)
	for y := 1; y < hub.gameMap.Height-1; y++ {
		for x := 1; x < hub.gameMap.Width-1; x++ {
			if hub.gameMap.TileAt(x, y) != TileEmpty {
				continue
			}
			spawnX, spawnY := float64(x)+0.5, float64(y)+0.5
			if hub.tankPositionBlockedLocked(playerID, spawnX, spawnY) {
				continue
			}
			occupied := false
			for _, bonus := range hub.bonuses {
				if bonus.X == x && bonus.Y == y {
					occupied = true
					break
				}
			}
			if !occupied {
				available = append(available, spawnPosition{x: x, y: y})
			}
		}
	}
	if len(available) == 0 {
		return 0, 0, false
	}
	position := available[hub.random.IntN(len(available))]
	return float64(position.x) + 0.5, float64(position.y) + 0.5, true
}

func (hub *Hub) tankPositionBlockedLocked(playerID string, x float64, y float64) bool {
	minimumDistance := 2*tankCollisionRadius - 1e-9
	for id, other := range hub.tanks {
		if id != playerID && other.Alive && math.Hypot(other.X-x, other.Y-y) < minimumDistance {
			return true
		}
	}
	return false
}

func (hub *Hub) reassignPlayerStateLocked(previousID string, playerID string) {
	for slot, weapon := range hub.weapons {
		if slot.PlayerID != previousID {
			continue
		}
		delete(hub.weapons, slot)
		slot.PlayerID = playerID
		hub.weapons[slot] = weapon
	}
	for id, bullet := range hub.bullets {
		if bullet.OwnerID == previousID {
			bullet.OwnerID = playerID
			hub.bullets[id] = bullet
		}
	}
}

func (hub *Hub) removeTankLocked(playerID string) {
	delete(hub.tanks, playerID)
	for slot := range hub.weapons {
		if slot.PlayerID == playerID {
			delete(hub.weapons, slot)
		}
	}
}

// weaponSettingsFor vrati nastavenia zbrane podla poradia v zozname.
func (hub *Hub) weaponSettingsFor(index int) WeaponSettings {
	if index < 0 || index >= len(hub.arsenal) {
		return weaponSettings()
	}
	return hub.arsenal[index]
}

func nextPosition(x float64, y float64, angle float64, distance float64) (float64, float64) {
	return x + math.Cos(angle)*distance, y - math.Sin(angle)*distance
}

func angleForDirection(direction Direction) float64 {
	switch direction {
	case Up:
		return math.Pi / 2
	case Down:
		return -math.Pi / 2
	case Left:
		return math.Pi
	case Right:
		return 0
	case UpLeft:
		return 3 * math.Pi / 4
	case UpRight:
		return math.Pi / 4
	case DownLeft:
		return -3 * math.Pi / 4
	case DownRight:
		return -math.Pi / 4
	default:
		return 0
	}
}

func (hub *Hub) broadcastLocked(message ServerMessage) {
	// Lekcia 18: Broadcast odosle snapshot vsetkym klientom rovnakej mapy.
	for _, out := range hub.clients {
		select {
		case out <- message:
		default:
		}
	}
}

func (hub *Hub) broadcastExceptLocked(playerID string, message ServerMessage) {
	for id, out := range hub.clients {
		if id == playerID {
			continue
		}
		select {
		case out <- message:
		default:
		}
	}
}

func (hub *Hub) snapshotLocked() Snapshot {
	// Lekcia 13: Snapshot kopiruje stav pod zamkom, aby ho mohol bezpecne precitat klient.
	now := time.Now()
	tanks := make(map[string]Tank, len(hub.tanks))
	for id, tank := range hub.tanks {
		tank.Bonuses = hub.activeBonusesLocked(tank, now)
		tank.Invulnerable = hub.tankInvulnerable(tank, now)
		tank.bonuses = nil
		armor := hub.armorSettingsFor(tank.Armor)
		tank.ArmorName = armor.Name
		tank.ArmorColor = armor.Color
		tank.ArmorDivisor = armor.Divisor
		tank.ArmorSpeed = armor.SpeedFactor
		if team, ok := hub.teamSettingsFor(tank.Team); ok {
			tank.TeamName = team.Name
			tank.TeamColor = team.Color
		} else {
			tank.TeamColor = tank.soloColor
		}
		tanks[id] = tank
	}
	bullets := make(map[string]Bullet, len(hub.bullets))
	for id, bullet := range hub.bullets {
		bullets[id] = bullet
	}
	explosions := make([]Explosion, 0, len(hub.explosions))
	for _, explosion := range hub.explosions {
		explosions = append(explosions, explosion)
	}
	bonuses := make([]BonusItem, 0, len(hub.bonuses))
	for _, item := range hub.bonuses {
		bonuses = append(bonuses, item)
	}
	weapons := make(map[string]WeaponStatus, len(hub.tanks))
	for id, tank := range hub.tanks {
		settings := hub.weaponSettingsFor(tank.Weapon)
		status := hub.weapons[weaponSlot{PlayerID: id, Weapon: tank.Weapon}].status(now, settings)
		status.Index = tank.Weapon
		status.Count = len(hub.arsenal)
		weapons[id] = status
	}
	return Snapshot{
		MapName:    hub.gameMap.Name,
		Width:      hub.gameMap.Width,
		Height:     hub.gameMap.Height,
		Tiles:      hub.gameMap.Rows(),
		Tanks:      tanks,
		Bullets:    bullets,
		Explosions: explosions,
		Weapons:    weapons,
		Bonuses:    bonuses,
		Teams:      hub.teamStatusesLocked(),
	}
}

// updateBonusesLocked zmaze stare bonusy, ukonci efekty a polozi novy bonus na mapu.
func (hub *Hub) updateBonusesLocked(now time.Time) bool {
	changed := false
	for id, item := range hub.bonuses {
		if !now.Before(item.ExpiresAt) {
			delete(hub.bonuses, id)
			changed = true
		}
	}
	for id, tank := range hub.tanks {
		for index, expiresAt := range tank.bonuses {
			if !now.Before(expiresAt) {
				delete(tank.bonuses, index)
			}
			// Aj bezici efekt meni odpocitavanie v prehliadaci.
			changed = true
		}
		hub.tanks[id] = tank
	}
	if !now.Before(hub.nextBonusAt) {
		hub.nextBonusAt = now.Add(bonusSpawnInterval)
		if hub.spawnBonusLocked(now) {
			changed = true
		}
	}
	return changed
}
