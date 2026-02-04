import { create } from "zustand";
import { GetConfig, GetCoreStatus, UpdateConfig, SetProxyMode } from "@wailsjs/go/app/App";
import { model } from "@wailsjs/go/models";

export type Page =
  | "profiles"
  | "subscriptions"
  | "routing"
  | "dns"
  | "settings"
  | "logs";

interface TrafficData {
  upload: number;
  download: number;
  totalUpload: number;
  totalDownload: number;
}

interface AppState {
  currentPage: Page;
  config: model.Config | null;
  coreStatus: model.CoreStatus | null;
  traffic: TrafficData;
  loading: boolean;

  setCurrentPage: (page: Page) => void;
  setCoreStatus: (status: model.CoreStatus) => void;
  setTraffic: (traffic: TrafficData) => void;

  loadConfig: () => Promise<void>;
  updateConfig: (config: model.Config) => Promise<void>;
  setProxyMode: (mode: number) => Promise<void>;
  loadCoreStatus: () => Promise<void>;
}

export const useAppStore = create<AppState>((set, get) => ({
  currentPage: "profiles",
  config: null,
  coreStatus: null,
  traffic: { upload: 0, download: 0, totalUpload: 0, totalDownload: 0 },
  loading: false,

  setCurrentPage: (page) => set({ currentPage: page }),
  setCoreStatus: (status) => set({ coreStatus: status }),
  setTraffic: (traffic) => set({ traffic }),

  loadConfig: async () => {
    const config = await GetConfig();
    set({ config });
  },

  updateConfig: async (config) => {
    await UpdateConfig(config);
    set({ config });
  },

  setProxyMode: async (mode) => {
    await SetProxyMode(mode);
    const config = get().config;
    if (config) {
      set({ config: { ...config, proxyMode: mode } as model.Config });
    }
  },

  loadCoreStatus: async () => {
    const status = await GetCoreStatus();
    set({ coreStatus: status });
  },
}));
