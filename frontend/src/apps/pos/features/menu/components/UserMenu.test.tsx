import { describe, it, expect, vi, beforeEach, type Mock } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import { AppConfig } from "@core/config/types/config.types";

vi.mock("@core/config/hooks/useConfig");
vi.mock("@core/auth/authProvider", () => {
  const mockLogout = vi.fn();
  return {
    createAuthProvider: vi.fn(() => ({
      logout: mockLogout,
    })),
  };
});

const mockToggleMode = vi.fn();
let mockComputedMode: "light" | "dark" = "light";
vi.mock("flowbite-react", async (importOriginal) => ({
  ...(await importOriginal<typeof import("flowbite-react")>()),
  useThemeMode: () => ({
    computedMode: mockComputedMode,
    toggleMode: mockToggleMode,
  }),
}));

vi.mock("@pos/components/Version", () => ({
  default: () => <span>Version 1.2.3</span>,
}));

// Import after mocks are set up
import { UserMenu } from "./UserMenu";
import useConfig from "@core/config/hooks/useConfig";
import { createAuthProvider } from "@core/auth/authProvider";
import { SumupReaderContext } from "@pos/features/sumup-reader/context/SumupReaderContext";

describe("UserMenu", () => {
  const mockUseConfig = vi.mocked(useConfig);
  const mockLogout = vi.mocked(createAuthProvider).mock.results[0]?.value
    .logout as Mock;

  const setAuthMode = (authMode: "proxy" | "oidc") =>
    mockUseConfig.mockReturnValue({ authMode } as unknown as AppConfig);

  const openMenu = () =>
    fireEvent.click(
      screen.getByRole("button", { name: "User menu for cashier" }),
    );

  beforeEach(() => {
    mockComputedMode = "light";
    setAuthMode("proxy");
  });

  it("is closed by default and opens on tap", () => {
    render(<UserMenu username="cashier" />);

    const trigger = screen.getByRole("button", {
      name: "User menu for cashier",
    });
    expect(trigger).toHaveAttribute("aria-expanded", "false");
    expect(screen.queryByRole("menu")).not.toBeInTheDocument();

    openMenu();

    expect(trigger).toHaveAttribute("aria-expanded", "true");
    expect(screen.getByRole("menu")).toBeInTheDocument();
    expect(screen.getByText("Version 1.2.3")).toBeInTheDocument();
  });

  it("closes again when the trigger is tapped a second time", () => {
    render(<UserMenu username="cashier" />);

    openMenu();
    openMenu();

    expect(screen.queryByRole("menu")).not.toBeInTheDocument();
  });

  it("links to admin and manual in a new tab and closes after choosing one", () => {
    render(<UserMenu username="cashier" />);
    openMenu();

    const manual = screen.getByRole("menuitem", { name: "Manual" });
    expect(manual).toHaveAttribute("href", "/manual.pdf");
    expect(manual).toHaveAttribute("target", "_blank");

    const admin = screen.getByRole("menuitem", { name: "Admin" });
    expect(admin).toHaveAttribute("href", "/admin");
    expect(admin).toHaveAttribute("target", "_blank");

    fireEvent.click(admin);
    expect(screen.queryByRole("menu")).not.toBeInTheDocument();
  });

  it("toggles dark mode and keeps the menu open", () => {
    render(<UserMenu username="cashier" />);
    openMenu();

    fireEvent.click(screen.getByRole("menuitem", { name: "Dark mode" }));

    expect(mockToggleMode).toHaveBeenCalledTimes(1);
    expect(screen.getByRole("menu")).toBeInTheDocument();
  });

  it("offers light mode while dark mode is active", () => {
    mockComputedMode = "dark";
    render(<UserMenu username="cashier" />);
    openMenu();

    expect(
      screen.getByRole("menuitem", { name: "Light mode" }),
    ).toBeInTheDocument();
  });

  it("closes on Escape", () => {
    render(<UserMenu username="cashier" />);
    openMenu();

    fireEvent.keyDown(document, { key: "Escape" });

    expect(screen.queryByRole("menu")).not.toBeInTheDocument();
  });

  it("closes when tapping outside, but not when tapping inside", () => {
    render(
      <div>
        <p>outside</p>
        <UserMenu username="cashier" />
      </div>,
    );
    openMenu();

    fireEvent.pointerDown(screen.getByRole("menu"));
    expect(screen.getByRole("menu")).toBeInTheDocument();

    fireEvent.pointerDown(screen.getByText("outside"));
    expect(screen.queryByRole("menu")).not.toBeInTheDocument();
  });

  it("does not offer logout in proxy auth mode", () => {
    render(<UserMenu username="cashier" />);
    openMenu();

    expect(
      screen.queryByRole("menuitem", { name: "Logout" }),
    ).not.toBeInTheDocument();
  });

  it("asks for confirmation before logging out in oidc mode", () => {
    setAuthMode("oidc");
    render(<UserMenu username="cashier" />);
    openMenu();

    fireEvent.click(screen.getByRole("menuitem", { name: "Logout" }));

    expect(screen.getByText("Really log out?")).toBeInTheDocument();
    expect(mockLogout).not.toHaveBeenCalled();

    fireEvent.click(screen.getByRole("button", { name: "Cancel" }));

    expect(screen.queryByText("Really log out?")).not.toBeInTheDocument();
    expect(
      screen.getByRole("menuitem", { name: "Logout" }),
    ).toBeInTheDocument();
  });

  it("logs out after confirming", async () => {
    setAuthMode("oidc");
    render(<UserMenu username="cashier" />);
    openMenu();

    fireEvent.click(screen.getByRole("menuitem", { name: "Logout" }));
    fireEvent.click(screen.getByRole("button", { name: "Log out" }));

    await waitFor(() => {
      expect(mockLogout).toHaveBeenCalledWith({});
    });
  });

  it("still leaves the POS when the logout request fails", async () => {
    mockLogout.mockRejectedValueOnce(new Error("Logout failed"));
    setAuthMode("oidc");
    render(<UserMenu username="cashier" />);
    openMenu();

    fireEvent.click(screen.getByRole("menuitem", { name: "Logout" }));
    fireEvent.click(screen.getByRole("button", { name: "Log out" }));

    await waitFor(() => {
      expect(mockLogout).toHaveBeenCalledWith({});
    });
  });

  it("opens the card reader selection when SumUp is enabled", () => {
    const openDialog = vi.fn();

    render(
      <SumupReaderContext
        value={{
          ...readerContextDefaults,
          enabled: true,
          selectedReaderId: "rdr_1",
          selectedReader: { id: "rdr_1", name: "Bar 1", status: "paired" },
          openDialog,
        }}
      >
        <UserMenu username="cashier" />
      </SumupReaderContext>,
    );
    openMenu();

    fireEvent.click(
      screen.getByRole("menuitem", { name: "Card reader: Bar 1" }),
    );

    expect(openDialog).toHaveBeenCalled();
    expect(screen.queryByRole("menu")).not.toBeInTheDocument();
  });

  it("has no card reader entry without SumUp", () => {
    render(<UserMenu username="cashier" />);
    openMenu();

    expect(screen.queryByText(/Card reader/)).not.toBeInTheDocument();
  });
});

const readerContextDefaults = {
  enabled: false,
  selectedReaderId: undefined,
  selectedReader: undefined,
  readers: null,
  loading: false,
  problem: null,
  selectReader: () => undefined,
  refresh: () => undefined,
  dialogOpen: false,
  openDialog: () => undefined,
  closeDialog: () => undefined,
};
