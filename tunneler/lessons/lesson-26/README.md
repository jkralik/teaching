# Lekcia 26: Testy

## Ciel

Napisat maly test pre herne pravidlo.

## Co si vysvetlime

- subory `*_test.go`
- `testing.T`
- co znamena ocakavanie
- prikaz `go test ./...`

## Kodovy krok

V `internal/game/map_test.go` pridaj test pre `nextPosition` alebo `TileAt`, ktory zlyha, ak sa zmeni ocakavane herne pravidlo.

Novy pojem: test je program, ktory automaticky overuje ocakavane spravanie funkcie.

## Overenie

Spusti `go test ./...` a potom test docasne uprav tak, aby zlyhal. Precitaj vystup a vrat ocakavanie spat.

## Uloha

Pridaj test pre funkciu `nextPosition`.

## Mini vyzva

Pridaj test, ze strela znici zem a zmizne.
