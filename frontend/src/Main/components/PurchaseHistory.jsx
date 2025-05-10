import { HiXCircle, HiOutlineExclamationCircle } from "react-icons/hi";
import React, { useState, useEffect, useRef } from "react";
import {
  Modal,
  ModalBody,
  Spinner,
  Table,
  TableCell,
  TableRow,
  TableHeadCell,
  TableHead,
  TableBody,
} from "flowbite-react";
import PropTypes from "prop-types";
import { useConfig } from "../../provider/ConfigProvider";
import "animate.css";
import MyButton from "./MyButton";

function PurchaseHistory({ history, removeFromPurchaseHistory }) {
  const [openModal, setOpenModal] = useState({ show: false, purchase: null });
  const [processing, setProcessing] = useState(false);

  const paymentMethods = useConfig().paymentMethods;

  const findPaymentMethod = (code) => {
    const paymentMethod = paymentMethods.find((method) => method.code === code);
    if (paymentMethod) {
      return paymentMethod.name;
    } else {
      return code;
    }
  };

  const confirmDelete = (purchase) => {
    if (processing) {
      return;
    }
    setProcessing(true);
    removeFromPurchaseHistory(purchase).then(() => {
      setProcessing(false);
      setOpenModal({ show: false });
    });
  };

  const [flash, setFlash] = useState(false);
  const flashCount = useRef(0);

  const triggerFlash = () => {
    setFlash(true);
    setTimeout(() => {
      setFlash(false);
    }, 500);
  };

  useEffect(() => {
    // not 100% sure why this is called three times
    if (flashCount.current < 3) {
      flashCount.current++;
      return;
    }
    triggerFlash();
  }, [history]);

  const currency = useConfig().currency;
  const dateLocale = useConfig().dateLocale;
  const dateOptions = useConfig().dateOptions;

  const formatDate = (date) => {
    return new Date(date).toLocaleString(dateLocale, dateOptions);
  };

  const compactTableTheme = {
    head: {
      cell: {
        base: "px-2 py-1",
      },
    },
    body: {
      cell: {
        base: "px-2 py-1",
      },
    },
  };

  return (
    <div className="mt-10">
      <Modal
        show={openModal.show}
        size="md"
        onClose={() => setOpenModal({ show: false })}
        popup
        dismissible
      >
        <ModalBody>
          <div className="text-center">
            <HiOutlineExclamationCircle className="mx-auto mb-4 h-14 w-14 text-gray-400 dark:text-gray-200" />
            <h3 className="mb-5 text-lg font-normal text-gray-500 dark:text-gray-200">
              Are you sure you want to delete this purchase?
            </h3>

            <Table className="mb-5">
              <TableBody className="dark:text-gray-100">
                {openModal?.purchase?.purchaseItems != null &&
                  openModal.purchase.purchaseItems.length > 0 &&
                  openModal.purchase.purchaseItems.map((purchaseItem) => (
                    <TableRow key={purchaseItem.id}>
                      <TableCell>{purchaseItem.quantity} x</TableCell>
                      <TableCell>{purchaseItem.product.name}</TableCell>
                      <TableCell className="text-right">
                        {currency.format(purchaseItem.totalGrossPrice)}
                      </TableCell>
                    </TableRow>
                  ))}
              </TableBody>
            </Table>

            <div className="mx-auto mb-4 text-gray-400 dark:text-gray-200">
              Selected payment method:{" "}
              <b>{findPaymentMethod(openModal?.purchase?.paymentMethod)}</b>
            </div>

            <div className="flex justify-center gap-4">
              <MyButton
                color="failure"
                disabled={processing}
                onClick={() => confirmDelete(openModal.purchase)}
              >
                Yes, I&apos;m sure
                {processing && <Spinner color="gray" className="ml-2" />}
              </MyButton>
              <MyButton
                color="black"
                className="bg-gray-200 dark:bg-gray-200"
                onClick={() => setOpenModal(false)}
              >
                No, cancel
              </MyButton>
            </div>
          </div>
        </ModalBody>
      </Modal>
      <Table
        striped
        theme={compactTableTheme}
        className={`table-fixed dark:text-gray-200 ${flash ? "animate__animated animate__pulse" : ""}`}
      >
        <TableHead>
          <TableRow>
            <TableHeadCell className="w-[55%]">Date</TableHeadCell>
            <TableHeadCell className="w-[15%] text-right">
              Total Price
            </TableHeadCell>
            <TableHeadCell className="w-[30%] text-right">Remove</TableHeadCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {history === null && (
            <TableRow>
              <TableCell colSpan={3} className="text-left">
                Purchases loading <Spinner className="ml-2" />
              </TableCell>
            </TableRow>
          )}
          {history !== null && history.length === 0 && (
            <TableRow>
              <TableCell colSpan={3} className="text-left">
                No purchases, yet.
              </TableCell>
            </TableRow>
          )}
          {history !== null &&
            history.slice(0, 3).map((purchase) => (
              <TableRow key={purchase.id}>
                <TableCell className="whitespace-nowrap">
                  {formatDate(purchase.createdAt)}
                </TableCell>
                <TableCell className="text-right">
                  {currency.format(purchase.totalGrossPrice)}
                </TableCell>
                <TableCell className="flex justify-end">
                  <MyButton
                    color="failure"
                    onClick={() => setOpenModal({ show: true, purchase })}
                  >
                    <HiXCircle />
                  </MyButton>
                </TableCell>
              </TableRow>
            ))}
        </TableBody>
      </Table>
    </div>
  );
}

PurchaseHistory.propTypes = {
  history: PropTypes.array,
  removeFromPurchaseHistory: PropTypes.func.isRequired,
};

export default PurchaseHistory;
