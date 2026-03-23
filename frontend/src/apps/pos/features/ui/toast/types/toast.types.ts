export type ToastType = "info" | "warning" | "error" | "success" | "debug";

export interface Toast {
  id: string;
  type: ToastType;
  message: string;
  autoClose?: boolean;
  duration?: number;
}
