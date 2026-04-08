export type ToastSeverity = "info" | "warning" | "error" | "success" | "debug";

export interface ToastData {
  id: string;
  severity: ToastSeverity;
  message: string;
  autoClose?: boolean;
  duration?: number;
  blocking?: boolean;
}
