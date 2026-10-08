# Lekcia 14: Pohyb tanku

## Ciel

Napogramovat pohyb podla smeru.

## Co si vysvetlime

- `switch`
- suradnice v mriezke
- smer tanku
- testovanie malej funkcie

## Kodovy krok

V `internal/game/hub.go` uprav `nextPosition`, aby pri prikaze `move` vratil novu poziciu podla smeru. Pridaj test pre vsetky styri smery.

Novy pojem: `switch` vybera jednu vetvu podla hodnoty a mala cista funkcia sa lahko testuje.

## Overenie

Spusti `go test ./internal/game` a potom sa pohni vo vsetkych styroch smeroch.

## Uloha

Uprav alebo otestuj funkciu `nextPosition`.

## Mini vyzva

Pridaj novy prikaz, ktory tank otoci bez pohybu.
