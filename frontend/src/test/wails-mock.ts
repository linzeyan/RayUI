import { vi } from "vitest";
import type { model, app, service } from "@wailsjs/go/models";

// Default mock data
export const mockConfig: model.Config = {
  activeProfileId: "",
  activeRoutingId: "",
  activeDnsPreset: "default",
  coreBasic: {
    logEnabled: true,
    logLevel: "warning",
    muxEnabled: false,
    allowInsecure: false,
    fingerprint: "chrome",
    enableFragment: false,
  },
  inbounds: [
    {
      protocol: "socks",
      listenAddr: "127.0.0.1",
      port: 10808,
      udpEnabled: true,
      sniffingEnabled: true,
      allowLAN: false,
    },
    {
      protocol: "http",
      listenAddr: "127.0.0.1",
      port: 10809,
      udpEnabled: false,
      sniffingEnabled: true,
      allowLAN: false,
    },
  ],
  proxyMode: 0,
  tun: {
    enabled: false,
    autoRoute: true,
    strictRoute: false,
    stack: "system",
    mtu: 9000,
    enableIPv6: false,
  },
  systemProxy: {
    exceptions: "localhost,127.0.0.1",
    notProxyLocal: true,
  },
  ui: {
    theme: "system",
    language: "en",
    fontFamily: "system-ui",
    fontSize: 14,
    autoHideOnStart: false,
    closeToTray: true,
    showInDock: true,
  },
  speedTest: {
    timeout: 10,
    url: "http://www.gstatic.com/generate_204",
    pingUrl: "https://www.google.com",
    concurrent: 5,
  },
} as unknown as model.Config;

export const mockCoreStatus: model.CoreStatus = {
  running: false,
  coreType: 0,
  version: "1.0.0",
  startTime: undefined,
  pid: undefined,
  profile: undefined,
} as unknown as model.CoreStatus;

export const mockSystemInfo: app.SystemInfo = {
  os: "darwin",
  arch: "arm64",
  appVersion: "0.1.0",
} as unknown as app.SystemInfo;

export const mockProfiles: model.ProfileItem[] = [];
export const mockSubscriptions: model.SubItem[] = [];
export const mockRoutings: model.RoutingItem[] = [];
export const mockStats: model.ServerStatItem[] = [];

// Event listeners store
const eventListeners: Map<string, ((...data: unknown[]) => void)[]> = new Map();

// Mock App API functions
export const mockAppApi = {
  AddProfile: vi.fn().mockResolvedValue(undefined),
  AddRouting: vi.fn().mockResolvedValue(undefined),
  AddSubscription: vi.fn().mockResolvedValue(undefined),
  CheckCoreUpdate: vi.fn().mockResolvedValue({
    coreType: 0,
    currentVersion: "1.0.0",
    latestVersion: "1.0.0",
    hasUpdate: false,
    downloadUrl: "",
    assetName: "",
  }),
  CheckGeoDataUpdate: vi.fn().mockResolvedValue(false),
  ClearLogs: vi.fn().mockResolvedValue(undefined),
  CloseAllConnections: vi.fn().mockResolvedValue(undefined),
  CloseConnection: vi.fn().mockResolvedValue(undefined),
  Context: vi.fn().mockResolvedValue({}),
  DeleteProfiles: vi.fn().mockResolvedValue(undefined),
  DeleteRouting: vi.fn().mockResolvedValue(undefined),
  DeleteSubscription: vi.fn().mockResolvedValue(undefined),
  DownloadCoreUpdate: vi.fn().mockResolvedValue(undefined),
  EnsureGeoData: vi.fn().mockResolvedValue(undefined),
  ExportShareLink: vi.fn().mockResolvedValue(""),
  GetConfig: vi.fn().mockResolvedValue(mockConfig),
  GetConnections: vi.fn().mockResolvedValue({
    downloadTotal: 0,
    uploadTotal: 0,
    connections: [],
  } as service.ConnectionsResponse),
  GetCoreStatus: vi.fn().mockResolvedValue(mockCoreStatus),
  GetDNSConfig: vi.fn().mockResolvedValue({
    remoteDns: "https://1.1.1.1/dns-query",
    directDns: "https://223.5.5.5/dns-query",
    bootstrapDns: "8.8.8.8",
    useSystemHosts: true,
    fakeIP: false,
    hosts: "",
    domainStrategy: "AsIs",
  }),
  GetGeoDataInfo: vi.fn().mockResolvedValue({
    geoipVersion: "",
    geositeVersion: "",
    geoipPath: "",
    geositePath: "",
    lastUpdated: 0,
  }),
  GetLogs: vi.fn().mockResolvedValue([]),
  GetProfiles: vi.fn().mockResolvedValue(mockProfiles),
  GetRoutings: vi.fn().mockResolvedValue(mockRoutings),
  GetStats: vi.fn().mockResolvedValue(mockStats),
  GetSubscriptions: vi.fn().mockResolvedValue(mockSubscriptions),
  GetSystemInfo: vi.fn().mockResolvedValue(mockSystemInfo),
  ImportFromText: vi.fn().mockResolvedValue(0),
  IsAutoStartEnabled: vi.fn().mockResolvedValue(false),
  IsCoreInstalled: vi.fn().mockResolvedValue(false),
  RequestQuit: vi.fn().mockResolvedValue(undefined),
  ResetAllStats: vi.fn().mockResolvedValue(undefined),
  ResetStats: vi.fn().mockResolvedValue(undefined),
  RestartCore: vi.fn().mockResolvedValue(undefined),
  SetActiveProfile: vi.fn().mockResolvedValue(undefined),
  SetActiveRouting: vi.fn().mockResolvedValue(undefined),
  SetAutoStart: vi.fn().mockResolvedValue(undefined),
  SetProxyMode: vi.fn().mockResolvedValue(undefined),
  ShouldCloseToTray: vi.fn().mockResolvedValue(true),
  Shutdown: vi.fn().mockResolvedValue(undefined),
  StartCore: vi.fn().mockResolvedValue(undefined),
  Startup: vi.fn().mockResolvedValue(undefined),
  StopCore: vi.fn().mockResolvedValue(undefined),
  SyncAllSubscriptions: vi.fn().mockResolvedValue({}),
  SyncSubscription: vi.fn().mockResolvedValue(0),
  TestAllProfiles: vi.fn().mockResolvedValue([]),
  TestProfiles: vi.fn().mockResolvedValue([]),
  ToggleCore: vi.fn().mockResolvedValue(undefined),
  UpdateConfig: vi.fn().mockResolvedValue(undefined),
  UpdateDNSConfig: vi.fn().mockResolvedValue(undefined),
  UpdateGeoData: vi.fn().mockResolvedValue(undefined),
  UpdateProfile: vi.fn().mockResolvedValue(undefined),
  UpdateRouting: vi.fn().mockResolvedValue(undefined),
  UpdateSubscription: vi.fn().mockResolvedValue(undefined),
};

