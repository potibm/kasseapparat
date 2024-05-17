
import { List, Datagrid, TextField, EmailField, DeleteButton, NumberField, NumberInput, Edit, SimpleForm, TextInput, Create, BooleanField, BooleanInput } from 'react-admin';
import InventoryIcon from '@mui/icons-material/Inventory';

export const ProductList = () => {
    return (
        <List sort={{ field: 'pos', order: 'ASC' }}>
            <Datagrid rowClick="edit">
              <NumberField source="id" />
              <TextField source="name" />
              <NumberField source="price" />
              <NumberField source="pos" />
              <BooleanField source="wrapAfter" sortable={false} />        
            </Datagrid>
        </List>
    )
}

export const ProductEdit = () => {
    return (
        <Edit>
            <SimpleForm>
                <NumberInput disabled source="id" />
                <TextInput source="name" />
                <NumberInput source="price" />
                <NumberInput source="pos" />
                <BooleanInput source="wrapAfter" />
                <BooleanInput source="apiExport" />
            </SimpleForm>
        </Edit>
    )
}

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
    )
}

export const ProductIcon = () => <InventoryIcon />;