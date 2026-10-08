# Lekcia 17: Klavesnica v JavaScripte

## Ciel

Posielat prikazy zo sipok a WASD na server.

## Co si vysvetlime

- udalost `keydown`
- objekt s mapovanim klavesov
- `event.preventDefault()`
- odoslanie JSON spravy

## Kodovy krok

V `web/app.js` dopln do mapy klaves aj WASD a pri odoslani prikazu zavolaj `event.preventDefault()`.

Novy pojem: udalost `keydown` opisuje stlacenie klavesu a objekt s mapovanim oddeli vstup od herneho prikazu.

## Overenie

Pohybuj tankom sipkami aj klavesmi WASD. Stranka sa pri pohybe nesmie posuvat.

## Uloha

Pridaj dalsie klavesy pre pohyb alebo zmen ovladanie.

## Mini vyzva

Pridaj klaves `r`, ktory neskor pouzijeme na respawn.
