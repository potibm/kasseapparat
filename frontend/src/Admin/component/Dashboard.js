// src/MyDashboard.js
import * as React from "react";
import { useEffect, useState } from "react";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import { Title, useDataProvider } from "react-admin";
import { Table } from "flowbite-react";
import { useConfig } from "../../provider/ConfigProvider";
import Decimal from "decimal.js";

const Dashboard = () => {
  return (
    <>
      <Title title="Kasseapparat" />
      <ProductStatsCard />
    </>
  );
};

const ProductStatsCard = () => {
  const [stats, setStats] = useState([]);
  const dataProvider = useDataProvider();
  const currency = useConfig().currency;

  useEffect(() => {
    // Fetch the stats
    dataProvider
      .getList("productStats", {
        pagination: { page: 1, perPage: 100 },
        sort: { field: "date", order: "DESC" },
        filter: {},
      })
      .then(({ data }) => {
        setStats(data);
      });
  }, [dataProvider]);

  if (stats.length === 0) {
    return <div>Loading...</div>;
  }

  const customCompactTheme = {
    head: {
      cell: {
        base: "px-3 py-2",
      },
    },
    body: {
      cell: {
        base: "px-3 py-2",
      },
    },
  };

  return (
    <Card>
      <CardContent>
        <Table className="mt-5" theme={customCompactTheme}>
          <Table.Head>
            <Table.HeadCell>Product</Table.HeadCell>
            <Table.HeadCell className="text-right">Units sold</Table.HeadCell>
            <Table.HeadCell className="text-right">Revenue</Table.HeadCell>
          </Table.Head>
          <Table.Body className="divide-y">
            {stats.map((stat) => (
              <Table.Row key={stat.id}>
                <Table.Cell>{stat.name}</Table.Cell>
                <Table.Cell className="text-right">{stat.soldItems}</Table.Cell>
                <Table.Cell className="text-right">
                  {currency.format(stat.totalPrice)}
                </Table.Cell>
              </Table.Row>
            ))}
            <Table.Row>
              <Table.Cell className="font-bold">Total</Table.Cell>
              <Table.Cell className="font-bold text-right">-</Table.Cell>
              <Table.Cell className="font-bold text-right">
                {currency.format(
                  stats.reduce(
                    (acc, stat) => acc.add(stat.totalPrice),
                    new Decimal(0),
                  ),
                )}
              </Table.Cell>
            </Table.Row>
          </Table.Body>
        </Table>
      </CardContent>
    </Card>
  );
};

export default Dashboard;
