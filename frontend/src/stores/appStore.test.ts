import { describe, it, expect, beforeEach } from "vitest";
import { useAppStore } from "./appStore";
import { mockAppApi, mockConfig, mockCoreStatus } from "@/test/wails-mock";

describe("appStore", () => {
  beforeEach(() => {
    useAppStore.setState({
      currentPage: "profiles",
      config: null,
      coreStatus: null,
      traffic: { upload: 0, download: 0, totalUpload: 0, totalDownload: 0 },
      loading: false,
    });
  });

  describe("setCurrentPage", () => {
    it("updates current page", () => {
      useAppStore.getState().setCurrentPage("settings");
      expect(useAppStore.getState().currentPage).toBe("settings");
    });
  });

  describe("setCoreStatus", () => {
    it("updates core status", () => {
      const status = { running: true, coreType: 0, version: "1.0.0" };
      useAppStore.getState().setCoreStatus(status as never);
      expect(useAppStore.getState().coreStatus).toEqual(status);
    });
  });

  describe("setTraffic", () => {
    it("updates traffic data", () => {
      const traffic = {
        upload: 1000,
        download: 2000,
        totalUpload: 5000,
        totalDownload: 10000,
      };
      useAppStore.getState().setTraffic(traffic);
      expect(useAppStore.getState().traffic).toEqual(traffic);
    });
  });

  describe("loadConfig", () => {
    it("loads config from backend", async () => {
      mockAppApi.GetConfig.mockResolvedValue(mockConfig);

      await useAppStore.getState().loadConfig();

      expect(mockAppApi.GetConfig).toHaveBeenCalled();
      expect(useAppStore.getState().config).toEqual(mockConfig);
    });
  });

  describe("updateConfig", () => {
    it("updates config in backend and store", async () => {
      const newConfig = { ...mockConfig, proxyMode: 1 };
      mockAppApi.UpdateConfig.mockResolvedValue(undefined);

      await useAppStore.getState().updateConfig(newConfig as never);

      expect(mockAppApi.UpdateConfig).toHaveBeenCalledWith(newConfig);
      expect(useAppStore.getState().config).toEqual(newConfig);
    });
  });

  describe("setProxyMode", () => {
    it("updates proxy mode", async () => {
      useAppStore.setState({ config: mockConfig as never });
      mockAppApi.SetProxyMode.mockResolvedValue(undefined);

      await useAppStore.getState().setProxyMode(2);

      expect(mockAppApi.SetProxyMode).toHaveBeenCalledWith(2);
      expect(useAppStore.getState().config?.proxyMode).toBe(2);
    });
  });

  describe("loadCoreStatus", () => {
    it("loads core status from backend", async () => {
      mockAppApi.GetCoreStatus.mockResolvedValue(mockCoreStatus);

      await useAppStore.getState().loadCoreStatus();

      expect(mockAppApi.GetCoreStatus).toHaveBeenCalled();
      expect(useAppStore.getState().coreStatus).toEqual(mockCoreStatus);
    });
  });
});
