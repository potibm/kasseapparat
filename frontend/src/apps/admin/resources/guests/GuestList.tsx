import React from "react";
import {
  List,
  Datagrid,
  TextField,
  NumberField,
  DateField,
  ListProps,
  SimpleList,
} from "react-admin";
import { useIsSmallScreen } from "@admin/ui/hooks/useIsSmallScreen";
import GuestActions from "./components/GuestActions";
import { GuestFilters } from "./components/GuestFilters";
import { useConfig } from "@core/config/hooks/useConfig";
import { ConditionalDeleteOnOwnershipButton } from "../../ui/buttons/ConditionalDeleteOnOwnershipButton";
import { AttendedGuestsBooleanField } from "./components/AttendedGuestsBooleanField";

export const GuestList: React.FC<ListProps> = (props) => {
  const { dateLocale } = useConfig();
  const isSmall = useIsSmallScreen();

  return (
    <List
      sort={{ field: "id", order: "ASC" }}
      filters={GuestFilters}
      actions={<GuestActions />}
      {...props}
    >
      {isSmall ? (
        <SimpleList
          linkType="edit"
          primaryText={(guest) => guest.name}
          secondaryText={(guest) => guest.guestlist?.name ?? ""}
          tertiaryText={(guest) =>
            guest.arrivedAt
              ? "arrived " +
                new Date(guest.arrivedAt).toLocaleTimeString(dateLocale, {
                  hour: "2-digit",
                  minute: "2-digit",
                })
              : ""
          }
        />
      ) : (
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
            locales={dateLocale}
            options={{ weekday: "short", hour: "2-digit", minute: "2-digit" }}
          />
          <ConditionalDeleteOnOwnershipButton mutationMode="pessimistic" />
        </Datagrid>
      )}
    </List>
  );
};
