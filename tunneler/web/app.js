const canvas = document.querySelector("#game");
const context = canvas.getContext("2d");
const mapSelect = document.querySelector("#mapSelect");
const playerName = document.querySelector("#playerName");
const teamSelect = document.querySelector("#teamSelect");
const teamHud = document.querySelector("#teamHud");
const teamList = document.querySelector("#teamList");
const joinButton = document.querySelector("#joinButton");
const joinForm = document.querySelector("#joinForm");
const playerInfo = document.querySelector("#playerInfo");
const playerSummary = document.querySelector("#playerSummary");
const leaveButton = document.querySelector("#leaveButton");
const gameArea = document.querySelector("#gameArea");
const shootButton = document.querySelector("#shootButton");
const weaponButton = document.querySelector("#weaponButton");
const armorButton = document.querySelector("#armorButton");
const statusText = document.querySelector("#status");
const weaponHud = document.querySelector("#weaponHud");
const ammoCount = document.querySelector("#ammoCount");
const weaponName = document.querySelector("#weaponName");
const weaponDamage = document.querySelector("#weaponDamage");
const armorInfo = document.querySelector("#armorInfo");
const bonusInfo = document.querySelector("#bonusInfo");
const reloadStatus = document.querySelector("#reloadStatus");
const reloadTrack = document.querySelector("#reloadTrack");
const reloadProgress = document.querySelector("#reloadProgress");

// Lekcia 05: Tento element je pripraveny na doplnenie vlastneho statusu alebo mena timu.

let socket = null;
let state = null;
let playerId = null;
let reloadDurationMs = 0;
let lastWeaponIndex = null;
const movementInputs = new Map();
let movementTimer = null;

playerName.disabled = false;
mapSelect.disabled = false;

// Farby policok mapy, pouziva ich drawTile.
const tileColors = {
  ".": "#17130f",
  "#": "#8b6339",
  "X": "#50493d",
};

async function loadMaps() {
  // Lekcia 09: API uz vracia nazvy map. Pridaj vlastny subor do maps/ a over ho vo vybere.
  const response = await fetch("/api/maps");
  const maps = await response.json();
  mapSelect.innerHTML = "";
  for (const map of maps) {
    const option = document.createElement("option");
    option.value = map;
    option.textContent = map;
    mapSelect.append(option);
  }
}

// Timy: server posiela zoznam timov, vlozime ich medzi "Automaticky" a "Bez timu".
async function loadTeams() {
  const response = await fetch("/api/teams");
  const teams = await response.json();
  const noTeam = teamSelect.querySelector('option[value="none"]');
  teams.forEach((team, index) => {
    const option = document.createElement("option");
    option.value = String(index + 1);
    option.textContent = team.name;
    option.style.color = team.color;
    teamSelect.insertBefore(option, noTeam);
  });
}

function updateTeamHud() {
  const teams = state?.teams || [];
  teamHud.hidden = teams.length === 0;
  teamList.innerHTML = "";
  for (const team of teams) {
    const item = document.createElement("li");
    item.style.color = team.color;
    item.textContent = `${team.name}: ${team.score} bodov (${team.players} hracov)`;
    teamList.append(item);
  }
}

// Prepina medzi uvodnou obrazovkou (vyber mena, mapy, timu) a hrou (mapa, zbran).
function showGameScreen(playing) {
  document.body.classList.toggle("is-playing", playing);
  joinForm.hidden = playing;
  playerInfo.hidden = !playing;
  gameArea.hidden = !playing;
  if (playing) {
    // Kurzor nesmie ostat v skrytom policku s menom, inak by sa klavesy ignorovali.
    document.activeElement?.blur();
  } else {
    weaponHud.hidden = true;
    teamHud.hidden = true;
  }
}

