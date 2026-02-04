import { describe, it, expect, vi } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useWailsEvent } from "./useWailsEvent";
import { mockRuntime, emitWailsEvent } from "@/test/wails-mock";

describe("useWailsEvent", () => {
  it("registers event listener on mount", () => {
    const callback = vi.fn();
    renderHook(() => useWailsEvent("test:event", callback));

    expect(mockRuntime.EventsOnMultiple).toHaveBeenCalledWith(
      "test:event",
      callback,
      -1
    );
  });

  it("calls callback when event fires", () => {
    const callback = vi.fn();
    renderHook(() => useWailsEvent("test:event", callback));

    act(() => {
      emitWailsEvent("test:event", { foo: "bar" });
    });

    expect(callback).toHaveBeenCalledWith({ foo: "bar" });
  });

  it("unregisters event listener on unmount", () => {
    const callback = vi.fn();
    const { unmount } = renderHook(() =>
      useWailsEvent("test:event", callback)
    );

    unmount();

    // After unmount, emitting should not call callback again
    // (the cleanup function was called)
    act(() => {
      emitWailsEvent("test:event", "after-unmount");
    });

    expect(callback).not.toHaveBeenCalledWith("after-unmount");
  });
});
