import { faker } from "@faker-js/faker";
import Decimal from "decimal.js";
import {
  User as UserType,
  Product as ProductType,
  Guest as GuestType,
  Purchase as PurchaseType,
  PurchaseItem as PurchaseItemType,
} from "./api.schemas";

/**
 * Generates a realistic mock Product object for testing purposes.
 *
 * @param overrides  Specific fields to override in the generated object.
 */
export const createMockProduct = (
  overrides?: Partial<ProductType>,
): ProductType => {
  // We generate a realistic base price
  const netPriceNum = faker.number.float({
    min: 5,
    max: 200,
    fractionDigits: 2,
  });
  const vatRateNum = faker.number.int({ min: 0, max: 25 });
  const vatAmountNum = netPriceNum * (vatRateNum / 100);
  const grossPriceNum = netPriceNum + vatAmountNum;

  return {
    id: faker.number.int({ min: 1, max: 99999 }),
    name: faker.commerce.productName(),

    netPrice: new Decimal(netPriceNum.toFixed(2)),
    grossPrice: new Decimal(grossPriceNum.toFixed(2)),
    vatRate: new Decimal(vatRateNum.toFixed(2)),
    vatAmount: new Decimal(vatAmountNum.toFixed(2)),

    wrapAfter: faker.datatype.boolean(),
    hidden: faker.datatype.boolean(0.1),
    soldOut: faker.datatype.boolean(0.2),
    apiExport: true,
    pos: faker.number.int({ min: 1, max: 100 }),
    totalStock: faker.number.int({ min: 10, max: 1000 }),
    guestlists: null,
    unitsSold: faker.number.int({ min: 0, max: 500 }),
    soldOutRequestCount: faker.number.int({ min: 0, max: 50 }),

    ...overrides,
  };
};

/**
 * Generates a realistic mock Guest object for testing purposes.
 *
 * @param overrides  Specific fields to override in the generated object.
 */
export const createMockGuest = (overrides?: Partial<GuestType>): GuestType => {
  const additionalGuestsNum = faker.number.int({ min: 0, max: 5 });

  return {
    id: faker.number.int({ min: 1, max: 99999 }),
    name: faker.person.fullName(),
    code:
      faker.helpers.maybe(
        () => faker.string.alphanumeric({ length: 8, casing: "upper" }),
        { probability: 0.2 },
      ) ?? null,
    listName: faker.commerce.department(),
    additionalGuests: additionalGuestsNum,
    arrivalNote:
      faker.helpers.maybe(() => faker.lorem.sentence(), { probability: 0.1 }) ??
      null,
    attendedGuests: faker.number.int({ min: 0, max: additionalGuestsNum }),

    ...overrides,
  };
};

/**
 * Generates a realistic mock Purchase object for testing purposes.
 *
 * @param overrides  Specific fields to override in the generated object.
 */
export const createMockPurchase = (
  overrides?: Partial<PurchaseType>,
): PurchaseType => {
  const user = createMockUser();

  const purchaseItemCount = faker.number.int({ min: 1, max: 5 });
  const purchaseItems = Array.from({ length: purchaseItemCount }, () =>
    createMockPurchaseItem(),
  );

  let totalNetPrice = new Decimal(0);
  let totalGrossPrice = new Decimal(0);
  let totalVatAmount = new Decimal(0);

  purchaseItems.forEach((item) => {
    totalNetPrice = totalNetPrice.add(item.totalNetPrice);
    totalGrossPrice = totalGrossPrice.add(item.totalGrossPrice);
    totalVatAmount = totalVatAmount.add(item.totalVatAmount);
  });

  return {
    id: faker.string.uuid(),
    createdAt: faker.date.recent().toISOString(),
    createdById: user.id,
    createdBy: user,
    paymentMethod: faker.helpers.arrayElement(["CASH", "CC"]),
    totalNetPrice: totalNetPrice,
    totalGrossPrice: totalGrossPrice,
    totalVatAmount: totalVatAmount,
    sumupTransactionId: null,
    sumupClientTransactionId: null,
    status: "pending",
    purchaseItems: purchaseItems,

    ...overrides,
  };
};

/**
 * Generates a realistic mock User object for testing purposes.
 *
 * @param overrides  Specific fields to override in the generated object.
 */
export const createMockUser = (overrides?: Partial<UserType>): UserType => {
  return {
    id: faker.number.int({ min: 1, max: 99999 }),
    username: faker.internet.displayName(),
    email: faker.internet.email(),
    admin: faker.datatype.boolean(0.1),

    ...overrides,
  };
};

/**
 * Generates a realistic mock PurchaseItem object for testing purposes.
 *
 * @param overrides  Specific fields to override in the generated object.
 */
export const createMockPurchaseItem = (
  overrides?: Partial<PurchaseItemType>,
): PurchaseItemType => {
  const quantity = faker.number.int({ min: 1, max: 5 });
  const vatRateNum = new Decimal(faker.number.int({ min: 0, max: 25 }));

  const netPrice = new Decimal(
    faker.number.float({ min: 5, max: 200, fractionDigits: 2 }),
  );
  const vatAmount = netPrice.mul(vatRateNum.div(100)).toDecimalPlaces(2);
  const grossPrice = netPrice.add(vatAmount).toDecimalPlaces(2);

  const totalVatAmountNum = vatAmount.mul(quantity).toDecimalPlaces(2);
  const totalGrossPrice = grossPrice.mul(quantity).toDecimalPlaces(2);
  const totalNetPrice = netPrice.mul(quantity).toDecimalPlaces(2);

  const product = createMockProduct();

  return {
    id: faker.number.int({ min: 1, max: 99999 }),
    purchaseID: faker.string.uuid(),
    productID: product.id,
    product: product,
    quantity: quantity,
    netPrice: netPrice,
    grossPrice: grossPrice,
    vatRate: vatRateNum,
    vatAmount: vatAmount,
    totalNetPrice: totalNetPrice,
    totalGrossPrice: totalGrossPrice,
    totalVatAmount: totalVatAmountNum,

    ...overrides,
  };
};
