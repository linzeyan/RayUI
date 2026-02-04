import { describe, it, expect, beforeEach } from "vitest";
import { screen, waitFor } from "@testing-library/react";
import { render, userEvent } from "@/test/test-utils";
import { ProfilesPage } from "./ProfilesPage";
import { useProfileStore } from "@/stores/profileStore";
import { useSubscriptionStore } from "@/stores/subscriptionStore";
import { useAppStore } from "@/stores/appStore";
import { mockAppApi, mockConfig } from "@/test/wails-mock";

describe("ProfilesPage", () => {
  beforeEach(() => {
    // Reset stores
    useProfileStore.setState({
      profiles: [],
      speedResults: new Map(),
      selectedIds: new Set(),
      filterSubId: "all",
      searchQuery: "",
      loading: false,
    });
    useSubscriptionStore.setState({
      subscriptions: [],
      loading: false,
    });
    useAppStore.setState({
      config: mockConfig,
      coreStatus: null,
    });
  });

  it("renders without crashing", async () => {
    render(<ProfilesPage />);

    // Wait for initial load
    await waitFor(() => {
      expect(mockAppApi.GetProfiles).toHaveBeenCalled();
    });

    // Should show empty state (text from en.json)
    expect(screen.getByText("No profiles yet")).toBeInTheDocument();
  });

  it("renders toolbar with all buttons", async () => {
    render(<ProfilesPage />);

    await waitFor(() => {
      expect(screen.getByText("Add Profile")).toBeInTheDocument();
    });

    expect(screen.getByText("Import")).toBeInTheDocument();
    expect(screen.getByText("Test All")).toBeInTheDocument();
    expect(
      screen.getByPlaceholderText("Search profiles...")
    ).toBeInTheDocument();
  });

  it("displays profiles when available", async () => {
    const mockProfiles = [
      {
        id: "1",
        remarks: "Server 1",
        address: "server1.example.com",
        port: 443,
        configType: 1, // VMess
        network: "tcp",
        subId: "",
      },
      {
        id: "2",
        remarks: "Server 2",
        address: "server2.example.com",
        port: 8443,
        configType: 2, // VLESS
        network: "ws",
        subId: "",
      },
    ];

    mockAppApi.GetProfiles.mockResolvedValue(mockProfiles);

    render(<ProfilesPage />);

    await waitFor(() => {
      expect(screen.getByText("Server 1")).toBeInTheDocument();
    });

    expect(screen.getByText("Server 2")).toBeInTheDocument();
    expect(screen.getByText("server1.example.com")).toBeInTheDocument();
    expect(screen.getByText("server2.example.com")).toBeInTheDocument();
  });

  it("filters profiles by search query", async () => {
    const mockProfiles = [
      {
        id: "1",
        remarks: "Japan Server",
        address: "jp.example.com",
        port: 443,
        configType: 1,
        network: "tcp",
        subId: "",
      },
      {
        id: "2",
        remarks: "US Server",
        address: "us.example.com",
        port: 443,
        configType: 1,
        network: "tcp",
        subId: "",
      },
    ];

    mockAppApi.GetProfiles.mockResolvedValue(mockProfiles);

    const user = userEvent.setup();
    render(<ProfilesPage />);

    await waitFor(() => {
      expect(screen.getByText("Japan Server")).toBeInTheDocument();
    });

    // Search for Japan
    const searchInput = screen.getByPlaceholderText("Search profiles...");
    await user.type(searchInput, "Japan");

    // Only Japan Server should be visible
    expect(screen.getByText("Japan Server")).toBeInTheDocument();
    expect(screen.queryByText("US Server")).not.toBeInTheDocument();
  });

  it("opens add profile dialog when clicking Add Profile button", async () => {
    const user = userEvent.setup();
    render(<ProfilesPage />);

    await waitFor(() => {
      expect(screen.getByText("Add Profile")).toBeInTheDocument();
    });

    await user.click(screen.getByText("Add Profile"));

    // Dialog should be open
    await waitFor(() => {
      expect(screen.getByRole("dialog")).toBeInTheDocument();
    });
  });

  it("opens import dialog when clicking Import button", async () => {
    const user = userEvent.setup();
    render(<ProfilesPage />);

    await waitFor(() => {
      expect(screen.getByText("Import")).toBeInTheDocument();
    });

    await user.click(screen.getByText("Import"));

    // Dialog should be open
    await waitFor(() => {
      expect(screen.getByRole("dialog")).toBeInTheDocument();
    });
  });

  it("calls TestAllProfiles when clicking Test All button", async () => {
    const mockProfiles = [
      {
        id: "1",
        remarks: "Server 1",
        address: "server1.example.com",
        port: 443,
        configType: 1,
        network: "tcp",
        subId: "",
      },
    ];

    mockAppApi.GetProfiles.mockResolvedValue(mockProfiles);
    mockAppApi.TestAllProfiles.mockResolvedValue([]);

    const user = userEvent.setup();
    render(<ProfilesPage />);

    await waitFor(() => {
      expect(screen.getByText("Server 1")).toBeInTheDocument();
    });

    await user.click(screen.getByText("Test All"));

    await waitFor(() => {
      expect(mockAppApi.TestAllProfiles).toHaveBeenCalled();
    });
  });

  it("shows selection bar when profiles are selected", async () => {
    const mockProfiles = [
      {
        id: "1",
        remarks: "Server 1",
        address: "server1.example.com",
        port: 443,
        configType: 1,
        network: "tcp",
        subId: "",
      },
    ];

    mockAppApi.GetProfiles.mockResolvedValue(mockProfiles);

    const user = userEvent.setup();
    render(<ProfilesPage />);

    await waitFor(() => {
      expect(screen.getByText("Server 1")).toBeInTheDocument();
    });

    // Click checkbox to select
    const checkboxes = screen.getAllByRole("checkbox");
    await user.click(checkboxes[1]); // First checkbox is "select all"

    // Selection bar should appear
    await waitFor(() => {
      expect(screen.getByText("1 selected")).toBeInTheDocument();
    });
    expect(screen.getByText("Test Selected")).toBeInTheDocument();
    expect(screen.getByText("Delete Selected")).toBeInTheDocument();
  });
});
