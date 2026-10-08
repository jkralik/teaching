# Lekcia 40: Rozhodovanie bota

## Ciel

Naucit bota nove pravidla: ustupit, ked ma malo zdravia, a zbierat bonusy.

## Co si vysvetlime

- `if` pravidla a ich **poradie** - prve pravidlo, ktore plati, vyhrava (`return action`)
- co bot vie o hre: `me.Health`, `me.MaxHealth`, `state.Bonuses`, `state.Weapons[me.ID]`
- ako pridat do `State` dalsiu informaciu zo servera (pole s rovnakym JSON nazvom ako v `internal/game/types.go`)

## Kodovy krok

Pridaj na zaciatok `Decide` pravidlo ustupu. Ked mam menej ako tretinu zdravia, utekam od nepriatela:

```go
enemy, found := nearestEnemy(me, state)
if found && me.Health*3 < me.MaxHealth {
	action.Move = true
	action.Angle = angleTo(enemy.X, enemy.Y, me.X, me.Y) // opacny smer: od nepriatela ku mne
	return action
}
```

Funkcia `angleTo` vrati smer z prveho bodu na druhy - ked ich vymenime, dostaneme smer od nepriatela.

Potom napis funkciu, ktora najde najblizsi bonus:

```go
func nearestBonus(me Tank, state State) (Bonus, bool) {
	// prejdi state.Bonuses a vrat ten s najmensou distance(me.X, me.Y, bonus.X+0.5, bonus.Y+0.5)
}
```

a pridaj pravidlo: ked nikoho nevidim a bonus je blizsie ako 10 policok, idem pren.

## Overenie

Napis test `TestRetreatsWhenHurt` podla vzoru v `brain_test.go`: tank so `Health: 20` nema strielat a ma ist od nepriatela (uhol `math.Pi`, ked je nepriatel vpravo).

Spusti `go test ./internal/bot/`.

## Uloha

Pozri `state.Weapons[me.ID]`. Ked ma zbran `Reloading == true`, nema zmysel strielat - nech bot radsej ustupi alebo zbiera bonus.

## Mini vyzva

Bot zatial nevie prepnut zbran. Pridaj do `Action` pole `Weapon *int` a v `send()` v `client.go` posli spravu:

```go
if action.Weapon != nil {
	conn.WriteJSON(clientMessage{Type: "weapon", Weapon: action.Weapon})
}
```

Je to ta ista sprava, ktoru bude posielat prehliadac po stlaceni klaves `1`-`9` (lekcia 43). Kedy sa oplati ktora zbran?
