package game

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
)

type Server struct {
	mu       sync.Mutex
	maps     map[string]*Map
	hubs     map[string]*Hub
	upgrader websocket.Upgrader
}

func NewServer(mapDir string) *Server {
	loadedMaps, err := LoadMaps(mapDir)
	if err != nil {
		log.Printf("mapy sa nepodarilo nacitat: %v", err)
		loadedMaps = make(map[string]*Map)
	}
	if len(loadedMaps) == 0 {
		loadedMaps["generated"] = GenerateMap("generated", 32, 20)
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
	writeJSON(w, MapNames(server.maps))
}

func (server *Server) HandleMap(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	hub := server.hubFor(name)
	writeJSON(w, hub.snapshot())
}

func (server *Server) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	mapName := r.URL.Query().Get("map")
	if mapName == "" {
		mapName = "generated"
	}
	playerName := strings.TrimSpace(r.URL.Query().Get("name"))
	if playerName == "" {
		playerName = "Hrac"
	}

	conn, err := server.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("websocket upgrade zlyhal: %v", err)
		return
	}
	defer conn.Close()

	playerID := r.RemoteAddr
	hub := server.hubFor(mapName)
	out := hub.Join(playerID, playerName)
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
		server.maps["generated"] = GenerateMap("generated", 32, 20)
	}
	gameMap, ok := server.maps[name]
	if !ok {
		gameMap = GenerateMap(name, 32, 20)
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
