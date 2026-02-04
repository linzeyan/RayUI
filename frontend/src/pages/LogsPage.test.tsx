import { describe, it, expect, beforeEach } from "vitest";
import { screen, waitFor } from "@testing-library/react";
import { render } from "@/test/test-utils";
import { LogsPage } from "./LogsPage";
import { useLogStore } from "@/stores/logStore";
import { mockAppApi } from "@/test/wails-mock";

describe("LogsPage", () => {
  beforeEach(() => {
    useLogStore.setState({
      logs: [],
      filterLevel: "all",
      searchQuery: "",
      autoScroll: true,
    });
    mockAppApi.GetConnections.mockResolvedValue({
      downloadTotal: 0,
      uploadTotal: 0,
      connections: [],
    });
  });

  it("renders without crashing", async () => {
    render(<LogsPage />);

    await waitFor(() => {
      expect(screen.getByText("Logs")).toBeInTheDocument();
    });
  });

  it("renders tabs for Logs and Connections", () => {
    render(<LogsPage />);

    expect(screen.getByText("Logs")).toBeInTheDocument();
    expect(screen.getByText("Connections")).toBeInTheDocument();
  });

  it("shows empty state when no logs", async () => {
    render(<LogsPage />);

    await waitFor(() => {
      expect(screen.getByText("No logs yet")).toBeInTheDocument();
    });
  });

  it("renders log toolbar elements", async () => {
    render(<LogsPage />);

    await waitFor(() => {
      expect(screen.getByText("Copy")).toBeInTheDocument();
      expect(screen.getByText("Clear")).toBeInTheDocument();
      expect(screen.getByPlaceholderText("Search logs...")).toBeInTheDocument();
    });
  });

  it("displays log lines", async () => {
    const mockLogs = [
      "[INFO] 2024-01-01 Connected to server",
      "[ERROR] 2024-01-01 Connection failed",
    ];
    mockAppApi.GetLogs.mockResolvedValue(mockLogs);

    render(<LogsPage />);

    await waitFor(() => {
      expect(
        screen.getByText("[INFO] 2024-01-01 Connected to server")
      ).toBeInTheDocument();
    });
    expect(
      screen.getByText("[ERROR] 2024-01-01 Connection failed")
    ).toBeInTheDocument();
  });
});
