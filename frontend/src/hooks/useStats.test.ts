import { describe, it, expect, beforeEach } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useStats } from "./useStats";
import { useAppStore } from "@/stores/appStore";
import { emitWailsEvent } from "@/test/wails-mock";

describe("useStats", () => {
  beforeEach(() => {
    useAppStore.setState({
      traffic: { upload: 0, download: 0, totalUpload: 0, totalDownload: 0 },
    });
  });

  it("updates traffic on stats:traffic event", () => {
    renderHook(() => useStats());

    const trafficData = {
      upload: 1024,
      download: 2048,
      totalUpload: 10000,
      totalDownload: 20000,
    };

    act(() => {
      emitWailsEvent("stats:traffic", trafficData);
    });

    expect(useAppStore.getState().traffic).toEqual(trafficData);
  });

  it("handles multiple traffic updates", () => {
    renderHook(() => useStats());

    act(() => {
      emitWailsEvent("stats:traffic", { upload: 100, download: 200, totalUpload: 100, totalDownload: 200 });
    });
    expect(useAppStore.getState().traffic.upload).toBe(100);

    act(() => {
      emitWailsEvent("stats:traffic", { upload: 500, download: 1000, totalUpload: 600, totalDownload: 1200 });
    });
    expect(useAppStore.getState().traffic.upload).toBe(500);
    expect(useAppStore.getState().traffic.download).toBe(1000);
  });

  it("handles zero traffic", () => {
    renderHook(() => useStats());

    act(() => {
      emitWailsEvent("stats:traffic", { upload: 0, download: 0, totalUpload: 0, totalDownload: 0 });
    });
    expect(useAppStore.getState().traffic.upload).toBe(0);
    expect(useAppStore.getState().traffic.download).toBe(0);
  });

  it("stops receiving events after unmount", () => {
    const { unmount } = renderHook(() => useStats());

    act(() => {
      emitWailsEvent("stats:traffic", { upload: 100, download: 200, totalUpload: 100, totalDownload: 200 });
    });
    expect(useAppStore.getState().traffic.upload).toBe(100);

    unmount();

    // Reset and emit again.
    useAppStore.setState({
      traffic: { upload: 0, download: 0, totalUpload: 0, totalDownload: 0 },
    });
    act(() => {
      emitWailsEvent("stats:traffic", { upload: 999, download: 999, totalUpload: 999, totalDownload: 999 });
    });
    // Should not have been updated after unmount.
    expect(useAppStore.getState().traffic.upload).toBe(0);
  });
});
