const root = document.documentElement;
const styles = getComputedStyle(root);
const terminal = new Terminal({
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
const fit = new FitAddon.FitAddon();

terminal.loadAddon(fit);
terminal.loadAddon(new WebLinksAddon.WebLinksAddon());
terminal.open(document.querySelector("#terminal"));
fit.fit();

const commands = {
  help: "about      who I am\r\nprojects   selected work\r\nclear      clear the terminal",
  about: "Harry Lawton. Software engineer exploring graphics, games, and tools from first principles.",
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
  if (data === "\r") {
    if (input === "clear") terminal.clear();
    else terminal.write(`\r\n${commands[input] || (input ? `${input}: command not found` : "")}`);
    input = "";
    prompt();
  } else if (data === "\u007f" && input) {
    input = input.slice(0, -1);
    terminal.write("\b \b");
  } else if (/^[\x20-\x7e]$/.test(data)) {
    input += data;
    terminal.write(data);
  }
});

window.addEventListener("resize", () => fit.fit());

for (const media of document.querySelectorAll(".project-media[data-demo]")) {
  function showDemo() {
    if (media.querySelector(".project-demo")) return;
    const demo = document.createElement("img");
    demo.className = "project-demo";
    demo.alt = "";
    demo.addEventListener("load", () => demo.classList.add("is-ready"), { once: true });
    media.append(demo);
    demo.src = media.dataset.demo;
  }

  function hideDemo() {
    media.querySelector(".project-demo")?.remove();
  }

  media.addEventListener("pointerenter", showDemo);
  media.addEventListener("pointerleave", hideDemo);
  media.addEventListener("focus", showDemo);
  media.addEventListener("blur", hideDemo);
}
