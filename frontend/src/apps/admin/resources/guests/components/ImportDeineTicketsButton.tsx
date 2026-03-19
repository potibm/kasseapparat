import React, { useRef } from "react";
import { useNotify, Button, useDataProvider } from "react-admin";
import UploadIcon from "@mui/icons-material/Upload";
import { createLogger } from "@core/logger/logger";

const log = createLogger("Admin");

export const ImportDeineTicketsButton: React.FC = () => {
  const fileInputRef = useRef<HTMLInputElement>(null);
  const dataProvider = useDataProvider();
  const notify = useNotify();

  const handleButtonClick = () => {
    fileInputRef.current?.click();
  };

  const handleFileChange = async (
    event: React.ChangeEvent<HTMLInputElement>,
  ) => {
    const file = event.target.files?.[0];
    if (!file) return;

    const formData = new FormData();
    formData.append("file", file);

    try {
      const response = await dataProvider.upload("guestsUpload", {
        data: formData,
      });
      if (response.data.createdGuests) {
        notify(
          `Success! ${response.data.createdGuests} entries have been created.`,
          { type: "info" },
        );
      } else {
        notify(
          "No entries have been created (as they might be dupes). Please check the file and try again.",
          { type: "warning" },
        );
      }
    } catch (error) {
      log.error("Error while uploading the file", error);
      notify("Error while uploading the file. Try again (later).", {
        type: "warning",
      });
    }

    event.target.value = "";
  };

  return (
    <>
      <input
        type="file"
        ref={fileInputRef}
        onChange={handleFileChange}
        style={{ display: "none" }}
        accept=".csv"
      />
      <Button
        label="Import DeineTickets.de"
        startIcon={<UploadIcon />}
        onClick={handleButtonClick}
      />
    </>
  );
};

export default ImportDeineTicketsButton;
