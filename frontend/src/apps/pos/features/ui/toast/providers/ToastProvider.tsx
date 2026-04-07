import React, { useState, useCallback } from "react";
import { Toast as ToastType } from "../types/toast.types";
import { ToastItem } from "../components/_internal/ToastItem";
import { ToastContext } from "../context/ToastContext";

export const ToastProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [toasts, setToasts] = useState<ToastType[]>([]);

  const isDev = import.meta.env.VITE_ENV === "development";

  const hasBlockingToast = toasts.some((t) => t.blocking);

  const removeToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const showToast = useCallback(
    (newToast: Omit<ToastType, "id">) => {
      if (newToast.type === "debug" && !isDev) return;

      const id = Math.random().toString(36).substring(2, 9);
      setToasts((prev) => [...prev, { ...newToast, id }]);
    },
    [isDev],
  );

  return (
    <ToastContext.Provider value={{ showToast }}>
      {children}
      {hasBlockingToast && (
        <div className="fixed inset-0 z-9998 bg-black/50 backdrop-blur-sm transition-opacity" />
      )}

      <div className="fixed bottom-5 right-5 z-9999 flex flex-col gap-3 w-80">
        {toasts.map((t) => (
          <ToastItem key={t.id} toast={t} onDismiss={removeToast} />
        ))}
      </div>
    </ToastContext.Provider>
  );
};
