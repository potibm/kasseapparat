export interface Product {
  id: number;
  name: string;
  price: string;
}

export const TEST_PRODUCTS = {
  REGULAR_TICKET: {
    id: 1,
    name: "🎟️ Regular",
    price: "40 €",
  },
  FREE_TICKET: {
    id: 3,
    name: "🎟️ Free",
    price: "0 €",
  },
  PREPAID_TICKET: {
    id: 4,
    name: "🎟️ Prepaid",
    price: "0 €",
  },
  COFFEE_MUG: {
    id: 16,
    name: "☕ Coffee Mug",
    price: "1 €",
  },
} as const;
