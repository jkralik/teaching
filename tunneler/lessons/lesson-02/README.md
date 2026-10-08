# Lekcia 02: Prvy Go program

## Ciel

Vediet spustit Go program a rozumiet funkcii `main`.

## Co si vysvetlime

- balicek `main`
- funkcia `main()`
- importy
- prikaz `go run`

## Kodovy krok

V `cmd/tunneler/main.go` vytvor funkciu `startupMessage(port string) string`, ktora vrati text `Tunneler server bezi na porte <port>`. V `main` jej vysledok vypis cez `log.Println`.

Novy pojem: funkcia s parametrom a navratovou hodnotou. Funkcia oddeli vypocet textu od spustenia servera.

## Overenie

Spusti `go run ./cmd/tunneler` a skontroluj, ze terminal vypise spravu s portom `8080`. Potom skus `PORT=3000 go run ./cmd/tunneler`.

## Uloha

V subore `cmd/tunneler/main.go` zmen text, ktory server vypise pri starte, a program znovu spusti.

## Git a GitHub

Po kazdej hotovej malej zmene ju uloz do GitHubu. V terminali otvorenom v priecinku `tunneler` pouzi tieto prikazy:

```bash
git status
git add cmd/tunneler/main.go
git commit -m "Zmenim uvodny vypis servera"
git push
```

- `git status` ukaze, ktore subory si zmenil.
- `git add` vyberie subor, ktory chceme ulozit.
- `git commit` vytvori ulozeny bod s kratkym popisom zmeny.
- `git push` posle tvoje commity na GitHub.

Ak menis viac suborov, napis ich za `git add`, napr.:

```bash
git add cmd/tunneler/main.go web/app.js
```

Na konci otvor svoj repozitar na GitHube a skontroluj, ze je tam novy commit.

## Praca s vetvami

Vetva je vlastna pracovna verzia projektu. Pre kazdu novu ulohu si vytvor vlastnu vetvu. Napriklad, ked pridavas vypis portu:

```bash
git switch -c vypis-portu
```

Teraz rob zmeny, vytvor commit a posli novu vetvu na GitHub:

```bash
git add cmd/tunneler/main.go
git commit -m "Pridam vypis portu"
git push -u origin vypis-portu
```

- `git switch -c vypis-portu` vytvori vetvu a hned do nej prepne.
- `git branch` ukaze zoznam vetiev; hviezdicka oznacuje tu, v ktorej prave pracujes.
- `git push -u origin vypis-portu` vytvori vetvu aj na GitHube a posle do nej commity.

Po odoslani otvor GitHub, vytvor pull request z vetvy `vypis-portu` do `main` a poziadaj ucitela o kontrolu. Ked je zmena spojena do `main`, vrat sa do nej a stiahni si najnovsi kod:

```bash
git switch main
git pull
```

## Mini vyzva

Pridaj do spravy aj adresu, ktoru mas otvorit v prehliadaci.
