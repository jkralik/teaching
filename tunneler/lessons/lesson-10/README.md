# Lekcia 10: Generovanie mapy

## Ciel

Vytvorit mapu pomocou nahody.

## Co si vysvetlime

- balicek `math/rand`
- pravdepodobnost
- hranice mapy
- preco generovana mapa potrebuje pevny okraj

## Kodovy krok

V `internal/game/map.go` zmen pravdepodobnost `TileDirt` v `GenerateMap`, aby sa v hre objavilo viac alebo menej tunelov.

Novy pojem: nahodny vyber je algoritmus, ktory moze mat iny vysledok pri kazdom spusteni.

## Overenie

Spusti hru viac razy a porovnaj vygenerovane mapy. Skontroluj, ze okraj zostava z kamena.

## Uloha

Uprav percento zeme a kamenov v `GenerateMap`.

## Mini vyzva

Vygeneruj mapu, ktora ma viac prazdneho miesta pre rychlejsiu hru.
