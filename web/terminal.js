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
  lineHeight: 1.35,
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

const terminalElement = document.querySelector("#terminal");
if (!(terminalElement instanceof HTMLElement)) throw new Error("Terminal element not found");

terminal.open(terminalElement);
fit.fit();

const protocol = location.protocol === "https:" ? "wss:" : "ws:";
const socket = new WebSocket(`${protocol}//${location.host}/shell`);

// socket.addEventListener("open", () => console.log("websocket connected"));
socket.addEventListener("message", (event) => {
    terminal.write(event.data)
});
socket.addEventListener("error", () => terminal.writeln("WebSocket error"));
socket.addEventListener("close", () => terminal.writeln("WebSocket closed"));

let buffer = "";

/**
 * @param {string} buf
 */
function sendInput(buf) {
  const payload = JSON.stringify({
    event: "input",
    command: buf
  });
  socket.send(payload);
}

terminal.onData((data) => {
  switch (data) {
    case "\r":
      if (buffer === "clear") terminal.clear();
      else {
         sendInput(buffer);
      }
      buffer = "";
      break;

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
