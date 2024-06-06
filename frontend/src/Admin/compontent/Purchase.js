import React from "react";
import {
  List,
  Datagrid,
  TextField,
  DeleteButton,
  NumberField,
  NumberInput,
  SimpleForm,
  TextInput,
  Create,
  BooleanInput,
  DateField,
  Show,
  SimpleShowLayout,
  ArrayField,
} from "react-admin";
import InventoryIcon from "@mui/icons-material/Inventory";

export const PurchaseList = () => {
  return (
    <List sort={{ field: "createdAt", order: "DESC" }}>
      <Datagrid rowClick="show" bulkActionButtons={false}>
        <NumberField source="id" />
        <DateField source="createdAt" showTime={true} />
        <NumberField source="totalPrice" />
        <TextField source="createdBy.username" />
        <DeleteButton mutationMode="pessimistic" />
      </Datagrid>
    </List>
  );
};

export const PurchaseShow = (props) => {
  return (
    <Show {...props}>
      <SimpleShowLayout>
        <NumberField source="id" />
        <DateField source="createdAt" showTime={true} />
        <NumberField source="totalPrice" />
        <ArrayField source="purchaseItems">
          <Datagrid bulkActionButtons={false}>
            <NumberField source="quantity" />
            <NumberField source="price" />
            <TextField source="product.name" />
            <NumberField source="totalPrice" />
          </Datagrid>
        </ArrayField>
      </SimpleShowLayout>
    </Show>
  );
};

export const ProductCreate = () => {
  return (
    <Create title="Create new product">
      <SimpleForm>
        <NumberInput disabled source="id" />
        <TextInput source="name" />
        <NumberInput source="price" />
        <NumberInput source="pos" />
        <BooleanInput source="wrapAfter" />
        <BooleanInput source="apiExport" />
      </SimpleForm>
    </Create>
  );
};

export const ProductIcon = () => <InventoryIcon />;
