import React from "react";
import {
  Modal,
  ModalBody,
  Spinner,
  Table,
  TableRow,
  TableCell,
  TableBody,
} from "flowbite-react";
import { HiOutlineExclamationCircle } from "react-icons/hi";
import { Purchase } from "../../../../utils/api.schemas";
import { useConfig } from "@core/config/hooks/useConfig";
import Button from "../../../../components/Button";

interface RefundModalProps {
  show: boolean;
  purchase: Purchase | null;
  processing: boolean;
  onClose: () => void;
  onConfirm: (purchase: Purchase) => void;
}

export const RefundModal: React.FC<RefundModalProps> = ({
  show,
  purchase,
  processing,
  onClose,
  onConfirm,
}) => {
  const { currency, paymentMethods } = useConfig();

  const getPaymentMethodName = (code?: string) => {
    return (
      paymentMethods.find((m) => m.code === code)?.name || code || "Unknown"
    );
  };

  if (!purchase) return null;

  return (
    <Modal show={show} size="md" onClose={onClose} popup dismissible>
      <ModalBody>
        <div className="text-center m-5">
          <HiOutlineExclamationCircle className="mx-auto mb-4 h-14 w-14 text-gray-400 dark:text-gray-200" />
          <h3 className="mb-5 text-lg font-normal text-gray-500 dark:text-gray-200">
            Are you sure you want to refund this purchase?
          </h3>

          <Table className="mb-5">
            <TableBody className="dark:text-gray-100">
              {purchase.purchaseItems?.map((item) => (
                <TableRow key={item.id}>
                  <TableCell>{item.quantity} x</TableCell>
                  <TableCell>{item.product.name}</TableCell>
                  <TableCell className="text-right">
                    {currency.format(item.totalGrossPrice.toNumber())}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>

          <div className="mx-auto mb-4 text-gray-400 dark:text-gray-200 text-sm">
            Payment method:{" "}
            <b>{getPaymentMethodName(purchase.paymentMethod)}</b>
          </div>

          <div className="flex justify-center gap-4">
            <Button
              disabled={processing}
              onClick={() => onConfirm(purchase)}
              aria-label="Confirm refund"
            >
              Yes, I&apos;m sure{" "}
              {processing && <Spinner size="sm" className="ml-2" />}
            </Button>
            <Button
              color="black"
              onClick={onClose}
              aria-label="Cancel refund"
              className="bg-gray-200 dark:bg-gray-200 dark:text-gray-800"
            >
              No, cancel
            </Button>
          </div>
        </div>
      </ModalBody>
    </Modal>
  );
};
