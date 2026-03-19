import { z } from "zod";
import { createLogger } from "@core/logger/logger";

const log = createLogger("Config");

export const COOKIE_NAME = "kasseapparat.sumup.reader-id";

const READER_ID_REGEX = /^rdr_[A-Z0-9]{26}$/;
const ReaderIdSchema = z
  .string()
  .regex(
    READER_ID_REGEX,
    "Invalid Reader ID format. Expected 'rdr_' followed by 26 characters.",
  );

export function getCurrentReaderId(): string | undefined {
  const cookies = document.cookie.split(";").map((c) => c.trim());
  const match = cookies.find((c) => c.startsWith(COOKIE_NAME + "="));
  if (!match) return undefined;

  let rawValue: string;
  try {
    rawValue = decodeURIComponent(match.slice(COOKIE_NAME.length + 1));
  } catch {
    return undefined;
  }

  const result = ReaderIdSchema.safeParse(rawValue);
  return result.success ? result.data : undefined;
}

export function setCurrentReaderId(readerId: string): void {
  const result = ReaderIdSchema.safeParse(readerId);

  if (!result.success) {
    log.error("Invalid Reader ID provided", result.error.format);
    return;
  }
  document.cookie = `${COOKIE_NAME}=${encodeURIComponent(readerId)}; path=/; SameSite=Lax`;
}

export function clearCurrentReaderId(): void {
  document.cookie = `${COOKIE_NAME}=; Max-Age=0; path=/; SameSite=Lax`;
}
