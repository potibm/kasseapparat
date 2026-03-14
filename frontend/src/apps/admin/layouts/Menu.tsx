import React from "react";
import { Menu as RaMenu, MenuItemLink, useSidebarState } from "react-admin";
import { useConfig } from "@core/config/providers/ConfigProvider";
import { MenuDivider, MenuMessage } from "./MenuComponents";
import { SumupReaderIcon } from "../components/SumupReader";
import { SumupTransactionIcon } from "../components/SumupTransaction";

export const Menu: React.FC = () => {
  const { environmentMessage, sumupEnabled } = useConfig();
  const [isSidebarOpen] = useSidebarState();

  return (
    <RaMenu>
      <RaMenu.DashboardItem />

      <MenuDivider name="POS" />
      <RaMenu.ResourceItem name="products" />
      <RaMenu.ResourceItem name="purchases" />
      <RaMenu.ResourceItem name="productInterests" />

      <MenuDivider name="Guestlist" />
      <RaMenu.ResourceItem name="guestlists" />
      <RaMenu.ResourceItem name="guests" />

      <MenuDivider name="Sumup" />
      {sumupEnabled
        ? [
            <RaMenu.ResourceItem name="sumupReaders" key="sumupReaders" />,
            <RaMenu.ResourceItem
              name="sumupTransactions"
              key="sumupTransactions"
            />,
          ]
        : [
            <MenuItemLink
              key="sumupReadersDisabled"
              to="#"
              primaryText="Readers (Disabled)"
              leftIcon={<SumupReaderIcon />}
              disabled
            />,
            <MenuItemLink
              key="sumupTransactionsDisabled"
              to="#"
              primaryText="Transactions (Disabled)"
              leftIcon={<SumupTransactionIcon />}
              disabled
            />,
          ]}

      <MenuDivider name="Admin" />
      <RaMenu.ResourceItem name="users" />

      {environmentMessage && isSidebarOpen && (
        <MenuMessage message={environmentMessage} />
      )}
    </RaMenu>
  );
};
