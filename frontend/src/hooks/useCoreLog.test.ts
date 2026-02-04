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
});
