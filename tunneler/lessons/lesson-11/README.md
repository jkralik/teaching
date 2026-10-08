# Lekcia 11: JSON

## Ciel

Pochopit, ako si Go a JavaScript posielaju data.

## Co si vysvetlime

- struktury v Go
- tagy `json`
- `JSON.stringify` a `JSON.parse`
- rozdiel medzi prikazom a stavom hry

## Kodovy krok

Do `Tank` v `internal/game/types.go` pridaj pole `Score int` s JSON tagom. V `web/app.js` zobraz skore vlastneho hraca vedla canvasu.

Novy pojem: JSON tag urcuje nazov vlastnosti, ktoru dostane JavaScript.

## Overenie

Odkop zem a skontroluj, ze sa skore v klientovi zvysi.

## Uloha

Pridaj do tanku novu vlastnost `Score` alebo zmen jej zobrazenie v klientovi.

## Mini vyzva

Zobraz skore hracov vedla hernej plochy.
