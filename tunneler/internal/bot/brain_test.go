package bot

import (
	"math"
	"strings"
	"testing"
)

// openMap vytvori prazdnu mapu so skalnym okrajom.
func openMap(width, height int) State {
	tiles := make([]string, height)
	for y := range tiles {
		if y == 0 || y == height-1 {
			tiles[y] = strings.Repeat("X", width)
			continue
		}
		tiles[y] = "X" + strings.Repeat(".", width-2) + "X"
	}
	return State{Width: width, Height: height, Tiles: tiles, Tanks: map[string]Tank{}}
}

func tank(id string, x, y float64, team int) Tank {
	return Tank{ID: id, Name: id, X: x, Y: y, Alive: true, Connected: true, Health: 100, MaxHealth: 100, Team: team}
}

func TestShootsEnemyInRange(t *testing.T) {
	state := openMap(20, 10)
	me := tank("me", 2.5, 5.5, 0)
	state.Tanks["me"] = me
	state.Tanks["enemy"] = tank("enemy", 7.5, 5.5, 0)

	action := NewBrain("").Decide(me, state)
	if !action.Shoot {
		t.Fatal("bot mal vystrelit na nepriatela")
	}
	if math.Abs(action.Angle) > 0.001 {
		t.Fatalf("nepriatel je vpravo, uhol mal byt 0, je %v", action.Angle)
	}
}

func TestDoesNotShootTeammate(t *testing.T) {
	state := openMap(20, 10)
	me := tank("me", 2.5, 5.5, 1)
	state.Tanks["me"] = me
	state.Tanks["friend"] = tank("friend", 5.5, 5.5, 1)

	if action := NewBrain("").Decide(me, state); action.Shoot {
		t.Fatal("bot nema strielat na spoluhraca")
	}
}

func TestDoesNotShootThroughRock(t *testing.T) {
	state := openMap(20, 10)
	row := []byte(state.Tiles[5])
	row[5] = 'X'
	state.Tiles[5] = string(row)
	me := tank("me", 2.5, 5.5, 0)
	state.Tanks["me"] = me
	state.Tanks["enemy"] = tank("enemy", 8.5, 5.5, 0)

	action := NewBrain("").Decide(me, state)
	if action.Shoot {
		t.Fatal("cez kamen sa strielat neda")
	}
	if !action.Move {
		t.Fatal("bot mal ist za nepriatelom")
	}
}

func TestMovesTowardFarEnemy(t *testing.T) {
	state := openMap(40, 10)
	me := tank("me", 2.5, 5.5, 0)
	state.Tanks["me"] = me
	state.Tanks["enemy"] = tank("enemy", 30.5, 5.5, 0)

	action := NewBrain("").Decide(me, state)
	if action.Shoot || !action.Move {
		t.Fatalf("daleky nepriatel: ocakavany pohyb bez strelby, je %+v", action)
	}
	if math.Abs(action.Angle) > 0.001 {
		t.Fatalf("mal ist doprava, uhol je %v", action.Angle)
	}
}

func TestCopilotFollowsOwner(t *testing.T) {
	state := openMap(40, 40)
	me := tank("me", 5.5, 30.5, 2)
	owner := tank("Jano", 5.5, 10.5, 2)
	state.Tanks["me"] = me
	state.Tanks["Jano"] = owner

	action := NewBrain("jano").Decide(me, state)
	if !action.Move || action.Shoot {
		t.Fatalf("copilot mal ist za hracom, je %+v", action)
	}
	if math.Abs(action.Angle-math.Pi/2) > 0.001 {
		t.Fatalf("hrac je hore, uhol mal byt Pi/2, je %v", action.Angle)
	}
}

func TestEscapesWhenStuck(t *testing.T) {
	state := openMap(40, 10)
	me := tank("me", 2.5, 5.5, 0)
	state.Tanks["me"] = me
	state.Tanks["enemy"] = tank("enemy", 30.5, 5.5, 0)

	brain := NewBrain("")
	var action Action
	for i := 0; i <= stuckTicks+1; i++ {
		action = brain.Decide(me, state) // tank sa nehybe
	}
	if brain.escape == 0 && brain.stuck != 0 {
		t.Fatalf("bot mal zistit, ze je zaseknuty (stuck=%d)", brain.stuck)
	}
	if !action.Move || action.Angle != brain.escapeAngle {
		t.Fatalf("zaseknuty bot mal utekat smerom %v, akcia %+v", brain.escapeAngle, action)
	}
}

func TestDirectionForAngle(t *testing.T) {
	cases := map[float64]string{0: "right", math.Pi / 2: "up", math.Pi: "left", -math.Pi / 2: "down", math.Pi / 4: "up-right"}
	for angle, want := range cases {
		if got := directionForAngle(angle); got != want {
			t.Errorf("directionForAngle(%v) = %q, chcem %q", angle, got, want)
		}
	}
}
