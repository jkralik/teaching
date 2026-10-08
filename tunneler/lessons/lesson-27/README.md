# Lekcia 27: Refaktoring

## Ciel

Zlepsit kod bez zmeny spravania hry.

## Co si vysvetlime

- male funkcie
- pomenovanie
- presun opakujuceho sa kodu
- preco refaktoring robime po testoch

## Kodovy krok

V `web/app.js` vyber jednu cast kreslenia a presun ju do malej funkcie s jasnym nazvom, napriklad `drawBullets`.

Novy pojem: refaktoring meni strukturu kodu bez zmeny jeho pozorovatelneho spravania.

## Overenie

Spusti hru pred zmenou a po nej. Kreslenie a ovladanie musia fungovat rovnako.

## Uloha

Najdi cast kodu, ktora by mohla mat lepsi nazov alebo samostatnu funkciu.

## Mini vyzva

Rozdel kreslenie v JavaScripte na funkcie `drawMap`, `drawTanks`, `drawBullets`.
