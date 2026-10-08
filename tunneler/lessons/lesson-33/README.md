# Lekcia 33: Druha zbran

## Ciel

Pridat do hry dalsiu zbran s inymi nastaveniami a prepinat medzi zbranami klavesom.

## Co si vysvetlime

- rez struktur `[]WeaponSettings`
- funkcia `weaponCatalog()`, ktora vrati zoznam zbrani
- index v rezi (prva zbran ma index `0`)
- sprava `weapon`, ktoru posiela klient pri stlaceni `Q`

## Kodovy krok

V `internal/game/weapon.go` najdi funkciu `weaponCatalog()`. Teraz obsahuje iba kanon. Pridaj do rezu dalsiu zbran, napriklad gulomet:

```go
func weaponCatalog() []WeaponSettings {
	return []WeaponSettings{
		weaponSettings(),
		{
			Name:            "Gulomet",
			Color:           "#7fd3ff",
			FireInterval:    100 * time.Millisecond,
			MagazineSize:    12,
			ReloadDuration:  2500 * time.Millisecond,
			BulletSpeed:     7,
			ExplosionRadius: 0.5,
			Damage:          10,
		},
	}
}
```

Server netreba inak menit. Uz vie:

- prepnut na dalsiu zbran (`Q`), alebo na konkretnu zbran podla cisla (klavesy `1`-`9` dorobime v lekcii 43),
- pocitat zasobnik kazdej zbrane zvlast,
- poslat strele farbu `Color`, aby ju prehliadac nakreslil inak,
- ukazat nazov zbrane a poskodenie v paneli Zbran.

Novy pojem: data namiesto kodu. Novu zbran sme pridali iba novymi hodnotami v zozname, nemuseli sme pisat novu logiku.

## Overenie

Restartuj server, pripoj sa a stlac `Q`. V paneli Zbran sa ma zmenit nazov, pocet nabojov a poskodenie. Vystrel a over, ze strela ma novu farbu. Spusti `go test ./...`.

## Uloha

Navrhni tretiu zbran, ktora je uplny opak gulometu, napriklad pomaly "Maziar" s velkym vybuchom a malym zasobnikom. Daj jej vlastne meno a farbu.

## Mini vyzva

Vymysli zbran, ktora je silna, ale ma nevyhodu. Napis jednou vetou, kedy sa ju oplati pouzit a kedy nie.
