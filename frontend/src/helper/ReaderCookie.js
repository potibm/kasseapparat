export const COOKIE_NAME = "sumup-reader-id";

export function getCurrentReaderId() {
  const cookies = document.cookie.split(";").map((c) => c.trim());
  const match = cookies.find((c) => c.startsWith(COOKIE_NAME + "="));
  return match ? decodeURIComponent(match.split("=")[1]) : undefined;
}

export function setCurrentReaderId(readerId) {
  document.cookie = `${COOKIE_NAME}=${encodeURIComponent(readerId)}; path=/; SameSite=Lax`;
}

export function clearCurrentReaderId() {
  document.cookie = `${COOKIE_NAME}=; Max-Age=0; path=/; SameSite=Lax`;
}