function connect() {
  // Lekcia 12: WebSocket udrziava trvale spojenie medzi prehliadacom a serverom.
  if (socket) {
    socket.close();
  }

  const params = new URLSearchParams({
    // Lekcia 24: Vybranu mapu a meno hraca posleme serveru v query parametroch.
    map: mapSelect.value,
    name: playerName.value || "Hrac",
    team: teamSelect.value,
  });
  const protocol = location.protocol === "https:" ? "wss" : "ws";
  const mySocket = new WebSocket(`${protocol}://${location.host}/ws?${params}`);
  socket = mySocket;
  joinButton.disabled = true;
  statusText.textContent = "Pripajam sa...";

  mySocket.addEventListener("open", () => {
	// Lekcia 06: Tu zmen text a CSS triedu statusu podla stavu spojenia.
    statusText.textContent = "Pripojene. Ovladanie najdes nizsie.";
  });

  mySocket.addEventListener("message", (event) => {
	// Lekcia 13: Snapshot je kopia pravdiveho stavu servera; klient ho iba zobrazuje.
    const message = JSON.parse(event.data);
    if (message.type === "state") {
      if (message.playerId) {
        playerId = message.playerId;
      }
      state = message.state;
      // Hra sa ukaze az ked pride prvy stav s nasim tankom.
      if (gameArea.hidden && state.tanks[playerId]) {
        showGameScreen(true);
      }
      updatePlayerSummary();
      updateWeaponHud();
      updateTeamHud();
      draw();
    }
    if (message.type === "error") {
      statusText.textContent = message.error;
    }
  });

  mySocket.addEventListener("close", () => {
    // Stare spojenie (napr. po opatovnom pripojeni) uz obrazovku nemeni.
    if (socket !== mySocket) {
      return;
    }
    socket = null;
    state = null;
    playerId = null;
    stopMovement();
    joinButton.disabled = false;
    if (!gameArea.hidden) {
      statusText.textContent = "Odpojene od servera.";
    } else if (statusText.textContent === "Pripajam sa...") {
      statusText.textContent = "Nepodarilo sa pripojit.";
    }
    showGameScreen(false);
  });
}

// Vedomym odpojenim hrac odstrani svoj tank aj herny stav zo servera.
function disconnect() {
  if (socket?.readyState === WebSocket.OPEN) {
    socket.send(JSON.stringify({ type: "leave" }));
    socket.close();
  }
}

// Kratky popis hraca nad zbranou: meno, tim a mapa.
function updatePlayerSummary() {
  const tank = state?.tanks?.[playerId];
  if (!tank) {
    return;
  }
  const team = tank.teamName || "Bez timu";
  playerSummary.textContent = `${tank.name} | ${team} | mapa ${state.mapName}`;
  playerSummary.style.color = tank.teamColor || "";
}

function updateWeaponHud() {
  const weapon = state?.weapons?.[playerId];
  if (!weapon) {
    weaponHud.hidden = true;
    return;
  }

  weaponHud.hidden = false;
  weaponName.textContent = weapon.count > 1
    ? `${weapon.name} (${weapon.index + 1}/${weapon.count})`
    : weapon.name;
  weaponName.style.color = weapon.color || "";
  weaponDamage.textContent = `Poskodenie: ${weapon.damage}`;
  const tank = state.tanks[playerId];
  if (tank) {
    const divisor = tank.armorDivisor > 1 ? `zasah / ${tank.armorDivisor}` : "plny zasah";
    const speed = Math.round((tank.armorSpeed || 1) * 100);
    armorInfo.textContent = `Pancier: ${tank.armorName} (${divisor}, rychlost ${speed} %)`;
    armorInfo.style.color = tank.armorColor || "";
    updateBonusInfo(tank.bonuses || []);
  }
  if (weapon.index !== lastWeaponIndex) {
    lastWeaponIndex = weapon.index;
    reloadDurationMs = 0;
  }
  ammoCount.replaceChildren();
  for (let index = 0; index < weapon.magazineSize; index++) {
    const round = document.createElement("span");
    round.className = `ammo-round${index < weapon.shotsRemaining ? " is-loaded" : ""}`;
    round.setAttribute("aria-hidden", "true");
    ammoCount.append(round);
  }
  ammoCount.setAttribute(
    "aria-label",
    `${weapon.shotsRemaining} z ${weapon.magazineSize} nabojov`,
  );

  if (!weapon.reloading) {
    reloadStatus.classList.remove("is-reloading");
    reloadStatus.textContent = "Pripravena";
    reloadTrack.hidden = true;
    reloadDurationMs = 0;
    return;
  }

  reloadStatus.classList.add("is-reloading");
  reloadDurationMs = Math.max(reloadDurationMs, weapon.reloadRemainingMs);
  const remaining = Math.max(0, weapon.reloadRemainingMs);
  const progress = reloadDurationMs > 0
    ? Math.min(100, ((reloadDurationMs - remaining) / reloadDurationMs) * 100)
    : 0;
  reloadStatus.textContent = `Prebijanie… ${(remaining / 1000).toFixed(1)} s`;
  reloadTrack.hidden = false;
  reloadTrack.setAttribute("aria-valuenow", String(Math.round(progress)));
  reloadProgress.style.width = `${progress}%`;
}

