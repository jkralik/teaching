package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// tickInterval: ako casto sa bot rozhoduje. Prehliadac posiela krok tiez kazdych 30 ms,
// rychlejsie to nema zmysel - server aj tak pusti len jeden krok za 30 ms.
const tickInterval = 30 * time.Millisecond

// Options su nastavenia z prikazoveho riadku (cmd/bot).
type Options struct {
	Server string // napr. http://localhost:8080
	Map    string // nazov mapy, napr. arena
	Name   string // meno tanku
	Team   string // "auto", "none" alebo cislo timu
	Follow string // meno hraca, ktoremu bot pomaha (copilot)
}

// Decider je "mozog" tanku. Brain z brain.go ho implementuje,
// ale v turnaji (lekcia 42) moze mat kazdy svoj vlastny.
type Decider interface {
	Decide(me Tank, state State) Action
}

// Run pripoji bota k serveru a hra, kym sa nezrusi ctx alebo nespadne spojenie.
func Run(ctx context.Context, options Options, brain Decider) error {
	team := options.Team
	if team == "" && options.Follow != "" {
		team = followTeam(ctx, options)
	}
	if team == "" {
		team = "auto"
	}

	wsURL, err := websocketURL(options.Server, options.Map, options.Name, team)
	if err != nil {
		return err
	}
	conn, _, err := websocket.DefaultDialer.DialContext(ctx, wsURL, nil)
	if err != nil {
		return fmt.Errorf("pripojenie na %s: %w", wsURL, err)
	}
	defer conn.Close()
	log.Printf("bot %q sa pripojil na mapu %q (tim %s)", options.Name, options.Map, team)

	var (
		mu       sync.Mutex
		state    State
		playerID string
		hasState bool
	)
	readErr := make(chan error, 1)
	go func() {
		for {
			var message serverMessage
			if err := conn.ReadJSON(&message); err != nil {
				readErr <- err
				return
			}
			if message.Type == "error" {
				readErr <- errors.New(message.Error)
				return
			}
			if message.Type != "state" {
				continue
			}
			mu.Lock()
			state = message.State
			hasState = true
			if message.PlayerID != "" {
				playerID = message.PlayerID
			}
			mu.Unlock()
		}
	}()

	ticker := time.NewTicker(tickInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			_ = conn.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, "bot konci"))
			return nil
		case err := <-readErr:
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("spojenie so serverom: %w", err)
		case <-ticker.C:
			mu.Lock()
			current, id, ready := state, playerID, hasState
			mu.Unlock()
			if !ready || id == "" {
				continue
			}
			me, ok := current.Tanks[id]
			if !ok || !me.Alive || !me.Connected {
				continue // mrtvy tank caka na ozivenie, server by prikazy aj tak ignoroval
			}
			action := brain.Decide(me, current)
			if err := send(conn, me, action); err != nil {
				return fmt.Errorf("posielanie prikazu: %w", err)
			}
		}
	}
}

// send prelozi Action na rovnake spravy, ake posiela prehliadac.
func send(conn *websocket.Conn, me Tank, action Action) error {
	// Strela leti tam, kam je tank otoceny. Otocit sa da len spravou "move",
	// preto pred vystrelom tank natocime (a pritom spravi aj krok).
	aim := action.Shoot && math.Abs(angleDifference(me.Angle, action.Angle)) > 0.05
	if action.Move || aim {
		angle := action.Angle
		message := clientMessage{Type: "move", Direction: directionForAngle(angle), Angle: &angle}
		if err := conn.WriteJSON(message); err != nil {
			return err
		}
	}
	if action.Shoot {
		if err := conn.WriteJSON(clientMessage{Type: "shoot"}); err != nil {
			return err
		}
	}
	return nil
}

// followTeam zisti cez HTTP API, v akom time je hrac, ktoremu bot pomaha.
func followTeam(ctx context.Context, options Options) string {
	base := strings.TrimRight(options.Server, "/")
	request, err := http.NewRequestWithContext(ctx, http.MethodGet,
		base+"/api/maps/"+url.PathEscape(options.Map), nil)
	if err != nil {
		return ""
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		log.Printf("copilot: neviem nacitat mapu: %v", err)
		return ""
	}
	defer response.Body.Close()
	var state State
	if err := json.NewDecoder(response.Body).Decode(&state); err != nil {
		log.Printf("copilot: zla odpoved servera: %v", err)
		return ""
	}
	for _, tank := range state.Tanks {
		if strings.EqualFold(tank.Name, options.Follow) {
			if tank.Team >= 1 {
				return fmt.Sprint(tank.Team)
			}
			log.Printf("copilot: hrac %q nema tim, bot ide do timu auto", options.Follow)
			return ""
		}
	}
	log.Printf("copilot: hrac %q nie je na mape %q, bot ide do timu auto", options.Follow, options.Map)
	return ""
}

// websocketURL z http://server spravi ws://server/ws?map=...&name=...&team=...
func websocketURL(server string, mapName string, name string, team string) (string, error) {
	parsed, err := url.Parse(server)
	if err != nil {
		return "", fmt.Errorf("zla adresa servera %q: %w", server, err)
	}
	switch parsed.Scheme {
	case "http", "ws":
		parsed.Scheme = "ws"
	case "https", "wss":
		parsed.Scheme = "wss"
	default:
		return "", fmt.Errorf("adresa servera musi zacinat http:// alebo https://, nie %q", server)
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/") + "/ws"
	query := url.Values{}
	query.Set("map", mapName)
	query.Set("name", name)
	query.Set("team", team)
	parsed.RawQuery = query.Encode()
	return parsed.String(), nil
}

// directionForAngle vyberie najblizsi z 8 smerov (server ho pouziva len na zobrazenie).
func directionForAngle(angle float64) string {
	directions := []string{"right", "up-right", "up", "up-left", "left", "down-left", "down", "down-right"}
	index := int(math.Round(normalizeAngle(angle)/(math.Pi/4))) % len(directions)
	return directions[index]
}

// normalizeAngle vrati uhol v rozsahu 0 az 2*Pi.
func normalizeAngle(angle float64) float64 {
	angle = math.Mod(angle, 2*math.Pi)
	if angle < 0 {
		angle += 2 * math.Pi
	}
	return angle
}

// angleDifference vrati najkratsi rozdiel dvoch uhlov (-Pi az Pi).
func angleDifference(a float64, b float64) float64 {
	diff := normalizeAngle(b - a)
	if diff > math.Pi {
		diff -= 2 * math.Pi
	}
	return diff
}
