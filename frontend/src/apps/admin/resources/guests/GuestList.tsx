import React from "react";
import { List, Datagrid, TextField, NumberField, DateField } from "react-admin";
import GuestActions from "./components/GuestActions";
import { GuestFilters } from "./components/GuestFilters";
import { useConfig } from "../../../../core/config/providers/ConfigProvider";
import { ConditionalDeleteOnOwnershipButton } from "../../ui/buttons/ConditionalDeleteOnOwnershipButton";
import { AttendedGuestsBooleanField } from "./components/AttendedGuestsBooleanField";
import { ListProps } from "flowbite-react";

export const GuestList: React.FC<ListProps> = (props) => {
  const { locale } = useConfig();

  return (
    <List
      sort={{ field: "id", order: "ASC" }}
      filters={GuestFilters}
      actions={<GuestActions />}
      {...props}
    >
      <Datagrid rowClick="edit" bulkActionButtons={false}>
        <NumberField source="id" />
        <TextField source="name" />
        <TextField source="guestlist.name" label="Guestlist" />
        <NumberField source="additionalGuests" sortable={false} />
        <AttendedGuestsBooleanField label="present" sortable={false} />
        <DateField
          source="arrivedAt"
          showTime={true}
          emptyText="-"
          locales={locale}
          options={{ weekday: "short", hour: "2-digit", minute: "2-digit" }}
        />
        <ConditionalDeleteOnOwnershipButton mutationMode="pessimistic" />
      </Datagrid>
    </List>
  );
};
