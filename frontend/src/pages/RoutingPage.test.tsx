import { describe, it, expect, beforeEach } from "vitest";
import { screen, waitFor } from "@testing-library/react";
import { render, userEvent } from "@/test/test-utils";
import { RoutingPage } from "./RoutingPage";
import { useRoutingStore } from "@/stores/routingStore";
import { useAppStore } from "@/stores/appStore";
import { mockAppApi, mockConfig } from "@/test/wails-mock";

describe("RoutingPage", () => {
  beforeEach(() => {
    useRoutingStore.setState({
      routings: [],
      loading: false,
    });
    useAppStore.setState({
      config: mockConfig,
    });
  });

  it("renders without crashing", async () => {
    render(<RoutingPage />);

    await waitFor(() => {
      expect(screen.getByText("Routing")).toBeInTheDocument();
    });
  });

  it("shows empty state when no routings", async () => {
    render(<RoutingPage />);

    await waitFor(() => {
      expect(screen.getByText("No routing rules")).toBeInTheDocument();
    });
  });

  it("renders add button", () => {
    render(<RoutingPage />);
    expect(screen.getByText("Add Rule Set")).toBeInTheDocument();
  });

  it("displays routing items", async () => {
    const mockRoutings = [
      {
        id: "r1",
        remarks: "Global Proxy",
        domainStrategy: "prefer_ipv4",
        rules: [{ id: "rule1" }],
        enabled: true,
        locked: true,
        sort: 0,
      },
      {
        id: "r2",
        remarks: "Custom Rules",
        domainStrategy: "prefer_ipv4",
        rules: [],
        enabled: true,
        locked: false,
        sort: 1,
      },
    ];
    mockAppApi.GetRoutings.mockResolvedValue(mockRoutings);

    render(<RoutingPage />);

    await waitFor(() => {
      expect(screen.getByText("Global Proxy")).toBeInTheDocument();
    });
    expect(screen.getByText("Custom Rules")).toBeInTheDocument();
    expect(screen.getByText("Built-in")).toBeInTheDocument();
  });

  it("opens add dialog when clicking Add Rule Set", async () => {
    const user = userEvent.setup();
    render(<RoutingPage />);

    await user.click(screen.getByText("Add Rule Set"));

    await waitFor(() => {
      expect(screen.getByRole("dialog")).toBeInTheDocument();
    });
  });
});
