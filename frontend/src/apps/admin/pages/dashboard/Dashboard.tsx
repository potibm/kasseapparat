import * as React from "react";
import { useEffect, useState } from "react";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import { Title, useDataProvider, RaRecord } from "react-admin";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeadCell,
  TableRow,
} from "flowbite-react";
import { useConfig } from "@core/config/providers/ConfigProvider";
import Decimal from "decimal.js";
import { Typography } from "@mui/material";

interface ProductStat extends RaRecord {
  name: string;
  soldItems: number;
  totalNetPrice: string | number;
  totalGrossPrice: string | number;
}

const Dashboard: React.FC = () => {
  return (
    <>
      <Title title="Kasseapparat" />
      <ProductStatsCard />
    </>
  );
};

const ProductStatsCard: React.FC = () => {
  const [stats, setStats] = useState<ProductStat[] | null>(null);
  const dataProvider = useDataProvider();
  const { currency } = useConfig(); // Destructuring ist cleaner

  useEffect(() => {
    dataProvider
      .getList<ProductStat>("productStats", {
        pagination: { page: 1, perPage: 100 },
        sort: { field: "name", order: "ASC" },
        filter: {},
      })
      .then(({ data }) => {
        setStats(data);
      })
      .catch((error) => {
        console.error("Dashboard fetch failed", error);
        setStats([]);
      });
  }, [dataProvider]);

  if (stats === null) return <Typography sx={{ p: 2 }}>Loading...</Typography>;
  if (stats.length === 0)
    return <Typography sx={{ p: 2 }}>No products yet.</Typography>;

  const customCompactTheme = {
    head: { cell: { base: "px-3 py-2" } },
    body: { cell: { base: "px-3 py-2" } },
  };

  // Summen-Berechnung mit Decimal.js
  const totalNet = stats.reduce(
    (acc, stat) => acc.add(new Decimal(stat.totalNetPrice)),
    new Decimal(0),
  );

  const totalGross = stats.reduce(
    (acc, stat) => acc.add(new Decimal(stat.totalGrossPrice)),
    new Decimal(0),
  );

  return (
    <Card sx={{ mt: 2 }}>
      <CardContent>
        <Typography variant="h6" gutterBottom>
          Product Sales Stats
        </Typography>
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
                <TableCell className="font-medium">{stat.name}</TableCell>
                <TableCell className="text-right">{stat.soldItems}</TableCell>
                <TableCell className="text-right">
                  {currency.format(new Decimal(stat.totalNetPrice).toNumber())}
                </TableCell>
                <TableCell className="text-right">
                  {currency.format(
                    new Decimal(stat.totalGrossPrice).toNumber(),
                  )}
                </TableCell>
              </TableRow>
            ))}
            <TableRow className="bg-gray-50 dark:bg-gray-800">
              <TableCell className="font-bold">Total</TableCell>
              <TableCell className="text-right">-</TableCell>
              <TableCell className="font-bold text-right">
                {currency.format(totalNet.toNumber())}
              </TableCell>
              <TableCell className="font-bold text-right">
                {currency.format(totalGross.toNumber())}
              </TableCell>
            </TableRow>
          </TableBody>
        </Table>
      </CardContent>
    </Card>
  );
};

export default Dashboard;
