import React, { useRef } from "react";
import {
  TopToolbar,
  ExportButton,
  useNotify,
  FilterButton,
  CreateButton,
  Button,
  useDataProvider,
} from "react-admin";
import UploadIcon from "@mui/icons-material/Upload";

const ImportDeineTicketsButton = () => {
  const fileInputRef = useRef(null);
  const dataProvider = useDataProvider();
  const notify = useNotify();

  const handleButtonClick = () => {
    fileInputRef.current.click();
  };

  const handleFileChange = async (event) => {
    const file = event.target.files[0];
    if (!file) {
      return;
    }

    const formData = new FormData();
    formData.append("file", file);

    try {
      const response = await dataProvider.upload("listEntriesUpload", {
        data: formData,
      });
      if (response.data.createdEntries) {
        notify(
          `Success! ${response.data.createdEntries} entries have been created.`,
          "info",
        );
      } else {
        notify(
          "No entries have been created (as they might be dupes). Please check the file and try again.",
          "warning",
        );
      }
    } catch (error) {
      notify("Error while uploading the file. Try again (later).", "warning");
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
      />
      <Button label="Import DeineTickets.de" onClick={handleButtonClick}>
        <UploadIcon />
      </Button>
    </>
  );
};

const ListEntryActions = (props) => (
  <TopToolbar>
    <FilterButton />
    <CreateButton />
    <ExportButton {...props} />
    <ImportDeineTicketsButton />
  </TopToolbar>
);

export default ListEntryActions;
