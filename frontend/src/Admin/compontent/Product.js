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
  ReferenceInput,
  SelectInput,
} from "react-admin";
import InventoryIcon from "@mui/icons-material/Inventory";
import { useConfig } from "../../provider/ConfigProvider";

export const ProductList = () => {
  const { permissions } = usePermissions();
  const currency = useConfig().currencyOptions;
  const locale = useConfig().Locale;

  return (
    <List sort={{ field: "pos", order: "ASC" }}>
      <Datagrid rowClick="edit" bulkActionButtons={false}>
        <NumberField source="id" />
        <TextField source="name" />
        <NumberField source="price" locales={locale} options={currency} />
        <NumberField source="pos" />
        <BooleanField source="wrapAfter" sortable={false} />
        <BooleanField source="hidden" sortable={false} />
        <TextField source="associatedList.name" sortable={false} />
        {permissions === "admin" && <DeleteButton mutationMode="pessimistic" />}
      </Datagrid>
    </List>
  );
};

export const ProductEdit = () => {
  return (
    <Edit>
      <SimpleForm
        toolbar={
          <Toolbar>
            <SaveButton />
          </Toolbar>
        }
      >
        <NumberInput disabled source="id" />
        <TextInput source="name" />
        <NumberInput source="price" min={0} />
        <NumberInput source="pos" />
        <BooleanInput source="wrapAfter" />
        <BooleanInput source="hidden" />
        <BooleanInput source="apiExport" />
        <ReferenceInput source="associatedListId" reference="lists">
          <SelectInput optionText="name" />
        </ReferenceInput>
      </SimpleForm>
    </Edit>
  );
};

export const ProductCreate = () => {
  return (
    <Create title="Create new product">
      <SimpleForm>
        <NumberInput disabled source="id" />
        <TextInput source="name" />
        <NumberInput source="price" min={0} />
        <NumberInput source="pos" />
        <BooleanInput source="wrapAfter" />
        <BooleanInput source="hidden" />
        <BooleanInput source="apiExport" />
        <ReferenceInput source="associatedListId" reference="lists">
          <SelectInput optionText="name" />
        </ReferenceInput>
      </SimpleForm>
    </Create>
  );
};

export const ProductIcon = () => <InventoryIcon />;
