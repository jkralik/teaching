# Lekcia 36: Pancier tanku

## Ciel

Pochopit, ako pancier zmensuje poskodenie a spomaluje tank, a pridat vlastny pancier.

## Co si vysvetlime

- zoznam pancierov `armorCatalog()` v `internal/game/armor.go`
- delenie poskodenia (`Divisor`) a spomalenie (`SpeedFactor`)
- ako klient kresli pancier ako farebny obrys tanku

## Kodovy krok

Otvor `internal/game/armor.go`. Kazdy pancier ma tieto nastavenia:

| Pole | Co robi | Priklad |
| --- | --- | --- |
| `Name` | nazov v paneli Zbran | `"Bronzovy"` |
| `Color` | farba obrysu tanku (prazdna = ziadny obrys) | `"#cd7f32"` |
| `Divisor` | poskodenie zbrane sa nim vydeli; `0` alebo `1` = plne poskodenie | `2` |
| `SpeedFactor` | nasobok rychlosti tanku (`1` = plna rychlost, `0.5` = polovicna) | `0.85` |

V hre stlac `E` a pancier sa prepne. Obrys tanku zmeni farbu a cim vacsi `Divisor`, tym je obrys hrubsi. Tank s tazsim pancierom jazdi pomalsie.

Funkcia `armoredDamage` vypocita zasah. Vzdy zoberie aspon 1 zivot, aby tank nebol nezranitelny. Funkcia `armoredSpeed` vypocita krok tanku.

Novy pojem: kompromis (trade-off) znamena, ze nieco ziskas a nieco ine stratis. Silny pancier chrani, ale spomaluje.

## Overenie

Spusti `go test ./...`. Test `TestArmorCatalogIsPlayable` skontroluje, ze kazdy pancier ma nazov a rozumne cisla.

Spusti server, pripoj sa v dvoch oknach a zapis si:

| Pancier | Kolko zasahov Kanonom vydrzi | Je tank prilis pomaly? |
| --- | --- | --- |
| Bez panciera | | |
| Bronzovy | | |
| Strieborny | | |
| Zlaty | | |

## Uloha

Pridaj do `armorCatalog()` vlastny pancier, napriklad:

```go
{Name: "Diamantovy", Color: "#6fe3ff", Divisor: 8, SpeedFactor: 0.35},
```

Restartuj server, prepni nan klavesom `E` a over farbu obrysu aj rychlost. Je tento pancier fer?

## Mini vyzva

Napis test, ktory overi, ze pancier s vacsim `Divisor` ma vzdy mensi alebo rovnaky `SpeedFactor` (silnejsia ochrana = pomalsi tank).
