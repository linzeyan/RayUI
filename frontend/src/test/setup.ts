import "@testing-library/jest-dom/vitest";
import { vi, beforeEach, afterEach } from "vitest";
import { setupWailsMocks, resetWailsMocks } from "./wails-mock";

// Setup Wails mocks before each test
beforeEach(() => {
  setupWailsMocks();
});

// Reset mocks after each test
afterEach(() => {
  resetWailsMocks();
});

// Mock window.matchMedia
Object.defineProperty(window, "matchMedia", {
  writable: true,
  value: vi.fn().mockImplementation((query) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: vi.fn(),
    removeListener: vi.fn(),
    addEventListener: vi.fn(),
    removeEventListener: vi.fn(),
    dispatchEvent: vi.fn(),
  })),
});

// Mock ResizeObserver as a proper class
class MockResizeObserver implements ResizeObserver {
  callback: ResizeObserverCallback;

  constructor(callback: ResizeObserverCallback) {
    this.callback = callback;
  }

  observe = vi.fn();
  unobserve = vi.fn();
  disconnect = vi.fn();
}

global.ResizeObserver = MockResizeObserver;

// Mock IntersectionObserver as a proper class
class MockIntersectionObserver implements IntersectionObserver {
  callback: IntersectionObserverCallback;
  root: Element | Document | null = null;
  rootMargin: string = "";
  thresholds: ReadonlyArray<number> = [];

  constructor(callback: IntersectionObserverCallback) {
    this.callback = callback;
  }

  observe = vi.fn();
  unobserve = vi.fn();
  disconnect = vi.fn();
  takeRecords = vi.fn().mockReturnValue([]);
}

global.IntersectionObserver = MockIntersectionObserver;
