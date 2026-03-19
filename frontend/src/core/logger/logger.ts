const IS_DEV = import.meta.env.DEV;

export const LOGGER_COMPONENTS = {
  Admin: { emoji: "🔑", color: "blue" },
  Auth: { emoji: "🔐", color: "green" },
  Config: { emoji: "⚙️", color: "teal" },
  Api: { emoji: "📡", color: "purple" },
  Cart: { emoji: "🛒", color: "orange" },
  Guestlist: { emoji: "📝", color: "pink" },
  Payment: { emoji: "💳", color: "red" },
  PaymentWebsocket: { emoji: "📡", color: "red" },
  Purchase: { emoji: "📜", color: "brown" },
  Core: { emoji: "⚪", color: "black" },
} as const;

type ComponentKey = keyof typeof LOGGER_COMPONENTS;

class Logger {
  private readonly emoji: string;
  private readonly color: string;

  constructor(private readonly component: ComponentKey) {
    this.emoji = LOGGER_COMPONENTS[component]?.emoji || "⚪";
    this.color = LOGGER_COMPONENTS[component]?.color || "black";
  }

  private log(
    level: "info" | "warn" | "error" | "debug",
    msg: string,
    ...args: any[]
  ) {
    const levelColors = {
      debug: { text: "gray", bg: "#f8f9fa" },
      info: { text: "darkgray", bg: "#e7f3ff" },
      warn: { text: "orange", bg: "#fff3cd" },
      error: { text: "red", bg: "#f8d7da" },
    };

    // Debug-Logs in Production unterdrücken
    if (!IS_DEV && level === "debug") return;

    const levelColor = levelColors[level] || levelColors.info;

    const sLevel = `color: ${levelColor.text}; font-weight: bold; display: inline-block; width: 60px;`;
    const sBadge = `background: ${this.color}; color: #fff; padding: 2px 6px; border-radius: 4px; font-weight: bold; margin-right: 5px;`;
    const sMsg = `color: ${levelColor.text};`;

    console.log(
      `%c${level.padEnd(8)} %c${this.emoji} ${this.component} %c${msg}`,
      sLevel,
      sBadge,
      sMsg,
      ...args,
    );
  }

  debug(msg: string, ...args: any[]) {
    this.log("debug", msg, ...args);
  }
  info(msg: string, ...args: any[]) {
    this.log("info", msg, ...args);
  }
  warn(msg: string, ...args: any[]) {
    this.log("warn", msg, ...args);
  }
  error(msg: string, ...args: any[]) {
    this.log("error", msg, ...args);
  }
}

export const createLogger = (component: ComponentKey) => new Logger(component);
