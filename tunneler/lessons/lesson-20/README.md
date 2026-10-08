# Lekcia 20: Mutex a bezpecny stav

## Ciel

Pochopit, preco viac gorutin nemoze menit stav naraz bez pravidiel.

## Co si vysvetlime

- `sync.Mutex`
- `Lock` a `Unlock`
- `defer`
- kriticka sekcia

## Kodovy krok

V kazdej metode, ktora cita alebo meni stav hubu, pouzi `hub.mu.Lock()` a `defer hub.mu.Unlock()`.

Novy pojem: mutex chrani kriticku sekciu, aby dve gorutiny nemenili mapu hracov naraz.

## Overenie

Spusti `go test -race ./...` a skontroluj, ze testy nehlasia datove preteky.

## Uloha

Najdi miesta, kde sa pouziva `hub.mu.Lock()`.

## Mini vyzva

Vysvetli vlastnymi slovami, co by sa mohlo stat bez mutexu.
