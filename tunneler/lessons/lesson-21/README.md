# Lekcia 21: Strely

## Ciel

Pridat objekt, ktory sa hybe sam.

## Co si vysvetlime

- struktura `Bullet`
- vlastnik strely
- tik hry cez `time.Ticker`
- zanik strely po naraze

## Kodovy krok

V `internal/game/hub.go` vytvor pri prikaze `shoot` novy `Bullet` a v `tick` ho posun kazdych 30 milisekund pomocou `time.Ticker`.

Novy pojem: ticker pravidelne vytvara udalosti, ktore mozu pohanat herny cas.

## Overenie

Vystrel a sleduj, ze strela sa sama pohybuje a po naraze zmizne.

## Uloha

Zmen rychlost striel upravou casu v `time.NewTicker`.

## Mini vyzva

Pridaj limit, aby hrac nemohol vystrelit nekonecne vela striel naraz.
