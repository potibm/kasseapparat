import { Modal, Spinner, ModalBody, ModalHeader } from "flowbite-react";
import MyButton from "./MyButton";
import React, { useState } from "react";
import PropTypes from "prop-types";

const ProductInterestModal = ({
  show,
  onClose,
  product,
  addProductInterest,
}) => {
  const [processing, setProcessing] = useState(false);

  const handleRegisterInterest = () => {
    setProcessing(true);

    // Simulate API call
    addProductInterest(product)
      .then(() => {
        setProcessing(false);
        onClose();
      })
      .catch((error) => {
        console.error("Error registering interest: ", error);
        setProcessing(false);
      });
  };

  return (
    <Modal show={show} onClose={onClose} dismissible>
      <ModalHeader>
        <div className="text-2xl font-bold text-center">Register interest</div>
      </ModalHeader>
      <ModalBody>
        <div className="text-center dark:text-white">
          <p className="mb-2">{product.name} is currently sold out.</p>

          <p className="mb-5">
            To improve our stock management, it would be nice to record
            people&apos;s interest in this product.
          </p>

          <div className="flex justify-center gap-4">
            <MyButton
              disabled={processing}
              onClick={() => handleRegisterInterest()}
            >
              Yes, register interest
              {processing && <Spinner color="gray" className="ml-2" />}
            </MyButton>
            <MyButton
              color="black"
              onClick={() => onClose()}
              className="bg-gray-200 dark:bg-gray-200 dark:text-gray-800"
            >
              No, cancel
            </MyButton>
          </div>
        </div>
      </ModalBody>
    </Modal>
  );
};

ProductInterestModal.propTypes = {
  show: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  product: PropTypes.object.isRequired,
  addProductInterest: PropTypes.func.isRequired,
};

export default ProductInterestModal;
