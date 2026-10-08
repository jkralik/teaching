# Lekcia 07: Rezy a mapa

## Ciel

Pochopit, ako Go uchovava mapu ako mriezku.

## Co si vysvetlime

- slice
- dvojrozmerna struktura `[][]Tile`
- suradnice `x` a `y`
- rozdiel medzi riadkom a stlpcom

## Kodovy krok

V `internal/game/map.go` dopln metodu `Inside(x, y int) bool`, ktora vrati `true`, ak suradnice patria do mapy. Pouzi ju v teste alebo pri overeni suradnice pred pristupom k policku.

Novy pojem: metoda patri konkretnemu typu a `bool` vyjadruje odpoved ano alebo nie.

## Overenie

V `internal/game/map_test.go` pridaj pripady pre vnutorne policko, okraj a suradnicu mimo mapy. Spusti `go test ./internal/game`.

## Uloha

Najdi typ `Map` a vysvetli, kde je ulozena sirka, vyska a policka mapy.

## Mini vyzva

Pouzi metodu `Inside` v `TileAt`, aby bolo rozhodnutie o hranici mapy na jednom mieste.