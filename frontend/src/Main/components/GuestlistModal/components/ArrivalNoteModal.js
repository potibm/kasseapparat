import React from "react";
import { Modal, Button } from "flowbite-react";
import PropTypes from "prop-types";

const GuestlistArrivalNoteModal = ({ isOpen, onClose, arrivalNote, name }) => {
  return (
    <Modal show={isOpen} onClose={onClose}>
      <Modal.Header>Arrival Note concerning {name}</Modal.Header>
      <Modal.Body>
        <div className="space-y-6">
          <p className="text-xl italic font-serif leading-relaxed text-gray-500 dark:text-gray-400">
            &raquo;{arrivalNote}&laquo;
          </p>
        </div>
      </Modal.Body>
      <Modal.Footer>
        <Button onClick={onClose}>OK</Button>
      </Modal.Footer>
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
