/**
 * @typedef {Object} InputMessage
 * @property {"input"} event
 * @property {string} command
 */

/**
 * @typedef {Object} ResizeMessage
 * @property {"resize"} event
 * @property {Resize} resize
 */

/**
 * @typedef {Object} Resize
 * @property {number} rows
 * @property {number} cols
 */

/** @typedef {InputMessage | ResizeMessage} ClientMessage */

/** @typedef {(args: string[]) => string} CommandHandler */

/**
 * @typedef {Object} Disposable
 * @property {() => void} dispose
 */

/**
 * @typedef {Object} TerminalTheme
 * @property {string} [background]
 * @property {string} [foreground]
 * @property {string} [cursor]
 * @property {string} [green]
 */

/**
 * @typedef {Object} TerminalOptions
 * @property {boolean} [cursorBlink]
 * @property {string} [fontFamily]
 * @property {number} [fontSize]
 * @property {number} [lineHeight]
 * @property {boolean} [convertEol]
 * @property {TerminalTheme} [theme]
 */

/**
 * @typedef {Object} TerminalAddon
 * @property {(terminal: TerminalInstance) => void} activate
 * @property {() => void} [dispose]
 */

/**
 * @typedef {Object} TerminalInstance
 * @property {TerminalOptions} options
 * @property {number} cols
 * @property {number} rows
 * @property {(addon: TerminalAddon) => void} loadAddon
 * @property {(element: HTMLElement) => void} open
 * @property {(data: string) => void} write
 * @property {(data: string) => void} writeln
 * @property {() => void} clear
 * @property {(handler: (event: KeyboardEvent) => boolean) => void} attachCustomKeyEventHandler
 * @property {(handler: (data: string) => void) => Disposable} onData
 */

/** @typedef {new (options?: TerminalOptions) => TerminalInstance} TerminalConstructor */

/**
 * @typedef {TerminalAddon & { fit: () => void }} FitAddonInstance
 */

/** @typedef {new () => FitAddonInstance} FitAddonConstructor */
/** @typedef {new () => TerminalAddon} WebLinksAddonConstructor */

/**
 * @typedef {Object} XTermGlobals
 * @property {TerminalConstructor} Terminal
 * @property {{ FitAddon: FitAddonConstructor }} FitAddon
 * @property {{ WebLinksAddon: WebLinksAddonConstructor }} WebLinksAddon
 */

export {};
