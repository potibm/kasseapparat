import React from "react";
import {
  usePermissions,
  List,
  Datagrid,
  TextField,
  DeleteButton,
  NumberField,
  NumberInput,
  Edit,
  SimpleForm,
  TextInput,
  Create,
  BooleanField,
  BooleanInput,
  SaveButton,
  Toolbar,
  required,
  TabbedForm,
  FormTab,
} from "react-admin";
import InventoryIcon from "@mui/icons-material/Inventory";
import { useConfig } from "../../provider/ConfigProvider";
import DecimalInput from "./DecimalInput";

export const ProductList = () => {
  const currency = useConfig().currencyOptions;
  const locale = useConfig().Locale;

  const { permissions, isLoading: permissionsLoading } = usePermissions();
  if (permissionsLoading) return <>Loading...</>;

  return (
    <List sort={{ field: "pos", order: "ASC" }}>
      <Datagrid rowClick="edit" bulkActionButtons={false}>
        <NumberField source="id" />
        <TextField source="name" />
        <NumberField source="price" locales={locale} options={currency} />
        <NumberField source="pos" />
        <BooleanField source="wrapAfter" sortable={false} />
        <BooleanField source="soldOut" sortable={false} />
        <BooleanField source="hidden" sortable={false} />
        {permissions === "admin" && <DeleteButton mutationMode="pessimistic" />}
      </Datagrid>
    </List>
  );
};

export const ProductEdit = () => {
  return (
    <Edit>
      <TabbedForm
        toolbar={
          <Toolbar>
            <SaveButton />
          </Toolbar>
        }
      >
        <FormTab label="General">
          <NumberInput disabled source="id" />
          <TextInput source="name" validate={required()} />
          <DecimalInput source="price" />
        </FormTab>
        <FormTab label="Layout">
          <NumberInput
            source="pos"
            helperText="The products will shown in this order"
          />
          <BooleanInput
            source="wrapAfter"
            helperText="Create a line break afther this product"
          />
          <BooleanInput
            source="hidden"
            helperText="Hide this product from the POS"
          />
        </FormTab>
        <FormTab label="Stock">
          <h3>Stock</h3>
          <NumberInput
            source="totalStock"
            min={0}
            step={1}
            helperText="Number of available items. Shown for informational purposes, only."
          />
          <NumberInput
            source="unitsSold"
            min={0}
            step={1}
            disabled={true}
            helperText="Number of sold items. Shown for informational purposes, only."
          />
        </FormTab>
        <FormTab label="Sold Out">
          <BooleanInput
            source="soldOut"
            helperText="Still show the product to collect information how big the interrest is"
          />
          <NumberInput source="soldOutRequestCount" disabled={true} />
        </FormTab>
        <FormTab label="API">
          <BooleanInput source="apiExport" />
        </FormTab>
      </TabbedForm>
    </Edit>
  );
};

export const ProductCreate = () => {
  return (
    <Create title="Create new product">
      <SimpleForm>
        <NumberInput disabled source="id" />
        <TextInput source="name" validate={required()} />
        <DecimalInput source="price" min={0} />
        <NumberInput
          source="pos"
          helperText="The products will shown in this order"
        />
        <BooleanInput
          source="wrapAfter"
          helperText="Create a line break afther this product"
        />
        <BooleanInput source="hidden" />
      </SimpleForm>
    </Create>
  );
};

export const ProductIcon = () => <InventoryIcon />;
