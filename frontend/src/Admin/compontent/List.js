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
} from "react-admin";
import GroupsIcon from "@mui/icons-material/Groups";

export const ListList = () => {
  const { permissions } = usePermissions();

  return (
    <List sort={{ field: "id", order: "ASC" }}>
      <Datagrid rowClick="edit" bulkActionButtons={false}>
        <NumberField source="id" />
        <TextField source="name" />
        <BooleanField source="typeCode" sortable={false} />
        {permissions === "admin" && <DeleteButton mutationMode="pessimistic" />}
      </Datagrid>
    </List>
  );
};

export const ListEdit = () => {
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
        <TextInput source="name" validate={required()} />
        <BooleanInput source="typeCode" />
      </SimpleForm>
    </Edit>
  );
};

export const ListCreate = () => {
  return (
    <Create title="Create new list">
      <SimpleForm>
        <NumberInput disabled source="id" />
        <TextInput source="name" validate={required()} />
        <BooleanInput source="typeCode" />
      </SimpleForm>
    </Create>
  );
};

export const ListIcon = () => <GroupsIcon />;
