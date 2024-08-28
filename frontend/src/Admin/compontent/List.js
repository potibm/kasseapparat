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
  ReferenceInput,
  SelectInput,
  useGetIdentity,
  useRecordContext,
  CreateButton,
} from "react-admin";
import GroupsIcon from "@mui/icons-material/Groups";
import { useNavigate } from "react-router-dom";

const ConditionalDeleteButton = (props) => {
  const record = useRecordContext(props);

  const { permissions, isLoading: permissionsLoading } = usePermissions();
  const { data: identity, isLoading: identityLoading } = useGetIdentity();
  if (permissionsLoading || identityLoading) return <>Loading...</>;

  const currentUserId = identity.id;
  const createdByCurrentUser = record && record.createdById === currentUserId;

  if (permissions === "admin" || createdByCurrentUser) {
    return <DeleteButton {...props} />;
  }
  return null;
};

const CreateListEntryButton = (props) => {
  const record = useRecordContext(props);
  const navigate = useNavigate();

  const handleCreateEntry = (guestlistId) => {
    navigate(`/admin/listEntries/create?list_id=${guestlistId}`);
  };

  return (
    <CreateButton
      {...props}
      label="Create Entry"
      onClick={(e) => {
        e.preventDefault();
        e.stopPropagation();
        handleCreateEntry(record?.id);
      }}
    />
  );
};

export const ListList = (props) => {
  return (
    <List {...props} sort={{ field: "id", order: "ASC" }}>
      <Datagrid rowClick="edit" bulkActionButtons={false}>
        <NumberField source="id" />
        <TextField source="name" />
        <BooleanField source="typeCode" sortable={false} />
        <TextField source="product.name" sortable={false} />
        <ConditionalDeleteButton mutationMode="pessimistic" />
        <CreateListEntryButton />
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
        <ReferenceInput source="productId" reference="products">
          <SelectInput optionText="name" validate={required()} />
        </ReferenceInput>
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
        <ReferenceInput source="productId" reference="products">
          <SelectInput optionText="name" validate={required()} />
        </ReferenceInput>
      </SimpleForm>
    </Create>
  );
};

export const ListIcon = () => <GroupsIcon />;
