# Lekcia 19: Gorutiny a kanaly

## Ciel

Rozumiet zakladu sucasneho behu v Go.

## Co si vysvetlime

- `go func()`
- kanal `chan`
- citanie a zapis sprav
- preco WebSocket potrebuje citanie aj zapis

## Kodovy krok

V `internal/game/hub.go` spusti herny loop cez `go hub.loop()` a pomocou kanala `stop` mu umozni skoncit.

Novy pojem: gorutina bezi sucasne a kanal je bezpecny sposob odovzdania signalu medzi gorutinami.

## Overenie

Spusti testy a sleduj, ze pripojenie klienta neblokuje server pri obsluhe dalsieho klienta.

## Uloha

Najdi kanal, cez ktory hub posiela spravy klientovi.

## Mini vyzva

Pridaj spravu pri pripojeni noveho hraca.
