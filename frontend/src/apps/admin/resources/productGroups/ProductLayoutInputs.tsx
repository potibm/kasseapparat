import React from "react";
import { ReferenceInput, SelectInput, TextInput } from "react-admin";
import { useWatch } from "react-hook-form";
import SectionInput from "../sections/SectionInput";
import { isInGroup, parseProductGroupId } from "./productGroup.utils";
import { ProductGroup } from "./types";

/**
 * The layout fields of a product: its group and variant, and its section. A product in
 * a group always takes the group's section (ADR 006), so the section select is only
 * shown for products without a group.
 */
export const ProductLayoutInputs: React.FC<{ preselectSection?: boolean }> = ({
  preselectSection = false,
}) => {
  const productGroupId = useWatch({ name: "productGroupId" });
  const grouped = isInGroup(productGroupId);

  return (
    <>
      <ReferenceInput
        source="productGroupId"
        reference="productGroups"
        sort={{ field: "name", order: "ASC" }}
        perPage={100}
      >
        <SelectInput
          label="Product group"
          optionText={(record: ProductGroup) => record.name}
          emptyText="None"
          parse={parseProductGroupId}
          helperText="Variants of one article, e.g. the sizes of a shirt, share a row in the POS"
        />
      </ReferenceInput>
      {grouped ? (
        <TextInput
          source="variant"
          label="Variant"
          helperText="Label of the button in the group's row, e.g. “XL”. Empty shows the product name."
        />
      ) : (
        <SectionInput preselectFirst={preselectSection} />
      )}
    </>
  );
};

export default ProductLayoutInputs;
