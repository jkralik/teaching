# Lekcia 34: Nova vlastnost zbrane

## Ciel

Pridat zbraniam uplne novu vlastnost a naprogramovat pravidlo, ktore ju pouzije.

## Co si vysvetlime

- pridanie noveho pola do struktury
- nulova hodnota v Go (`0`, `""`, `false`), ked pole nevyplnis
- ako strela "pamata", z ktorej zbrane bola vystrelena (`Bullet.Weapon`)
- funkcia `weaponSettingsFor(index)` v `hub.go`

## Kodovy krok

Pridame dostrel: strela po urcitej vzdialenosti vybuchne sama.

1. V `internal/game/weapon.go` pridaj do `WeaponSettings` pole `Range float64` a nastav ho kazdej zbrani (napriklad kanon `12`, gulomet `6`).
2. V `internal/game/types.go` pridaj do `Bullet` pole `Traveled float64` s tagom `json:"-"` (klient ho nepotrebuje).
3. V `internal/game/hub.go` vo funkcii `moveBullet` (pri komentari Lekcia 34) tesne pred riadok `bullet.X = nextX` pripocitaj prejdenu vzdialenost a pri prekroceni dostrelu nechaj strelu vybuchnut:

```go
bullet.Traveled += bulletStep
if settings.Range > 0 && bullet.Traveled >= settings.Range {
	hub.explode(bullet, nextX, nextY, settings)
	return false
}
```

Novy pojem: ak ma pole nulovu hodnotu, mozeme ju pouzit ako "vypnute". Preto podmienka `settings.Range > 0` necha zbrane bez dostrelu letiet ako doteraz.

## Overenie

Vystrel do prazdneho tunela. Strela gulometu ma vybuchnut skor ako strela kanona.

## Uloha

Vyber si jednu dalsiu vlastnost a naprogramuj ju:

- `SelfDamage bool` - vybuch zrani aj strelca,
- `Knockback float64` - vybuch odsunie zasiahnuty tank,
- `Pellets int` - jeden vystrel vytvori viac striel do vejara (brokovnica).

## Mini vyzva

Pridaj do `WeaponStatus` novu vlastnost a zobraz ju v paneli Zbran v `web/app.js`, napriklad dostrel zbrane.
