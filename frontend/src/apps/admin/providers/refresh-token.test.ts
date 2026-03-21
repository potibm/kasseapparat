import { describe, it, expect, beforeEach, vi, afterEach } from "vitest";
import { refreshToken, TOKEN_REFRESH_THRESHOLD } from "./refresh-token"; // Adjust path if necessary
import { getSession, updateToken } from "../utils/auth-utils";

function addSeconds(date: Date, seconds: number): Date {
  date.setSeconds(date.getSeconds() + seconds);
  return date;
}

// --- 1. MOCKS ---
vi.mock("../utils/auth-utils", () => ({
  getSession: vi.fn(),
  updateToken: vi.fn(),
}));

// Helper to create a type-safe session mock without 'any'
// We only need the expiryDate for these tests, so we cast it safely via unknown
type SessionData = NonNullable<ReturnType<typeof getSession>>;
const createMockSession = (expiryDate: string): SessionData =>
  ({ expiryDate }) as unknown as SessionData;

// --- 2. TESTS ---
describe("refreshToken Logic", () => {
  const now = new Date("2024-01-01T12:00:00.000Z");
  const refrehThresholdSeconds = TOKEN_REFRESH_THRESHOLD / 1000;

  beforeEach(() => {
    vi.clearAllMocks();

    // Freeze time to a fixed point: Jan 1st, 2024 at 12:00:00.000Z
    vi.useFakeTimers();
    vi.setSystemTime(now);
  });

  afterEach(() => {
    // Always clean up timers so we don't break other tests
    vi.useRealTimers();
  });

  it("should throw an error if no active session is found", async () => {
    vi.mocked(getSession).mockReturnValue(null);

    await expect(refreshToken()).rejects.toThrow(
      "No active session found. Please log in.",
    );
    expect(updateToken).not.toHaveBeenCalled();
  });

  it("should return early and do nothing if the token is valid for MORE than 10 seconds", async () => {
    // Token expires at 12:00:15 (15 seconds in the future)
    vi.mocked(getSession).mockReturnValue(
      createMockSession(
        addSeconds(now, refrehThresholdSeconds + 5).toISOString(),
      ),
    );

    await refreshToken();

    // Since it's > 10s, updateToken should NOT be called
    expect(updateToken).not.toHaveBeenCalled();
  });

  it("should call updateToken if the token expires in EXACTLY 10 seconds", async () => {
    // Token expires at 12:00:10 (exactly 10 seconds in the future)
    // 10000 > 10000 is false, so it should proceed to updateToken
    vi.mocked(getSession).mockReturnValue(
      createMockSession(addSeconds(now, refrehThresholdSeconds).toISOString()),
    );

    await refreshToken();

    expect(updateToken).toHaveBeenCalledTimes(1);
  });

  it("should call updateToken if the token expires in LESS than 10 seconds", async () => {
    // Token expires at 12:00:05 (only 5 seconds left)
    vi.mocked(getSession).mockReturnValue(
      createMockSession(
        addSeconds(now, refrehThresholdSeconds - 5).toISOString(),
      ),
    );

    await refreshToken();

    expect(updateToken).toHaveBeenCalledTimes(1);
  });

  it("should call updateToken if the token is ALREADY expired (in the past)", async () => {
    // Token expired 1 hour ago
    vi.mocked(getSession).mockReturnValue(
      createMockSession(addSeconds(now, -3600).toISOString()),
    );

    await refreshToken();

    expect(updateToken).toHaveBeenCalledTimes(1);
  });
});
