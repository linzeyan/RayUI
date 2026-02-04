import { describe, it, expect, beforeEach } from "vitest";
import { screen, waitFor } from "@testing-library/react";
import { render, userEvent } from "@/test/test-utils";
import { SubscriptionsPage } from "./SubscriptionsPage";
import { useSubscriptionStore } from "@/stores/subscriptionStore";
import { useProfileStore } from "@/stores/profileStore";
import { mockAppApi } from "@/test/wails-mock";

describe("SubscriptionsPage", () => {
  beforeEach(() => {
    useSubscriptionStore.setState({
      subscriptions: [],
      syncingIds: new Set(),
      loading: false,
    });
    useProfileStore.setState({
      profiles: [],
      loading: false,
    });
  });

  it("renders without crashing", async () => {
    render(<SubscriptionsPage />);

    await waitFor(() => {
      expect(screen.getByText("Subscriptions")).toBeInTheDocument();
    });
  });

  it("shows empty state when no subscriptions", async () => {
    render(<SubscriptionsPage />);

    await waitFor(() => {
      expect(screen.getByText("No subscriptions yet")).toBeInTheDocument();
    });
  });

  it("renders toolbar buttons", async () => {
    render(<SubscriptionsPage />);

    await waitFor(() => {
      expect(screen.getByText("Add Subscription")).toBeInTheDocument();
      expect(screen.getByText("Update All")).toBeInTheDocument();
    });
  });

  it("displays subscriptions when available", async () => {
    const mockSubs = [
      {
        id: "s1",
        remarks: "My Subscription",
        url: "https://example.com/sub",
        enabled: true,
        sort: 0,
        autoUpdateInterval: 60,
        updateTime: 0,
      },
    ];
    mockAppApi.GetSubscriptions.mockResolvedValue(mockSubs);

    render(<SubscriptionsPage />);

    await waitFor(() => {
      expect(screen.getByText("My Subscription")).toBeInTheDocument();
    });
    expect(screen.getByText("https://example.com/sub")).toBeInTheDocument();
  });

  it("opens add subscription dialog", async () => {
    const user = userEvent.setup();
    render(<SubscriptionsPage />);

    await waitFor(() => {
      expect(screen.getByText("Add Subscription")).toBeInTheDocument();
    });

    await user.click(screen.getByText("Add Subscription"));

    await waitFor(() => {
      expect(screen.getByRole("dialog")).toBeInTheDocument();
    });
  });

  it("shows server count per subscription", async () => {
    const mockSubs = [
      {
        id: "s1",
        remarks: "Sub 1",
        url: "https://example.com",
        enabled: true,
        sort: 0,
        autoUpdateInterval: 60,
        updateTime: 0,
      },
    ];
    const mockProfiles = [
      { id: "p1", subId: "s1", remarks: "Server 1" },
      { id: "p2", subId: "s1", remarks: "Server 2" },
    ];
    mockAppApi.GetSubscriptions.mockResolvedValue(mockSubs);
    mockAppApi.GetProfiles.mockResolvedValue(mockProfiles);

    render(<SubscriptionsPage />);

    await waitFor(() => {
      expect(screen.getByText(/2 servers/)).toBeInTheDocument();
    });
  });
});
