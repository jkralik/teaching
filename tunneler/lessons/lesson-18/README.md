# Lekcia 18: Viac hracov

## Ciel

Pochopit, ako server drzi zoznam pripojenych hracov.

## Co si vysvetlime

- mapa `map[string]Tank`
- ID hraca
- pripojenie a odpojenie
- broadcast spravy vsetkym klientom

## Kodovy krok

V `internal/game/hub.go` nechaj `Join` pridat novy `Tank` do mapy `tanks` a po pripojeni posli snapshot vsetkym klientom.

Novy pojem: server udrzuje zdielany stav hracov a broadcast posle jednu spravu viacerym klientom.

## Overenie

Otvor hru v dvoch oknach. Po pripojeni druheho hraca musia oba klienty vidiet oba tanky.

## Uloha

Otvor hru v dvoch oknach prehliadaca a sleduj, ako sa tanky navzajom vidia.

## Mini vyzva

Zmen startovaciu poziciu novych hracov.
