import React from "react";
import {
  Modal,
  Button,
  ModalFooter,
  ModalBody,
  ModalHeader,
} from "flowbite-react";

interface GuestlistArrivalNoteModalProps {
  isOpen: boolean;
  onClose: () => void;
  arrivalNote: string;
  name: string;
}

const GuestlistArrivalNoteModal: React.FC<GuestlistArrivalNoteModalProps> = ({
  isOpen,
  onClose,
  arrivalNote,
  name,
}) => {
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

export default GuestlistArrivalNoteModal;
