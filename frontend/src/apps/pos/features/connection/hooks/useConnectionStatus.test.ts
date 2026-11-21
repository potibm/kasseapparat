import { describe, it, expect, vi, beforeEach, afterEach } from "vitest";
import { renderHook, act } from "@testing-library/react";
import {
  useConnectionStatus,
  DEFAULT_INTERVAL_MS,
  DEFAULT_TIMEOUT_MS,
} from "./useConnectionStatus";

const URL = "http://localhost:3001/api/v3/config";

/** fetch mock that answers after `ms` with the given ok flag */
const respondAfter = (ms: number, ok = true) =>
  vi.fn(
    () =>
      new Promise<Response>((resolve) =>
        setTimeout(() => resolve({ ok } as Response), ms),
      ),
  );

/** fetch mock that never answers, but rejects once its request is aborted */
const neverRespond = () =>
  vi.fn(
    (_url: string, init?: RequestInit) =>
      new Promise<Response>((_resolve, reject) => {
        init?.signal?.addEventListener("abort", () =>
          reject(new DOMException("Aborted", "AbortError")),
        );
      }),
  );

const setOnLine = (value: boolean) =>
  Object.defineProperty(globalThis.navigator, "onLine", {
    configurable: true,
    get: () => value,
  });

const advance = async (ms: number) => {
  await act(async () => {
    await vi.advanceTimersByTimeAsync(ms);
  });
};

describe("useConnectionStatus", () => {
  beforeEach(() => {
    vi.useFakeTimers();
    setOnLine(true);
  });

  afterEach(() => {
    vi.useRealTimers();
    vi.unstubAllGlobals();
    setOnLine(true);
  });

  it("is online after a fast, successful check", async () => {
    const fetchMock = respondAfter(100);
    vi.stubGlobal("fetch", fetchMock);

    const { result } = renderHook(() => useConnectionStatus(URL));
    await advance(100);

    expect(result.current).toBe("online");
    expect(fetchMock).toHaveBeenCalledWith(
      URL,
      expect.objectContaining({ method: "GET", cache: "no-store" }),
    );
  });

  it("is slow when the check takes longer than the threshold", async () => {
    vi.stubGlobal("fetch", respondAfter(2_500));

    const { result } = renderHook(() => useConnectionStatus(URL));
    await advance(2_500);

    expect(result.current).toBe("slow");
  });

  it("is offline when the backend answers with an error", async () => {
    vi.stubGlobal("fetch", respondAfter(10, false));

    const { result } = renderHook(() => useConnectionStatus(URL));
    await advance(10);

    expect(result.current).toBe("offline");
  });

  it("is offline when the request fails", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockRejectedValue(new TypeError("Failed to fetch")),
    );

    const { result } = renderHook(() => useConnectionStatus(URL));
    await advance(0);

    expect(result.current).toBe("offline");
  });

  it("is offline when the check times out", async () => {
    vi.stubGlobal("fetch", neverRespond());

    const { result } = renderHook(() => useConnectionStatus(URL));
    await advance(DEFAULT_TIMEOUT_MS);

    expect(result.current).toBe("offline");
  });

  it("does not even try while the browser reports offline", async () => {
    const fetchMock = respondAfter(10);
    vi.stubGlobal("fetch", fetchMock);
    setOnLine(false);

    const { result } = renderHook(() => useConnectionStatus(URL));
    await advance(0);

    expect(result.current).toBe("offline");
    expect(fetchMock).not.toHaveBeenCalled();
  });

  it("checks again on every interval", async () => {
    const fetchMock = respondAfter(10);
    vi.stubGlobal("fetch", fetchMock);

    renderHook(() => useConnectionStatus(URL));
    await advance(10);
    expect(fetchMock).toHaveBeenCalledTimes(1);

    await advance(DEFAULT_INTERVAL_MS);
    expect(fetchMock).toHaveBeenCalledTimes(2);
  });

  it("does not start a second check while one is still running", async () => {
    const fetchMock = respondAfter(20_000);
    vi.stubGlobal("fetch", fetchMock);

    renderHook(() =>
      useConnectionStatus(URL, { intervalMs: 5_000, timeoutMs: 30_000 }),
    );
    await advance(12_000);

    expect(fetchMock).toHaveBeenCalledTimes(1);
  });

  it("switches to offline on the browser's offline event and re-checks when back online", async () => {
    const fetchMock = respondAfter(10);
    vi.stubGlobal("fetch", fetchMock);

    const { result } = renderHook(() => useConnectionStatus(URL));
    await advance(10);
    expect(result.current).toBe("online");

    await act(async () => {
      globalThis.dispatchEvent(new Event("offline"));
    });
    expect(result.current).toBe("offline");

    await act(async () => {
      globalThis.dispatchEvent(new Event("online"));
    });
    await advance(10);

    expect(fetchMock).toHaveBeenCalledTimes(2);
    expect(result.current).toBe("online");
  });

  it("re-checks when the page becomes visible again", async () => {
    const fetchMock = respondAfter(10);
    vi.stubGlobal("fetch", fetchMock);

    renderHook(() => useConnectionStatus(URL));
    await advance(10);

    Object.defineProperty(document, "visibilityState", {
      configurable: true,
      get: () => "visible",
    });
    await act(async () => {
      document.dispatchEvent(new Event("visibilitychange"));
    });
    await advance(10);

    expect(fetchMock).toHaveBeenCalledTimes(2);
  });

  it("stops checking after unmount", async () => {
    const fetchMock = respondAfter(10);
    vi.stubGlobal("fetch", fetchMock);

    const { unmount } = renderHook(() => useConnectionStatus(URL));
    await advance(10);
    unmount();

    await advance(DEFAULT_INTERVAL_MS * 2);
    expect(fetchMock).toHaveBeenCalledTimes(1);
  });
});
