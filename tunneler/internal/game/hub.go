package game

import (
	"fmt"
	"sync"
	"time"
)

type Command struct {
	PlayerID   string
	PlayerName string
	Message    ClientMessage
}

type Hub struct {
	mu      sync.Mutex
	gameMap *Map
	tanks   map[string]Tank
	bullets map[string]Bullet
	clients map[string]chan ServerMessage
	stop    chan struct{}
}

func NewHub(gameMap *Map) *Hub {
	hub := &Hub{
		gameMap: gameMap.Clone(),
		tanks:   make(map[string]Tank),
		bullets: make(map[string]Bullet),
		clients: make(map[string]chan ServerMessage),
		stop:    make(chan struct{}),
	}
	go hub.loop()
	return hub
}

func (hub *Hub) Join(playerID string, playerName string) chan ServerMessage {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	out := make(chan ServerMessage, 8)
	hub.clients[playerID] = out
	hub.tanks[playerID] = Tank{
		ID:        playerID,
		Name:      playerName,
		X:         2 + len(hub.tanks)%max(1, hub.gameMap.Width-4),
		Y:         2,
		Direction: Right,
		Alive:     true,
	}
	out <- ServerMessage{Type: "state", State: hub.snapshotLocked()}
	return out
}

func (hub *Hub) Leave(playerID string) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	if out, ok := hub.clients[playerID]; ok {
		close(out)
	}
	delete(hub.clients, playerID)
	delete(hub.tanks, playerID)
	hub.broadcastLocked(ServerMessage{Type: "state", State: hub.snapshotLocked()})
}

func (hub *Hub) Handle(command Command) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	tank, ok := hub.tanks[command.PlayerID]
	if !ok || !tank.Alive {
		return
	}

	switch command.Message.Type {
	case "move":
		tank.Direction = command.Message.Direction
		hub.moveTank(&tank)
		hub.tanks[command.PlayerID] = tank
	case "shoot":
		bullet := Bullet{
			ID:        fmt.Sprintf("%s-%d", tank.ID, time.Now().UnixNano()),
			OwnerID:   tank.ID,
			X:         tank.X,
			Y:         tank.Y,
			Direction: tank.Direction,
		}
		hub.moveBullet(&bullet)
		hub.bullets[bullet.ID] = bullet
	}

	hub.broadcastLocked(ServerMessage{Type: "state", State: hub.snapshotLocked()})
}

func (hub *Hub) loop() {
	ticker := time.NewTicker(120 * time.Millisecond)
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
	hub.mu.Lock()
	defer hub.mu.Unlock()

	changed := false
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

func (hub *Hub) moveTank(tank *Tank) {
	nextX, nextY := nextPosition(tank.X, tank.Y, tank.Direction)
	tile := hub.gameMap.TileAt(nextX, nextY)
	if tile == TileRock {
		return
	}
	if tile == TileDirt {
		hub.gameMap.SetTile(nextX, nextY, TileEmpty)
		tank.Score++
		return
	}
	tank.X = nextX
	tank.Y = nextY
}

func (hub *Hub) moveBullet(bullet *Bullet) bool {
	nextX, nextY := nextPosition(bullet.X, bullet.Y, bullet.Direction)
	tile := hub.gameMap.TileAt(nextX, nextY)
	if tile == TileRock {
		return false
	}
	if tile == TileDirt {
		hub.gameMap.SetTile(nextX, nextY, TileEmpty)
		return false
	}
	for id, tank := range hub.tanks {
		if id == bullet.OwnerID || !tank.Alive {
			continue
		}
		if tank.X == nextX && tank.Y == nextY {
			tank.Alive = false
			hub.tanks[id] = tank
			if owner, ok := hub.tanks[bullet.OwnerID]; ok {
				owner.Score += 5
				hub.tanks[bullet.OwnerID] = owner
			}
			return false
		}
	}
	bullet.X = nextX
	bullet.Y = nextY
	return true
}

func nextPosition(x int, y int, direction Direction) (int, int) {
	switch direction {
	case Up:
		return x, y - 1
	case Down:
		return x, y + 1
	case Left:
		return x - 1, y
	case Right:
		return x + 1, y
	default:
		return x, y
	}
}

func (hub *Hub) broadcastLocked(message ServerMessage) {
	for _, out := range hub.clients {
		select {
		case out <- message:
		default:
		}
	}
}

func (hub *Hub) snapshotLocked() Snapshot {
	tanks := make(map[string]Tank, len(hub.tanks))
	for id, tank := range hub.tanks {
		tanks[id] = tank
	}
	bullets := make(map[string]Bullet, len(hub.bullets))
	for id, bullet := range hub.bullets {
		bullets[id] = bullet
	}
	return Snapshot{
		MapName: hub.gameMap.Name,
		Width:   hub.gameMap.Width,
		Height:  hub.gameMap.Height,
		Tiles:   hub.gameMap.Rows(),
		Tanks:   tanks,
		Bullets: bullets,
	}
}
