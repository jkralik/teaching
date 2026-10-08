# Lekcia 04: HTTP server

## Ciel

Pochopit, ako Go odpoveda prehliadacu.

## Co si vysvetlime

- `net/http`
- cesta URL
- handler funkcia
- JSON odpoved

## Kodovy krok

V `cmd/tunneler/main.go` pridaj handler `helloHandler`, ktory odpovie na `GET /api/hello` JSON objektom `{"message":"Ahoj z Tunneleru"}`. Handler zaregistruj do `mux`.

Novy pojem: HTTP handler prijme poziadavku a vytvori odpoved. JSON je format, ktory vie precitat aj JavaScript.

## Overenie

Spusti server a otvor `http://localhost:8080/api/hello`. V prehliadaci alebo cez `curl` musis vidiet JSON odpoved.

## Uloha

Pridaj novu cestu `/api/hello`, ktora vrati kratku JSON spravu.

## Mini vyzva

Pridaj do JSON odpovede aj pole `lesson` s cislom lekcie.