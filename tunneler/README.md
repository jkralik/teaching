# Tunneler

Tunneler je kurzovy projekt pre deti, ktore sa ucia programovat v Go. Cielom je postupne vytvorit multiplayer hru, v ktorej sa tanky pohybuju po mape, kopu tunely, zbieraju body a strielaju po sebe.

Projekt ma priblizne 31 lekcii po 1 hodine. Kazda lekcia je samostatny priecinok v `lessons/lesson-XX` s README zadanim.

## Co projekt obsahuje

- Go HTTP server
- WebSocket komunikaciu medzi prehliadacom a serverom
- jednoduchy HTML, CSS a JavaScript klient
- mapy nacitane zo suboru alebo vygenerovane serverom
- zakladnu hernu logiku: hraci, pohyb, kopanie, strely a viac map naraz

## Spustenie

```bash
cd tunneler
go mod tidy
go run ./cmd/tunneler
```

Potom otvor v prehliadaci:

```text
http://localhost:8080
```

Port sa da zmenit cez premennu prostredia:

```bash
PORT=3000 go run ./cmd/tunneler
```

## Git a GitHub

Po kazdej dokoncenej malej zmene si uloz pracu do repozitara:

```bash
git status
git add cesta/k/suboru
git commit -m "Kratky popis zmeny"
git push
```

Pre novu ulohu si najprv vytvor vlastnu vetvu. Napriklad pre upravu vypisu portu:

```bash
git switch -c vypis-portu
git push -u origin vypis-portu
```

Potom vytvor pull request z vlastnej vetvy do `main`. Ked je zmena spojena, vrat sa do `main` a stiahni nove zmeny:

```bash
git switch main
git pull
```

Podrobny postup a vysvetlenie prikazov je v [druhej lekcii](lessons/lesson-02/README.md).

## Zakladne pravidla hry

- Hraci si vyberu mapu v prehliadaci.
- Kazdy hrac je tank na spolocnej mape.
- Sipkami alebo WASD sa tank pohybuje.
- Klavesom medzernik tank vystreli.
- Zem sa da odkopat pohybom do blokov `#`.
- Kamene `X` su pevna prekazka.
- Server je pravda: klient len posiela prikazy a kresli stav.

## Struktura

```text
cmd/tunneler/       hlavny program servera
internal/game/      herna logika a WebSocket hub
web/                HTML, CSS a JavaScript klient
maps/               textove mapy
lessons/            lekcie kurzu
```

## Format mapy

Mapa je textovy subor. Kazdy znak je jedno policko:

- `.` prazdne miesto
- `#` zem, ktoru tank odkope
- `X` pevny kamen

Priklad:

```text
XXXXXXXXXXXXXXXX
X....####......X
X..######..##..X
X..............X
XXXXXXXXXXXXXXXX
```

## Ako ucit kurz

Lekcie su navrhnute tak, aby deti najprv pochopili male casti programu: premenne, funkcie, HTTP server, JSON, herny stav. Neskor sa postupne pridaju WebSockety, viac hracov, mapy, strely, kolizie a jednoduche vylepsenia hry.

Odporucany rytmus hodiny:

1. 10 minut: vysvetlenie novej myslienky
2. 35 minut: spolocne programovanie
3. 10 minut: vlastna uprava alebo mini vyzva
4. 5 minut: spustenie hry a zhrnutie
