import { z } from "zod";
import Decimal from "decimal.js";

const DecimalSchema = z.string().transform((val) => new Decimal(val));

export const ProductSchema = z.object({
  id: z.number(),
  name: z.string(),
  netPrice: DecimalSchema,
  grossPrice: DecimalSchema,
  vatRate: DecimalSchema,
  vatAmount: DecimalSchema,
  wrapAfter: z.boolean(),
  hidden: z.boolean(),
  soldOut: z.boolean(),
  apiExport: z.boolean(),
  pos: z.number(),
  totalStock: z.number(),
  guestlists: z
    .array(
      z.object({
        id: z.number(),
        name: z.string(),
        typeCode: z.boolean(),
        productId: z.number(),
      }),
    )
    .nullable(),
  unitsSold: z.number(),
  soldOutRequestCount: z.number(),
});

export type Product = z.infer<typeof ProductSchema>;

const PurchaseStatusSchema = z.enum([
  "pending",
  "confirmed",
  "refunded",
  "failed",
]);

export const PurchaseSchema = z.object({
  id: z.uuid(),
  createdAt: z.string(),
  createdById: z.number(),
  createdBy: z.object({
    id: z.number(),
    username: z.string(),
    email: z.string(),
    admin: z.boolean(),
  }),
  paymentMethod: z.string(),
  totalNetPrice: DecimalSchema,
  totalGrossPrice: DecimalSchema,
  totalVatAmount: DecimalSchema,
  sumupTransactionId: z.uuid().nullable(),
  sumupClientTransactionId: z.uuid().nullable(),
  status: PurchaseStatusSchema,
  purchaseItems: z.array(
    z.object({
      id: z.number(),
      purchaseID: z.uuid(),
      productID: z.number(),
      product: ProductSchema,
      quantity: z.number(),
      netPrice: DecimalSchema,
      grossPrice: DecimalSchema,
      vatRate: DecimalSchema,
      vatAmount: DecimalSchema,
      totalNetPrice: DecimalSchema,
      totalGrossPrice: DecimalSchema,
      totalVatAmount: DecimalSchema,
    }),
  ),
});

export type Purchase = z.infer<typeof PurchaseSchema>;

export const GuestSchema = z.object({
  id: z.number(),
  name: z.string(),
  code: z.string().nullable(),
  listName: z.string(),
  additionalGuests: z.number().min(0).max(10),
  arrivalNote: z.string().nullable(),
});

export type Guest = z.infer<typeof GuestSchema>;
