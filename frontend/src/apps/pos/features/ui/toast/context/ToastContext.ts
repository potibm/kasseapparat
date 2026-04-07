import { createContext } from "react";
import { ToastData as ToastType } from "../types/toast.types";

interface ToastContextType {
  showToast: (toast: Omit<ToastType, "id">) => void;
}

export const ToastContext = createContext<ToastContextType | undefined>(
  undefined,
);

export default ToastContext;
