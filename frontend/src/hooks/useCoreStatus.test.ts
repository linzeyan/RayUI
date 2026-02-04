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
});
