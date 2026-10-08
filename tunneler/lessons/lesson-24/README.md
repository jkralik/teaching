# Lekcia 24: Vyber mapy

## Ciel

Umoznit hracom rozhodnut, do akej mapy sa pripoja.

## Co si vysvetlime

- endpoint `/api/maps`
- `select` v HTML
- query parameter `map`
- samostatny hub pre kazdu mapu

## Kodovy krok

V `web/index.html` pridaj `select` pre mapu. V `web/app.js` nacitaj `/api/maps` a vybranu mapu posli v query parametri WebSocketu.

Novy pojem: query parameter je cast URL, ktorou klient vyberie data alebo spravanie servera.

## Overenie

Vyber dve rozne mapy v dvoch oknach a over, ze hraci sa navzajom nevidia.

## Uloha

Pridaj novu mapu a over, ze hraci v roznych mapach sa nevidia.

## Mini vyzva

Zobraz aktualny nazov mapy nad canvasom.
