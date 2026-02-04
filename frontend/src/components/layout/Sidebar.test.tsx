import { describe, it, expect, vi, beforeEach } from "vitest";
import { screen, waitFor } from "@testing-library/react";
import { render, userEvent } from "@/test/test-utils";
import { Sidebar } from "./Sidebar";
import { useAppStore } from "@/stores/appStore";
import { mockAppApi } from "@/test/wails-mock";

describe("Sidebar", () => {
  beforeEach(() => {
    // Reset store state
    useAppStore.setState({
      currentPage: "profiles",
      coreStatus: null,
      config: null,
    });
  });

  it("renders without crashing", () => {
    render(<Sidebar />);
    expect(screen.getByText("RayUI")).toBeInTheDocument();
  });

  it("renders all navigation items", () => {
    render(<Sidebar />);

    // Check for navigation items by their translated text
    expect(screen.getByText("Profiles")).toBeInTheDocument();
    expect(screen.getByText("Subscriptions")).toBeInTheDocument();
    expect(screen.getByText("Routing")).toBeInTheDocument();
    expect(screen.getByText("DNS")).toBeInTheDocument();
    expect(screen.getByText("Settings")).toBeInTheDocument();
    expect(screen.getByText("Logs")).toBeInTheDocument();
  });

  it("highlights current page", () => {
    useAppStore.setState({ currentPage: "settings" });
    render(<Sidebar />);

    const settingsButton = screen.getByRole("button", { name: /settings/i });
    expect(settingsButton).toHaveClass("bg-sidebar-accent");
  });

  it("changes page on navigation click", async () => {
    const user = userEvent.setup();
    render(<Sidebar />);

    const subscriptionsButton = screen.getByRole("button", {
      name: /subscriptions/i,
    });
    await user.click(subscriptionsButton);

    expect(useAppStore.getState().currentPage).toBe("subscriptions");
  });

  it("shows disconnected status when core is not running", () => {
    useAppStore.setState({
      coreStatus: { running: false, coreType: 0, version: "1.0.0" },
    });
    render(<Sidebar />);

    expect(screen.getByText("Disconnected")).toBeInTheDocument();
  });

  it("shows connected status when core is running", () => {
    useAppStore.setState({
      coreStatus: {
        running: true,
        coreType: 0,
        version: "1.0.0",
        profile: "Test Server",
      },
    });
    render(<Sidebar />);

    expect(screen.getByText("Connected")).toBeInTheDocument();
    expect(screen.getByText("Test Server")).toBeInTheDocument();
  });

  it("calls StartCore when clicking toggle button while disconnected", async () => {
    const user = userEvent.setup();
    useAppStore.setState({
      coreStatus: { running: false, coreType: 0, version: "1.0.0" },
    });

    render(<Sidebar />);

    const toggleButton = screen.getByRole("button", { name: /disconnected/i });
    await user.click(toggleButton);

    await waitFor(() => {
      expect(mockAppApi.StartCore).toHaveBeenCalled();
    });
  });

  it("calls StopCore when clicking toggle button while connected", async () => {
    const user = userEvent.setup();
    useAppStore.setState({
      coreStatus: { running: true, coreType: 0, version: "1.0.0" },
    });

    render(<Sidebar />);

    const toggleButton = screen.getByRole("button", { name: /connected/i });
    await user.click(toggleButton);

    await waitFor(() => {
      expect(mockAppApi.StopCore).toHaveBeenCalled();
    });
  });
});
