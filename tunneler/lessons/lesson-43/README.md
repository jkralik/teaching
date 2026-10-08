# Lekcia 43: Vyber zbrane cislami

## Ciel

Dorobit v prehliadaci klavesy `1`-`9`, ktore hned vyberu konkretnu zbran, a napisat ich do ovladania na stranke.

## Co si vysvetlime

- udalost `keydown` a vlastnost `event.key` (napr. `"3"`)
- regularny vyraz `/^[1-9]$/` - pasuje iba na jednu cislicu 1 az 9
- prevod textu na cislo funkciou `Number()`
- preco klaves `1` posiela index `0` (prva zbran v rezi `weaponCatalog()` ma index `0`)

## Kodovy krok

Server uz spravu so zbranou pozna - posiela mu ju aj bot (lekcia 40). Treba ju poslat aj z prehliadaca.

V `web/app.js` najdi v obsluhe `keydown` komentar `Lekcia 43` a nahrad ho tymto kodom:

```js
if (/^[1-9]$/.test(event.key)) {
  event.preventDefault();
  send({ type: "weapon", weapon: Number(event.key) - 1 });
}
```

Potom v `web/index.html` v casti `Ovladanie` pridaj za riadok s klavesom `Q`:

```html
<dt><kbd>1</kbd> &ndash; <kbd>9</kbd></dt>
<dd>Vybrat konkretnu zbran</dd>
```

## Over to

Obnov stranku, pripoj sa a stlac `2`. V paneli Zbran sa ma zmenit nazov zbrane. Stlac `1` a mas zase kanon.

Co sa stane, ked stlacis `9`, ale v `weaponCatalog()` je menej zbrani? Najdi v `internal/game/hub.go` riadok `case "weapon":` a zisti, preco hra nespadne.

## Vyzva

- Zobraz v paneli Zbran pri nazve aj cislo klavesu, napr. `[2] Gulomet`.
- Vymysli, ako by sa dala vybrat predchadzajuca zbran (napr. klaves `Z`).
