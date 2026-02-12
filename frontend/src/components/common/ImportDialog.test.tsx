import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@/test/test-utils";
import { ImportDialog } from "./ImportDialog";

describe("ImportDialog", () => {
  const defaultProps = {
    open: true,
    onOpenChange: vi.fn(),
    onImport: vi.fn().mockResolvedValue(3),
  };

  it("renders when open", () => {
    render(<ImportDialog {...defaultProps} />);
    // Should show the import title
    expect(screen.getByRole("dialog")).toBeInTheDocument();
  });

  it("does not render dialog content when closed", () => {
    render(<ImportDialog {...defaultProps} open={false} />);
    expect(screen.queryByRole("dialog")).not.toBeInTheDocument();
  });

  it("has a textarea for paste input", () => {
    render(<ImportDialog {...defaultProps} />);
    const textarea = screen.getByPlaceholderText(/vmess:\/\//);
    expect(textarea).toBeInTheDocument();
  });

  it("has clipboard paste button", () => {
    render(<ImportDialog {...defaultProps} />);
    // There should be a button for clipboard import
    const buttons = screen.getAllByRole("button");
    expect(buttons.length).toBeGreaterThanOrEqual(2);
  });

  it("import button is disabled when textarea is empty", () => {
    render(<ImportDialog {...defaultProps} />);
    // Find footer buttons (excluding the dialog close X button).
    const buttons = screen.getAllByRole("button").filter(
      (btn) => !btn.hasAttribute("data-slot") || btn.getAttribute("data-slot") !== "dialog-close"
    );
    // The last non-close button should be the import/submit button.
    const importBtn = buttons[buttons.length - 1];
    expect(importBtn).toBeDisabled();
  });
});
