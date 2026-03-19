import { Modal, Spinner, ModalBody, ModalHeader } from "flowbite-react";
import Button from "../../../../components/Button";
import React, { useState } from "react";
import { Product as ProductType } from "../../../../utils/api.schemas";
import { createLogger } from "@core/logger/logger";

const log = createLogger("Product");

interface ProductInterestModalProps {
  show: boolean;
  onClose: () => void;
  product: ProductType;
  addProductInterest: (product: ProductType) => Promise<void>;
}

const ProductInterestModal: React.FC<ProductInterestModalProps> = ({
  show,
  onClose,
  product,
  addProductInterest,
}) => {
  const [processing, setProcessing] = useState(false);

  const handleRegisterInterest = async () => {
    setProcessing(true);
    try {
      await addProductInterest(product);
      log.info("Interest registered for product", { productId: product.id });
      onClose();
    } catch (error) {
      log.error("Error registering interest", { productId: product.id }, error);
    } finally {
      setProcessing(false);
    }
  };

  return (
    <Modal show={show} onClose={onClose} dismissible>
      <ModalHeader>
        <div className="text-2xl font-bold text-center">Register interest</div>
      </ModalHeader>
      <ModalBody>
        <div className="text-center dark:text-white">
          <p className="mb-2">{product.name} is currently sold out.</p>

          <p className="mb-5">
            To improve our stock management, it would be nice to record
            people&apos;s interest in this product.
          </p>

          <div className="flex justify-center gap-4">
            <Button
              disabled={processing}
              onClick={() => handleRegisterInterest()}
            >
              Yes, register interest
              {processing && <Spinner color="gray" className="ml-2" />}
            </Button>
            <Button
              color="black"
              onClick={() => onClose()}
              className="bg-gray-200 dark:bg-gray-200 dark:text-gray-800"
            >
              No, cancel
            </Button>
          </div>
        </div>
      </ModalBody>
    </Modal>
  );
};

export default ProductInterestModal;
