export enum PurchaseErrorType {
  ReaderBusy = "READER_BUSY",
  OutOfStock = "OUT_OF_STOCK",
  Generic = "GENERIC",
}

// The backend rejects a purchase that exceeds the available stock with one of these
// messages in the error details, which the API client puts on the error message.
const outOfStockMarkers = ["left in stock", "is sold out"];

export const getPurchaseErrorType = (error: unknown): PurchaseErrorType => {
  if (error instanceof Error) {
    if (error.message.includes("Reader Busy")) {
      return PurchaseErrorType.ReaderBusy;
    }

    if (outOfStockMarkers.some((marker) => error.message.includes(marker))) {
      return PurchaseErrorType.OutOfStock;
    }
  }

  return PurchaseErrorType.Generic;
};

export const getErrorMessage = (type: PurchaseErrorType): string => {
  if (type === PurchaseErrorType.ReaderBusy) {
    return "The SumUp reader is currently busy. Please complete or cancel the ongoing transaction.";
  }

  if (type === PurchaseErrorType.OutOfStock) {
    return "Not enough stock left for one of the products in the cart.";
  }

  return "An error occurred while processing the purchase.";
};
