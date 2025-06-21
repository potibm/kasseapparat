import React, { useEffect, useState, useRef } from "react";
import { Modal, ModalBody, ModalFooter, Spinner } from "flowbite-react";
import { HiClock, HiCheckCircle, HiXCircle } from "react-icons/hi";
import { getCurrentReaderId } from "../../../helper/ReaderCookie";
import PropTypes from "prop-types";
import MyButton from "../MyButton";
import { useAuth } from "../../../Auth/provider/AuthProvider";

const PollingModal = ({ show, purchase, onClose, onConfirmed, onComplete }) => {
  const [status, setStatus] = useState(purchase.status);
  const [lastUpdate, setLastUpdate] = useState(Date.now());
  const [now, setNow] = useState(Date.now()); // for age display
  const [error, setError] = useState(null);
  const [flash, setFlash] = useState(false);
  const [processing, setProcessing] = useState(false);
  const wsRef = useRef(null);
  const { token: jwtToken } = useAuth();

  const sumUpReaderId = getCurrentReaderId();
  const closeModalTimeout = 2000;
  const ageInSeconds = Math.max(0, Math.round((now - lastUpdate) / 1000));

  const cancelPayment = () => {
    // disable the button to prevent multiple clicks
    if (processing) return;
    setProcessing(true);

    // send a message to the websocket to cancel the payment
    if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
      wsRef.current.send(
        JSON.stringify({ type: "cancel_payment", reader_id: sumUpReaderId }),
      );
      console.log("Cancel payment message sent to WebSocket");
    } else {
      console.error("WebSocket is not open, cannot send cancel message");
      setError("WebSocket is not open, cannot cancel payment");
      setProcessing(false);
    }
  };

  // Update "now" every second
  useEffect(() => {
    const interval = setInterval(() => setNow(Date.now()), 1000);
    return () => clearInterval(interval);
  }, []);

  // map status to color + icon
  const statusInfo = {
    pending: {
      color: "text-yellow-600",
      icon: <HiClock className="inline mr-1 h-5 w-5" />,
      spinnerColor: "warning",
    },
    confirmed: {
      color: "text-green-600",
      icon: <HiCheckCircle className="inline mr-1 h-5 w-5" />,
      spinnerColor: null,
    },
    failed: {
      color: "text-red-600",
      icon: <HiXCircle className="inline mr-1 h-5 w-5" />,
      spinnerColor: null,
    },
  };
  const current = statusInfo[status] || {
    color: "text-gray-600",
    icon: null,
    spinnerColor: "info",
  };

  useEffect(() => {
    const ws = new WebSocket(
      `ws://localhost:3001/api/v2/purchases/${purchase.id}/ws?token=${jwtToken}`,
    );
    wsRef.current = ws;

    ws.onopen = () => {
      console.log("WebSocket connected");
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === "status_update") {
          setFlash(true);
          setTimeout(() => setFlash(false), 500);
          setStatus(data.status);
          setLastUpdate(Date.now());

          if (data.status === "confirmed") {
            ws.close();
            onComplete(true);
            setTimeout(() => onConfirmed(data), closeModalTimeout);
          } else if (data.status === "failed") {
            ws.close();
            setError("Purchase failed.");
            onComplete(false);
            setTimeout(onClose, closeModalTimeout);
          } else {
            console.log("Purchase status update:", data.status);
          }
        } else if (data.type === "cancel_ack") {
          // close the WebSocket connection after cancel acknowledgment
          ws.close();

          setProcessing(false);
          setStatus("failed");
          setError("Purchase cancelled by user.");
          setTimeout(onClose, closeModalTimeout);
          onComplete(false);
        }
      } catch (err) {
        console.error("WebSocket parsing error:", err);
      }
    };

    ws.onerror = (err) => {
      console.error("WebSocket error", err);
      setError("WebSocket error occurred");
    };

    ws.onclose = () => {
      console.log("WebSocket closed");
    };

    return () => {
      if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
        wsRef.current.close();
      }
    };
  }, [purchase.id, onConfirmed, onClose, onComplete, jwtToken]);

  return (
    <Modal show={show} size="md" popup dismissible={false}>
      <ModalBody className="m-5">
        <div className="text-center">
          <h3 className="text-lg font-semibold mb-8">Purchase Status</h3>
          <div
            className={`text-center space-y-4 ${flash ? "animate__animated animate__headShake" : ""}`}
          >
            {status === "pending" && (
              <Spinner color={current.spinnerColor} size="xl" />
            )}
            <div
              className={`text-lg font-medium flex items-center justify-center gap-2 ${current.color}`}
            >
              {current.icon}
              <span>
                {status} ({ageInSeconds} sec ago)
              </span>
            </div>
          </div>
          {error && <p className="text-red-600 font-medium">{error}</p>}
        </div>
      </ModalBody>
      <ModalFooter>
        <MyButton
          color="failure"
          disabled={processing}
          onClick={() => cancelPayment()}
        >
          Abort Purchase
          {processing && <Spinner color="gray" className="ml-2" />}
        </MyButton>
      </ModalFooter>
    </Modal>
  );
};

PollingModal.propTypes = {
  show: PropTypes.bool.isRequired,
  purchase: PropTypes.shape({
    id: PropTypes.number.isRequired,
    status: PropTypes.string.isRequired,
  }),
  onClose: PropTypes.func.isRequired,
  onConfirmed: PropTypes.func.isRequired,
  onComplete: PropTypes.func.isRequired,
};

export default PollingModal;
