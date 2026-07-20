import React, { useState } from "react";
import { Button, useNotify } from "react-admin";
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
import { useConfig } from "@core/config/hooks/useConfig";

interface PaymentMethod {
  code: string;
  name: string;
}

interface PurchaseExportButtonProps {
  paymentMethods: PaymentMethod[];
}

export const PurchaseExportButton: React.FC<PurchaseExportButtonProps> = ({
  paymentMethods,
}) => {
  const [open, setOpen] = useState(false);
  const [selected, setSelected] = useState<string[]>([]);
  const notify = useNotify();

  const { apiBaseUrl } = useConfig();

  const togglePaymentMethod = (code: string) => {
    setSelected((prev) =>
      prev.includes(code) ? prev.filter((c) => c !== code) : [...prev, code],
    );
  };

  const handleExport = async () => {
    try {
      const params = new URLSearchParams();
      if (selected.length > 0) {
        params.set("paymentMethods", selected.join(","));
      }

      const response = await fetch(
        `${apiBaseUrl}/purchases/export?${params.toString()}`,
        {
          method: "GET",
        },
      );

      if (!response.ok) {
        throw new Error(`Status ${response.status}: ${response.statusText}`);
      }

      // Handling the file download
      const blob = await response.blob();
      const url = globalThis.URL.createObjectURL(blob);

      // Extract filename from header or fallback
      const contentDisposition = response.headers.get("Content-Disposition");
      const filenameMatch =
        contentDisposition &&
        new RegExp(/filename="([^"]+)"/).exec(contentDisposition);
      const filename = filenameMatch ? filenameMatch[1] : "purchases.csv";

      const a = document.createElement("a");
      a.href = url;
      a.download = filename;
      document.body.appendChild(a);
      a.click();
      a.remove();
      globalThis.URL.revokeObjectURL(url);

      setOpen(false);
      notify("Export started successfully", { type: "info" });
    } catch (error: unknown) {
      notify(
        `Export failed: ${error instanceof Error ? error.message : "Unknown error"}`,
        { type: "error" },
      );
    }
  };

  return (
    <>
      <Button
        label="Export CSV"
        startIcon={<FileDownloadIcon />}
        onClick={() => setOpen(true)}
      />
      <Dialog
        open={open}
        onClose={() => setOpen(false)}
        fullWidth
        maxWidth="xs"
      >
        <DialogTitle>Export Purchases</DialogTitle>
        <DialogContent>
          <FormGroup sx={{ mt: 1 }}>
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
        <DialogActions sx={{ px: 3, pb: 2 }}>
          <Button
            label="Download"
            onClick={handleExport}
            variant="contained"
            disabled={selected.length === 0}
          />
          <Button label="ra.action.cancel" onClick={() => setOpen(false)} />
        </DialogActions>
      </Dialog>
    </>
  );
};
