# Lekcia 12: WebSocket spojenie

## Ciel

Rozumiet, preco hra potrebuje trvale spojenie.

## Co si vysvetlime

- rozdiel HTTP request a WebSocket
- otvorenie spojenia v JavaScripte
- upgrade na serveri
- spravy zo servera ku klientovi

## Kodovy krok

V `web/app.js` zmen text statusu na `Spojene so serverom` po otvoreni WebSocketu a na `Spojenie zatvorene` pri jeho zatvoreni.

Novy pojem: WebSocket je trvale obojsmerne spojenie, cez ktore server moze posielat stav aj bez noveho HTTP requestu.

## Overenie

Otvor hru, pozoruj status a potom zastav server. Status sa musi zmenit po zavreti spojenia.

## Uloha

Najdi `HandleWebSocket` a popis, co sa stane pri pripojeni hraca.

## Mini vyzva

Zmen text statusu pri otvoreni a zatvoreni spojenia.
