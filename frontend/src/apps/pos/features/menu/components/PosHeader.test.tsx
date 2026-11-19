import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import PosHeader from "./PosHeader";

vi.mock("./UserMenu", () => ({
  default: ({ username }: { username: string }) => (
    <span>menu for {username}</span>
  ),
}));

describe("PosHeader", () => {
  it("shows the app name and the user menu", () => {
    render(<PosHeader username="cashier" />);

    expect(screen.getByTestId("pos-header")).toHaveTextContent("Kasseapparat");
    expect(screen.getByText("menu for cashier")).toBeInTheDocument();
    expect(screen.queryByRole("status")).not.toBeInTheDocument();
  });

  it("shows the environment message when one is configured", () => {
    render(
      <PosHeader username="cashier" environmentMessage="Test environment" />,
    );

    expect(screen.getByRole("status")).toHaveTextContent("Test environment");
  });
});
