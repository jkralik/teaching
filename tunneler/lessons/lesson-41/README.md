# Lekcia 41: Hladanie cesty

## Ciel

Bot nepojde rovno do steny, ale najde cestu okolo kamenov pomocou algoritmu **BFS**.

## Co si vysvetlime

- mapa ako mriezka policok: `state.TileAt(x, y)` vrati `.`, `#` alebo `X`
- **BFS** (prehladavanie do sirky) - rovnaky napad ako v [lekcii 08](../lesson-08/README.md): fronta (`slice`) a navstivene policka (`map`)
- preco BFS najde najkratsiu cestu: najprv skusi vsetky policka vzdialene 1, potom 2, potom 3...

## Kodovy krok

Vytvor subor `internal/bot/path.go`:

```go
package bot

type point struct{ X, Y int }

// nextStep vrati policko, na ktore mam ist, aby som sa dostal k cielu.
func nextStep(state State, from point, to point) (point, bool) {
	previous := map[point]point{from: from}
	queue := []point{from}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		if current == to {
			// ideme od ciela spat, kym nenajdeme prve policko za startom
			for previous[current] != from {
				current = previous[current]
			}
			return current, true
		}
		for _, d := range []point{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
			next := point{current.X + d.X, current.Y + d.Y}
			if _, seen := previous[next]; seen || state.TileAt(next.X, next.Y) == TileRock {
				continue
			}
			previous[next] = current
			queue = append(queue, next)
		}
	}
	return point{}, false
}
```

V `Decide` namiesto `angleTo(me.X, me.Y, enemy.X, enemy.Y)` pouzi:

```go
from := point{int(me.X), int(me.Y)}
to := point{int(enemy.X), int(enemy.Y)}
if step, ok := nextStep(state, from, to); ok {
	action.Angle = angleTo(me.X, me.Y, float64(step.X)+0.5, float64(step.Y)+0.5)
}
```

`+0.5` je stred policka - tanky stoja v stredoch policok.

## Overenie

Napis test: mapa so stenou z kamenov medzi botom a nepriatelom, v stene je jedna diera. Prvy krok musi viest k diere, nie rovno do steny.

## Uloha

Zem `#` sa da prekopat, ale je pomala. Uprav BFS tak, aby zem bola "drahsia" - napriklad ju preskoc, ak existuje cesta len cez prazdne `.` policka (spusti BFS dvakrat: najprv bez zeme, potom so zemou).

## Mini vyzva

Pouzi `nextStep` aj na cestu k bonusu z lekcie 40 a pre copilota k jeho hracovi.
