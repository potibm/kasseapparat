import { useEffect, useRef, useState, useCallback } from "react";
import { useAuth } from "@pos/features/auth/providers/auth-provider";
import { useConfig } from "@core/config/providers/ConfigProvider";
import {
  UsePaymentWebSocketReturn,
  PaymentStatus,
} from "../types/payment.types";

export const usePaymentWebSocket = (
  purchaseId: string,
): UsePaymentWebSocketReturn => {
  const [status, setStatus] = useState<PaymentStatus>("pending");
  const [error, setError] = useState<string | null>(null);
  const [lastMessageAt, setLastMessageAt] = useState<number>(Date.now());
  const [isConnected, setIsConnected] = useState(false);

  const wsRef = useRef<WebSocket | null>(null);
  const { getToken } = useAuth();
  const { websocketHost } = useConfig();

  /**
   * Send message over websocket, when connection is open.
   */
  const sendMessage = useCallback((type: string, payload: any) => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      const message = JSON.stringify({ type, ...payload });
      console.log("WS Sending:", message);
      wsRef.current.send(message);
    } else {
      console.warn(
        "WS: Attempted to send message while connection was not open.",
      );
    }
  }, []);

  /**
   * Cancel the payment by sending a cancel message to the server. The server should respond with a cancel_ack which we handle in onmessage.
   */
  const cancel = useCallback(
    (readerId: string | undefined) => {
      sendMessage("cancel_payment", { reader_id: readerId });
    },
    [sendMessage],
  );

  useEffect(() => {
    let isMounted = true;
    let connectionTimeout: ReturnType<typeof setTimeout> | undefined;

    const initialize = async () => {
      try {
        const token = await getToken();
        if (!token) throw new Error("No authentication token available");

        const wsUrl = `${websocketHost}/api/v2/purchases/${purchaseId}/ws`;
        const ws = new WebSocket(wsUrl, [token]);
        wsRef.current = ws;

        // security timeout
        connectionTimeout = globalThis.setTimeout(() => {
          if (ws.readyState !== WebSocket.OPEN && isMounted) {
            console.error("WS: Connection timeout");
            setError("Could not reach the payment server.");
            setStatus("timeout");
          }
        }, 5000);

        ws.onopen = () => {
          if (!isMounted) return;
          console.log(`WS: Connected to purchase ${purchaseId}`);
          setIsConnected(true);
          clearTimeout(connectionTimeout);
          setError(null);
        };

        ws.onmessage = (event) => {
          if (!isMounted) return;
          try {
            const data = JSON.parse(event.data);
            setLastMessageAt(Date.now());

            // central logic for handling different message types from the server
            if (data.type === "status_update") {
              if (data.status === "confirmed") setStatus("confirmed");
              else if (data.status === "failed") setStatus("failed");
              else console.log("WS: Status Update:", data.status);
            } else if (data.type === "cancel_ack") {
              setStatus("cancelled");
            }
          } catch (e) {
            console.error("WS: Message parse error", e);
          }
        };

        ws.onclose = (event) => {
          if (!isMounted) return;
          setIsConnected(false);
          console.log("WS: Connection closed", event.code);

          // when still pending, a close usually indicates an error or network issue
          if (status === "pending" && !event.wasClean) {
            setStatus("connection_lost");
          }
        };

        ws.onerror = (err) => {
          if (!isMounted) return;
          console.error("WS: Error event", err);
          setError("A connection error occurred.");
        };
      } catch (err: any) {
        if (isMounted) {
          setError(err.message || "Failed to initialize WebSocket");
          setStatus("failed");
        }
      }
    };

    initialize();

    // Cleanup on unmount or id-change
    return () => {
      isMounted = false;
      clearTimeout(connectionTimeout);
      if (wsRef.current) {
        wsRef.current.close();
        wsRef.current = null;
      }
    };
  }, [purchaseId, getToken, websocketHost]);

  return {
    status,
    error,
    lastMessageAt,
    cancel,
    isConnected,
    sendMessage,
  };
};

export default usePaymentWebSocket;
