import { useEffect, useRef, useState, useCallback } from "react";
import { useAuth } from "@pos/features/auth/hooks/useAuth";
import { useConfig } from "@core/config/hooks/useConfig";
import {
  UsePaymentWebSocketReturn,
  PaymentStatus,
  PaymentCommand,
  PaymentMessagePayload,
} from "../types/payment.types";
import { createLogger } from "@core/logger/logger";

const log = createLogger("PaymentWebsocket");

export const usePaymentWebSocket = (
  purchaseId: string,
): UsePaymentWebSocketReturn => {
  const [status, setStatus] = useState<PaymentStatus>("pending");
  const [error, setError] = useState<string | null>(null);
  const [lastMessageAt, setLastMessageAt] = useState<number>(Date.now());
  const [isConnected, setIsConnected] = useState(false);

  const statusRef = useRef<PaymentStatus>(status);

  useEffect(() => {
    statusRef.current = status;
  }, [status]);

  const wsRef = useRef<WebSocket | null>(null);
  const { getToken } = useAuth();
  const { websocketHost } = useConfig();

  /**
   * Send message over websocket, when connection is open.
   */
  const sendMessage = useCallback(
    (type: PaymentCommand, payload: PaymentMessagePayload) => {
      if (wsRef.current?.readyState === WebSocket.OPEN) {
        const message = JSON.stringify({ type, ...payload });
        log.debug("Sending", purchaseId, message);
        wsRef.current.send(message);
      } else {
        log.warn("Attempted to send message while connection was not open.");
      }
    },
    [purchaseId],
  );

  /**
   * Cancel the payment by sending a cancel message to the server. The server should respond with a cancel_ack which we handle in onmessage.
   */
  const cancel = useCallback(
    (readerId: string | undefined) => {
      if (!readerId) {
        log.warn(
          "Cannot cancel payment: no reader ID provided",
          purchaseId,
          readerId,
        );
        return;
      }
      sendMessage("cancel_payment", { reader_id: readerId });
    },
    [sendMessage, purchaseId],
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
            log.error("Connection timeout", purchaseId);
            setError("Could not reach the payment server.");
            setStatus("timeout");
          }
        }, 5000);

        ws.onopen = () => {
          if (!isMounted) return;
          log.info("Connected to purchase", purchaseId);
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
              if (data.status === "confirmed") {
                log.info("Payment confirmed", purchaseId, data);
                setStatus("confirmed");
              } else if (data.status === "failed") {
                log.warn("Payment failed", purchaseId, data);
                setStatus("failed");
              } else {
                log.debug("Status Update", purchaseId, data.status);
              }
            } else if (data.type === "cancel_ack") {
              log.info("Cancel acknowledged", purchaseId, data);
              setStatus("cancelled");
            }
          } catch (e) {
            log.error("Message parse error", purchaseId, e);
          }
        };

        ws.onclose = (event) => {
          if (!isMounted) return;
          setIsConnected(false);
          log.info("Connection closed", purchaseId, event.code);

          // when still pending, a close usually indicates an error or network issue
          if (statusRef.current === "pending" && !event.wasClean) {
            setStatus("connection_lost");
          }
        };

        ws.onerror = (err) => {
          if (!isMounted) return;
          log.error("Error event", purchaseId, err);
          setError("A connection error occurred.");
        };
      } catch (err: unknown) {
        if (isMounted) {
          log.error("Initialization failed", purchaseId, err);
          setError((err as Error).message || "Failed to initialize WebSocket");
          setStatus("failed");
        }
      }
    };

    initialize();

    // Cleanup on unmount or id-change
    return () => {
      isMounted = false;
      if (connectionTimeout) clearTimeout(connectionTimeout);
      if (wsRef.current) {
        wsRef.current.onclose = null;
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
