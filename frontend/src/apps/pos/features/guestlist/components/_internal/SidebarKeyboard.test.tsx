import { describe, expect, it, vi } from "vitest";
import { fireEvent, render, screen } from "@testing-library/react";
import SidebarKeyboard from "./SidebarKeyboard";

describe("SidebarKeyboard", () => {
  it("offers A–Z only, without umlauts or digits", () => {
    render(<SidebarKeyboard term="" setTerm={vi.fn()} />);

    const letters = screen
      .getAllByRole("button", { name: /^Add . to search term$/ })
      .map((button) => button.textContent);

    expect(letters).toHaveLength(26);
    expect(letters[0]).toBe("A");
    expect(letters[25]).toBe("Z");
  });

  it("types letters and spaces and deletes the last character", () => {
    const setTerm = vi.fn();
    render(<SidebarKeyboard term="JE" setTerm={setTerm} />);

    fireEvent.click(
      screen.getByRole("button", { name: "Add A to search term" }),
    );
    expect(setTerm).toHaveBeenLastCalledWith("JEA");

    fireEvent.click(
      screen.getByRole("button", { name: "Add space to search term" }),
    );
    expect(setTerm).toHaveBeenLastCalledWith("JE ");

    fireEvent.click(
      screen.getByRole("button", {
        name: "Remove last character from search term",
      }),
    );
    expect(setTerm).toHaveBeenLastCalledWith("J");
  });
});
