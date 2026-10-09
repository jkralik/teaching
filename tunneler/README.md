# Tunneler

Tunneler je kurzovy projekt pre deti, ktore sa ucia programovat v Go. Cielom je postupne vytvorit multiplayer hru, v ktorej sa tanky pohybuju po mape, kopu tunely, zbieraju body a strielaju po sebe.

Projekt ma priblizne 31 lekcii po 1 hodine a bonusove bloky "Zbrojna dielna" (lekcie 32-38) a "AI tank" (lekcie 39-42). Kazda lekcia je samostatny priecinok v `lessons/lesson-XX` s README zadanim.

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
http://localhost:58080
```

Port sa da zmenit cez premennu prostredia:

```bash
PORT=3000 go run ./cmd/tunneler
```

### AI tank (bot)

Bot je samostatny program, ktory sa k serveru pripoji cez WebSocket ako bezny hrac (rovnake API ako prehliadac, ziadne specialne prava). V druhom terminali:

```bash
go run ./cmd/bot -map arena                          # bot proti vsetkym
go run ./cmd/bot -name Pomocnik -follow Jano         # copilot v time hraca Jano
go run ./cmd/bot -server http://192.168.1.20:58080 -name Robot2 -team none
```

Spravanie bota (mozog) je v `internal/bot/brain.go`, pozri lekcie 39-42.

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
- Tanky do seba nenarazaju prejazdom; aj odpojeny tank zostava na mape a blokuje cestu.
- Rychlost pohybu strazi server: tank spravi najviac jeden krok za 30 ms, aj keby klient posielal spravy `move` castejsie.
- Pri vypadku spojenia sa tank oznaci ako offline a nemoze konat ani strielat. Pripojenie s rovnakym menom do 5 minut obnovi jeho stav; potom sa tank odstrani. Tlacidlo „Odpojit sa“ odstrani tank aj jeho herny stav hned.
- Sipkami alebo WASD sa tank pohybuje aj sikmo; pri drzeni klavesu sa pohybuje plynulo.
- Na mobile sa tank ovlada dotykovymi sipkami; podrzanim viac sipiek naraz sa pohybuje sikmo. Tlacidla pod mapou sluzia na strelbu, zmenu zbrane a panciera.
- Klavesom medzernik tank vystreli.
- Strela leti smerom, ktorym je tank otoceny.
- Kazdy tank ma zivoty (`tankMaxHealth`, predvolene 100), ktore ukazuje farebny ukazovatel nad tankom. Zasah da strelcovi 1 bod, zniceny tank 5 bodov.
- Tank moze mat pancier zo zoznamu `armorCatalog()` v `internal/game/armor.go` (Bez panciera, Bronzovy, Strieborny, Zlaty). Klavesom `E` sa prepina. Pancier je vidiet ako farebny obrys tanku (silnejsi = hrubsi). Poskodenie zbrane sa vydeli cislom `Divisor` (0 alebo 1 = plne poskodenie, zasah vzdy zoberie aspon 1 zivot) a `SpeedFactor` spomali tank.
- Nastavenia zbrane (kadencia, zasobnik, prebijanie, rychlost strely, polomer vybuchu a poskodenie `Damage`) su pokope vo funkcii `weaponSettings()` v `internal/game/weapon.go`, kde ich mozete upravit ako cvicenie.
- Zbrane su v zozname `weaponCatalog()`. Klavesom `Q` sa prepne na dalsiu zbran. Vyber konkretnej zbrane klavesmi `1`-`9` si deti dorobia v lekcii 43. Panel Zbran ukazuje nazov, poskodenie, naboje a prebijanie; kazda zbran ma vlastny zasobnik a farbu strely.
- Panel zbrane zobrazuje pocet nabojov a pocas prebijania aj zostavajuci cas s ukazovatelom priebehu.
- Zem sa da odkopat pohybom do blokov `#`. Kopanie chvilu trva, takze tank je v zemi pomalsi; tvrdost zeme nastavuje `dirtHardness` v `internal/game/hub.go`.
- Kamene `X` su pevna prekazka.
- Na volnych polickach sa nahodne objavuju bonusy zo zoznamu `bonusCatalog()` v `internal/game/bonus.go` (napr. `R` = okamzity reload). Tank ich zbiera prejdenim; ten isty docasny bonus moze mat len raz naraz. Fialove policko `?` da nahodny efekt, aj negativny (napr. Blato spomali tank). Aktivne bonusy ukazuje panel Zbran.
- Hrac si pri pripojeni vyberie tim (`Automaticky` = tim s najmenej hracmi; `Bez timu` vytvori hracovi samostatny tim pomenovany podla neho s nahodnou farbou, ktoru nastavuju `soloColorSaturation` a `soloColorLightness`). Vyberatelne timy su v zozname `teamCatalog()` v `internal/game/team.go` (meno a farba tanku), prehliadac ich nacita z `/api/teams`. Strela spoluhraca nezrani a preleti cez neho; prepinac je konstanta `friendlyFire`. Skore timu je sucet bodov jeho hracov a ukazuje ho panel Timy. Po znovupripojeni hrac zostava vo svojom povodnom time.
- Server je pravda: klient len posiela prikazy a kresli stav.

## Struktura

```text
cmd/tunneler/       hlavny program servera
cmd/bot/            AI tank - klient, ktory hra cez WebSocket
internal/bot/       mozog bota (brain.go) a pripojenie k serveru
internal/game/      herna logika a WebSocket hub
web/                HTML, CSS a JavaScript klient
maps/               textove mapy
lessons/            lekcie kurzu
```

## Format mapy

Mapa je textovy subor. Kazdy znak sa pri nacitani rozsiri na 2 x 2 mensie herne policka:

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

## Navrhy

### Misko

- viacej zivotov
- noc (tank vidi len nejaky okruh)
- lepsi stit (nesmrtenlnost na kratky cas 10-20s)

### Janko

- skraslit tanky trosku 3D aby boli