import React, { useState } from "react";
import { Button } from "react-admin";
import {
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormGroup,
  FormControlLabel,
  Checkbox,
} from "@mui/material";
import FileDownloadIcon from "@mui/icons-material/FileDownload";
import { useConfig } from "../../provider/ConfigProvider";
import { getAdminData } from "../authProvider";

export const PurchaseExportButton = ({ paymentMethods }) => {
  const [open, setOpen] = useState(false);
  const [selected, setSelected] = useState([]);
  const adminData = getAdminData();
  const token = adminData?.token;

  const { apiHost } = useConfig();

  const togglePaymentMethod = (code) => {
    setSelected((prev) =>
      prev.includes(code) ? prev.filter((c) => c !== code) : [...prev, code],
    );
  };

  const handleExport = async () => {
    const params = new URLSearchParams();
    if (selected.length > 0) {
      params.set("paymentMethods", selected.join(","));
    }

    const response = await fetch(
      `${apiHost}/api/v2/purchases/export?${params.toString()}`,
      {
        method: "GET",
        headers: {
          Authorization: `Bearer ${token}`,
        },
      },
    );

    if (!response.ok) {
      alert(
        `An error occurred while exporting purchases (Status ${response.status}): ${response.statusText}. Please try again later or contact support if the issue persists.`,
      );
      return;
    }

    const blob = await response.blob();
    const url = window.URL.createObjectURL(blob);
    const filename = response.headers
      .get("Content-Disposition")
      ?.match(/filename="([^"]+)"/)?.[1];

    const a = document.createElement("a");
    a.href = url;
    a.download = filename || "purchases.csv";
    document.body.appendChild(a);
    a.click();
    a.remove();
    window.URL.revokeObjectURL(url);
    setOpen(false);
  };

  return (
    <>
      <Button
        label="Export CSV"
        startIcon={<FileDownloadIcon />}
        onClick={() => setOpen(true)}
      />
      <Dialog open={open} onClose={() => setOpen(false)}>
        <DialogTitle>Export Purchases</DialogTitle>
        <DialogContent>
          <FormGroup>
            {paymentMethods.map((pm) => (
              <FormControlLabel
                key={pm.code}
                control={
                  <Checkbox
                    checked={selected.includes(pm.code)}
                    onChange={() => togglePaymentMethod(pm.code)}
                  />
                }
                label={pm.name}
              />
            ))}
          </FormGroup>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpen(false)}>Cancel</Button>
          <Button onClick={handleExport} disabled={selected.length === 0}>
            Download
          </Button>
        </DialogActions>
      </Dialog>
    </>
  );
};
