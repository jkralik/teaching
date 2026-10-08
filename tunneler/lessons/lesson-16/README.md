# Lekcia 16: Kopanie tunelov

## Ciel

Spravit z tanku kopaca, ktory meni mapu.

## Co si vysvetlime

- mutovanie herneho stavu
- rozdiel medzi zemou a prazdnym miestom
- skore za odkopanu zem
- preco sa tank pri kopani najprv nepohne
- tvrdost zeme `dirtHardness` - kopanie trva niekolko pohybov

## Kodovy krok

V `moveTank` zmen `TileDirt` na `TileEmpty` a zvys `tank.Score`, ked tank kope tunel.

Tank nekope policko naraz. Kazdy pohyb prida do `tank.digProgress` tolko, aky je krok tanku. Ked progres dosiahne `dirtHardness`, zem zmizne. Ak sa tank otoci na ine policko, kope odznova.

Novy pojem: mutacia meni stav mapy; hodnota sa po akcii zachova aj v dalsom snimku.

## Overenie

Prejdi tankom do zeme a over, ze policko zmizne a skore narastie.

## Uloha

Zmen pocet bodov, ktore hrac dostane za odkopanie zeme.

Potom zmen `dirtHardness` (napriklad `0.2`, `1.2`, `3`) a porovnaj, ako rychlo tank kope tunel. Ktora hodnota je najzabavnejsia?

## Mini vyzva

Sprav pravidlo, ze niektora zem da viac bodov.
