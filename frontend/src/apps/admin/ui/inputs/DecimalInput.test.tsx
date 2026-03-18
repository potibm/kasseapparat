import { it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { AdminContext, SimpleForm } from "react-admin";
import DecimalInput from "./DecimalInput";

it("should integrate with react-admin context", () => {
  render(
    <AdminContext>
      <SimpleForm record={{ price: "100.5" }} toolbar={false}>
        <DecimalInput source="price" label="Price" />
      </SimpleForm>
    </AdminContext>
  );

  expect(screen.getByLabelText("Price")).toHaveValue("100,5");
});