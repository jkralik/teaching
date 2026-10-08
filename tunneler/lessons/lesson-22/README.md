# Lekcia 22: Zasah tanku

## Ciel

Vyhodnotit, ci strela trafila ineho hraca.

## Co si vysvetlime

- porovnanie suradnic
- preskocenie vlastnika strely
- zmena `Alive`
- body za zasah

## Kodovy krok

V `moveBullet` pridaj kontrolu, ktora pri rovnakej pozicii strely a supera nastavi `Alive` na `false` a vlastnikovi prida body.

Novy pojem: herne pravidlo je podmienka, ktora meni viac casti stavu naraz.

## Overenie

Otvor dve okna, vystrel na druheho hraca a over, ze zasah zmeni jeho stav.

## Uloha

Uprav pocet bodov za trafenie supera.

## Mini vyzva

Pridaj zivoty namiesto okamziteho vyradenia.
