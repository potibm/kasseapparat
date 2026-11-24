import { deleteErrorMessage } from "@admin/ui/buttons/deleteErrorMessage";

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
  "This section still contains products or product groups. Move them to another section first.";

export const SECTION_DELETE_FAILED_MESSAGE =
  "The section could not be deleted.";

/** Turns a failed section delete into a message for the admin. */
export const sectionDeleteErrorMessage = (error: unknown): string =>
  deleteErrorMessage(
    error,
    SECTION_NOT_EMPTY_MESSAGE,
    SECTION_DELETE_FAILED_MESSAGE,
  );
