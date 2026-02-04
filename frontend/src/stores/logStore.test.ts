import { describe, it, expect, beforeEach } from "vitest";
import { useLogStore } from "./logStore";
import { mockAppApi } from "@/test/wails-mock";

describe("logStore", () => {
  beforeEach(() => {
    useLogStore.setState({
      logs: [],
      filterLevel: "all",
      searchQuery: "",
      autoScroll: true,
    });
  });

  describe("loadLogs", () => {
    it("loads logs from backend with default limit", async () => {
      const mockLogs = ["[INFO] Started", "[DEBUG] Connected"];
      mockAppApi.GetLogs.mockResolvedValue(mockLogs);

      await useLogStore.getState().loadLogs();

      expect(mockAppApi.GetLogs).toHaveBeenCalledWith(500);
      expect(useLogStore.getState().logs).toEqual(mockLogs);
    });

    it("loads logs with custom limit", async () => {
      mockAppApi.GetLogs.mockResolvedValue([]);

      await useLogStore.getState().loadLogs(100);

      expect(mockAppApi.GetLogs).toHaveBeenCalledWith(100);
    });

    it("handles null response", async () => {
      mockAppApi.GetLogs.mockResolvedValue(null);

      await useLogStore.getState().loadLogs();

      expect(useLogStore.getState().logs).toEqual([]);
    });
  });

  describe("addLog", () => {
    it("appends a log line", () => {
      useLogStore.getState().addLog("[INFO] Test");

      expect(useLogStore.getState().logs).toEqual(["[INFO] Test"]);
    });

    it("caps at 2000 lines", () => {
      const existing = Array.from({ length: 2000 }, (_, i) => `Line ${i}`);
      useLogStore.setState({ logs: existing });

      useLogStore.getState().addLog("New line");

      const logs = useLogStore.getState().logs;
      expect(logs).toHaveLength(2000);
      expect(logs[logs.length - 1]).toBe("New line");
      expect(logs[0]).toBe("Line 1"); // Line 0 dropped
    });
  });

  describe("clearLogs", () => {
    it("clears logs on backend and in store", async () => {
      useLogStore.setState({ logs: ["[INFO] Old log"] });
      mockAppApi.ClearLogs.mockResolvedValue(undefined);

      await useLogStore.getState().clearLogs();

      expect(mockAppApi.ClearLogs).toHaveBeenCalled();
      expect(useLogStore.getState().logs).toEqual([]);
    });
  });

  describe("filter and search", () => {
    it("setFilterLevel updates level", () => {
      useLogStore.getState().setFilterLevel("error");
      expect(useLogStore.getState().filterLevel).toBe("error");
    });

    it("setSearchQuery updates query", () => {
      useLogStore.getState().setSearchQuery("connection");
      expect(useLogStore.getState().searchQuery).toBe("connection");
    });

    it("setAutoScroll updates auto scroll", () => {
      useLogStore.getState().setAutoScroll(false);
      expect(useLogStore.getState().autoScroll).toBe(false);
    });
  });
});
