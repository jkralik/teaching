# Lekcia 06: CSS a herna obrazovka

## Ciel

Upravit vzhlad hry bez zasahu do servera.

## Co si vysvetlime

- selektory v CSS
- farby cez premenne
- responzivne rozlozenie
- preco ma canvas pevny pomer stran

## Kodovy krok

V `web/style.css` pridaj triedu `.connected`, ktora zmeni farbu statusu na zelenu. V `web/app.js` tuto triedu nastav elementu `server-message` po otvoreni WebSocket spojenia.

Novy pojem: CSS trieda opisuje vzhlad a JavaScript ju moze zapnut alebo vypnut cez `classList`.

## Overenie

Po pripojeni do hry sa status zmeni na zeleny. Pri zatvorenom spojeni nesmie zelena predstierat, ze server funguje.

## Uloha

Zmen farby tanku, zeme alebo pozadia v `web/style.css` a `web/app.js`.

## Mini vyzva

Vymysli farebnu temu pre svoj tim a pridaj triedu `.warning` pre stav chyby.