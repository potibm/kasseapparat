import React, { useEffect, useState, useRef } from "react";
import { Modal, ModalBody, ModalFooter, Spinner } from "flowbite-react";
import { HiClock, HiCheckCircle, HiXCircle } from "react-icons/hi";
import { getCurrentReaderId } from "../../../helper/ReaderCookie";
import PropTypes from "prop-types";
import MyButton from "../MyButton";
import { useAuth } from "../../../Auth/provider/AuthProvider";
import { useConfig } from "../../../provider/ConfigProvider";

const PollingModal = ({ show, purchase, onClose, onConfirmed, onComplete }) => {
  const [status, setStatus] = useState(purchase.status);
  const [lastUpdate, setLastUpdate] = useState(Date.now());
  const [now, setNow] = useState(Date.now()); // for age display
  const [error, setError] = useState(null);
  const [flash, setFlash] = useState(false);
  const [processing, setProcessing] = useState(false);
  const wsRef = useRef(null);
  const { token: jwtToken } = useAuth();
  const { websocketHost } = useConfig();

  const sumUpReaderId = getCurrentReaderId();
  const closeModalTimeout = 3000;
  const ageInSeconds = Math.max(0, Math.round((now - lastUpdate) / 1000));

  const statusRef = useRef(status);
  const wasHandledRef = useRef(false);

  useEffect(() => {
    statusRef.current = status;
  }, [status]);

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
      fontSize: "text-lg",
    },
    confirmed: {
      color: "text-green-600",
      icon: <HiCheckCircle className="inline mr-1 h-8 w-8" />,
      spinnerColor: null,
      fontSize: "text-3xl",
    },
    failed: {
      color: "text-red-600",
      icon: <HiXCircle className="inline mr-1 h-8 w-8" />,
      spinnerColor: null,
      fontSize: "text-3xl",
    },
  };
  const current = statusInfo[status] || {
    color: "text-gray-600",
    icon: null,
    spinnerColor: "info",
    fontSize: "text-lg",
  };

  useEffect(() => {
    const handleFailure = (message) => {
      wasHandledRef.current = true;
      setStatus("failed");
      setError(message);
      setProcessing(false);
      onComplete(false);
      setTimeout(onClose, closeModalTimeout);
      if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
        wsRef.current.close();
      }
    };

    const handleSuccess = (data) => {
      wasHandledRef.current = true;
      setStatus("confirmed");
      setError(null);
      setProcessing(false);
      onComplete(true);
      setTimeout(() => onConfirmed(data), closeModalTimeout);
      if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
        wsRef.current.close();
      }
    };
    const ws = new WebSocket(
      `${websocketHost}/api/v2/purchases/${purchase.id}/ws?token=${jwtToken}`,
    );
    wsRef.current = ws;

    const connectionTimeout = setTimeout(() => {
      if (ws.readyState !== WebSocket.OPEN) {
        console.warn("WebSocket did not connect in time.");
        handleFailure("Could not fetch the status of the payment terminal.");
      }
    }, 3000); // 3 seconds

    ws.onopen = () => {
      clearTimeout(connectionTimeout);
      console.log("WebSocket connected");
    };

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        if (data.type === "status_update") {
          console.log("WebSocket message received:", data);
          setFlash(true);
          setTimeout(() => setFlash(false), 500);
          setLastUpdate(Date.now());

          if (data.status === "confirmed") {
            handleSuccess(data);
          } else if (data.status === "failed") {
            handleFailure("Purchase failed.");
          } else {
            console.log("Purchase status update:", data.status);
          }
        } else if (data.type === "cancel_ack") {
          // close the WebSocket connection after cancel acknowledgment
          handleFailure("Purchase cancelled by user.");
        }
      } catch (err) {
        console.error("WebSocket parsing error:", err, event.data);
      }
    };

    ws.onerror = (err) => {
      console.error("WebSocket error", err);
      setError("WebSocket error occurred");
    };

    ws.onclose = () => {
      console.log("WebSocket closed");
      if (!wasHandledRef.current && statusRef.current === "pending") {
        handleFailure("Connection lost.");
      }
    };

    return () => {
      clearTimeout(connectionTimeout);
      if (wsRef.current && wsRef.current.readyState === WebSocket.OPEN) {
        wsRef.current.close();
      }
    };
  }, [purchase.id, onConfirmed, onClose, onComplete, jwtToken, websocketHost]);

  if (!purchase?.id) return null;

  return (
    <Modal show={show} size="md" popup dismissible={false}>
      <ModalBody className="m-5">
        <div className="text-center">
          <h3 className="text-lg font-semibold mb-8">Purchase Status</h3>
          <div
            key={status}
            className={`text-center space-y-4 ${flash ? "animate__animated animate__headShake" : ""}`}
          >
            {status === "pending" && (
              <Spinner
                color={current.spinnerColor}
                size="xl"
                className="mb-3"
              />
            )}
            <div
              className={`text-lg font-medium flex items-center justify-center gap-2 ${current.color}`}
            >
              {current.icon}
              <span>
                {status}
                {status === "pending" && <> ({ageInSeconds} sec ago)</>}
              </span>
            </div>
            {error && <p className="text-red-600 font-medium">{error}</p>}
          </div>
        </div>
      </ModalBody>
      <ModalFooter>
        {status === "pending" && (
          <MyButton
            color="failure"
            disabled={processing}
            onClick={() => cancelPayment()}
          >
            {processing ? "Cancelling..." : "Abort Purchase"}
            {processing && <Spinner color="gray" className="ml-2" />}
          </MyButton>
        )}
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
