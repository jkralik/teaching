// Prikaz bot spusti AI tank, ktory sa pripoji k serveru cez WebSocket ako hrac.
//
//	go run ./cmd/bot -map arena
//	go run ./cmd/bot -map arena -name Pomocnik -follow Jano   (copilot v time hraca Jano)
package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/jkralik/teaching/tunneler/internal/bot"
)

func main() {
	options := bot.Options{}
	flag.StringVar(&options.Server, "server", "http://localhost:58080", "adresa servera")
	flag.StringVar(&options.Map, "map", "arena", "nazov mapy")
	flag.StringVar(&options.Name, "name", "Robot", "meno tanku")
	flag.StringVar(&options.Team, "team", "", "tim: auto, none alebo cislo (prazdne = auto alebo tim hraca z -follow)")
	flag.StringVar(&options.Follow, "follow", "", "meno hraca, ktoremu bot pomaha (copilot)")
	flag.Parse()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	brain := bot.NewBrain(options.Follow)
	for {
		err := bot.Run(ctx, options, brain)
		if ctx.Err() != nil {
			log.Print("bot konci")
			return
		}
		// Ked spojenie spadne, bot sa po chvili pripoji znova s rovnakym menom
		// a server mu vrati jeho povodny tank.
		log.Printf("%v - skusim znova o 3 s", err)
		select {
		case <-ctx.Done():
			return
		case <-time.After(3 * time.Second):
		}
	}
}
