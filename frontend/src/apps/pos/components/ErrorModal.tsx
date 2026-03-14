import React from "react";
import { Modal, ModalHeader, ModalBody } from "flowbite-react";

interface ErrorModalProps {
  message: string;
  onClose: () => void;
}

const ErrorModal: React.FC<ErrorModalProps> = ({ message, onClose }) => {
  return (
    <Modal show={message !== ""} onClose={onClose} dismissible={true}>
      <ModalHeader>Error</ModalHeader>
      <ModalBody className="dark:text-gray-200">
        <p>{message}</p>
      </ModalBody>
    </Modal>
  );
};

export default ErrorModal;
