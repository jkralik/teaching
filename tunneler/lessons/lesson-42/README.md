# Lekcia 42: Copilot a turnaj botov

## Ciel

Mat bota ako spoluhraca (copilota) a usporiadat turnaj, v ktorom proti sebe hraju boty roznych deti.

## Co si vysvetlime

- prepinac `-follow`: bot cez HTTP API (`/api/maps/{mapa}`) zisti, v akom time je hrac, a pripoji sa do neho
- prepinace `-name`, `-team` a `-server`
- rozhranie `Decider` v `client.go` - hocico, co ma funkciu `Decide`, moze riadit tank

## Kodovy krok

Pripoj sa v prehliadaci ako `Jano` do timu `Modri`. Potom spusti copilota:

```bash
go run ./cmd/bot -name Pomocnik -follow Jano
```

Pomocnik bude modry, bude sa drzat pri tebe (`followDistance`) a strielat na cervenych. Tvoje strely ho nezrania, lebo je tvoj spoluhrac.

Bot sa da pripojit aj na server kamarata:

```bash
go run ./cmd/bot -server http://192.168.1.20:8080 -map arena -name Terminator -team 1
```

Ked spadne spojenie, bot sa po 3 sekundach pripoji znova s rovnakym menom a pokracuje so svojim tankom.

## Overenie

Spusti `go test ./internal/bot/`. Test `TestCopilotJoinsOwnersTeam` overi, ze copilot sa dostane do timu svojho hraca.

## Uloha

Zlepsi copilota v `brain.go`:

- ked hrac striela, nech copilot striela tym istym smerom (`owner.Angle`)
- ked ma hrac malo zdravia, nech sa copilot postavi medzi hraca a nepriatela

## Turnaj botov

1. Kazde dieta si vylepsi svojho bota (lekcie 39-41) na svojom pocitaci.
2. Lektor spusti jeden server, vsetci sa pripoja cez `-server http://<ip-lektora>:8080 -team none` (kazdy hrac dostane vlastny tim).
3. Kazdy bot ma vlastne meno. Po 5 minutach vyhrava bot s najvacsim skore.
4. Potom turnaj timov: dvaja hraci a ich copiloti (`-follow`) proti dalsim dvom.

## Mini vyzva

Ktore pravidlo spravilo najvacsi rozdiel? Vymente si kod s kamaratom a skuste skombinovat najlepsie napady.
