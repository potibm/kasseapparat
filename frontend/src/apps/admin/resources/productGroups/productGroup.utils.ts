export const PRODUCT_GROUP_NOT_EMPTY_MESSAGE =
  "This product group still contains products. Take them out of the group first.";

export const PRODUCT_GROUP_DELETE_FAILED_MESSAGE =
  "The product group could not be deleted.";

/**
 * The product form sends 0 for "no group": the backend reads a missing or null
 * productGroupId as "unchanged" and 0 as "take the product out of its group".
 */
export const parseProductGroupId = (value: unknown): number => {
  if (value === "" || value === null || value === undefined) return 0;

  const id = Number(value);

  return Number.isFinite(id) ? id : 0;
};

/** Whether a form value for productGroupId puts the product into a group. */
export const isInGroup = (value: unknown): boolean =>
  parseProductGroupId(value) > 0;
