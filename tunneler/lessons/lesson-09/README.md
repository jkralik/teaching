# Lekcia 09: Nacitanie mapy zo suboru

## Ciel

Vediet precitat textovy subor s mapou.

## Co si vysvetlime

- `os.Open`
- `bufio.Scanner`
- osetrenie chyby
- preco musia mat riadky rovnaku dlzku

## Kodovy krok

V `maps/` vytvor subor `training.txt` s vlastnou mapou. V `internal/game/map.go` dopln kontrolu, ktora odmietne riadok s neplatnym znakom.

Novy pojem: citanie suboru prebieha postupne cez scanner a chyba sa vracia volajucemu.

## Overenie

Spusti `go test ./...`, potom server a vyber novu mapu v prehliadaci.

## Uloha

Vytvor novu mapu v priecinku `maps/` a over, ze sa zobrazi vo vybere map.

## Mini vyzva

Urob mapu, ktora ma v strede pevnu prekazku z `X`.
