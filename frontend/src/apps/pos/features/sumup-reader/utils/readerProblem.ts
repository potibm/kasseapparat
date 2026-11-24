import { SumupReader } from "@pos/api/schemas";

/** Statuses in which a reader can take a payment. */
const USABLE_STATUSES = new Set(["paired", "processing"]);

export type ReaderProblem =
  | { kind: "none-selected" }
  | { kind: "missing"; readerId: string }
  | { kind: "unusable"; reader: SumupReader };

/**
 * What keeps card payments from working on this device, or null when nothing does.
 * While the reader list is unknown (still loading, or SumUp unreachable) only a
 * missing selection counts: the payment itself will report other problems.
 */
export const findReaderProblem = (
  selectedReaderId: string | undefined,
  readers: SumupReader[] | null,
): ReaderProblem | null => {
  if (!selectedReaderId) return { kind: "none-selected" };
  if (readers === null) return null;

  const reader = readers.find((candidate) => candidate.id === selectedReaderId);
  if (!reader) return { kind: "missing", readerId: selectedReaderId };
  if (!USABLE_STATUSES.has(reader.status)) return { kind: "unusable", reader };

  return null;
};

export const describeReaderProblem = (problem: ReaderProblem): string => {
  switch (problem.kind) {
    case "none-selected":
      return "No card reader selected";
    case "missing":
      return "The selected card reader is no longer paired";
    case "unusable":
      return `Card reader „${problem.reader.name}“ is ${problem.reader.status}`;
  }
};