// Lekcia 37: Aktivne bonusy sa ukazu v paneli aj s odpocitavanim.
function updateBonusInfo(bonuses) {
  bonusInfo.replaceChildren("Bonusy: ");
  if (bonuses.length === 0) {
    bonusInfo.append("ziadne");
    return;
  }
  bonuses.forEach((bonus, index) => {
    const item = document.createElement("span");
    item.style.color = bonus.color;
    item.textContent = `${bonus.name} ${(bonus.remainingMs / 1000).toFixed(1)} s`;
    bonusInfo.append(index > 0 ? ", " : "", item);
  });
}

function send(message) {
  if (!socket || socket.readyState !== WebSocket.OPEN) {
    return;
  }
  socket.send(JSON.stringify(message));
}

// ============================================================
// Lekcia 23/27: KRESLENIE OBJEKTOV
// Kazdy objekt v hre ma vlastnu funkciu drawXxx. Ked chces zmenit,
// ako nieco vyzera, uprav len tu jednu funkciu alebo farby vo "vzhlad".
// Vsetky funkcie dostanu stred objektu v pixeloch (cx, cy) a velkost
// jedneho policka (size), takze nemusis riesit posun mapy.
// ============================================================

// Tu si mozes menit farby a velkosti bez toho, aby si menil kod funkcii.
const vzhlad = {
  pozadie: "#222",
  kamera: {
    zvacsenie: 2,            // 1 = cela mapa na obrazovke, 2 = dvakrat vacsia a kamera sleduje tvoj tank
  },
  tank: {
    farba: "#e3a72f",        // zivy pripojeny tank
    farbaOffline: "#766b60", // hrac sa odpojil
    farbaZniceny: "#5d5143", // zniceny tank caka na ozivenie
    sirka: 0.86,             // v nasobkoch policka
    vyska: 0.72,
    farbaHlavne: "#fff3d6",
    dlzkaHlavne: 0.58,
    hrubkaHlavne: 0.14,
    farbaMena: "#fff3d6",
    pismoMena: "14px Georgia",
  },
  strela: {
    farba: "#e45d3d",        // pouzije sa, ked zbran nema vlastnu farbu
    velkost: 0.14,
  },
  bonus: {
    velkost: 0.42,
    farbaObrysu: "#fff3d6",
    farbaZnacky: "#1b1610",
  },
  vybuch: {
    vypln: "rgba(255, 112, 42, 0.4)",
    obrys: "rgba(255, 204, 87, 0.9)",
  },
  stit: {
    velkost: 0.75,
  },
  zivot: {
    pozadie: "#1b1712",
    vela: "#5fbf4a",
    stredne: "#e3c22f",
    malo: "#e45d3d",
  },
};

function draw() {
  context.clearRect(0, 0, canvas.width, canvas.height);
  if (!state) {
    return;
  }

  // Kamera: ked hras, mapa je zvacsena a stred obrazovky sleduje tvoj tank.
  // Ked este nehras, vidis celu mapu.
  const myTank = state.tanks[playerId];
  const zoom = myTank ? vzhlad.kamera.zvacsenie : 1;
  const size = zoom * Math.min(canvas.width / state.width, canvas.height / state.height);
  const cameraX = myTank ? myTank.x : state.width / 2;
  const cameraY = myTank ? myTank.y : state.height / 2;
  const offsetX = cameraOffset(cameraX, state.width, size, canvas.width);
  const offsetY = cameraOffset(cameraY, state.height, size, canvas.height);
  // Prepocita poziciu z mapy (v polickach) na pixely na obrazovke.
  const screenX = (mapX) => offsetX + mapX * size;
  const screenY = (mapY) => offsetY + mapY * size;

  // Poradie je dolezite: co sa kresli neskor, je navrchu.
  for (let y = 0; y < state.height; y++) {
    for (let x = 0; x < state.width; x++) {
      drawTile(state.tiles[y][x], screenX(x), screenY(y), size);
    }
  }
  for (const bonus of state.bonuses || []) {
    drawBonus(bonus, screenX(bonus.x + 0.5), screenY(bonus.y + 0.5), size);
  }
  for (const explosion of state.explosions || []) {
    drawExplosion(explosion, screenX(explosion.x), screenY(explosion.y), size);
  }
  for (const tank of Object.values(state.tanks)) {
    drawTank(tank, screenX(tank.x), screenY(tank.y), size);
  }
  for (const bullet of Object.values(state.bullets)) {
    drawBullet(bullet, screenX(bullet.x), screenY(bullet.y), size);
  }
}

