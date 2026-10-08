package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jkralik/teaching/tunneler/internal/game"
)

func main() {
	// Lekcia 02: Tento uvodny vypis je pripraveny na oddelenie do funkcie startupMessage(port string).
	port := os.Getenv("PORT")
	if port == "" {
		// Lekcia 03: Predvolenu hodnotu 8080 nahrad konstantou defaultPort.
		port = "8080"
	}

	server := game.NewServer("maps")

	mux := http.NewServeMux()
	// Lekcia 04: Dopln vlastny GET endpoint a vrat z neho JSON odpoved.
	mux.HandleFunc("GET /api/maps", server.HandleMaps)
	mux.HandleFunc("GET /api/maps/{name}", server.HandleMap)
	mux.HandleFunc("GET /api/teams", server.HandleTeams)
	mux.HandleFunc("GET /ws", server.HandleWebSocket)
	mux.Handle("/", http.FileServer(http.Dir("web")))

	addr := "0.0.0.0:" + port
	log.Printf("Tunneler server bezi na http://%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
