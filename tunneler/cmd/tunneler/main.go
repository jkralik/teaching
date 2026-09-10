package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jkralik/teaching/tunneler/internal/game"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := game.NewServer("maps")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/maps", server.HandleMaps)
	mux.HandleFunc("GET /api/maps/{name}", server.HandleMap)
	mux.HandleFunc("GET /ws", server.HandleWebSocket)
	mux.Handle("/", http.FileServer(http.Dir("web")))

	addr := ":" + port
	log.Printf("Tunneler server bezi na http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal(err)
	}
}
