# Lekcia 15: Kolizie

## Ciel

Zabranit tanku prejst cez kamen alebo okraj mapy.

## Co si vysvetlime

- pevna prekazka
- kontrola cieloveho policka
- navrat z funkcie pomocou `return`
- preco su okraje mapy dolezite

## Kodovy krok

V `moveTank` dopln pravidlo, ktore nepovoli prejst cez `TileRock` ani mimo mapy.

Novy pojem: kontrola kolizie musi prebehnut pred zmenou stavu; `return` ukonci pohyb hned, ked je ciel obsadeny.

## Overenie

Pridaj test s kamenom a test s okrajom mapy. Spusti `go test ./internal/game`.

## Uloha

V `moveTank` najdi pravidlo pre kamen `X` a vyskusaj ho zmenit.

## Mini vyzva

Pridaj novy typ policka, ktory spomali tank alebo ho zastavi.
