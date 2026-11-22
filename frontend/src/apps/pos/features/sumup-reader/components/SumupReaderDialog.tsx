import React from "react";
import { Modal, ModalBody, ModalHeader, Spinner } from "flowbite-react";
import { HiCheck, HiRefresh } from "react-icons/hi";
import { useSumupReader } from "../hooks/useSumupReader";

const statusClass = (status: string) =>
  status === "paired" || status === "processing"
    ? "bg-green-100 text-green-800 dark:bg-green-900 dark:text-green-100"
    : "bg-red-100 text-red-800 dark:bg-red-900 dark:text-red-100";

/**
 * Lets the cashier pick the card reader for this iPad without going to the admin.
 * The choice is stored in this browser only, so every till keeps its own reader.
 */
const SumupReaderDialog: React.FC = () => {
  const {
    readers,
    loading,
    selectedReaderId,
    selectReader,
    refresh,
    closeDialog,
  } = useSumupReader();

  const handleSelect = (id: string) => {
    selectReader(id);
    closeDialog();
  };

  return (
    <Modal show onClose={closeDialog} dismissible size="md">
      <ModalHeader>Card reader for this till</ModalHeader>
      <ModalBody>
        <div className="flex flex-col gap-3 dark:text-white">
          {readers === null && !loading && (
            <p className="text-red-700 dark:text-red-300">
              The card readers could not be loaded.
            </p>
          )}
          {readers?.length === 0 && (
            <p className="text-gray-700 dark:text-gray-300">
              No card reader is paired yet. Pair one in the admin under
              “SumUp Readers”.
            </p>
          )}
          {readers?.map((reader) => {
            const selected = reader.id === selectedReaderId;

            return (
              <button
                key={reader.id}
                type="button"
                data-testid={"sumup-reader-" + reader.id}
                aria-pressed={selected}
                onClick={() => handleSelect(reader.id)}
                className={`pos-pressable flex min-h-14 items-center gap-3 rounded-xl border-2 px-4 text-left ${
                  selected
                    ? "border-cyan-700 bg-cyan-50 dark:border-cyan-500 dark:bg-cyan-950"
                    : "border-gray-200 dark:border-gray-700"
                }`}
              >
                <span className="flex-1">
                  <span className="block text-lg font-semibold">
                    {reader.name || reader.id}
                  </span>
                  <span className="block font-mono text-xs text-gray-500 dark:text-gray-400">
                    {reader.id}
                  </span>
                </span>
                <span
                  className={`rounded-full px-2 py-0.5 text-sm font-semibold ${statusClass(reader.status)}`}
                >
                  {reader.status}
                </span>
                {selected && (
                  <HiCheck
                    className="h-6 w-6 text-cyan-700 dark:text-cyan-300"
                    aria-label="Selected"
                  />
                )}
              </button>
            );
          })}
          <button
            type="button"
            onClick={refresh}
            disabled={loading}
            className="flex h-11 items-center justify-center gap-2 text-sm font-semibold text-cyan-800 disabled:opacity-60 dark:text-cyan-300"
          >
            {loading ? (
              <Spinner size="sm" />
            ) : (
              <HiRefresh className="h-5 w-5" aria-hidden="true" />
            )}
            Refresh
          </button>
        </div>
      </ModalBody>
    </Modal>
  );
};

export default SumupReaderDialog;
