import React from "react";
import { HiExclamation } from "react-icons/hi";
import { useSumupReader } from "../hooks/useSumupReader";
import { describeReaderProblem } from "../utils/readerProblem";

/**
 * One line under the checkout buttons, only when card payments cannot work. Tapping
 * it opens the reader selection. While everything is fine nothing is shown.
 */
const SumupReaderHint: React.FC = () => {
  const { enabled, problem, openDialog } = useSumupReader();

  if (!enabled || !problem) return null;

  return (
    <button
      type="button"
      data-testid="sumup-reader-hint"
      onClick={openDialog}
      className="col-span-2 flex min-h-11 items-center gap-2 rounded-lg bg-amber-50 px-3 text-left text-sm font-medium text-amber-900 dark:bg-amber-950 dark:text-amber-100"
    >
      <HiExclamation className="h-5 w-5 shrink-0" aria-hidden="true" />
      <span className="flex-1">{describeReaderProblem(problem)}</span>
      <span className="font-semibold underline">Choose</span>
    </button>
  );
};

export default SumupReaderHint;
