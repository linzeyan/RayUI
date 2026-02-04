import { create } from "zustand";
import { GetLogs, ClearLogs } from "@wailsjs/go/app/App";

interface LogState {
  logs: string[];
  filterLevel: string;
  searchQuery: string;
  autoScroll: boolean;

  loadLogs: (limit?: number) => Promise<void>;
  addLog: (line: string) => void;
  clearLogs: () => Promise<void>;
  setFilterLevel: (level: string) => void;
  setSearchQuery: (query: string) => void;
  setAutoScroll: (auto: boolean) => void;
}

export const useLogStore = create<LogState>((set, get) => ({
  logs: [],
  filterLevel: "all",
  searchQuery: "",
  autoScroll: true,

  loadLogs: async (limit = 500) => {
    const logs = await GetLogs(limit);
    set({ logs: logs || [] });
  },

  addLog: (line) =>
    set((s) => ({
      logs: [...s.logs, line].slice(-2000),
    })),

  clearLogs: async () => {
    await ClearLogs();
    set({ logs: [] });
  },

  setFilterLevel: (level) => set({ filterLevel: level }),
  setSearchQuery: (query) => set({ searchQuery: query }),
  setAutoScroll: (auto) => set({ autoScroll: auto }),
}));
