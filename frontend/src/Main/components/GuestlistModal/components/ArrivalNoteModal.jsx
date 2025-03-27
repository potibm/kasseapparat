import React from "react";
import {
  Modal,
  Button,
  ModalFooter,
  ModalBody,
  ModalHeader,
} from "flowbite-react";
import PropTypes from "prop-types";

const GuestlistArrivalNoteModal = ({ isOpen, onClose, arrivalNote, name }) => {
  return (
    <Modal show={isOpen} onClose={onClose}>
      <ModalHeader>Arrival Note concerning {name}</ModalHeader>
      <ModalBody>
        <div className="space-y-6">
          <p className="text-xl italic font-serif leading-relaxed text-gray-500 dark:text-gray-400">
            &raquo;{arrivalNote}&laquo;
          </p>
        </div>
      </ModalBody>
      <ModalFooter>
        <Button onClick={onClose}>OK</Button>
      </ModalFooter>
    </Modal>
  );
};

GuestlistArrivalNoteModal.propTypes = {
  isOpen: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
  arrivalNote: PropTypes.string.isRequired,
  name: PropTypes.string.isRequired,
};

export default GuestlistArrivalNoteModal;
