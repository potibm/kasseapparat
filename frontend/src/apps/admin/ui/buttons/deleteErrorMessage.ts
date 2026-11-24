import { HttpError } from "react-admin";

/**
 * Turns a failed delete into a message for the admin. The backend refuses to delete a
 * section or product group that still has content with 409, which react-admin would
 * only show as "Conflict" (ADR 006).
 */
export const deleteErrorMessage = (
  error: unknown,
  conflictMessage: string,
  fallbackMessage: string,
): string => {
  if (error instanceof HttpError) {
    if (error.status === 409) return conflictMessage;

    const details = (error.body as { details?: unknown } | undefined)?.details;
    if (typeof details === "string" && details !== "") return details;
  }

  if (error instanceof Error && error.message !== "") return error.message;

  return fallbackMessage;
};
