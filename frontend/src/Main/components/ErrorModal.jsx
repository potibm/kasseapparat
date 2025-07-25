import React from "react";
import { Modal, ModalHeader, ModalBody } from "flowbite-react";
import PropTypes from "prop-types";

const ErrorModal = ({ message, onClose }) => {
  return (
    <Modal show={message !== ""} onClose={onClose}>
      <ModalHeader>Error</ModalHeader>
      <ModalBody className="dark:text-gray-200">
        <p>{message}</p>
      </ModalBody>
    </Modal>
  );
};

ErrorModal.propTypes = {
  message: PropTypes.string.isRequired,
  onClose: PropTypes.func.isRequired,
};

export default ErrorModal;
