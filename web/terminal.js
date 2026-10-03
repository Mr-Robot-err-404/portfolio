// @ts-check

/** @typedef {import("./types.js").XTermGlobals} XTermGlobals */

const root = document.documentElement;
const styles = getComputedStyle(root);

const xtermGlobals = /** @type {XTermGlobals} */ (/** @type {unknown} */ (window));
const XTerm = xtermGlobals.Terminal;
const XTermFitAddon = xtermGlobals.FitAddon.FitAddon;
const XTermWebLinksAddon = xtermGlobals.WebLinksAddon.WebLinksAddon;

const terminal = new XTerm({
  cursorBlink: true,
  fontFamily: '"IBM Plex Mono", monospace',
  fontSize: 14,
  lineHeight: 1,
  theme: {
    background: styles.getPropertyValue("--surface").trim(),
    foreground: styles.getPropertyValue("--text").trim(),
    cursor: styles.getPropertyValue("--accent").trim(),
    green: styles.getPropertyValue("--accent").trim(),
  },
});
const fit = new XTermFitAddon();

terminal.loadAddon(fit);
terminal.loadAddon(new XTermWebLinksAddon());
terminal.options.convertEol = true;

const terminalElement = document.querySelector("#terminal");
if (!(terminalElement instanceof HTMLElement)) throw new Error("Terminal element not found");

terminal.open(terminalElement);
fit.fit();

const protocol = location.protocol === "https:" ? "wss:" : "ws:";
const socket = new WebSocket(`${protocol}//${location.host}/shell`);

socket.binaryType = "arraybuffer";

socket.addEventListener("open", () => {
  socket.send(JSON.stringify({
    event: "resize",
    resize: { cols: terminal.cols, rows: terminal.rows },
  }));
});

socket.addEventListener("message", async (event) => {
  if (typeof event.data === "string") {
    terminal.write(event.data);
    return;
  }
  const bytes = new Uint8Array(event.data);
  terminal.write(new TextDecoder().decode(bytes));
});
socket.addEventListener("error", () => terminal.writeln("WebSocket error"));
socket.addEventListener("close", () => terminal.writeln("WebSocket closed"));

const promptText = "\x1b[32mvisitor@portfolio\x1b[0m:~ \x1b[34m❯\x1b[0m ";
let buffer = "";

class CommandHistory {
  constructor() {
    /** @type {string[]} */
    this.list = [];
    this.idx = 0;
    this.draft = "";
  }

  /** @param {string} command */
  append(command) {
    if (command && this.list.at(-1) !== command) this.list.push(command);
    this.idx = this.list.length;
    this.draft = "";
  }

  /** @param {string} current */
  previous(current) {
    if (!this.list.length) return null;
    if (this.idx === this.list.length) this.draft = current;
    this.idx = Math.max(0, this.idx - 1);
    return this.list[this.idx];
  }

  next() {
    if (!this.list.length || this.idx === this.list.length) return null;
    this.idx++;
    return this.idx === this.list.length ? this.draft : this.list[this.idx];
  }
}

const commandHistory = new CommandHistory();

/**
 * @param {string} buf
 */
function sendInput(buf) {
  const payload = JSON.stringify({
    event: "input",
    command: buf,
  });
  socket.send(payload);
}

function redrawInput() {
  terminal.write(`\r\x1b[2K${promptText}${buffer}`);
}

terminal.attachCustomKeyEventHandler((event) => {
  if (event.type !== "keydown") return true;

  switch (event.key) {
    case "ArrowUp":
      buffer = commandHistory.previous(buffer) ?? buffer;
      redrawInput();
      return false;

    case "ArrowDown":
      buffer = commandHistory.next() ?? buffer;
      redrawInput();
      return false;

    case "Backspace":
      if (buffer) {
        buffer = buffer.slice(0, -1);
        terminal.write("\b \b");
      }
      return false;

    default:
      return true;
  }
});

terminal.onData((data) => {
  switch (data) {
    case "\r":
      commandHistory.append(buffer);
      sendInput(buffer);
      buffer = "";
      break;

    case "\b":
    case "\u007f":
      if (!buffer) return;
      buffer = buffer.slice(0, -1);
      terminal.write("\b \b");
      break;

    default:
      if (!/^[\x20-\x7e]$/.test(data)) return;
      buffer += data;
      terminal.write(data);
  }
});

window.addEventListener("resize", () => fit.fit());

const projectMedia = /** @type {NodeListOf<HTMLElement>} */ (
  document.querySelectorAll(".project-media[data-demo], .project-media[data-video]")
);

for (const media of projectMedia) {
  function showDemo() {
    if (media.querySelector(".project-demo")) return;
    const source = media.dataset.video || media.dataset.demo;
    if (!source) return;
    const isVideo = Boolean(media.dataset.video);

    /** @type {HTMLImageElement | HTMLVideoElement} */
    let demo;

    if (isVideo) {
      const video = document.createElement("video");
      demo = video;
      demo.muted = true;
      demo.loop = true;
      demo.playsInline = true;
      demo.addEventListener(
        "loadeddata",
        () => {
          if (!demo.isConnected) return;
          demo.classList.add("is-ready");
          video.play().catch(() => {});
        },
        { once: true },
      );
    } else {
      const image = document.createElement("img");
      demo = image;
      demo.alt = "";
      demo.addEventListener("load", () => demo.classList.add("is-ready"), { once: true });
    }
    demo.className = "project-demo";
    media.append(demo);
    demo.src = source;
  }

  function hideDemo() {
    const demo = media.querySelector(".project-demo");
    if (demo instanceof HTMLVideoElement) demo.pause();
    demo?.remove();
  }

  media.addEventListener("pointerenter", showDemo);
  media.addEventListener("pointerleave", hideDemo);
  media.addEventListener("focus", showDemo);
  media.addEventListener("blur", hideDemo);
}
