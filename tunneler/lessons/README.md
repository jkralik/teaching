# Lekcie

Kazda lekcia je planovana na priblizne 1 hodinu. README v priecinku lekcie opisuje ciel, vysvetlenie, programatorsku ulohu a mini vyzvu.

Odporucanie pre lektora: nechaj deti casto spustat program. Hra je motivacia, ale hlavny ciel je pochopit male kroky v Go.

## Ako lekcie postupujeme

V kazdej lekcii urobime jeden maly krok v rovnakom projekte:

1. najprv program spustime a pozrieme sa, co uz funguje,
2. dopiseme kratky kus kodu, ktory prida viditelnu funkcionalitu,
3. pomenujeme novy Go alebo webovy pojem,
4. overime vysledok spustenim hry, prikazom alebo testom.

Ak je v zadani napisane iba "najdi", znamena to pripravu pred kodovym krokom, nie koniec lekcie.

## Lekcia 08

- [Datove struktury a rychlost](lesson-08/README.md) - porovnanie `slice` a `map`, algoritmy hladania a meranie pomocou benchmarkov

## Zbrojna dielna (bonus)

Po turnaji v lekcii 31 si deti mozu hrat so zbranami. Vsetko sa deje hlavne v `internal/game/weapon.go`.

- [Lekcia 32: Laboratorium zbrani](lesson-32/README.md) - zmena kadencie, zasobnika, rychlosti strely, vybuchu, poskodenia a zivotov
- [Lekcia 33: Druha zbran](lesson-33/README.md) - nova zbran v `weaponCatalog()` a prepinanie klavesom `Q`
- [Lekcia 34: Nova vlastnost zbrane](lesson-34/README.md) - nove pole v `WeaponSettings`, napriklad dostrel
- [Lekcia 35: Test a balans zbrani](lesson-35/README.md) - testy pre kazdu zbran a vyvazenie zbrani
- [Lekcia 36: Pancier tanku](lesson-36/README.md) - farebne pancieri v `armorCatalog()`, delenie poskodenia a spomalenie tanku
- [Lekcia 37: Bonusy na mape](lesson-37/README.md) - `bonusCatalog()`, okamzity reload, Turbo, Silna strela a Nesmrtelnost
- [Lekcia 38: Prekvapenie ?](lesson-38/README.md) - nahodne policko `?`, negativne efekty a `mysteryChance`

## AI tank (bonus)

AI tank je samostatny program v `cmd/bot`, ktory sa k serveru pripaja cez WebSocket ako hrac. Deti upravuju hlavne jeho "mozog" v `internal/bot/brain.go`.

- [Lekcia 39: Prvy bot](lesson-39/README.md) - spustenie bota, funkcia `Decide` a ladenie konstant
- [Lekcia 40: Rozhodovanie bota](lesson-40/README.md) - ustup pri malom zdravi, zbieranie bonusov, prepinanie zbrani
- [Lekcia 41: Hladanie cesty](lesson-41/README.md) - BFS po mriezke mapy okolo kamenov
- [Lekcia 42: Copilot a turnaj botov](lesson-42/README.md) - bot v time hraca (`-follow`) a turnaj botov
- [Lekcia 43: Vyber zbrane cislami](lesson-43/README.md) - klavesy `1`-`9` v prehliadaci a sprava `weapon` s indexom
