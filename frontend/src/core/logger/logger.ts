const IS_DEV = import.meta.env.DEV;

interface LoggerComponentBadgeStyle {
  emoji: string;
  color: string;
}

export const LOGGER_COMPONENTS: Record<string, LoggerComponentBadgeStyle> = {
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
  Product: { emoji: "📦", color: "cyan" },
} as const;

type ComponentKey = keyof typeof LOGGER_COMPONENTS;

type LogLevel = "debug" | "info" | "warn" | "error";
type ConsoleMethod = "log" | "warn" | "error";

class Logger {
  private readonly badgeStyle: LoggerComponentBadgeStyle;

  constructor(private readonly component: ComponentKey) {
    this.badgeStyle = LOGGER_COMPONENTS[component] || {
      emoji: "⚪",
      color: "black",
    };
  }

  private log(level: LogLevel, msg: string, ...args: unknown[]) {
    // prevent debug logs in production
    if (!IS_DEV && level === "debug") return;

    const levelColor = this.colorForLevel(level);
    const consoleMethod = this.consoleMethodForLevel(level);

    const sLevel = `color: ${levelColor}; font-weight: bold; display: inline-block; width: 60px;`;
    const sBadge = `background: ${this.badgeStyle.color}; color: #fff; padding: 2px 6px; border-radius: 4px; font-weight: bold; margin-right: 5px; display: inline-block; width: 150px;`;
    const sMsg = `color: ${levelColor};`;

    // eslint-disable-next-line no-console
    console[consoleMethod](
      `%c${level.padEnd(8)} %c${this.badgeStyle.emoji} ${this.component} %c${msg}`,
      sLevel,
      sBadge,
      sMsg,
      ...args,
    );
  }

  private consoleMethodForLevel(level: LogLevel): ConsoleMethod {
    const methodMap: Record<LogLevel, ConsoleMethod> = {
      debug: "log",
      info: "log",
      warn: "warn",
      error: "error",
    };
    return methodMap[level] || methodMap.info;
  }

  private colorForLevel(level: LogLevel): string {
    const levelColors: Record<LogLevel, string> = {
      debug: "gray",
      info: "darkgray",
      warn: "orange",
      error: "red",
    };
    return levelColors[level] || levelColors.info;
  }

  debug(msg: string, ...args: unknown[]) {
    this.log("debug", msg, ...args);
  }
  info(msg: string, ...args: unknown[]) {
    this.log("info", msg, ...args);
  }
  warn(msg: string, ...args: unknown[]) {
    this.log("warn", msg, ...args);
  }
  error(msg: string, ...args: unknown[]) {
    this.log("error", msg, ...args);
  }
}

export const createLogger = (component: ComponentKey) => new Logger(component);
