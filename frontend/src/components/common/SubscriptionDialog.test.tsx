import { describe, it, expect, vi } from "vitest";
import { render, screen } from "@/test/test-utils";
import { SubscriptionDialog } from "./SubscriptionDialog";

describe("SubscriptionDialog", () => {
  const defaultProps = {
    open: true,
    onOpenChange: vi.fn(),
    subscription: null,
    onSave: vi.fn(),
  };

  it("renders add mode when subscription is null", () => {
    render(<SubscriptionDialog {...defaultProps} />);
    expect(screen.getByRole("dialog")).toBeInTheDocument();
  });

  it("does not render when closed", () => {
    render(<SubscriptionDialog {...defaultProps} open={false} />);
    expect(screen.queryByRole("dialog")).not.toBeInTheDocument();
  });

  it("renders form fields", () => {
    render(<SubscriptionDialog {...defaultProps} />);
    // Should have name and URL inputs
    expect(screen.getByPlaceholderText("My Subscription")).toBeInTheDocument();
    expect(screen.getByPlaceholderText(/https:\/\/example.com/)).toBeInTheDocument();
  });

  it("renders filter and user agent fields", () => {
    render(<SubscriptionDialog {...defaultProps} />);
    expect(screen.getByPlaceholderText(/regex filter/)).toBeInTheDocument();
    expect(screen.getByPlaceholderText("(optional)")).toBeInTheDocument();
  });

  it("save button is disabled when name or URL is empty", () => {
    render(<SubscriptionDialog {...defaultProps} />);
    const buttons = screen.getAllByRole("button");
    const saveBtn = buttons.find(
      (btn) => btn.textContent?.includes("Save") || btn.textContent?.includes("儲存")
    );
    if (saveBtn) {
      expect(saveBtn).toBeDisabled();
    }
  });

  it("renders in edit mode when subscription is provided", () => {
    const sub = {
      id: "sub-1",
      remarks: "My Sub",
      url: "https://example.com/sub",
      enabled: true,
      sort: 0,
      autoUpdateInterval: 60,
      updateTime: 1700000000,
    };
    render(
      <SubscriptionDialog
        {...defaultProps}
        subscription={sub as never}
      />
    );
    // Should have the name pre-filled
    const nameInput = screen.getByPlaceholderText("My Subscription") as HTMLInputElement;
    expect(nameInput.value).toBe("My Sub");
  });
});
