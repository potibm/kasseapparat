import { it, expect, vi, describe, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { AdminContext, RecordContextProvider } from "react-admin";
import { MemoryRouter } from "react-router";
import CreateGuestlistEntryButton from "./CreateGuestlistEntryButton";

const mockNavigate = vi.fn();
vi.mock("react-router", async () => {
  const actual = await vi.importActual("react-router");
  return {
    ...actual,
    useNavigate: () => mockNavigate,
  };
});

describe("CreateGuestlistEntryButton", () => {
  beforeEach(() => {
    mockNavigate.mockClear();
  });

  it("renders null when no record context", () => {
    const { container } = render(
      <MemoryRouter>
        <AdminContext>
          <CreateGuestlistEntryButton />
        </AdminContext>
      </MemoryRouter>,
    );

    expect(container.firstChild).toBeNull();
  });

  it("renders button with label when record exists", () => {
    render(
      <MemoryRouter>
        <AdminContext>
          <RecordContextProvider value={{ id: 1 }}>
            <CreateGuestlistEntryButton />
          </RecordContextProvider>
        </AdminContext>
      </MemoryRouter>,
    );

    expect(screen.getByText("Add Guest")).toBeInTheDocument();
  });

  it("navigates to create guest page with guestlist_id on click", () => {
    render(
      <MemoryRouter>
        <AdminContext>
          <RecordContextProvider value={{ id: 1 }}>
            <CreateGuestlistEntryButton />
          </RecordContextProvider>
        </AdminContext>
      </MemoryRouter>,
    );

    const button = screen.getByText("Add Guest").closest("button");
    if (button) {
      fireEvent.click(button);
    }

    expect(mockNavigate).toHaveBeenCalledWith(
      "/admin/guests/create?guestlist_id=1",
    );
  });
});
