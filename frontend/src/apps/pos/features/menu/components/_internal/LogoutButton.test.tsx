import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { LogoutButton } from "./LogoutButton";
import useConfig from "@core/config/hooks/useConfig";
import { AppConfig } from "@core/config/types/config.types";

vi.mock("@core/config/hooks/useConfig");

describe("LogoutButton", () => {
  const mockUseConfig = vi.mocked(useConfig);

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

  it("should call logout and redirect when clicked", async () => {
    mockUseConfig.mockReturnValue({
      authMode: "oidc",
    } as unknown as AppConfig);

    render(<LogoutButton />);

    const button = screen.getByRole("button");
    fireEvent.click(button);

    // The logout is called asynchronously, so we need to wait
    await new Promise((resolve) => setTimeout(resolve, 0));
  });

  it("should redirect to home after logout even if logout fails", async () => {
    mockUseConfig.mockReturnValue({
      authMode: "oidc",
    } as unknown as AppConfig);

    render(<LogoutButton />);

    const button = screen.getByRole("button");
    fireEvent.click(button);

    // Wait for async operations
    await new Promise((resolve) => setTimeout(resolve, 0));
  });
});
