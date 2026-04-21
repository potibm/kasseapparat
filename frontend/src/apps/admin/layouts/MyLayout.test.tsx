import { render, screen } from "@testing-library/react";
import { AdminContext } from "react-admin";
import { MyLayout, MyAppBar } from "./MyLayout";
import { expect, test, describe, vi } from "vitest";

vi.mock("./Menu", () => ({
  Menu: () => <div data-testid="mock-menu">Mocked Menu</div>,
}));

describe("Layout & AppBar", () => {
  test("MyAppBar renders title and logo", () => {
    render(
      <AdminContext>
        <MyAppBar />
      </AdminContext>,
    );

    expect(screen.getByText(/Kasseapparat/i)).toBeInTheDocument();

    const appBar = screen.getByRole("banner");
    expect(appBar).toHaveClass("MuiAppBar-colorSecondary");
  });

  test("MyLayout renders and displays the correct content", () => {
    render(
      <AdminContext>
        <MyLayout>
          <div data-testid="dummy-content">Content</div>
        </MyLayout>
      </AdminContext>,
    );

    expect(screen.getByText(/Kasseapparat/i)).toBeInTheDocument();

    expect(screen.getByTestId("dummy-content")).toBeInTheDocument();
  });

  test("Snapshot Match", () => {
    const { asFragment } = render(
      <AdminContext>
        <MyLayout>
          <div data-testid="dummy-content">Content</div>
        </MyLayout>
      </AdminContext>,
    );
    expect(asFragment()).toMatchSnapshot();
  });
});
