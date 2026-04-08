export enum PurchaseErrorType {
  ReaderBusy = "READER_BUSY",
  Generic = "GENERIC",
}

export const getPurchaseErrorType = (error: unknown): PurchaseErrorType => {
  if (error instanceof Error && error.message.includes("Reader Busy")) {
    return PurchaseErrorType.ReaderBusy;
  }
  return PurchaseErrorType.Generic;
};

export const getErrorMessage = (type: PurchaseErrorType): string => {
  switch (type) {
    case PurchaseErrorType.ReaderBusy:
      return "The SumUp reader is currently busy. Please complete or cancel the ongoing transaction.";
    default:
      return "An error occurred while processing the purchase.";
  }
};
