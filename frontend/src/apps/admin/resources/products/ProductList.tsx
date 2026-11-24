import React from "react";
import {
  usePermissions,
  List,
  Datagrid,
  TextField,
  DeleteButton,
  NumberField,
  BooleanField,
  ReferenceField,
  SimpleList,
} from "react-admin";
import { useIsSmallScreen } from "@admin/ui/hooks/useIsSmallScreen";
import { useConfig } from "@core/config/hooks/useConfig";
import { ProductFilters } from "./ProductFilters";
import ProductStockField from "./ProductStockField";
import { summarizeStock } from "./stock.utils";
import { Product } from "./types";

export const ProductList: React.FC = () => {
  const { currencyOptions, currencyLocale } = useConfig();
  const { permissions, isLoading: permissionsLoading } = usePermissions();
  const isSmall = useIsSmallScreen();
  const price = new Intl.NumberFormat(currencyLocale, currencyOptions);

  if (permissionsLoading) return <>Loading...</>;

  return (
    <List<Product>
      sort={{ field: "pos", order: "ASC" }}
      filters={ProductFilters}
    >
      {isSmall ? (
        <SimpleList<Product>
          linkType="edit"
          primaryText={(product) => product.name}
          secondaryText={(product) => {
            const stock = summarizeStock(product);
            const stockText =
              stock.total === null
                ? `${stock.sold} sold`
                : `${stock.left} / ${stock.total} left`;
            return `${price.format(product.grossPrice)} · ${stockText}`;
          }}
          tertiaryText={(product) =>
            [product.soldOut && "Sold out", product.hidden && "Hidden"]
              .filter(Boolean)
              .join(" · ")
          }
        />
      ) : (
        <Datagrid rowClick="edit" bulkActionButtons={false}>
          <NumberField source="id" />
          <TextField source="name" />
          <NumberField source="vatRate" />
          <NumberField
            source="grossPrice"
            locales={currencyLocale}
            options={currencyOptions}
          />
          <ReferenceField
            source="sectionId"
            reference="sections"
            label="Section"
            link={false}
          />
          <ReferenceField
            source="productGroupId"
            reference="productGroups"
            label="Group"
            link={false}
            emptyText="–"
          />
          <TextField source="variant" sortable={false} emptyText="–" />
          <NumberField source="pos" />
          <ProductStockField label="Stock" />
          <BooleanField source="soldOut" sortable={false} />
          <BooleanField source="hidden" sortable={false} />
          {permissions === "admin" && (
            <DeleteButton mutationMode="pessimistic" />
          )}
        </Datagrid>
      )}
    </List>
  );
};
