# Lekcia 13: Herny stav

## Ciel

Pochopit, preco ma server uchovavat pravdivy stav hry.

## Co si vysvetlime

- mapa, tanky a strely ako stav
- preco klient neposiela svoju polohu priamo
- co znamena snapshot
- kopirovanie dat pre odoslanie

## Kodovy krok

Do `Snapshot` v `internal/game/types.go` pridaj pole `PlayerCount int` a napln ho v `snapshotLocked`. V klientovi zobraz pocet hracov.

Novy pojem: snapshot je kopia stavu urcena na odoslanie; klient nema menit stav servera priamo.

## Overenie

Otvor hru v dvoch oknach a over, ze obe ukazuju pocet pripojenych hracov.

## Uloha

Preskumaj strukturu `Snapshot` a najdi, kde sa vytvara.

## Mini vyzva

Pridaj do snapshotu pocet hracov na mape.
