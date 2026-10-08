package game

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Map struct {
	Name   string
	Width  int
	Height int
	Tiles  [][]Tile
}

func LoadMap(path string) (*Map, error) {
	// Lekcia 09: Tento scanner je pripraveny na nacitanie vlastnej mapy zo suboru.
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var rows [][]Tile
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Lekcia 25: Pokaz mapovy subor a precitaj chybu, ktoru vrati validacia.
		line := strings.TrimRight(scanner.Text(), "\r")
		if line == "" {
			continue
		}
		row := make([]Tile, 0, len(line)*2)
		for index := range line {
			row = append(row, Tile(line[index]), Tile(line[index]))
		}
		rows = append(rows, append([]Tile(nil), row...), row)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("mapa %s je prazdna", path)
	}

	width := len(rows[0])
	for _, row := range rows {
		if len(row) != width {
			return nil, fmt.Errorf("mapa %s nema v kazdom riadku rovnaku sirku", path)
		}
	}

	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	return &Map{Name: name, Width: width, Height: len(rows), Tiles: rows}, nil
}

func LoadMaps(dir string) (map[string]*Map, error) {
	paths, err := filepath.Glob(filepath.Join(dir, "*.txt"))
	if err != nil {
		return nil, err
	}
	maps := make(map[string]*Map)
	for _, path := range paths {
		loadedMap, err := LoadMap(path)
		if err != nil {
			return nil, err
		}
		maps[loadedMap.Name] = loadedMap
	}
	return maps, nil
}

func MapNames(maps map[string]*Map) []string {
	names := make([]string, 0, len(maps)+1)
	for name := range maps {
		names = append(names, name)
	}
	names = append(names, "generated")
	sort.Strings(names)
	return names
}

func GenerateMap(name string, width int, height int) *Map {
	// Lekcia 10: Zmen tieto pravdepodobnosti a sleduj, ako sa meni vygenerovana mapa.
	tiles := make([][]Tile, height)
	for y := range height {
		tiles[y] = make([]Tile, width)
		for x := range width {
			switch {
			case x == 0 || y == 0 || x == width-1 || y == height-1:
				tiles[y][x] = TileRock
			case rand.Intn(100) < 12:
				tiles[y][x] = TileRock
			case rand.Intn(100) < 55:
				tiles[y][x] = TileDirt
			default:
				tiles[y][x] = TileEmpty
			}
		}
	}
	return &Map{Name: name, Width: width, Height: height, Tiles: tiles}
}

func (gameMap *Map) Clone() *Map {
	tiles := make([][]Tile, gameMap.Height)
	for y := range gameMap.Height {
		tiles[y] = append([]Tile(nil), gameMap.Tiles[y]...)
	}
	return &Map{Name: gameMap.Name, Width: gameMap.Width, Height: gameMap.Height, Tiles: tiles}
}

func (gameMap *Map) TileAt(x int, y int) Tile {
	// Lekcia 07: Dopln metodu Inside(x, y int) bool a pouzi ju pre kontrolu hranic mapy.
	if x < 0 || y < 0 || x >= gameMap.Width || y >= gameMap.Height {
		return TileRock
	}
	return gameMap.Tiles[y][x]
}

func (gameMap *Map) SetTile(x int, y int, tile Tile) {
	if x < 0 || y < 0 || x >= gameMap.Width || y >= gameMap.Height {
		return
	}
	gameMap.Tiles[y][x] = tile
}

func (gameMap *Map) Rows() []string {
	rows := make([]string, gameMap.Height)
	for y := range gameMap.Height {
		builder := strings.Builder{}
		builder.Grow(gameMap.Width)
		for x := range gameMap.Width {
			builder.WriteByte(byte(gameMap.Tiles[y][x]))
		}
		rows[y] = builder.String()
	}
	return rows
}
