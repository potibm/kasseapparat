import React from 'react'
import { Modal } from 'flowbite-react'
import PropTypes from 'prop-types'

function ErrorModal ({ message, onClose }) {
  return (
        <Modal show={message !== ''} onClose={onClose}>
            <Modal.Header>Error</Modal.Header>
            <Modal.Body>
                <p>{message}</p>
            </Modal.Body>
        </Modal>
  )
}

ErrorModal.propTypes = {
  message: PropTypes.string.isRequired,
  onClose: PropTypes.func.isRequired
}

export default ErrorModal
