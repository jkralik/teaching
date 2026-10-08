package bot

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"github.com/jkralik/teaching/tunneler/internal/game"
)

func startServer(t *testing.T) *httptest.Server {
	t.Helper()
	server := game.NewServer(t.TempDir())
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/maps/{name}", server.HandleMap)
	mux.HandleFunc("GET /ws", server.HandleWebSocket)
	httpServer := httptest.NewServer(mux)
	t.Cleanup(httpServer.Close)
	return httpServer
}

func mapState(t *testing.T, server string, name string) State {
	t.Helper()
	response, err := http.Get(server + "/api/maps/" + name)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	var state State
	if err := json.NewDecoder(response.Body).Decode(&state); err != nil {
		t.Fatal(err)
	}
	return state
}

func findTank(state State, name string) (Tank, bool) {
	for _, tank := range state.Tanks {
		if tank.Name == name {
			return tank, true
		}
	}
	return Tank{}, false
}

func TestBotJoinsAndPlays(t *testing.T) {
	server := startServer(t)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	options := Options{Server: server.URL, Map: "test", Name: "Robot"}
	if err := Run(ctx, options, NewBrain("")); err != nil {
		t.Fatalf("Run: %v", err)
	}
	if _, ok := findTank(mapState(t, server.URL, "test"), "Robot"); !ok {
		t.Fatal("tank Robot sa na mape neobjavil")
	}
}

func TestCopilotJoinsOwnersTeam(t *testing.T) {
	server := startServer(t)
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws?map=test&name=Jano&team=2"
	human, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer human.Close()
	var first serverMessage
	if err := human.ReadJSON(&first); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	options := Options{Server: server.URL, Map: "test", Name: "Pomocnik", Follow: "jano"}
	if err := Run(ctx, options, NewBrain(options.Follow)); err != nil {
		t.Fatalf("Run: %v", err)
	}
	copilot, ok := findTank(mapState(t, server.URL, "test"), "Pomocnik")
	if !ok {
		t.Fatal("copilot sa na mape neobjavil")
	}
	if copilot.Team != 2 {
		t.Fatalf("copilot mal byt v time 2 ako Jano, je v time %d", copilot.Team)
	}
}

func TestWebsocketURL(t *testing.T) {
	got, err := websocketURL("https://hra.sk/", "arena", "Robot 1", "auto")
	if err != nil {
		t.Fatal(err)
	}
	want := "wss://hra.sk/ws?map=arena&name=Robot+1&team=auto"
	if got != want {
		t.Fatalf("websocketURL = %q, chcem %q", got, want)
	}
}
