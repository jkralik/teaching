# Lekcia 01: Co ideme postavit

## Ciel

Pochopit, aka hra je Tunneler a z coho sa sklada webova multiplayer hra.

## Co si vysvetlime

- rozdiel medzi serverom a klientom
- preco je Go vhodne na server
- co je HTTP a WebSocket
- ako bude vyzerat tank, mapa a prikaz hraca

## Kodovy krok

V `web/index.html` dopln nad hernu plochu nadpis `<h1>Tunneler</h1>` a pod neho kratky popis hry v `<p>`. Tym pridame prvy viditelny prvok nasej hry.

Novy pojem: HTML element, jeho otvaraci a zatvaraci tag.

## Overenie

Spusti `go run ./cmd/tunneler`, otvor stranku a skontroluj, ze sa nad canvasom zobrazi nadpis a popis.

## Uloha

Spusti hotovy zaklad projektu prikazom `go run ./cmd/tunneler` a otvor `http://localhost:8080`. Vyskusaj sa pripojit do mapy.

## Mini vyzva

Nakresli si na papier vlastnu mapu z troch znakov: `.`, `#`, `X`, a potom ju skus zapisat do noveho suboru v `maps/`.