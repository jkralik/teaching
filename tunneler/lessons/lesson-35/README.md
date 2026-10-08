# Lekcia 35: Test a balans zbrani

## Ciel

Overit testom, ze kazda zbran je hratelna, a vyladit zbrane tak, aby ziadna nebola najlepsia vo vsetkom.

## Co si vysvetlime

- test `TestWeaponCatalogHasPlayableWeapons` v `internal/game/weapon_test.go`
- cyklus `for index, weapon := range catalog`
- `t.Fatalf` a zrozumitelna chybova sprava
- balans: kazda vyhoda ma mat nevyhodu

## Kodovy krok

Otvor `internal/game/weapon_test.go`. Test uz kontroluje, ze kazda zbran ma meno, kladne poskodenie a ze sa mena neopakuju. Pri komentari Lekcia 35 dopln vlastne pravidlo, napriklad:

```go
if weapon.Damage >= tankMaxHealth {
	t.Fatalf("weapon %q destroys a tank with one hit", weapon.Name)
}
```

Tank bez panciera dostane plne poskodenie, preto staci kontrolovat `weapon.Damage`.

Potom vymysli jednoduche cislo "sila zbrane" a porovnaj zbrane medzi sebou:

```go
damagePerSecond := float64(weapon.Damage*weapon.MagazineSize) /
	(float64(weapon.MagazineSize-1)*weapon.FireInterval.Seconds() + weapon.ReloadDuration.Seconds())
t.Logf("%s: %.1f poskodenia za sekundu", weapon.Name, damagePerSecond)
```

Spusti `go test -v ./internal/game -run Weapon` a pozri sa na vypis.

Novy pojem: balans hry znamena upravovat cisla tak, aby kazda volba mala zmysel.

## Overenie

Test prejde pre vsetky zbrane. Ak niektoru zbran schvalne pokazis (napriklad `Damage: 0`), test zlyha a chybova sprava povie, ktora zbran je zla.

## Uloha

Usporiadaj medzi sebou turnaj: kazdy tim prinesie jednu zbran, vlozite ich do `weaponCatalog()` a zahrate si. Po hre upravte cisla zbrane, ktora vyhravala prilis casto.

## Mini vyzva

Vypis poskodenie za sekundu pre kazdu kombinaciu zbrane a panciera (dva vnorene cykly cez `weaponCatalog()` a `armorCatalog()`, poskodenie prepocitaj cez `armoredDamage`).

Dalsia vyzva: pridaj test, ktory overi, ze ziadna zbran nema viac ako dvojnasobne poskodenie za sekundu oproti najslabsej zbrani.
