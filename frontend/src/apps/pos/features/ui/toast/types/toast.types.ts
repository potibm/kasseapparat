export type ToastSeverity = "info" | "warning" | "error" | "success" | "debug";

export interface ToastAction {
  label: string;
  onClick: () => void;
}

export interface ToastData {
  id: string;
  severity: ToastSeverity;
  message: string;
  autoClose?: boolean;
  duration?: number;
  blocking?: boolean;
  /** Optional button next to the message, e.g. "Undo". Clicking it also dismisses the toast. */
  action?: ToastAction;
}
