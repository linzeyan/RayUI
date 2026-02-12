import { describe, it, expect, beforeEach } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useCoreLog } from "./useCoreLog";
import { useLogStore } from "@/stores/logStore";
import { emitWailsEvent } from "@/test/wails-mock";

describe("useCoreLog", () => {
  beforeEach(() => {
    useLogStore.setState({ logs: [] });
  });

  it("adds log line on core:log event", () => {
    renderHook(() => useCoreLog());

    act(() => {
      emitWailsEvent("core:log", "[INFO] Connection established");
    });

    expect(useLogStore.getState().logs).toContain(
      "[INFO] Connection established"
    );
  });

  it("accumulates multiple log events", () => {
    renderHook(() => useCoreLog());

    act(() => {
      emitWailsEvent("core:log", "Line 1");
      emitWailsEvent("core:log", "Line 2");
    });

    expect(useLogStore.getState().logs).toEqual(["Line 1", "Line 2"]);
  });

  it("stops receiving events after unmount", () => {
    const { unmount } = renderHook(() => useCoreLog());

    act(() => {
      emitWailsEvent("core:log", "Before unmount");
    });
    expect(useLogStore.getState().logs).toHaveLength(1);

    unmount();

    act(() => {
      emitWailsEvent("core:log", "After unmount");
    });
    // After unmount the listener should be cleaned up,
    // so no new log should be added.
    expect(useLogStore.getState().logs).toHaveLength(1);
  });
});
