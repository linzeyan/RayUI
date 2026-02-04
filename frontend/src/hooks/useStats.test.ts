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
});
