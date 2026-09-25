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

/** @type {Record<string, string>} */
const commands = {
  help: "about      who I am\r\nprojects   selected work\r\nclear      clear the terminal",
  about:
    "Harry Lawton. Software engineer exploring graphics, games, and tools from first principles.",
  projects: "terminal-wireframe   perkins   tinyrenderer",
};
let input = "";

function prompt() {
  terminal.write("\r\n\x1b[32mvisitor@portfolio\x1b[0m:\x1b[34m~\x1b[0m$ ");
}

terminal.writeln("\x1b[32mPORTFOLIO / SYSTEM ONLINE\x1b[0m");
terminal.writeln("A real sandboxed shell is coming next. Type 'help'.");
prompt();

terminal.onData((data) => {
  switch (data) {
    case "\r":
      if (input === "clear") terminal.clear();
      else {
        const output = commands[input] ?? (input ? `${input}: command not found` : "");
        terminal.write(`\r\n${output}`);
      }
      input = "";
      break;

    case "\u007f":
      if (!input) return;
      input = input.slice(0, -1);
      terminal.write("\b \b");
      break;

    default:
      if (!/^[\x20-\x7e]$/.test(data)) return;
      input += data;
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
