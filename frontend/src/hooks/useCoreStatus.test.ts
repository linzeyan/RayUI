import { describe, it, expect, beforeEach } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useCoreStatus } from "./useCoreStatus";
import { useAppStore } from "@/stores/appStore";
import { emitWailsEvent } from "@/test/wails-mock";

describe("useCoreStatus", () => {
  beforeEach(() => {
    useAppStore.setState({ coreStatus: null });
  });

  it("updates core status on core:status event", () => {
    renderHook(() => useCoreStatus());

    const statusData = {
      running: true,
      coreType: 1,
      version: "1.8.0",
      profile: "Test Server",
    };

    act(() => {
      emitWailsEvent("core:status", statusData);
    });

    expect(useAppStore.getState().coreStatus).toEqual(statusData);
  });

  it("handles multiple status updates", () => {
    renderHook(() => useCoreStatus());

    act(() => {
      emitWailsEvent("core:status", { running: true, coreType: 1, version: "1.0" });
    });
    expect(useAppStore.getState().coreStatus?.running).toBe(true);

    act(() => {
      emitWailsEvent("core:status", { running: false, coreType: 1, version: "1.0" });
    });
    expect(useAppStore.getState().coreStatus?.running).toBe(false);
  });

  it("stops receiving events after unmount", () => {
    const { unmount } = renderHook(() => useCoreStatus());

    act(() => {
      emitWailsEvent("core:status", { running: true, coreType: 1, version: "1.0" });
    });
    expect(useAppStore.getState().coreStatus?.running).toBe(true);

    unmount();

    // Reset and emit again - should not update since hook is unmounted.
    useAppStore.setState({ coreStatus: null });
    act(() => {
      emitWailsEvent("core:status", { running: true, coreType: 2, version: "2.0" });
    });
    expect(useAppStore.getState().coreStatus).toBeNull();
  });
});
