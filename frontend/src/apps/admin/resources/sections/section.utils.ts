import { HttpError } from "react-admin";

export const UNNAMED_SECTION_LABEL = "Row break";

/**
 * How a section is named in lists and selects. An unnamed section is shown as a
 * plain row break in the POS (ADR 006), so it is labelled as one here as well.
 */
export const sectionLabel = (
  section: { id?: number | string; name?: string } | undefined | null,
): string => {
  if (!section) return "";

  const name = section.name?.trim();
  if (name) return name;

  return `${UNNAMED_SECTION_LABEL} #${section.id}`;
};

export const SECTION_NOT_EMPTY_MESSAGE =
  "This section still contains products. Move them to another section first.";

/**
 * Turns a failed section delete into a message for the admin. The backend answers a
 * section that still has products with 409, which react-admin would only show as
 * "Conflict".
 */
export const sectionDeleteErrorMessage = (error: unknown): string => {
  if (error instanceof HttpError) {
    if (error.status === 409) return SECTION_NOT_EMPTY_MESSAGE;

    const details = (error.body as { details?: unknown } | undefined)?.details;
    if (typeof details === "string" && details !== "") return details;
  }

  if (error instanceof Error && error.message !== "") return error.message;

  return "The section could not be deleted.";
};
