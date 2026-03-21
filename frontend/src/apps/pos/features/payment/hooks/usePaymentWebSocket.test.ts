import { describe, it, expect, beforeEach, vi, afterEach } from "vitest";
import { renderHook, act, waitFor } from "@testing-library/react";
import { usePaymentWebSocket } from "./usePaymentWebSocket";
import { useAuth } from "@pos/features/auth/providers/AuthProvider";
import { useConfig } from "@core/config/providers/ConfigProvider";

// dependency mocks
vi.mock("@pos/features/auth/providers/AuthProvider", () => ({
  useAuth: vi.fn(),
}));

vi.mock("@core/config/providers/ConfigProvider", () => ({
  useConfig: vi.fn(),
}));

vi.mock("@core/logger/logger", () => ({
  createLogger: () => ({
    info: vi.fn(),
    warn: vi.fn(),
    error: vi.fn(),
    debug: vi.fn(),
  }),
}));

// websocket mock
class MockWebSocket {
  url: string;
  protocols: string | string[];
  readyState: number = 0; // 0 = CONNECTING, 1 = OPEN

  send = vi.fn();
  close = vi.fn();

  onopen: (() => void) | null = null;
  onmessage: ((event: { data: string }) => void) | null = null;
  onclose: ((event: { wasClean: boolean; code: number }) => void) | null = null;
  onerror: ((event: unknown) => void) | null = null;

  constructor(url: string, protocols: string | string[]) {
    this.url = url;
    this.protocols = protocols;
  }
}

describe("usePaymentWebSocket Hook", () => {
  let mockWebSocketInstance: MockWebSocket;

  beforeEach(() => {
    vi.clearAllMocks();

    vi.mocked(useAuth).mockReturnValue({
      getToken: vi.fn(async () => "fake-ws-token"),
    } as unknown as ReturnType<typeof useAuth>);
    vi.mocked(useConfig).mockReturnValue({
      websocketHost: "wss://test.com",
    } as unknown as ReturnType<typeof useConfig>);

    const MockWSConstructor = Object.assign(
      vi.fn(function (url: string, protocols: string | string[]) {
        mockWebSocketInstance = new MockWebSocket(url, protocols);
        return mockWebSocketInstance;
      }),
      {
        CONNECTING: 0,
        OPEN: 1,
        CLOSING: 2,
        CLOSED: 3,
      },
    );

    vi.stubGlobal("WebSocket", MockWSConstructor);
  });

  afterEach(() => {
    vi.unstubAllGlobals();
    mockWebSocketInstance = undefined as unknown as MockWebSocket;
  });

  describe("Initialization & Connection", () => {
    it("should initialize and connect successfully", async () => {
      const { result } = renderHook(() => usePaymentWebSocket("purchase-123"));

      await waitFor(() => {
        expect(mockWebSocketInstance).toBeDefined();
      });

      expect(globalThis.WebSocket).toHaveBeenCalledWith(
        "wss://test.com/api/v2/purchases/purchase-123/ws",
        ["fake-ws-token"],
      );

      act(() => {
        mockWebSocketInstance.readyState = 1; // OPEN
        mockWebSocketInstance.onopen?.();
      });

      expect(result.current.isConnected).toBe(true);
      expect(result.current.status).toBe("pending");
      expect(result.current.error).toBeNull();
    });

    it("should handle connection timeouts", async () => {
      vi.useFakeTimers();
      const { result } = renderHook(() => usePaymentWebSocket("purchase-123"));

      // magic time travel to trigger the timeout
      await act(async () => {
        await vi.advanceTimersByTimeAsync(5000);
      });

      expect(result.current.status).toBe("timeout");
      expect(result.current.error).toBe("Could not reach the payment server.");
      expect(result.current.isConnected).toBe(false);

      vi.useRealTimers();
    });
  });

  describe("Message Handling", () => {
    it("should handle a successful payment confirmation", async () => {
      const { result } = renderHook(() => usePaymentWebSocket("purchase-123"));

      await waitFor(() => {
        expect(mockWebSocketInstance).toBeDefined();
      });

      act(() => {
        mockWebSocketInstance.onmessage?.({
          data: JSON.stringify({ type: "status_update", status: "confirmed" }),
        });
      });

      expect(result.current.status).toBe("confirmed");
    });

    it("should handle cancel acknowledgement", async () => {
      const { result } = renderHook(() => usePaymentWebSocket("purchase-123"));

      await waitFor(() => {
        expect(mockWebSocketInstance).toBeDefined();
      });

      act(() => {
        mockWebSocketInstance.onmessage?.({
          data: JSON.stringify({ type: "cancel_ack" }),
        });
      });

      expect(result.current.status).toBe("cancelled");
    });
  });

  describe("Sending Commands", () => {
    it("should send a cancel command if readerId is provided", async () => {
      const { result } = renderHook(() => usePaymentWebSocket("purchase-123"));

      await waitFor(() => {
        expect(mockWebSocketInstance).toBeDefined();
      });

      act(() => {
        mockWebSocketInstance.readyState = 1; // WebSocket.OPEN
        mockWebSocketInstance.onopen?.();
      });

      act(() => {
        result.current.cancel("reader-999");
      });

      expect(mockWebSocketInstance.send).toHaveBeenCalledWith(
        JSON.stringify({ type: "cancel_payment", reader_id: "reader-999" }),
      );
    });

    it("should NOT send anything if the connection is not open", async () => {
      const { result } = renderHook(() => usePaymentWebSocket("purchase-123"));

      await waitFor(() => {
        expect(mockWebSocketInstance).toBeDefined();
      });

      act(() => {
        result.current.cancel("reader-999");
      });

      expect(mockWebSocketInstance.send).not.toHaveBeenCalled();
    });
  });

  describe("Cleanup & Disconnect", () => {
    it("should correctly handle unexpected disconnects", async () => {
      const { result } = renderHook(() => usePaymentWebSocket("purchase-123"));

      await waitFor(() => {
        expect(mockWebSocketInstance).toBeDefined();
      });

      act(() => {
        mockWebSocketInstance.onclose?.({ wasClean: false, code: 1006 });
      });

      expect(result.current.isConnected).toBe(false);
      expect(result.current.status).toBe("connection_lost");
    });

    it("should close the connection when the component unmounts", async () => {
      const { unmount } = renderHook(() => usePaymentWebSocket("purchase-123"));

      await waitFor(() => {
        expect(mockWebSocketInstance).toBeDefined();
      });

      unmount();

      expect(mockWebSocketInstance.close).toHaveBeenCalledTimes(1);
      expect(mockWebSocketInstance.onclose).toBeNull();
    });
  });
});
