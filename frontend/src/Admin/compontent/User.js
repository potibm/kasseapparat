
import { List, Datagrid, TextField, EmailField, DeleteButton, NumberField, NumberInput, Edit, SimpleForm, TextInput, Create, BooleanField, BooleanInput, UpdateButton, SaveButton, Toolbar, DeleteWithConfirmButton, PasswordInput } from 'react-admin';
import PersonIcon from '@mui/icons-material/Person';

export const UserList = () => {
    return (
        <List sort={{ field: 'id', order: 'ASC' }}>
            <Datagrid rowClick="edit" bulkActionButtons={false}>
              <NumberField source="id" />
              <TextField source="username" />
              <DeleteButton mutationMode="pessimistic" /> 
            </Datagrid>
        </List>
    )
}

export const UserEdit = () => {

    const equalToPassword = (value, allValues) => {
        if (value !== allValues.password) {
            return 'The two passwords must match';
        }
    }    

    return (
        <Edit>
            <SimpleForm toolbar={<Toolbar><SaveButton /></Toolbar>}>
                <NumberInput disabled source="id" />
                <TextInput source="username" />
                <PasswordInput source="password" />
                <PasswordInput source="confirm_password" validate={equalToPassword} />
            </SimpleForm>
        </Edit>
    )
}

export const UserCreate = () => {

    const equalToPassword = (value, allValues) => {
        if (value !== allValues.password) {
            return 'The two passwords must match';
        }
    }    

    return (
        <Create title="Create new user">
            <SimpleForm>
            <NumberInput disabled source="id" />
                <TextInput source="username" />
                <PasswordInput source="password" />
                <PasswordInput source="confirm_password" validate={equalToPassword} />
            </SimpleForm>
        </Create>
    )
}

export const UserIcon = () => <PersonIcon />;