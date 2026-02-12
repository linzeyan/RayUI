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

  describe("addLog - boundary", () => {
    it("adding to exactly 1999 logs should NOT drop any", () => {
      const existing = Array.from({ length: 1999 }, (_, i) => `Line ${i}`);
      useLogStore.setState({ logs: existing });

      useLogStore.getState().addLog("Line 1999");

      const logs = useLogStore.getState().logs;
      expect(logs).toHaveLength(2000);
      expect(logs[0]).toBe("Line 0"); // Nothing dropped
      expect(logs[1999]).toBe("Line 1999");
    });

    it("adding multiple logs quickly preserves order", () => {
      useLogStore.getState().addLog("First");
      useLogStore.getState().addLog("Second");
      useLogStore.getState().addLog("Third");

      const logs = useLogStore.getState().logs;
      expect(logs).toEqual(["First", "Second", "Third"]);
    });

    it("handles empty string log lines", () => {
      useLogStore.getState().addLog("");
      expect(useLogStore.getState().logs).toEqual([""]);
    });

    it("handles very long log lines", () => {
      const longLine = "x".repeat(10000);
      useLogStore.getState().addLog(longLine);
      expect(useLogStore.getState().logs[0]).toHaveLength(10000);
    });
  });

  describe("filter and search - combined", () => {
    it("filter and search can both be set independently", () => {
      useLogStore.getState().setFilterLevel("warning");
      useLogStore.getState().setSearchQuery("timeout");

      expect(useLogStore.getState().filterLevel).toBe("warning");
      expect(useLogStore.getState().searchQuery).toBe("timeout");
    });

    it("resetting filter does not affect search", () => {
      useLogStore.getState().setFilterLevel("error");
      useLogStore.getState().setSearchQuery("test");
      useLogStore.getState().setFilterLevel("all");

      expect(useLogStore.getState().filterLevel).toBe("all");
      expect(useLogStore.getState().searchQuery).toBe("test");
    });

    it("clearing search preserves filter level", () => {
      useLogStore.getState().setFilterLevel("debug");
      useLogStore.getState().setSearchQuery("query");
      useLogStore.getState().setSearchQuery("");

      expect(useLogStore.getState().filterLevel).toBe("debug");
      expect(useLogStore.getState().searchQuery).toBe("");
    });
  });
});
