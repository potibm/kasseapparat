import React from "react";
import {
  List,
  Datagrid,
  TextField,
  NumberField,
  DateField,
  FunctionField,
  SimpleList,
} from "react-admin";
import { useIsSmallScreen } from "@admin/ui/hooks/useIsSmallScreen";
import { useConfig } from "@core/config/hooks/useConfig";
import { PurchaseListToolbar } from "./components/PurchaseListToolbar";
import { PurchaseStatusField } from "./components/PurchaseStatusField";
import { PaymentMethodField } from "./components/PaymentMethodField";
import { ConditionalDeleteOnAdminButton } from "../../ui/buttons/ConditionalDeleteOnAdminButton";
import { ConditionalRefundButton } from "./components/ConditionalRefundButton";
import {
  PurchaseItemLike,
  summarizePurchaseItems,
} from "./purchaseItems.utils";

export const PurchaseList: React.FC = () => {
  const { currencyOptions, currencyLocale, dateLocale } = useConfig();
  const isSmall = useIsSmallScreen();
  const price = new Intl.NumberFormat(currencyLocale, currencyOptions);

  return (
    <List sort={{ field: "createdAt", order: "DESC" }} actions={false}>
      <PurchaseListToolbar />
      {isSmall ? (
        <SimpleList
          linkType="show"
          primaryText={(record) => price.format(Number(record.totalGrossPrice))}
          secondaryText={(record) =>
            summarizePurchaseItems(record.purchaseItems)
          }
          tertiaryText={(record) =>
            new Date(record.createdAt).toLocaleTimeString(dateLocale, {
              weekday: "short",
              hour: "2-digit",
              minute: "2-digit",
            })
          }
        />
      ) : (
        <Datagrid rowClick="show" bulkActionButtons={false}>
          <FunctionField
            source="id"
            label="ID"
            render={(record) =>
              record.id ? (
                <span title={record.id}>
                  {record.id.slice(0, 4)}…{record.id.slice(-4)}
                </span>
              ) : null
            }
          />
          <DateField
            source="createdAt"
            showTime={true}
            locales={dateLocale}
            options={{ weekday: "short", hour: "2-digit", minute: "2-digit" }}
          />
          <FunctionField
            label="Items"
            sortable={false}
            render={(record: { purchaseItems?: PurchaseItemLike[] }) =>
              summarizePurchaseItems(record.purchaseItems)
            }
          />
          <NumberField
            source="totalGrossPrice"
            locales={currencyLocale}
            options={currencyOptions}
          />
          <TextField source="createdBy.username" />
          <PaymentMethodField source="paymentMethod" />
          <PurchaseStatusField source="status" label="Status" />
          <ConditionalDeleteOnAdminButton mutationMode="pessimistic" />
          <ConditionalRefundButton />
        </Datagrid>
      )}
    </List>
  );
};