// Kde na obrazovke zacina mapa (v pixeloch) v jednom smere (x alebo y).
// center je pozicia kamery na mape, mapLength dlzka mapy v polickach,
// screenLength velkost obrazovky v pixeloch.
function cameraOffset(center, mapLength, size, screenLength) {
  const mapPixels = mapLength * size;
  if (mapPixels <= screenLength) {
    // Mapa sa zmesti cela, tak ju dame do stredu.
    return (screenLength - mapPixels) / 2;
  }
  // Kamera sa pri okraji mapy zastavi, aby sme nevideli prazdne miesto za mapou.
  const offset = screenLength / 2 - center * size;
  return Math.min(0, Math.max(screenLength - mapPixels, offset));
}

// Policko mapy. (x, y) je LAVY HORNY roh policka, nie stred.
// tile je znak z mapy: "." tunel, "#" hlina, "X" skala.
function drawTile(tile, x, y, size) {
  context.fillStyle = tileColors[tile] || vzhlad.pozadie;
  context.fillRect(x, y, size, size);
}

// Cely tank: telo, pancier, hlaven, stit, zivot a meno.
function drawTank(tank, cx, cy, size) {
  context.save();
  context.translate(cx, cy);
  // Po otoceni kreslime tank, akoby vzdy mieril doprava (smer +x).
  context.rotate(-tank.angle);
  drawTankBody(tank, size);
  drawArmorOutline(tank, size);
  drawTankBarrel(tank, size);
  context.restore();

  if (tank.invulnerable && tank.alive) {
    drawShield(cx, cy, size);
  }
  const barY = drawHealthBar(tank, cx, cy, size);
  drawTankName(tank, cx, barY - 4, size);
}

// Telo tanku. Bod (0, 0) je stred tanku, tank mieri doprava.
function drawTankBody(tank, size) {
  context.fillStyle = tankColor(tank);
  const width = size * vzhlad.tank.sirka;
  const height = size * vzhlad.tank.vyska;
  context.fillRect(-width / 2, -height / 2, width, height);
}

// Farba tanku podla toho, ci je zivy, ci je hrac pripojeny a v akom je time.
function tankColor(tank) {
  if (!tank.alive) {
    return vzhlad.tank.farbaZniceny;
  }
  if (!tank.connected) {
    return vzhlad.tank.farbaOffline;
  }
  return tank.teamColor || vzhlad.tank.farba;
}

// Hlaven tanku: ciara zo stredu smerom doprava.
function drawTankBarrel(tank, size) {
  context.strokeStyle = vzhlad.tank.farbaHlavne;
  context.lineWidth = Math.max(2, size * vzhlad.tank.hrubkaHlavne);
  context.beginPath();
  context.moveTo(0, 0);
  context.lineTo(size * vzhlad.tank.dlzkaHlavne, 0);
  context.stroke();
}

// Meno hraca nad tankom. y je spodok textu.
function drawTankName(tank, cx, y, size) {
  context.fillStyle = tank.teamColor || vzhlad.tank.farbaMena;
  context.font = vzhlad.tank.pismoMena;
  context.fillText(`${tank.name}${tank.connected ? "" : " (offline)"}`, cx - size / 2, y);
}

// Strela. bullet.color je farba zbrane, bullet.weapon jej nazov.
function drawBullet(bullet, cx, cy, size) {
  context.fillStyle = bullet.color || vzhlad.strela.farba;
  context.beginPath();
  context.arc(cx, cy, Math.max(2, size * vzhlad.strela.velkost), 0, Math.PI * 2);
  context.fill();
}

// Vybuch. explosion.radius je polomer v polickach.
function drawExplosion(explosion, cx, cy, size) {
  context.fillStyle = vzhlad.vybuch.vypln;
  context.strokeStyle = vzhlad.vybuch.obrys;
  context.lineWidth = Math.max(1, size * 0.06);
  context.beginPath();
  context.arc(cx, cy, explosion.radius * size, 0, Math.PI * 2);
  context.fill();
  context.stroke();
}

