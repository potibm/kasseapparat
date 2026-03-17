import { z } from "zod";

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

  const rawValue = decodeURIComponent(match.split("=")[1]);

  const result = ReaderIdSchema.safeParse(rawValue);
  return result.success ? result.data : undefined;
}

export function setCurrentReaderId(readerId: string): void {
  const result = ReaderIdSchema.safeParse(readerId);

  if (!result.success) {
    console.error("Invalid Reader ID provided:", result.error.format);
    return;
  }
  document.cookie = `${COOKIE_NAME}=${encodeURIComponent(readerId)}; path=/; SameSite=Lax`;
}

export function clearCurrentReaderId(): void {
  document.cookie = `${COOKIE_NAME}=; Max-Age=0; path=/; SameSite=Lax`;
}