// Mock Wails runtime functions
export const mockRuntime = {
  LogPrint: vi.fn(),
  LogTrace: vi.fn(),
  LogDebug: vi.fn(),
  LogInfo: vi.fn(),
  LogWarning: vi.fn(),
  LogError: vi.fn(),
  LogFatal: vi.fn(),
  EventsOnMultiple: vi.fn(
    (
      eventName: string,
      callback: (...data: unknown[]) => void,
      _maxCallbacks: number
    ) => {
      const listeners = eventListeners.get(eventName) || [];
      listeners.push(callback);
      eventListeners.set(eventName, listeners);
      return () => {
        const currentListeners = eventListeners.get(eventName) || [];
        const index = currentListeners.indexOf(callback);
        if (index > -1) {
          currentListeners.splice(index, 1);
        }
      };
    }
  ),
  EventsOff: vi.fn((eventName: string) => {
    eventListeners.delete(eventName);
  }),
  EventsOffAll: vi.fn(() => {
    eventListeners.clear();
  }),
  EventsEmit: vi.fn((eventName: string, ...data: unknown[]) => {
    const listeners = eventListeners.get(eventName) || [];
    listeners.forEach((cb) => cb(...data));
  }),
  WindowReload: vi.fn(),
  WindowReloadApp: vi.fn(),
  WindowSetAlwaysOnTop: vi.fn(),
  WindowSetSystemDefaultTheme: vi.fn(),
  WindowSetLightTheme: vi.fn(),
  WindowSetDarkTheme: vi.fn(),
  WindowCenter: vi.fn(),
  WindowSetTitle: vi.fn(),
  WindowFullscreen: vi.fn(),
  WindowUnfullscreen: vi.fn(),
  WindowIsFullscreen: vi.fn().mockResolvedValue(false),
  WindowGetSize: vi.fn().mockResolvedValue({ w: 1024, h: 768 }),
  WindowSetSize: vi.fn(),
  WindowSetMaxSize: vi.fn(),
  WindowSetMinSize: vi.fn(),
  WindowSetPosition: vi.fn(),
  WindowGetPosition: vi.fn().mockResolvedValue({ x: 0, y: 0 }),
  WindowHide: vi.fn(),
  WindowShow: vi.fn(),
  WindowMaximise: vi.fn(),
  WindowToggleMaximise: vi.fn(),
  WindowUnmaximise: vi.fn(),
  WindowIsMaximised: vi.fn().mockResolvedValue(false),
  WindowMinimise: vi.fn(),
  WindowUnminimise: vi.fn(),
  WindowSetBackgroundColour: vi.fn(),
  ScreenGetAll: vi.fn().mockResolvedValue([]),
  WindowIsMinimised: vi.fn().mockResolvedValue(false),
  WindowIsNormal: vi.fn().mockResolvedValue(true),
  BrowserOpenURL: vi.fn(),
  Environment: vi.fn().mockResolvedValue({
    buildType: "production",
    platform: "darwin",
    arch: "arm64",
  }),
  Quit: vi.fn(),
  Hide: vi.fn(),
  Show: vi.fn(),
  ClipboardGetText: vi.fn().mockResolvedValue(""),
  ClipboardSetText: vi.fn().mockResolvedValue(true),
  OnFileDrop: vi.fn(),
  OnFileDropOff: vi.fn(),
  CanResolveFilePaths: vi.fn().mockReturnValue(false),
  ResolveFilePaths: vi.fn(),
};

// Setup global window mocks
export function setupWailsMocks() {
  // Reset all mocks
  Object.values(mockAppApi).forEach((fn) => fn.mockClear());
  Object.values(mockRuntime).forEach((fn) => {
    if (typeof fn.mockClear === "function") {
      fn.mockClear();
    }
  });
  eventListeners.clear();

  // Setup window.go.app.App
  (window as Record<string, unknown>).go = {
    app: {
      App: mockAppApi,
    },
  };

  // Setup window.runtime
  (window as Record<string, unknown>).runtime = mockRuntime;
}

// Helper to emit events in tests
export function emitWailsEvent(eventName: string, ...data: unknown[]) {
  mockRuntime.EventsEmit(eventName, ...data);
}

// Helper to reset mocks between tests
export function resetWailsMocks() {
  Object.values(mockAppApi).forEach((fn) => fn.mockClear());
  Object.values(mockRuntime).forEach((fn) => {
    if (typeof fn.mockClear === "function") {
      fn.mockClear();
    }
  });
  eventListeners.clear();
}
