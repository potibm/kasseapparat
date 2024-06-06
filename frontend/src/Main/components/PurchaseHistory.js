import { HiXCircle, HiOutlineExclamationCircle } from "react-icons/hi";
import React, { useState } from "react";
import { Button, Modal, Table } from "flowbite-react";
import PropTypes from "prop-types";

function PurchaseHistory({ history, currency, removeFromPurchaseHistory }) {
  const [openModal, setOpenModal] = useState({ show: false, purchase: null });

  const confirmDelete = (purchase) => {
    setOpenModal({ show: false });
    removeFromPurchaseHistory(purchase);
  };

  return (
    <div className="mt-10">
      <Modal
        show={openModal.show}
        size="md"
        onClose={() => setOpenModal({ show: false })}
        popup
      >
        <Modal.Header />
        <Modal.Body>
          <div className="text-center">
            <HiOutlineExclamationCircle className="mx-auto mb-4 h-14 w-14 text-gray-400 dark:text-gray-200" />
            <h3 className="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
              Are you sure you want to delete this purchase?
            </h3>
            <div className="flex justify-center gap-4">
              <Button
                color="failure"
                onClick={() => confirmDelete(openModal.purchase)}
              >
                {"Yes, I'm sure"}
              </Button>
              <Button color="gray" onClick={() => setOpenModal(false)}>
                No, cancel
              </Button>
            </div>
          </div>
        </Modal.Body>
      </Modal>
      <Table striped>
        <Table.Head>
          <Table.HeadCell>Date</Table.HeadCell>
          <Table.HeadCell className="text-right">Total Price</Table.HeadCell>
          <Table.HeadCell>Remove</Table.HeadCell>
        </Table.Head>
        <Table.Body>
          {history.slice(0, 2).map((purchase) => (
            <Table.Row key={purchase.id}>
              <Table.Cell className="text-right">
                {new Date(purchase.createdAt).toLocaleString("de-DE", {
                  weekday: "long",
                  hour: "2-digit",
                  minute: "2-digit",
                })}
              </Table.Cell>
              <Table.Cell className="text-right">
                {currency.format(purchase.totalPrice)}
              </Table.Cell>
              <Table.Cell>
                <Button
                  color="failure"
                  onClick={() => setOpenModal({ show: true, purchase })}
                >
                  <HiXCircle />
                </Button>
              </Table.Cell>
            </Table.Row>
          ))}
        </Table.Body>
      </Table>
    </div>
  );
}

PurchaseHistory.propTypes = {
  history: PropTypes.array.isRequired,
  currency: PropTypes.object.isRequired,
  removeFromPurchaseHistory: PropTypes.func.isRequired,
};

export default PurchaseHistory;
