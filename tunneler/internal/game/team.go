package game

import (
	"fmt"
	"math"
	"net/http"
	"strings"
)

// TeamSettings popisuje jeden tim.
type TeamSettings struct {
	Name    string `json:"name"`
	Color   string `json:"color"` // farba tela tanku v prehliadaci
	OwnerID string `json:"-"`
}

// TeamStatus je skore timu, ktore posielame prehliadacu.
type TeamStatus struct {
	Team    int    `json:"team"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	Score   int    `json:"score"`
	Players int    `json:"players"`
}

// Specialne cisla timu pri pripajani. Skutocne timy su 1, 2, 3, ... (poradie v teamCatalog() + 1).
const (
	TeamNone = 0  // bez timu: hra kazdy proti kazdemu
	TeamAuto = -1 // server vyberie tim s najmenej hracmi
)

// Timy: Ma strela zranit aj spoluhraca? false = spoluhraca strela preleti a nezrani ho.
const friendlyFire = false

// Timy: Zoznam timov. Pridaj vlastny tim s menom a farbou.
// Farba by sa nemala bit s farbou panciera (obrys tanku).
func teamCatalog() []TeamSettings {
	return []TeamSettings{
		{Name: "Cerveni", Color: "#d9534f"},
		{Name: "Modri", Color: "#4a90d9"},
	}
}

// teamSettingsFor vrati nastavenia timu; ok je false pre "bez timu" alebo neplatne cislo.
func (hub *Hub) teamSettingsFor(team int) (TeamSettings, bool) {
	if team < 1 || team > len(hub.teams) {
		return TeamSettings{}, false
	}
	return hub.teams[team-1], true
}

// Osobne timy dostanu nahodnu farbu. Odtien (hue) je nahodny 0-359,
// sytost a svetlost su pevne, aby bola farba vzdy pekne vyrazna.
const (
	soloColorSaturation = 0.65 // 0 = siva, 1 = najsytejsia
	soloColorLightness  = 0.55 // 0 = cierna, 1 = biela
)

// randomSoloColorLocked vyberie nahodnu farbu osobneho timu, napr. "#3fbf6a".
func (hub *Hub) randomSoloColorLocked() string {
	hue := hub.random.Float64() * 360
	return hslToHex(hue, soloColorSaturation, soloColorLightness)
}

// createSoloTeamLocked vytvori vlastny tim pre hraca, ktory si zvolil "bez timu".
func (hub *Hub) createSoloTeamLocked(playerID string, playerName string) int {
	name := strings.TrimSpace(playerName)
	if name == "" {
		name = "Hrac"
	}
	hub.teams = append(hub.teams, TeamSettings{
		Name:    name,
		Color:   hub.randomSoloColorLocked(),
		OwnerID: playerID,
	})
	return len(hub.teams)
}

func (hub *Hub) removeSoloTeamLocked(playerID string) {
	teamIndex := -1
	for index := hub.playableTeamCount; index < len(hub.teams); index++ {
		if hub.teams[index].OwnerID == playerID {
			teamIndex = index + 1
			break
		}
	}
	if teamIndex == -1 {
		return
	}

	hub.teams = append(hub.teams[:teamIndex-1], hub.teams[teamIndex:]...)
	for id, tank := range hub.tanks {
		if tank.Team == teamIndex {
			tank.Team = TeamNone
			hub.tanks[id] = tank
		} else if tank.Team > teamIndex {
			tank.Team--
			hub.tanks[id] = tank
		}
	}
	for id, bullet := range hub.bullets {
		if bullet.Team > teamIndex {
			bullet.Team--
			hub.bullets[id] = bullet
		} else if bullet.Team == teamIndex {
			bullet.Team = TeamNone
			hub.bullets[id] = bullet
		}
	}
}

// hslToHex prevedie farbu z HSL (odtien, sytost, svetlost) na zapis "#rrggbb".
func hslToHex(hue float64, saturation float64, lightness float64) string {
	chroma := (1 - math.Abs(2*lightness-1)) * saturation
	x := chroma * (1 - math.Abs(math.Mod(hue/60, 2)-1))
	var r, g, b float64
	switch {
	case hue < 60:
		r, g = chroma, x
	case hue < 120:
		r, g = x, chroma
	case hue < 180:
		g, b = chroma, x
	case hue < 240:
		g, b = x, chroma
	case hue < 300:
		r, b = x, chroma
	default:
		r, b = chroma, x
	}
	m := lightness - chroma/2
	toByte := func(value float64) int { return int(math.Round((value + m) * 255)) }
	return fmt.Sprintf("#%02x%02x%02x", toByte(r), toByte(g), toByte(b))
}

// sameTeam povie, ci su dvaja hraci spoluhraci. Hraci bez timu nie su spoluhraci nikoho.
func sameTeam(a int, b int) bool {
	return a != TeamNone && a == b
}

// canHurt povie, ci strela timu shooterTeam smie zranit tank timu targetTeam.
func canHurt(shooterTeam int, targetTeam int) bool {
	return friendlyFire || !sameTeam(shooterTeam, targetTeam)
}

// chooseTeamLocked zmeni volbu hraca na skutocne cislo timu.
func (hub *Hub) chooseTeamLocked(requested int) int {
	if requested == TeamNone || hub.playableTeamCount == 0 {
		return TeamNone
	}
	if requested >= 1 && requested <= hub.playableTeamCount {
		return requested
	}
	// Automaticky vyber: tim s najmenej hracmi, pri rovnosti ten prvy.
	counts := make([]int, hub.playableTeamCount+1)
	for _, tank := range hub.tanks {
		if tank.Team >= 1 && tank.Team <= hub.playableTeamCount {
			counts[tank.Team]++
		}
	}
	best := 1
	for team := 2; team <= hub.playableTeamCount; team++ {
		if counts[team] < counts[best] {
			best = team
		}
	}
	return best
}

// teamStatusesLocked spocita skore a pocet hracov kazdeho timu.
func (hub *Hub) teamStatusesLocked() []TeamStatus {
	statuses := make([]TeamStatus, len(hub.teams))
	for index, settings := range hub.teams {
		statuses[index] = TeamStatus{Team: index + 1, Name: settings.Name, Color: settings.Color}
	}
	for _, tank := range hub.tanks {
		if tank.Team < 1 || tank.Team > len(statuses) {
			continue
		}
		statuses[tank.Team-1].Score += tank.Score
		statuses[tank.Team-1].Players++
	}
	return statuses
}

// HandleTeams posle prehliadacu zoznam timov pre vyber pri pripojeni.
func (server *Server) HandleTeams(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, teamCatalog())
}
