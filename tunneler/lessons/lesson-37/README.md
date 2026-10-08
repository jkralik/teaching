# Lekcia 37: Bonusy na mape

## Ciel

Pochopit, ako sa na mape objavuju docasne bonusy, a pridat vlastne bonusy: Turbo, Silnu strelu a Nesmrtelnost.

## Co si vysvetlime

- zoznam bonusov `bonusCatalog()` v `internal/game/bonus.go`
- konstanty, ktore riadia objavovanie bonusov na mape
- jednorazovy efekt (`Duration` je 0) a docasny efekt (`Duration` je napr. 5 sekund)
- preco ten isty tank nemoze mat ten isty bonus dvakrat naraz

## Kodovy krok

Otvor `internal/game/bonus.go`. Kazdy bonus ma tieto nastavenia (co nevyplnis, nic nemeni):

| Pole | Co robi | Priklad |
| --- | --- | --- |
| `Name` | nazov v paneli Bonusy | `"Turbo"` |
| `Symbol` | pismeno na policku | `"T"` |
| `Color` | farba policka | `"#7dff8a"` |
| `Duration` | ako dlho efekt trva; `0` = jednorazovy | `5 * time.Second` |
| `InstantReload` | hned nabije vsetky zbrane | `true` |
| `SpeedFactor` | nasobok rychlosti tanku | `1.6` |
| `DamageFactor` | nasobok poskodenia striel | `2` |
| `Invulnerable` | tank nedostava poskodenie | `true` |
| `MysteryOnly` | bonus padne iba z policka `?` (lekcia 38) | `true` |

Na zaciatku je v katalogu iba `Okamzity reload` (modre `R`). Ked tank prejde cez policko s bonusom, ziska ho. Docasny bonus sa ukaze v paneli `Bonusy` s odpocitavanim. Ak tank uz taky bonus ma, bonus zostane lezat na mape pre ostatnych.

Nad katalogom su konstanty `bonusSpawnInterval` (ako casto sa objavi bonus), `bonusLifetime` (kedy zmizne), `bonusMaxOnMap` (kolko ich moze byt naraz) a `bonusPickupRadius`.

Novy pojem: docasny efekt je zmena, ktora ma cas konca. Server si pre tank pamata, kedy efekt skonci, a potom ho sam zrusi.

## Overenie

Spusti `go test ./...`. Test `TestBonusCatalogIsPlayable` skontroluje, ze kazdy bonus ma nazov, znacku, farbu a ze netrva dlhsie ako 20 sekund.

Spusti server a pockaj par sekund, kym sa na mape objavi `R`. Vystrel cely zasobnik a prejdi cez `R`. Zbran je hned nabita?

## Uloha

Pridaj do `bonusCatalog()` tieto tri bonusy:

```go
{Name: "Turbo", Symbol: "T", Color: "#7dff8a", Duration: 5 * time.Second, SpeedFactor: 1.6},
{Name: "Silna strela", Symbol: "D", Color: "#ff7a59", Duration: 5 * time.Second, DamageFactor: 2},
{Name: "Nesmrtelnost", Symbol: "N", Color: "#fff27a", Duration: 5 * time.Second, Invulnerable: true},
```

Restartuj server a vyskusaj kazdy bonus vo dvojici:

| Bonus | Co sa stalo | Je to fer? |
| --- | --- | --- |
| Turbo | | |
| Silna strela | | |
| Nesmrtelnost | (okolo tanku je zlty stit) | |

Skus zmenit `bonusSpawnInterval` alebo `bonusMaxOnMap`. Je hra zabavnejsia s viac bonusmi alebo s menej?

## Mini vyzva

Vymysli kombinovany bonus, napriklad `Berserk`, ktory na 4 sekundy zdvojnasobi poskodenie a zaroven zrychli tank. Staci jeden riadok s viacerymi poliami.
