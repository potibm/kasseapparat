import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";
import { CriticalError } from "./CriticalError";

describe("CriticalError", () => {
  it("should render title and message", () => {
    render(<CriticalError title="Test Title" message="Test Message" />);

    expect(screen.getByText("Test Title")).toBeDefined();
    expect(screen.getByText("Test Message")).toBeDefined();
  });

  it("should render details when provided", () => {
    render(
      <CriticalError
        title="Test Title"
        message="Test Message"
        details="Test Details"
      />,
    );

    expect(screen.getByText("Test Details")).toBeDefined();
  });

  it("should not render details when not provided", () => {
    render(<CriticalError title="Test Title" message="Test Message" />);

    expect(screen.queryByRole("pre")).toBeNull();
  });

  it("should render with correct structure", () => {
    const { container } = render(
      <CriticalError title="Error Title" message="Error Message" />,
    );

    // Check for the outer container with flex layout
    const outerDiv = container.firstChild as HTMLElement;
    expect(outerDiv.className).toContain("flex");
    expect(outerDiv.className).toContain("h-screen");
    expect(outerDiv.className).toContain("items-center");
    expect(outerDiv.className).toContain("justify-center");

    // Check for Alert component
    expect(screen.getByRole("alert")).toBeDefined();
  });

  it("should render multiple lines in details", () => {
    const multiLineDetails = "Line 1\nLine 2\nLine 3";
    const { container } = render(
      <CriticalError
        title="Error"
        message="Message"
        details={multiLineDetails}
      />,
    );

    // Check that the pre element contains the multi-line text
    const preElement = container.querySelector("pre");
    expect(preElement).not.toBeNull();
    expect(preElement?.textContent).toBe(multiLineDetails);
  });

  it("should handle empty strings gracefully", () => {
    render(<CriticalError title="" message="" />);

    // Should still render the component structure
    expect(screen.getByRole("alert")).toBeDefined();
  });

  it("should handle special characters in message", () => {
    const specialMessage = "Error: <script>alert('xss')</script>";
    render(
      <CriticalError title="Title" message={specialMessage} />,
    );

    // React should escape the HTML
    expect(screen.getByText(specialMessage)).toBeDefined();
  });
});
