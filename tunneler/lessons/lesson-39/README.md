# Lekcia 39: Prvy bot

## Ciel

Spustit AI tank, ktory hra proti nam, a zmenit jeho spravanie pomocou konstant.

## Co si vysvetlime

- bot je **samostatny program** v `cmd/bot` - k serveru sa pripaja cez WebSocket presne ako prehliadac
- server nevie, ze hra robot; bot moze len to, co hrac s klavesnicou (ziadne podvadzanie)
- "mozog" bota je v `internal/bot/brain.go`, funkcia `Decide(me, state)` sa vola kazdych 30 ms
- `Decide` dostane moj tank a celu mapu a vrati `Action`: kam ist (`Move`, `Angle`) a ci strielat (`Shoot`)

## Kodovy krok

Spusti server v jednom terminali:

```bash
go run ./cmd/tunneler
```

V druhom terminali spusti bota:

```bash
go run ./cmd/bot -map arena
```

Otvor `http://localhost:8080`, pripoj sa na mapu `arena` a bojuj proti tanku `Robot`.

Bot sa rozhoduje v tomto poradi (pozri `Decide`):

1. vidi nepriatela blizsie ako `shootDistance` a v ceste nie je kamen -> otoci sa a striela
2. dlho sa nepohol (`stuckTicks`) -> chvilu ide nahodnym smerom
3. je copilot a jeho hrac je daleko -> ide za hracom (lekcia 42)
4. pozna nepriatela -> ide za nim
5. inak sa tula po mape

Uhol: `0` = doprava, `math.Pi/2` = hore, `math.Pi` = dolava. Pocita ho funkcia `angleTo`.

Novy pojem: **konstanta** je hodnota, ktoru program nemeni, ale my ju mozeme lahko prepisat a vyskusat, co sa stane.

## Overenie

Spusti `go test ./internal/bot/`. Testy overia, ze bot striela na nepriatela, nestriela na spoluhraca a cez kamen.

## Uloha

Zmen konstanty v `brain.go`, restartuj bota (`Ctrl+C` a znova `go run ./cmd/bot`) a zapis si, co sa stalo:

| Zmena | Co bot robi |
| --- | --- |
| `shootDistance = 3.0` | |
| `shootDistance = 20.0` | |
| `wanderTicks = 5` | |
| `stuckTicks = 3` | |

Preco pri `stuckTicks = 3` bot nevie prekopat zem? (Napoveda: kopanie jedneho policka trva asi 12 krokov.)

## Mini vyzva

Spusti dvoch botov naraz s inymi menami: `go run ./cmd/bot -name Robot2 -team none`. Kto vyhra?
