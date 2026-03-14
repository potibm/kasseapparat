export interface ApiCreatePayloadPurchase {
  paymentMethod: string;
  cart: CartItem[];
  totalGrossPrice: string;
  totalNetPrice: string;
  sumupReaderId?: string;
}

export interface CartItem {
  id: number;
  quantity: number;
  lists: null;
  guestlists: null;
}

export interface ApiCreateResponseProductInterest {
  id: number;
  productID: number;
}

export interface AuthApiError {
  message: string;
  status?: number;
  details?: string;
  data?: {
    details?: string;
    message?: string;
  };
}
