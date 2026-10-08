# Lekcia 03: Premenne a konstanty

## Ciel

Pouzit premennu a konstantu v skutocnom programe.

## Co si vysvetlime

- `var` a kratky zapis `:=`
- konstanty `const`
- textove hodnoty
- preco nechceme magicke cisla

## Kodovy krok

V `cmd/tunneler/main.go` vytvor konstantu `defaultPort = "8080"` a pouzi ju namiesto textu `"8080"`, ked nastavujes predvoleny port.

Novy pojem: konstanta pomenuva hodnotu, ktora sa pocas behu programu nemeni. Jedno miesto je vdaka nej zdrojom pravdy.

## Overenie

Spusti hru bez premennej `PORT` a potom s `PORT=3000`. V oboch pripadoch skontroluj uvodny vypis.

## Uloha

Najdi v kode velkost generovanej mapy a uloz ju do pomenovanych konstant.

## Mini vyzva

Vytvor konstanty `generatedMapWidth` a `generatedMapHeight` a zmen rozmery generovanej mapy.