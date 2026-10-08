# Lekcia 25: Chyby a logovanie

## Ciel

Nebat sa chyb a vediet ich citat.

## Co si vysvetlime

- `error`
- `if err != nil`
- `log.Printf`
- rozdiel medzi chybou pre programatora a hraca

## Kodovy krok

V `internal/game/map.go` osetri chybu z `os.Open`, `scanner.Err` a neplatnej sirky riadku. V `main` ju zapis cez `log.Printf`.

Novy pojem: `error` je navratova hodnota, ktoru program musi skontrolovat a zrozumitelne zaznamenat.

## Overenie

Docasne pokaz mapovy subor, spusti server a precitaj chybovu spravu v terminali.

## Uloha

Skus pokazit subor mapy tak, aby mal riadky roznej dlzky, a precitaj chybovu spravu.

## Mini vyzva

Sprav chybovu spravu zrozumitelnejsiu pre deti v triede.
