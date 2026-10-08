package game

// ArmorSettings popisuje jeden typ panciera tanku.
type ArmorSettings struct {
	Name        string
	Color       string  // farba obrysu tanku v prehliadaci, prazdna = bez obrysu
	Divisor     int     // poskodenie zbrane sa vydeli tymto cislom, 0 alebo 1 = plne poskodenie
	SpeedFactor float64 // nasobok rychlosti tanku, 1 = plna rychlost, 0.5 = polovicna
}

// Lekcia 36: Zoznam pancierov. Pridaj vlastny pancier s inou farbou, ochranou a rychlostou.
// Tank zacina s prvym pancierom v zozname a klavesom E prepina na dalsi.
func armorCatalog() []ArmorSettings {
	return []ArmorSettings{
		{Name: "Bez panciera", Color: "", Divisor: 0, SpeedFactor: 1},
		{Name: "Bronzovy", Color: "#cd7f32", Divisor: 2, SpeedFactor: 0.85},
		{Name: "Strieborny", Color: "#c9d3dc", Divisor: 3, SpeedFactor: 0.7},
		{Name: "Zlaty", Color: "#ffd84a", Divisor: 5, SpeedFactor: 0.5},
	}
}

// armoredDamage vrati poskodenie po zapocitani panciera.
// Zasah vzdy zoberie aspon 1 zivot, aby tank nebol nezranitelny.
func armoredDamage(damage int, divisor int) int {
	if divisor <= 1 {
		return damage
	}
	return max(1, damage/divisor)
}

// armoredSpeed vrati dlzku kroku tanku po zapocitani panciera.
func armoredSpeed(step float64, armor ArmorSettings) float64 {
	if armor.SpeedFactor <= 0 {
		return step
	}
	return step * armor.SpeedFactor
}

func (hub *Hub) armorSettingsFor(index int) ArmorSettings {
	if index < 0 || index >= len(hub.armors) {
		return ArmorSettings{Name: "Bez panciera", SpeedFactor: 1}
	}
	return hub.armors[index]
}
