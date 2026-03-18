export type PaymentStatus =
  | "pending"
  | "confirmed"
  | "failed"
  | "cancelled"
  | "connection_lost"
  | "timeout";

export type PaymentCommand = "cancel_payment";

interface PaymentMessageCancelPayload {
  reader_id: string;
}

export type PaymentMessagePayload = PaymentMessageCancelPayload;

export interface UsePaymentWebSocketReturn {
  status: PaymentStatus;
  error: string | null;
  lastMessageAt: number;
  cancel: (readerId: string | undefined) => void;
  isConnected: boolean;
  sendMessage: (type: PaymentCommand, payload: PaymentMessagePayload) => void;
}
