# Lekcia 38: Prekvapenie ?

## Ciel

Pridat na mapu riziko: policko `?` moze dat dobry bonus, ale aj negativny efekt.

## Co si vysvetlime

- pole `MysteryOnly` v `bonusCatalog()`
- konstantu `mysteryChance` (pravdepodobnost, ze sa objavi `?`)
- nahodu v hre a preco robi hru napinavejsou

## Kodovy krok

Fialove policko `?` neprezradi, co v nom je. Ked nan tank prejde, server nahodne vyberie jeden bonus z `bonusCatalog()`, ktory tank este nema. Moze to byt dobry bonus aj taky, ktory ma `MysteryOnly: true`.

Bonus s `MysteryOnly: true` sa nikdy neobjavi na mape priamo, len cez `?`. Uz je pripraveny negativny efekt:

```go
{Name: "Blato", Symbol: "B", Color: "#8a6b4a", Duration: 5 * time.Second, SpeedFactor: 0.5, MysteryOnly: true},
```

Blato na 5 sekund spomali tank na polovicu.

Konstanta `mysteryChance = 0.3` znamena, ze priblizne 3 z 10 novych bonusov budu `?`.

Novy pojem: pravdepodobnost je cislo od 0 do 1, ktore hovori, ako casto sa nieco stane. 0 = nikdy, 1 = vzdy.

## Overenie

Spusti `go test ./...`. Test `TestMysteryTileGivesRandomEffect` overi, ze policko `?` da tanku nejaky efekt.

Spusti server a zbieraj `?`. Zapis si, co si dostal:

| Pokus | Efekt |
| --- | --- |
| 1 | |
| 2 | |
| 3 | |

## Uloha

Pridaj vlastne negativne efekty, napriklad:

```go
{Name: "Slaba strela", Symbol: "S", Color: "#9a9a9a", Duration: 5 * time.Second, DamageFactor: 0.5, MysteryOnly: true},
{Name: "Zaseknuty pas", Symbol: "Z", Color: "#555555", Duration: 2 * time.Second, SpeedFactor: 0.1, MysteryOnly: true},
```

Potom zmen `mysteryChance` na `0.8`. Chces este zbierat `?`, ked je tam vela zlych efektov?

## Mini vyzva

Napis test, ktory overi, ze kazdy bonus s `MysteryOnly: true` ma `Duration` vacsie ako 0, aby negativny efekt raz skoncil.
