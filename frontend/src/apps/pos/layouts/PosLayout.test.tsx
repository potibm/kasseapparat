import { describe, it, expect, afterEach } from "vitest";
import { render, screen } from "@testing-library/react";
import PosLayout, { POS_TOUCH_CLASS } from "./PosLayout";

describe("PosLayout", () => {
  afterEach(() => {
    document.documentElement.classList.remove(POS_TOUCH_CLASS);
  });

  it("renders content, sidebar, top alert and overlays", () => {
    render(
      <PosLayout
        sidebar={<div>sidebar content</div>}
        topAlert={<div>top alert</div>}
        overlays={<div>overlay content</div>}
      >
        <div>main content</div>
      </PosLayout>,
    );

    expect(screen.getByText("main content")).toBeInTheDocument();
    expect(screen.getByText("sidebar content")).toBeInTheDocument();
    expect(screen.getByText("top alert")).toBeInTheDocument();
    expect(screen.getByText("overlay content")).toBeInTheDocument();
  });

  it("sets the touch class on <html> while mounted", () => {
    expect(document.documentElement).not.toHaveClass(POS_TOUCH_CLASS);

    const { unmount } = render(
      <PosLayout sidebar={null}>
        <div>main content</div>
      </PosLayout>,
    );

    expect(document.documentElement).toHaveClass(POS_TOUCH_CLASS);

    unmount();

    expect(document.documentElement).not.toHaveClass(POS_TOUCH_CLASS);
  });
});
