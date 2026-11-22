import { Guest } from "@pos/api/schemas";

/**
 * The guest a scanned ticket belongs to: exactly one result whose code equals the
 * search term (ignoring case and surrounding spaces). Works whether or not the
 * scanner sends Enter after the code – the results of the search are enough.
 */
export const findScannedGuest = (
  entries: Guest[],
  query: string,
): Guest | null => {
  const scanned = query.trim().toLowerCase();
  if (scanned === "") return null;

  const matches = entries.filter(
    (entry) => entry.code?.toLowerCase() === scanned,
  );

  return matches.length === 1 ? matches[0] : null;
};