// Lekcia 37: Bonusova krabicka je farebny kruzok so znackou, prekvapenie ma otaznik.
// bonus.color je farba, bonus.symbol znacka a bonus.name nazov bonusu.
function drawBonus(bonus, cx, cy, size) {
  const pulse = 1 + 0.12 * Math.sin(Date.now() / 180);
  const radius = Math.max(5, size * vzhlad.bonus.velkost * pulse);
  context.fillStyle = bonus.color;
  context.strokeStyle = vzhlad.bonus.farbaObrysu;
  context.lineWidth = Math.max(1, size * 0.08);
  context.beginPath();
  context.arc(cx, cy, radius, 0, Math.PI * 2);
  context.fill();
  context.stroke();
  context.fillStyle = vzhlad.bonus.farbaZnacky;
  context.font = `bold ${Math.max(9, Math.round(radius * 1.3))}px Georgia`;
  context.textAlign = "center";
  context.textBaseline = "middle";
  context.fillText(bonus.symbol, cx, cy + 1);
  context.textAlign = "start";
  context.textBaseline = "alphabetic";
}

// Lekcia 37: Nesmrtelny tank ma okolo seba blikajuci stit.
function drawShield(cx, cy, size) {
  context.strokeStyle = `rgba(255, 242, 122, ${0.55 + 0.35 * Math.sin(Date.now() / 120)})`;
  context.lineWidth = Math.max(2, size * 0.1);
  context.beginPath();
  context.arc(cx, cy, size * vzhlad.stit.velkost, 0, Math.PI * 2);
  context.stroke();
}

// Lekcia 36: Pancier sa kresli ako farebny obrys. Silnejsi pancier ma hrubsi obrys.
// Kresli sa v otocenych suradniciach tanku, (0, 0) je stred tanku.
function drawArmorOutline(tank, size) {
  if (!tank.armorColor) {
    return;
  }
  const bodyWidth = size * vzhlad.tank.sirka;
  const bodyHeight = size * vzhlad.tank.vyska;
  const thickness = Math.max(2, size * (0.04 + 0.025 * Math.max(1, tank.armorDivisor)));
  const width = bodyWidth + thickness;
  const height = bodyHeight + thickness;
  context.strokeStyle = tank.alive ? tank.armorColor : "#6f665b";
  context.lineWidth = thickness;
  context.strokeRect(-width / 2, -height / 2, width, height);
  context.strokeStyle = "rgba(20, 16, 12, 0.85)";
  context.lineWidth = Math.max(1, thickness / 3);
  context.strokeRect(-bodyWidth / 2, -bodyHeight / 2, bodyWidth, bodyHeight);
}

// Pruh zivota nad tankom. Vrati y horneho okraja pruhu (nad nim je meno).
function drawHealthBar(tank, cx, cy, size) {
  const width = Math.max(24, size * 1.2);
  const height = Math.max(4, size * 0.16);
  const x = cx - width / 2;
  const y = cy - size * 0.6 - height;
  const maxHealth = tank.maxHealth || 1;
  const ratio = Math.max(0, Math.min(1, (tank.health ?? maxHealth) / maxHealth));

  context.fillStyle = vzhlad.zivot.pozadie;
  context.fillRect(x - 1, y - 1, width + 2, height + 2);
  context.fillStyle = ratio > 0.6 ? vzhlad.zivot.vela : ratio > 0.3 ? vzhlad.zivot.stredne : vzhlad.zivot.malo;
  context.fillRect(x, y, width * ratio, height);
  return y;
}

playerName.addEventListener("keydown", (event) => {
  event.stopPropagation();
});

mapSelect.addEventListener("keydown", (event) => {
  event.stopPropagation();
});

const movementByKey = {
  ArrowUp: [0, -1], w: [0, -1], W: [0, -1],
  ArrowDown: [0, 1], s: [0, 1], S: [0, 1],
  ArrowLeft: [-1, 0], a: [-1, 0], A: [-1, 0],
  ArrowRight: [1, 0], d: [1, 0], D: [1, 0],
};

function sendMovement() {
  let dx = 0;
  let dy = 0;
  for (const vector of movementInputs.values()) {
    dx += vector[0];
    dy += vector[1];
  }
  if (dx === 0 && dy === 0) {
    return;
  }

  const magnitude = Math.hypot(dx, dy);
  dx /= magnitude;
  dy /= magnitude;
  const direction = dy < 0
    ? (dx < 0 ? "up-left" : dx > 0 ? "up-right" : "up")
    : dy > 0
      ? (dx < 0 ? "down-left" : dx > 0 ? "down-right" : "down")
      : (dx < 0 ? "left" : "right");
  send({ type: "move", direction, angle: Math.atan2(-dy, dx) });
}

