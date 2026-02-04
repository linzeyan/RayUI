import { describe, it, expect, beforeEach } from "vitest";
import { useRoutingStore } from "./routingStore";
import { mockAppApi } from "@/test/wails-mock";

describe("routingStore", () => {
  beforeEach(() => {
    useRoutingStore.setState({
      routings: [],
      loading: false,
    });
  });

  describe("loadRoutings", () => {
    it("loads routings from backend", async () => {
      const mockRoutings = [
        { id: "r1", remarks: "Global", enabled: true },
        { id: "r2", remarks: "BypassCN", enabled: false },
      ];
      mockAppApi.GetRoutings.mockResolvedValue(mockRoutings);

      await useRoutingStore.getState().loadRoutings();

      expect(mockAppApi.GetRoutings).toHaveBeenCalled();
      expect(useRoutingStore.getState().routings).toEqual(mockRoutings);
    });

    it("handles null response", async () => {
      mockAppApi.GetRoutings.mockResolvedValue(null);

      await useRoutingStore.getState().loadRoutings();

      expect(useRoutingStore.getState().routings).toEqual([]);
    });

    it("sets loading state correctly", async () => {
      mockAppApi.GetRoutings.mockImplementation(
        () =>
          new Promise((resolve) => {
            expect(useRoutingStore.getState().loading).toBe(true);
            resolve([]);
          })
      );

      await useRoutingStore.getState().loadRoutings();

      expect(useRoutingStore.getState().loading).toBe(false);
    });
  });

  describe("addRouting", () => {
    it("adds routing and reloads", async () => {
      const routing = { id: "r1", remarks: "Custom" };
      mockAppApi.AddRouting.mockResolvedValue(undefined);
      mockAppApi.GetRoutings.mockResolvedValue([routing]);

      await useRoutingStore.getState().addRouting(routing as never);

      expect(mockAppApi.AddRouting).toHaveBeenCalledWith(routing);
      expect(mockAppApi.GetRoutings).toHaveBeenCalled();
    });
  });

  describe("updateRouting", () => {
    it("updates routing and reloads", async () => {
      const routing = { id: "r1", remarks: "Updated" };
      mockAppApi.UpdateRouting.mockResolvedValue(undefined);
      mockAppApi.GetRoutings.mockResolvedValue([routing]);

      await useRoutingStore.getState().updateRouting(routing as never);

      expect(mockAppApi.UpdateRouting).toHaveBeenCalledWith(routing);
    });
  });

  describe("deleteRouting", () => {
    it("deletes routing and reloads", async () => {
      mockAppApi.DeleteRouting.mockResolvedValue(undefined);
      mockAppApi.GetRoutings.mockResolvedValue([]);

      await useRoutingStore.getState().deleteRouting("r1");

      expect(mockAppApi.DeleteRouting).toHaveBeenCalledWith("r1");
    });
  });

  describe("setActiveRouting", () => {
    it("sets active routing and reloads", async () => {
      mockAppApi.SetActiveRouting.mockResolvedValue(undefined);
      mockAppApi.GetRoutings.mockResolvedValue([]);

      await useRoutingStore.getState().setActiveRouting("r1");

      expect(mockAppApi.SetActiveRouting).toHaveBeenCalledWith("r1");
      expect(mockAppApi.GetRoutings).toHaveBeenCalled();
    });
  });
});
