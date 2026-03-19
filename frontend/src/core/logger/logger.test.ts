import { describe, it, expect, vi, beforeEach } from "vitest";
import { createLogger } from "./logger";

describe("Logger Utility", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("should format info logs correctly (matching the format string)", () => {
    const logSpy = vi.spyOn(console, "log").mockImplementation(() => {});
    const logger = createLogger("Auth");
    const message = "User logged in";
    const meta = { userId: 123 };

    logger.info(message, meta);

    const firstCallArgs = logSpy.mock.calls[0];

    expect(firstCallArgs[0]).toMatch(/info.*🔐.*Auth/);

    expect(firstCallArgs[1]).toContain("font-weight: bold"); // Level style
    expect(firstCallArgs[2]).toContain("background: green"); // Badge style

    expect(firstCallArgs[4]).toEqual(meta);
  });

  it("should use console.warn for warnings", () => {
    const warnSpy = vi.spyOn(console, "warn").mockImplementation(() => {});
    const logger = createLogger("Payment");

    logger.warn("Attention");

    expect(warnSpy).toHaveBeenCalled();
  });

  it("should use console.error for errors", () => {
    const errorSpy = vi.spyOn(console, "error").mockImplementation(() => {});
    const logger = createLogger("Api");

    logger.error("Failed");

    expect(errorSpy).toHaveBeenCalled();
  });
});

describe("Logger Production Mode", () => {
  beforeEach(() => {
    vi.restoreAllMocks();
    vi.resetModules();
  });

  it("should NOT log debug messages when DEV is false", async () => {
    vi.stubEnv("DEV", false);

    const { createLogger } = await import("./logger");
    const logSpy = vi.spyOn(console, "log").mockImplementation(() => {});

    const logger = createLogger("Api");
    logger.debug("This should be silent");

    expect(logSpy).not.toHaveBeenCalled();
  });

  it("should STILL log warn and error messages in production", async () => {
    vi.stubEnv("DEV", false);

    const { createLogger } = await import("./logger");
    const warnSpy = vi.spyOn(console, "warn").mockImplementation(() => {});
    const errorSpy = vi.spyOn(console, "error").mockImplementation(() => {});

    const logger = createLogger("Payment");
    logger.warn("Production warning");
    logger.error("Production error");

    expect(warnSpy).toHaveBeenCalled();
    expect(errorSpy).toHaveBeenCalled();
  });
});