function startMovement() {
  if (!movementTimer) {
    sendMovement();
    movementTimer = window.setInterval(sendMovement, 30);
  }
}

document.addEventListener("keydown", (event) => {
  // Lekcia 17: Mapovanie klavesov oddeluje ovladanie od sprav posielanych serveru.
  if (event.target instanceof HTMLInputElement || event.target instanceof HTMLSelectElement || event.target instanceof HTMLTextAreaElement || event.target instanceof HTMLButtonElement) {
    return;
  }

  if (movementByKey[event.key]) {
    event.preventDefault();
    movementInputs.set(`keyboard:${event.key}`, movementByKey[event.key]);
    // Rychlost tanku urcuje len casovac (kazdych 30 ms jeden krok).
    // Drzany klaves posiela opakovane keydown udalosti, tie by tank zrychlili,
    // preto krok posleme hned len pri prvom stlaceni.
    startMovement();
  }
  if (event.key === " ") {
    // Lekcia 21: Medzernik vytvori na serveri novu strelu, ktoru pohana ticker.
    event.preventDefault();
    send({ type: "shoot" });
  }
  // Lekcia 33: Q prepne na dalsiu zbran.
  if (event.key === "q" || event.key === "Q") {
    event.preventDefault();
    send({ type: "weapon" });
  }
  // Lekcia 36: E prepne na dalsi pancier.
  if (event.key === "e" || event.key === "E") {
    event.preventDefault();
    send({ type: "armor" });
  }
  // Lekcia 43: Sem doplnime klavesy 1-9, ktore vyberu konkretnu zbran.
});

document.addEventListener("keyup", (event) => {
  movementInputs.delete(`keyboard:${event.key}`);
  if (movementInputs.size === 0 && movementTimer) {
    window.clearInterval(movementTimer);
    movementTimer = null;
  }
});

const directionVectors = {
  up: [0, -1],
  down: [0, 1],
  left: [-1, 0],
  right: [1, 0],
};

document.querySelectorAll(".direction-button").forEach((button) => {
  const vector = directionVectors[button.dataset.direction];

  button.addEventListener("pointerdown", (event) => {
    event.preventDefault();
    button.setPointerCapture(event.pointerId);
    movementInputs.set(`pointer:${event.pointerId}`, vector);
    startMovement();
  });

  const stopPointerMovement = (event) => {
    movementInputs.delete(`pointer:${event.pointerId}`);
    if (movementInputs.size === 0 && movementTimer) {
      window.clearInterval(movementTimer);
      movementTimer = null;
    }
  };
  button.addEventListener("pointerup", stopPointerMovement);
  button.addEventListener("pointercancel", stopPointerMovement);
  button.addEventListener("lostpointercapture", stopPointerMovement);
});

function stopMovement() {
  movementInputs.clear();
  if (movementTimer) {
    window.clearInterval(movementTimer);
    movementTimer = null;
  }
}

window.addEventListener("blur", stopMovement);

function bindActionButton(button, message) {
  button.addEventListener("pointerdown", (event) => {
    event.preventDefault();
    send(message);
  });
  button.addEventListener("click", (event) => {
    // Pointer akcie sa odosielaju hned pri dotyku; click obsluhuje klavesnicu a asistivne technologie.
    if (event.detail === 0) {
      send(message);
    }
  });
}

bindActionButton(shootButton, { type: "shoot" });
bindActionButton(weaponButton, { type: "weapon" });
bindActionButton(armorButton, { type: "armor" });
joinButton.addEventListener("click", connect);
leaveButton.addEventListener("click", disconnect);

// Lekcia 30: Pod canvas mozes doplnit HUD so zoznamom hracov a ich score.
// Lekcia 31: Sem pridaj vlastnu funkcionalitu a priprav jej kratku ukazku.

loadMaps().catch((error) => {
  statusText.textContent = `Mapy sa nepodarilo nacitat: ${error.message}`;
});
loadTeams().catch((error) => {
  statusText.textContent = `Timy sa nepodarilo nacitat: ${error.message}`;
});
draw();