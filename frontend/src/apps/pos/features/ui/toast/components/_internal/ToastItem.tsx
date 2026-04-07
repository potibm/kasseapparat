import { useEffect } from "react";
import { Toast, ToastToggle } from "flowbite-react";
import {
  HiCheck,
  HiExclamation,
  HiX,
  HiCode,
  HiInformationCircle,
} from "react-icons/hi";
import { IconType } from "react-icons";
import {
  Toast as ToastType,
  ToastType as ToastTypeType,
} from "../../types/toast.types";

export const ToastItem = ({
  toast,
  onDismiss,
}: {
  toast: ToastType;
  onDismiss: (id: string) => void;
}) => {
  const styles: Record<ToastTypeType, { icon: IconType; color: string }> = {
    success: { icon: HiCheck, color: "text-green-500 bg-green-100" },
    error: { icon: HiX, color: "text-red-500 bg-red-100" },
    warning: { icon: HiExclamation, color: "text-yellow-500 bg-yellow-100" },
    info: { icon: HiInformationCircle, color: "text-blue-500 bg-blue-100" },
    debug: {
      icon: HiCode,
      color: "text-purple-500 bg-purple-100 dark:bg-purple-900",
    },
  };

  useEffect(() => {
    if (toast.autoClose === false) return;

    const timer = setTimeout(() => {
      onDismiss(toast.id);
    }, toast.duration || 10000);

    return () => clearTimeout(timer);
  }, [toast, onDismiss]);

  const { icon: Icon, color } = styles[toast.type];

  return (
    <Toast className="max-w-xs border border-gray-200 shadow-lg mb-2">
      <div
        className={`inline-flex h-8 w-8 shrink-0 items-center justify-center rounded-lg ${color}`}
      >
        <Icon className="h-5 w-5" />
      </div>
      <div className="ml-3 text-sm font-normal">
        {toast.type === "debug" && (
          <span className="font-bold mr-1">[DEBUG]</span>
        )}
        {toast.message}
      </div>
      <ToastToggle onDismiss={() => onDismiss(toast.id)} />
    </Toast>
  );
};
