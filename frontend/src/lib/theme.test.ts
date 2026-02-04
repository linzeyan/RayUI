import { describe, it, expect, vi, beforeEach } from "vitest";
import { applyTheme, watchSystemTheme } from "./theme";

describe("applyTheme", () => {
  beforeEach(() => {
    document.documentElement.classList.remove("dark");
  });

  it("applies dark theme", () => {
    applyTheme("dark");
    expect(document.documentElement.classList.contains("dark")).toBe(true);
  });

  it("applies light theme", () => {
    document.documentElement.classList.add("dark");
    applyTheme("light");
    expect(document.documentElement.classList.contains("dark")).toBe(false);
  });

  it("applies system theme based on matchMedia", () => {
    // Default mock returns matches: false (light)
    applyTheme("system");
    expect(document.documentElement.classList.contains("dark")).toBe(false);
  });
});

describe("watchSystemTheme", () => {
  it("returns no-op cleanup for non-system theme", () => {
    const cleanup = watchSystemTheme("dark");
    expect(typeof cleanup).toBe("function");
    cleanup(); // should not throw
  });

  it("registers and cleans up change listener for system theme", () => {
    // Create a stable mock for matchMedia that tracks addEventListener
    const addListener = vi.fn();
    const removeListener = vi.fn();
    const mockMq = {
      matches: false,
      media: "(prefers-color-scheme: dark)",
      onchange: null,
      addListener: vi.fn(),
      removeListener: vi.fn(),
      addEventListener: addListener,
      removeEventListener: removeListener,
      dispatchEvent: vi.fn(),
    };

    const origMatchMedia = window.matchMedia;
    window.matchMedia = vi.fn().mockReturnValue(mockMq);

    const cleanup = watchSystemTheme("system");

    expect(addListener).toHaveBeenCalledWith("change", expect.any(Function));

    cleanup();

    expect(removeListener).toHaveBeenCalledWith("change", expect.any(Function));

    // Restore
    window.matchMedia = origMatchMedia;
  });

  it("does not add listener for non-system theme", () => {
    const addListener = vi.fn();
    const mockMq = {
      matches: false,
      media: "(prefers-color-scheme: dark)",
      onchange: null,
      addListener: vi.fn(),
      removeListener: vi.fn(),
      addEventListener: addListener,
      removeEventListener: vi.fn(),
      dispatchEvent: vi.fn(),
    };

    const origMatchMedia = window.matchMedia;
    window.matchMedia = vi.fn().mockReturnValue(mockMq);

    watchSystemTheme("light");

    expect(addListener).not.toHaveBeenCalled();

    // Restore
    window.matchMedia = origMatchMedia;
  });
});
