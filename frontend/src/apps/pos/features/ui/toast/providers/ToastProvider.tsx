import React, { createContext, useState, useCallback } from "react";
import { Toast as ToastType } from "../types/toast.types";
import { ToastItem } from "../components/_internal/ToastItem";

interface ToastContextType {
  showToast: (toast: Omit<ToastType, "id">) => void;
}

export const ToastContext = createContext<ToastContextType | undefined>(
  undefined,
);

export const ToastProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [toasts, setToasts] = useState<ToastType[]>([]);

  const isDev = process.env.NODE_ENV === "development";

  const showToast = useCallback(
    (newToast: Omit<ToastType, "id">) => {
      if (newToast.type === "debug" && !isDev) return;

      const id = Math.random().toString(36).substring(2, 9);
      setToasts((prev) => [...prev, { ...newToast, id }]);

      if (newToast.autoClose !== false) {
        setTimeout(() => {
          setToasts((prev) => prev.filter((t) => t.id !== id));
        }, newToast.duration || 10000);
      }
    },
    [isDev],
  );

  return (
    <ToastContext.Provider value={{ showToast }}>
      {children}
      <div className="fixed bottom-5 right-5 z-9999 flex flex-col gap-3 w-80">
        {toasts.map((t) => (
          <ToastItem
            key={t.id}
            toast={t}
            onDismiss={(id: string) =>
              setToasts((prev) => prev.filter((toast) => toast.id !== id))
            }
          />
        ))}
      </div>
    </ToastContext.Provider>
  );
};
