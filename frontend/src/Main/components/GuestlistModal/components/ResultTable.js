import React from "react";
import { Table, Alert, Spinner } from "flowbite-react";
import { HiInformationCircle, HiXCircle } from "react-icons/hi";
import PropTypes from "prop-types";
import GuestListResultTableRow from "./ResultTableRow";

const GuestListResultTable = ({
  loading,
  error,
  guestListEntries,
  onAddToCart,
  hasListItem,
  loadedSearchQuery,
}) => {
  return (
    <>
      {loading && (
        <div className="absolute inset-0 flex items-center justify-center bg-white bg-opacity-75 z-10">
          <Spinner size="xl" />
        </div>
      )}

      {error && (
        <Alert className="my-3" color="failure" icon={HiInformationCircle}>
          {error}
        </Alert>
      )}

      {!loading && guestListEntries.length === 0 && (
        <Alert className="my-3" color="warning" icon={HiXCircle}>
          No entries found
        </Alert>
      )}

      {guestListEntries.length > 0 && (
        <div className="space-y-4">
          <Table hoverable>
            <Table.Head>
              <Table.HeadCell className="w-1/12"></Table.HeadCell>
              <Table.HeadCell className="w-5/12">Name</Table.HeadCell>
              <Table.HeadCell className="w-6/12">Action</Table.HeadCell>
            </Table.Head>
            <Table.Body className="divide-y">
              {guestListEntries.map((entry) => (
                <GuestListResultTableRow
                  key={entry.id}
                  entry={entry}
                  onAddToCart={onAddToCart}
                  hasListItem={hasListItem}
                  loadedSearchQuery={loadedSearchQuery}
                />
              ))}
            </Table.Body>
          </Table>
        </div>
      )}
    </>
  );
};

GuestListResultTable.propTypes = {
  loading: PropTypes.bool,
  error: PropTypes.string,
  guestListEntries: PropTypes.array.isRequired,
  onAddToCart: PropTypes.func.isRequired,
  hasListItem: PropTypes.func.isRequired,
  loadedSearchQuery: PropTypes.string.isRequired,
};

export default GuestListResultTable;
