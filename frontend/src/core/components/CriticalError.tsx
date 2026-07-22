import React from "react";
import { Alert } from "flowbite-react";

interface CriticalErrorProps {
  title: string;
  message: string;
  details?: string;
}

export const CriticalError: React.FC<CriticalErrorProps> = ({
  title,
  message,
  details,
}) => {
  return (
    <div className="flex h-screen items-center justify-center p-4">
      <div className="max-w-2xl">
        <Alert color="failure">
          <h2 className="text-lg font-semibold mb-2">{title}</h2>
          <p className="mb-2">{message}</p>
          {details && (
            <pre className="text-sm bg-red-50 p-2 rounded overflow-auto">
              {details}
            </pre>
          )}
        </Alert>
      </div>
    </div>
  );
};
