import React from 'react';
import { Modal } from 'flowbite-react';

export default function ErrorModal({ message, onClose }) {
    return (
        <Modal show={message !== ''} onClose={onClose}>
            <Modal.Header>Error</Modal.Header>
            <Modal.Body>
                <p>{message}</p>
            </Modal.Body>
        </Modal>
    );
}
