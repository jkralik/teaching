package game

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

type Server struct {
	mu           sync.Mutex
	maps         map[string]*Map
	hubs         map[string]*Hub
	upgrader     websocket.Upgrader
	nextPlayerID atomic.Uint64
}

func NewServer(mapDir string) *Server {
	loadedMaps, err := LoadMaps(mapDir)
	if err != nil {
		log.Printf("mapy sa nepodarilo nacitat: %v", err)
		loadedMaps = make(map[string]*Map)
	}
	if len(loadedMaps) == 0 {
		loadedMaps["generated"] = GenerateMap("generated", 64, 40)
	}

	return &Server{
		maps: loadedMaps,
		hubs: make(map[string]*Hub),
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool { return true },
		},
	}
}

func (server *Server) HandleMaps(w http.ResponseWriter, r *http.Request) {
	// Lekcia 24: Tento endpoint naplni select v klientovi zoznamom dostupnych map.
	writeJSON(w, MapNames(server.maps))
}

func (server *Server) HandleMap(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	hub := server.hubFor(name)
	writeJSON(w, hub.snapshot())
}

func (server *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Lekcia 24: Query parameter map vybera hub, ku ktoremu sa hrac pripoji.
	mapName := r.URL.Query().Get("map")
	if mapName == "" {
		mapName = "generated"
	}
	playerName := strings.TrimSpace(r.URL.Query().Get("name"))
	if playerName == "" {
		playerName = "Hrac"
	}
	// Timy: "auto" alebo prazdne = server vyberie tim, "none" = vlastny tim, cislo = konkretny tim.
	team := TeamAuto
	switch value := r.URL.Query().Get("team"); value {
	case "", "auto":
	case "none":
		team = TeamNone
	default:
		if number, err := strconv.Atoi(value); err == nil && number >= 1 {
			team = number
		}
	}

	conn, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade zlyhal: %v", err)
		return
	}
	defer conn.Close()

	playerID := fmt.Sprintf("player-%d", server.nextPlayerID.Add(1))
	hub := server.hubFor(mapName)
	out := hub.JoinTeam(playerID, playerName, team)
	defer hub.Leave(playerID)

	done := make(chan struct{})
	go func() {
		defer close(done)
		for message := range out {
			if err := conn.WriteJSON(message); err != nil {
				return
			}
		}
	}()

	for {
		var message ClientMessage
		if err := conn.ReadJSON(&message); err != nil {
			return
		}
		if message.Type == "leave" {
			hub.RemovePlayer(playerID)
			return
		}
		hub.Handle(Command{PlayerID: playerID, PlayerName: playerName, Message: message})
		select {
		case <-done:
			return
		default:
		}
	}
}

func (server *Server) hubFor(name string) *Hub {
	server.mu.Lock()
	defer server.mu.Unlock()

	if name == "generated" || name == "" {
		server.maps["generated"] = GenerateMap("generated", 64, 40)
	}
	gameMap, ok := server.maps[name]
	if !ok {
		gameMap = GenerateMap(name, 64, 40)
		server.maps[name] = gameMap
	}
	if hub, ok := server.hubs[name]; ok {
		return hub
	}
	hub := NewHub(gameMap)
	server.hubs[name] = hub
	return hub
}

func (hub *Hub) snapshot() Snapshot {
	hub.mu.Lock()
	defer hub.mu.Unlock()
	return hub.snapshotLocked()
}

func writeJSON(w http.ResponseWriter, value any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(value); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
