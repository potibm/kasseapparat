import React from "react";
import {
  Table,
  Alert,
  Spinner,
  TableHead,
  TableHeadCell,
  TableBody,
  TableRow,
} from "flowbite-react";
import { HiInformationCircle, HiXCircle } from "react-icons/hi";
import GuestlistResultTableRow from "./GuestlistResultTableRow";
import { Guest as GuestType } from "@pos/utils/api.schemas";

interface GuestlistResultTableProps {
  loading: boolean;
  error: string | null;
  guestlistEntries: GuestType[];
  onAddToCart: (listEntry: GuestType, additionalGuests: number) => void;
  hasListItem: (guest: GuestType) => boolean;
  loadedSearchQuery: string;
}

const GuestlistResultTable: React.FC<GuestlistResultTableProps> = ({
  loading,
  error,
  guestlistEntries,
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

      {!loading && guestlistEntries.length === 0 && (
        <Alert className="my-3" color="warning" icon={HiXCircle}>
          No matching guests in this guestlist
        </Alert>
      )}

      {guestlistEntries.length > 0 && (
        <div className="space-y-4">
          <Table hoverable className="dark:text-white">
            <TableHead>
              <TableRow>
                <TableHeadCell className="w-1/12"></TableHeadCell>
                <TableHeadCell className="w-5/12">Name</TableHeadCell>
                <TableHeadCell className="w-6/12">Action</TableHeadCell>
              </TableRow>
            </TableHead>
            <TableBody className="divide-y">
              {guestlistEntries.map((entry) => (
                <GuestlistResultTableRow
                  key={entry.id}
                  entry={entry}
                  onAddToCart={onAddToCart}
                  hasListItem={hasListItem}
                  loadedSearchQuery={loadedSearchQuery}
                />
              ))}
            </TableBody>
          </Table>
        </div>
      )}
    </>
  );
};

export default GuestlistResultTable;
