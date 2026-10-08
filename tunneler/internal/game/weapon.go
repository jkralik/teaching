package game

import "time"

// WeaponSettings obsahuje hodnoty, ktore mozno menit pri cviceni.
// Lekcia 34: Sem mozes pridat novu vlastnost zbrane, napriklad dostrel.
type WeaponSettings struct {
	Name            string
	Color           string
	FireInterval    time.Duration
	MagazineSize    int
	ReloadDuration  time.Duration
	BulletSpeed     float64
	ExplosionRadius float64
	Damage          int
}

// Pocet zivotov, s ktorymi tank zacina hru.
const tankMaxHealth = 100

type WeaponStatus struct {
	Name              string `json:"name"`
	Color             string `json:"color"`
	Damage            int    `json:"damage"`
	Index             int    `json:"index"`
	Count             int    `json:"count"`
	MagazineSize      int    `json:"magazineSize"`
	ShotsRemaining    int    `json:"shotsRemaining"`
	Reloading         bool   `json:"reloading"`
	ReloadRemainingMs int64  `json:"reloadRemainingMs"`
}

type weaponState struct {
	shotsFired     int
	lastShotAt     time.Time
	reloadingUntil time.Time
}

// Lekcia 32: Zmen hodnoty kanonu, spusti hru a sleduj, ako sa zmeni suboj.
func weaponSettings() WeaponSettings {
	return WeaponSettings{
		Name:            "Kanon",
		Color:           "#e45d3d",              // farba strely v prehliadaci
		FireInterval:    350 * time.Millisecond, // cas medzi vystrelmi
		MagazineSize:    3,                      // pocet vystrelov pred prebijanim
		ReloadDuration:  1500 * time.Millisecond,
		BulletSpeed:     3.75, // herne policka za sekundu
		ExplosionRadius: 0.9,  // polomer v hernych polickach
		Damage:          35,   // kolko zivotov zoberie jeden zasah
	}
}

// weaponCatalog vrati vsetky zbrane, medzi ktorymi hrac prepina klavesom Q (a od lekcie 43 aj 1-9).
// Lekcia 33: Pridaj do zoznamu dalsiu zbran s inym menom, farbou a nastaveniami.
func weaponCatalog() []WeaponSettings {
	return []WeaponSettings{
		weaponSettings(),
	}
}

func (state *weaponState) tryShoot(now time.Time, settings WeaponSettings) bool {
	if !state.reloadingUntil.IsZero() {
		if now.Before(state.reloadingUntil) {
			return false
		}
		state.shotsFired = 0
		state.reloadingUntil = time.Time{}
	}
	if !state.lastShotAt.IsZero() && now.Sub(state.lastShotAt) < settings.FireInterval {
		return false
	}

	state.shotsFired++
	state.lastShotAt = now
	if state.shotsFired == settings.MagazineSize {
		state.reloadingUntil = now.Add(settings.ReloadDuration)
	}
	return true
}

func (state weaponState) status(now time.Time, settings WeaponSettings) WeaponStatus {
	status := WeaponStatus{
		Name:           settings.Name,
		Color:          settings.Color,
		Damage:         settings.Damage,
		MagazineSize:   settings.MagazineSize,
		ShotsRemaining: max(0, settings.MagazineSize-state.shotsFired),
	}
	if !state.reloadingUntil.IsZero() {
		remaining := state.reloadingUntil.Sub(now)
		if remaining > 0 {
			status.Reloading = true
			status.ShotsRemaining = 0
			status.ReloadRemainingMs = remaining.Milliseconds()
			if status.ReloadRemainingMs == 0 {
				status.ReloadRemainingMs = 1
			}
		} else {
			status.ShotsRemaining = settings.MagazineSize
		}
	}
	return status
}
