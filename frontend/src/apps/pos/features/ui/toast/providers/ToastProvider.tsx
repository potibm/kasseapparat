import React, { useState, useCallback, useMemo } from "react";
import { ToastData } from "../types/toast.types";
import { ToastItem } from "../components/_internal/ToastItem";
import { ToastContext } from "../context/ToastContext";

export const ToastProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [toasts, setToasts] = useState<ToastData[]>([]);

  const isDev = import.meta.env.VITE_ENV === "development";

  const hasBlockingToast = toasts.some((t) => t.blocking);

  const removeToast = useCallback((id: string) => {
    setToasts((prev) => prev.filter((t) => t.id !== id));
  }, []);

  const showToast = useCallback(
    (newToast: Omit<ToastData, "id">) => {
      if (newToast.severity === "debug" && !isDev) return;

      const id = Math.random().toString(36).substring(2, 9);

      const resolvedAutoClose = newToast.blocking ? false : newToast.autoClose;

      setToasts((prev) => [
        ...prev,
        { ...newToast, id, autoClose: resolvedAutoClose },
      ]);
    },
    [isDev],
  );

  const contextValue = useMemo(() => ({ showToast }), [showToast]);

  return (
    <ToastContext value={contextValue}>
      {children}
      {hasBlockingToast && (
        <div className="fixed inset-0 z-9998 bg-black/50 backdrop-blur-sm transition-opacity" />
      )}

      {/* Bottom centre of the product area: the cart, checkout and purchase
          history on the right stay visible and usable. */}
      <div
        data-testid="toast-container"
        className="pointer-events-none fixed right-[380px] bottom-5 left-5 z-9999 flex flex-col items-center gap-3"
      >
        {toasts.map((t) => (
          <div key={t.id} className="pointer-events-auto w-80">
            <ToastItem toast={t} onDismiss={removeToast} />
          </div>
        ))}
      </div>
    </ToastContext>
  );
};
