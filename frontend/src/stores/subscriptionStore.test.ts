import { describe, it, expect, beforeEach } from "vitest";
import { useSubscriptionStore } from "./subscriptionStore";
import { mockAppApi } from "@/test/wails-mock";

describe("subscriptionStore", () => {
  beforeEach(() => {
    useSubscriptionStore.setState({
      subscriptions: [],
      syncingIds: new Set(),
      loading: false,
    });
  });

  describe("loadSubscriptions", () => {
    it("loads subscriptions from backend", async () => {
      const mockSubs = [
        { id: "s1", remarks: "Sub 1", url: "https://example.com/sub1" },
        { id: "s2", remarks: "Sub 2", url: "https://example.com/sub2" },
      ];
      mockAppApi.GetSubscriptions.mockResolvedValue(mockSubs);

      await useSubscriptionStore.getState().loadSubscriptions();

      expect(mockAppApi.GetSubscriptions).toHaveBeenCalled();
      expect(useSubscriptionStore.getState().subscriptions).toEqual(mockSubs);
    });

    it("handles null response", async () => {
      mockAppApi.GetSubscriptions.mockResolvedValue(null);

      await useSubscriptionStore.getState().loadSubscriptions();

      expect(useSubscriptionStore.getState().subscriptions).toEqual([]);
    });

    it("sets loading state correctly", async () => {
      mockAppApi.GetSubscriptions.mockImplementation(
        () =>
          new Promise((resolve) => {
            expect(useSubscriptionStore.getState().loading).toBe(true);
            resolve([]);
          })
      );

      await useSubscriptionStore.getState().loadSubscriptions();

      expect(useSubscriptionStore.getState().loading).toBe(false);
    });
  });

  describe("addSubscription", () => {
    it("adds subscription and reloads", async () => {
      const sub = { id: "s1", remarks: "New Sub", url: "https://example.com" };
      mockAppApi.AddSubscription.mockResolvedValue(undefined);
      mockAppApi.GetSubscriptions.mockResolvedValue([sub]);

      await useSubscriptionStore.getState().addSubscription(sub as never);

      expect(mockAppApi.AddSubscription).toHaveBeenCalledWith(sub);
      expect(mockAppApi.GetSubscriptions).toHaveBeenCalled();
    });
  });

  describe("updateSubscription", () => {
    it("updates subscription and reloads", async () => {
      const sub = { id: "s1", remarks: "Updated Sub", url: "https://example.com" };
      mockAppApi.UpdateSubscription.mockResolvedValue(undefined);
      mockAppApi.GetSubscriptions.mockResolvedValue([sub]);

      await useSubscriptionStore.getState().updateSubscription(sub as never);

      expect(mockAppApi.UpdateSubscription).toHaveBeenCalledWith(sub);
    });
  });

  describe("deleteSubscription", () => {
    it("deletes subscription and reloads", async () => {
      mockAppApi.DeleteSubscription.mockResolvedValue(undefined);
      mockAppApi.GetSubscriptions.mockResolvedValue([]);

      await useSubscriptionStore.getState().deleteSubscription("s1");

      expect(mockAppApi.DeleteSubscription).toHaveBeenCalledWith("s1");
    });
  });

  describe("syncSubscription", () => {
    it("syncs subscription and returns count", async () => {
      mockAppApi.SyncSubscription.mockResolvedValue(10);
      mockAppApi.GetSubscriptions.mockResolvedValue([]);

      const count = await useSubscriptionStore.getState().syncSubscription("s1");

      expect(mockAppApi.SyncSubscription).toHaveBeenCalledWith("s1");
      expect(count).toBe(10);
    });

    it("tracks syncing state per subscription", async () => {
      let resolveSync: (value: number) => void;
      mockAppApi.SyncSubscription.mockImplementation(
        () =>
          new Promise<number>((resolve) => {
            resolveSync = resolve;
          })
      );
      mockAppApi.GetSubscriptions.mockResolvedValue([]);

      const promise = useSubscriptionStore.getState().syncSubscription("s1");

      expect(useSubscriptionStore.getState().syncingIds.has("s1")).toBe(true);

      resolveSync!(5);
      await promise;

      expect(useSubscriptionStore.getState().syncingIds.has("s1")).toBe(false);
    });
  });

  describe("syncAllSubscriptions", () => {
    it("syncs all subscriptions", async () => {
      useSubscriptionStore.setState({
        subscriptions: [
          { id: "s1", remarks: "Sub 1" },
          { id: "s2", remarks: "Sub 2" },
        ] as never[],
      });

      const results = { s1: 5, s2: 3 };
      mockAppApi.SyncAllSubscriptions.mockResolvedValue(results);
      mockAppApi.GetSubscriptions.mockResolvedValue([]);

      const ret = await useSubscriptionStore.getState().syncAllSubscriptions();

      expect(mockAppApi.SyncAllSubscriptions).toHaveBeenCalled();
      expect(ret).toEqual(results);
    });

    it("clears syncing state after completion", async () => {
      useSubscriptionStore.setState({
        subscriptions: [{ id: "s1" }] as never[],
      });
      mockAppApi.SyncAllSubscriptions.mockResolvedValue({});
      mockAppApi.GetSubscriptions.mockResolvedValue([]);

      await useSubscriptionStore.getState().syncAllSubscriptions();

      expect(useSubscriptionStore.getState().syncingIds.size).toBe(0);
    });
  });
});
