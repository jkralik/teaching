# Lekcia 32: Laboratorium zbrani

## Ciel

Pochopit, ako cisla v nastaveniach zbrane menia hru, a vyskusat to ako maly vedecky pokus.

## Co si vysvetlime

- struktura `WeaponSettings` a jej polia
- typ `time.Duration` (`350 * time.Millisecond`)
- konstanta `tankMaxHealth`
- delenie celych cisel v Go (`35 / 2` je `17`, nie `17.5`)
- pokus: zmenim jednu vec, ostatne necham rovnake

## Kodovy krok

Otvor `internal/game/weapon.go` a najdi funkciu `weaponSettings()`. Kazde pole ma komentar, co znamena:

| Pole | Co robi | Predvolene |
| --- | --- | --- |
| `FireInterval` | najkratsi cas medzi dvoma vystrelmi | `350 * time.Millisecond` |
| `MagazineSize` | kolko striel vystrelis pred prebijanim | `3` |
| `ReloadDuration` | ako dlho trva prebijanie | `1500 * time.Millisecond` |
| `BulletSpeed` | kolko policok preleti strela za sekundu | `3.75` |
| `ExplosionRadius` | polomer vybuchu v polickach | `0.9` |
| `Damage` | kolko zivotov zoberie jeden zasah | `35` |

Pod funkciou je konstanta `tankMaxHealth` (predvolene `100`) - s kolkymi zivotmi tank zacina.

Pancier tanku je v subore `internal/game/armor.go`. Klavesom `E` v hre prepnes pancier a v paneli Zbran uvidis, kolkokrat zmensi zasah. Pancier si podrobne upravime v lekcii 36.

Zmen vzdy iba jednu hodnotu, restartuj server (`go run ./cmd/tunneler`) a vyskusaj hru v dvoch oknach prehliadaca.

Novy pojem: experiment znamena zmenit jednu premennu a pozorovat, co sa stane. Ak zmenis vela veci naraz, nevies, ktora zmena sposobila rozdiel.

## Overenie

Do tabulky si zapis pokus a vysledok, napriklad:

| Zmena | Kolko zasahov treba na znicenie tanku | Je to zabavnejsie? |
| --- | --- | --- |
| `Damage: 35` | 3 | |
| `Damage: 50` | 2 | |
| `Damage: 100` | 1 | |
| `Damage: 35`, pancier Bronzovy (klaves `E`) | | |

Spusti aj `go test ./...`. Ak test zlyha, precitaj si, ktoru hodnotu nema rad.

## Uloha

Zmen `tankMaxHealth` a `Damage` tak, aby tank vydrzal presne 4 zasahy. Vysvetli, ako si to vypocital.

Potom vyskusaj, kolko zasahov vydrzi tank s kazdym pancierom. Pozor, Go pri deleni celych cisel zahodi desatinnu cast (`35 / 2` je `17`).

## Mini vyzva

Najdi nastavenie, pri ktorom je hra uplne nehratelna (napriklad strela takmer stoji alebo sa prebija 30 sekund), a vysvetli, preco.
