const canvas = document.querySelector("#game");
const context = canvas.getContext("2d");
const mapSelect = document.querySelector("#mapSelect");
const playerName = document.querySelector("#playerName");
const joinButton = document.querySelector("#joinButton");
const statusText = document.querySelector("#status");

let socket = null;
let state = null;

playerName.disabled = false;
mapSelect.disabled = false;

const tileColors = {
  ".": "#17130f",
  "#": "#8b6339",
  "X": "#50493d",
};

async function loadMaps() {
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

function connect() {
  if (socket) {
    socket.close();
  }

  const params = new URLSearchParams({
    map: mapSelect.value,
    name: playerName.value || "Hrac",
  });
  const protocol = location.protocol === "https:" ? "wss" : "ws";
  socket = new WebSocket(`${protocol}://${location.host}/ws?${params}`);
  statusText.textContent = "Pripajam sa...";

  socket.addEventListener("open", () => {
    statusText.textContent = "Pripojene. Pohyb: sipky alebo WASD. Strela: medzernik.";
  });

  socket.addEventListener("message", (event) => {
    const message = JSON.parse(event.data);
    if (message.type === "state") {
      state = message.state;
      draw();
    }
    if (message.type === "error") {
      statusText.textContent = message.error;
    }
  });

  socket.addEventListener("close", () => {
    statusText.textContent = "Odpojene od servera.";
  });
}

function send(message) {
  if (!socket || socket.readyState !== WebSocket.OPEN) {
    return;
  }
  socket.send(JSON.stringify(message));
}

function draw() {
  context.clearRect(0, 0, canvas.width, canvas.height);
  if (!state) {
    return;
  }

  const tileSize = Math.min(canvas.width / state.width, canvas.height / state.height);
  const offsetX = (canvas.width - state.width * tileSize) / 2;
  const offsetY = (canvas.height - state.height * tileSize) / 2;

  for (let y = 0; y < state.height; y++) {
    for (let x = 0; x < state.width; x++) {
      const tile = state.tiles[y][x];
      context.fillStyle = tileColors[tile] || "#222";
      context.fillRect(offsetX + x * tileSize, offsetY + y * tileSize, tileSize, tileSize);
    }
  }

  for (const tank of Object.values(state.tanks)) {
    context.fillStyle = tank.alive ? "#e3a72f" : "#5d5143";
    context.fillRect(offsetX + tank.x * tileSize + 3, offsetY + tank.y * tileSize + 3, tileSize - 6, tileSize - 6);
    context.fillStyle = "#fff3d6";
    context.font = "14px Georgia";
    context.fillText(tank.name, offsetX + tank.x * tileSize, offsetY + tank.y * tileSize - 4);
  }

  for (const bullet of Object.values(state.bullets)) {
    context.fillStyle = "#e45d3d";
    context.beginPath();
    context.arc(offsetX + (bullet.x + 0.5) * tileSize, offsetY + (bullet.y + 0.5) * tileSize, Math.max(3, tileSize * 0.18), 0, Math.PI * 2);
    context.fill();
  }
}

playerName.addEventListener("keydown", (event) => {
  event.stopPropagation();
});

mapSelect.addEventListener("keydown", (event) => {
  event.stopPropagation();
});

document.addEventListener("keydown", (event) => {
  if (event.target instanceof HTMLInputElement || event.target instanceof HTMLSelectElement || event.target instanceof HTMLTextAreaElement) {
    return;
  }

  const directionByKey = {
    ArrowUp: "up",
    w: "up",
    W: "up",
    ArrowDown: "down",
    s: "down",
    S: "down",
    ArrowLeft: "left",
    a: "left",
    A: "left",
    ArrowRight: "right",
    d: "right",
    D: "right",
  };
  if (directionByKey[event.key]) {
    event.preventDefault();
    send({ type: "move", direction: directionByKey[event.key] });
  }
  if (event.key === " ") {
    event.preventDefault();
    send({ type: "shoot" });
  }
});

joinButton.addEventListener("click", connect);

loadMaps().catch((error) => {
  statusText.textContent = `Mapy sa nepodarilo nacitat: ${error.message}`;
});
draw();