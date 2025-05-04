// src/MyDashboard.js
import * as React from "react";
import { useEffect, useState } from "react";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import { Title, useDataProvider } from "react-admin";
import {
  Table,
  TableCell,
  TableHead,
  TableHeadCell,
  TableBody,
  TableRow,
} from "flowbite-react";
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
  const [stats, setStats] = useState(null);
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

  if (stats === null) {
    return <div>Loading...</div>;
  }
  else if (stats.length === 0) {
    return <div>No products, yet.</div>;
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
          <TableHead>
            <TableRow>
              <TableHeadCell>Product</TableHeadCell>
              <TableHeadCell className="text-right">Units sold</TableHeadCell>
              <TableHeadCell className="text-right">Revenue Net</TableHeadCell>
              <TableHeadCell className="text-right">
                Revenue Gross
              </TableHeadCell>
            </TableRow>
          </TableHead>
          <TableBody className="divide-y">
            {stats.map((stat) => (
              <TableRow key={stat.id}>
                <TableCell>{stat.name}</TableCell>
                <TableCell className="text-right">{stat.soldItems}</TableCell>
                <TableCell className="text-right">
                  {currency.format(stat.totalNetPrice)}
                </TableCell>
                <TableCell className="text-right">
                  {currency.format(stat.totalGrossPrice)}
                </TableCell>
              </TableRow>
            ))}
            <TableRow>
              <TableCell className="font-bold">Total</TableCell>
              <TableCell className="font-bold text-right">-</TableCell>
              <TableCell className="font-bold text-right">
                {currency.format(
                  stats.reduce(
                    (acc, stat) => acc.add(stat.totalNetPrice),
                    new Decimal(0),
                  ),
                )}
              </TableCell>
              <TableCell className="font-bold text-right">
                {currency.format(
                  stats.reduce(
                    (acc, stat) => acc.add(stat.totalGrossPrice),
                    new Decimal(0),
                  ),
                )}
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
};

export default Dashboard;
