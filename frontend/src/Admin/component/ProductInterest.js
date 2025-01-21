import React from "react";
import {
  List,
  Datagrid,
  TextField,
  DeleteButton,
  NumberField,
  DateField,
} from "react-admin";
import ThumbUpIcon from "@mui/icons-material/ThumbUp";

export const ProductInterestList = () => {
  return (
    <List sort={{ field: "pos", order: "ASC" }}>
      <Datagrid rowClick="" bulkActionButtons={false}>
        <NumberField source="id" />
        <DateField showTime source="createdAt" />
        <NumberField source="product.id" />
        <TextField source="product.name" />
        <DeleteButton mutationMode="pessimistic" />
      </Datagrid>
    </List>
  );
};

export const ProductInterestIcon = () => <ThumbUpIcon />;
