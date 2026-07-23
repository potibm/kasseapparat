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
    __mockLogout: mockLogout,
  };
});

// Import after mocks are set up
import { LogoutButton } from "./LogoutButton";
import useConfig from "@core/config/hooks/useConfig";
import { createAuthProvider } from "@core/auth/authProvider";

describe("LogoutButton", () => {
  const mockUseConfig = vi.mocked(useConfig);
  const mockedAuthProvider = vi.mocked(createAuthProvider);
  const mockLogout = mockedAuthProvider.mock.results[0]?.value.logout as Mock;

  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("should not render when authMode is not oidc", () => {
    mockUseConfig.mockReturnValue({
      authMode: "proxy",
    } as unknown as AppConfig);

    const { container } = render(<LogoutButton />);
    expect(container.firstChild).toBeNull();
  });

  it("should render when authMode is oidc", () => {
    mockUseConfig.mockReturnValue({
      authMode: "oidc",
    } as unknown as AppConfig);

    render(<LogoutButton />);
    expect(screen.getByRole("button")).toBeDefined();
    expect(screen.getAllByText("Logout").length).toBeGreaterThan(0);
  });

  it("should call logout when clicked", async () => {
    mockUseConfig.mockReturnValue({
      authMode: "oidc",
    } as unknown as AppConfig);

    render(<LogoutButton />);

    const button = screen.getByRole("button");
    fireEvent.click(button);

    await waitFor(() => {
      expect(mockLogout).toHaveBeenCalledWith({});
    });
  });

  it("should redirect to home after logout even if logout fails", async () => {
    mockLogout.mockRejectedValueOnce(new Error("Logout failed"));
    mockUseConfig.mockReturnValue({
      authMode: "oidc",
    } as unknown as AppConfig);

    render(<LogoutButton />);

    const button = screen.getByRole("button");
    fireEvent.click(button);

    await waitFor(() => {
      expect(mockLogout).toHaveBeenCalledWith({});
    });
  });
});
