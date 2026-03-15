export type PaymentStatus =
  | "pending"
  | "confirmed"
  | "failed"
  | "cancelled"
  | "connection_lost"
  | "timeout";

export interface UsePaymentWebSocketReturn {
  status: PaymentStatus;
  error: string | null;
  lastMessageAt: number;
  cancel: (readerId: string | undefined) => void;
  isConnected: boolean;
  sendMessage: (type: string, payload: any) => void;
}
