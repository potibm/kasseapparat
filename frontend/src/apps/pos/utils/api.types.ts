export interface ApiCreatePayloadPurchase {
  paymentMethod: string;
  cart: CartItem[];
  totalGrossPrice: string;
  totalNetPrice: string;
  sumupReaderId?: string;
}

interface CartItem {
  id: number;
  quantity: number;
  lists: null;
  guestlists: null;
}

export interface ApiCreateResponseProductInterest {
  id: number;
  productID: number;
}
