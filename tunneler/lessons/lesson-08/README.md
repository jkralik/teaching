# Lekcia 08: Datove struktury a rychlost

## Ciel

Porovnat dve jednoduche datove struktury a zistit, preco rovnaky problem nemusi mat vzdy rovnako rychle riesenie.

## Co si vysvetlime

- `slice` je usporiadany zoznam hodnot
- `map` uklada hodnoty pod klucom, napriklad hraca pod jeho ID
- hladanie v `slice` kontroluje hodnoty postupne
- hladanie v `map` najde hodnotu priamo podla kluca
- rychlejsia struktura moze potrebovat viac pamate

## Priklad z hry

Predstav si zoznam hracov:

```go
players := []Tank{
    {ID: "ada"},
    {ID: "ben"},
    {ID: "cyril"},
}
```

Ak hladame hraca `cyril`, pri `slice` sa kontroluju ID hracov jedno po druhom. Pri mnohych hracoch to trva dlhsie.

Mapa pouziva ID ako kluc:

```go
playersByID := map[string]Tank{
    "ada":   {ID: "ada"},
    "ben":   {ID: "ben"},
    "cyril": {ID: "cyril"},
}

player := playersByID["cyril"]
```

Pomocou `map` vieme ziskat hraca podla ID rychlo aj vtedy, ked ich je vela. Mapa si vsak musi pamatat aj kluce a svoju vnutornu organizaciu, preto zaberie viac pamate ako jednoduchy `slice` s rovnakymi hodnotami.

## Algoritmus hladania

Hladanie v `slice` ma cas $O(n)$: pri dvojnasobnom pocte hracov moze skontrolovat priblizne dvojnasobok hodnot.

Hladanie v `map` ma priemerne cas $O(1)$: pocet krokov sa pri pridani dalsich hracov velmi nemeni.

## Uloha

Vytvor subor `internal/game/structures_benchmark_test.go`. Pridaj do neho dve benchmark funkcie: jednu, ktora hlada posledneho z 1 000 hracov v `slice`, a druhu, ktora rovnakeho hraca hlada v `map` podla ID.

Spusti meranie:

```bash
go test ./internal/game -bench=. -benchmem
```

Vo vysledku porovnaj:

- `ns/op`: priemerny cas jednej operacie; mensie cislo je rychlejsie
- `B/op`: kolko bajtov pamate vznikne pri jednej operacii
- `allocs/op`: kolko novych alokacii potrebuje jedna operacia

Benchmark meria opakovane, aby nahodny kratky cas neovplyvnil vysledok. Cisla sa mozu lisit medzi pocitacmi; porovnavaj len dva riadky z jedneho spustenia.

## Otazky po merani

1. Ktore hladanie bolo rychlejsie?
2. Preco sa `slice` hodi, ked chceme hracov kreslit v poradi?
3. Preco sa `map` hodi, ked potrebujeme casto najst konkretneho hraca podla ID?
4. Ake data by si v hre ulozil do `slice` a ake do `map`?

## Mini vyzva

Skus meranie s 10, 100, 1 000 a 10 000 hracmi. Nakresli si graf poctu hracov a casu hladania. Pozoruj, pri ktorej strukture cas rastie najviac.
