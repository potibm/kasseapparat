import { use } from "react";
import { ToastContext } from "./../context/ToastContext";

export const useToast = () => {
  const context = use(ToastContext);

  if (context === undefined) {
    throw new Error("useToast must be used within a ToastProvider");
  }

  return context;
};
