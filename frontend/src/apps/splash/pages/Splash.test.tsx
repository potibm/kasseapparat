import { describe, expect, it, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router";
import Splash from "./Splash";

vi.mock("../../pos/components/Version", () => ({
  default: () => <span>v-test</span>,
}));

describe("Splash", () => {
  it("links to the POS and the admin", () => {
    render(
      <MemoryRouter>
        <Splash />
      </MemoryRouter>,
    );

    expect(screen.getByRole("link", { name: /POS/ })).toHaveAttribute(
      "href",
      "/pos",
    );
    expect(screen.getByRole("link", { name: /Admin/ })).toHaveAttribute(
      "href",
      "/admin",
    );
  });
});
