import { Modal, Spinner } from "flowbite-react";
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
      <Modal.Header>
        <div className="text-2xl font-bold text-center">Register interest</div>
      </Modal.Header>
      <Modal.Body>
        <div className="text-center">
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
            <MyButton color="black" onClick={() => onClose()}>
              No, cancel
            </MyButton>
          </div>
        </div>
      </Modal.Body>
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
