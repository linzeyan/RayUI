import { create } from "zustand";
import {
  GetProfiles,
  AddProfile,
  UpdateProfile,
  DeleteProfiles,
  SetActiveProfile,
  ImportFromText,
  ExportShareLink,
  TestProfiles,
  TestAllProfiles,
} from "@wailsjs/go/app/App";
import { model } from "@wailsjs/go/models";

interface ProfileState {
  profiles: model.ProfileItem[];
  speedResults: Map<string, model.SpeedTestResult>;
  selectedIds: Set<string>;
  filterSubId: string;
  searchQuery: string;
  loading: boolean;

  loadProfiles: (subId?: string) => Promise<void>;
  addProfile: (profile: model.ProfileItem) => Promise<void>;
  updateProfile: (profile: model.ProfileItem) => Promise<void>;
  deleteProfiles: (ids: string[]) => Promise<void>;
  setActiveProfile: (id: string) => Promise<void>;
  importFromText: (text: string) => Promise<number>;
  exportShareLink: (id: string) => Promise<string>;
  testProfiles: (ids: string[]) => Promise<void>;
  testAllProfiles: () => Promise<void>;
  setSelectedIds: (ids: Set<string>) => void;
  toggleSelected: (id: string) => void;
  setFilterSubId: (subId: string) => void;
  setSearchQuery: (query: string) => void;
}

export const useProfileStore = create<ProfileState>((set, get) => ({
  profiles: [],
  speedResults: new Map(),
  selectedIds: new Set(),
  filterSubId: "all",
  searchQuery: "",
  loading: false,

  loadProfiles: async (subId = "all") => {
    set({ loading: true });
    try {
      // Convert "all" to empty string for backend API
      const profiles = await GetProfiles(subId === "all" ? "" : subId);
      set({ profiles: profiles || [] });
    } finally {
      set({ loading: false });
    }
  },

  addProfile: async (profile) => {
    await AddProfile(profile);
    await get().loadProfiles(get().filterSubId);
  },

  updateProfile: async (profile) => {
    await UpdateProfile(profile);
    await get().loadProfiles(get().filterSubId);
  },

  deleteProfiles: async (ids) => {
    await DeleteProfiles(ids);
    set((s) => {
      const newSelected = new Set(s.selectedIds);
      ids.forEach((id) => newSelected.delete(id));
      return { selectedIds: newSelected };
    });
    await get().loadProfiles(get().filterSubId);
  },

  setActiveProfile: async (id) => {
    await SetActiveProfile(id);
  },

  importFromText: async (text) => {
    const count = await ImportFromText(text);
    await get().loadProfiles(get().filterSubId);
    return count;
  },

  exportShareLink: async (id) => {
    return await ExportShareLink(id);
  },

  testProfiles: async (ids) => {
    const results = await TestProfiles(ids);
    set((s) => {
      const map = new Map(s.speedResults);
      results.forEach((r) => map.set(r.profileId, r));
      return { speedResults: map };
    });
  },

  testAllProfiles: async () => {
    const results = await TestAllProfiles();
    set((s) => {
      const map = new Map(s.speedResults);
      results.forEach((r) => map.set(r.profileId, r));
      return { speedResults: map };
    });
  },

  setSelectedIds: (ids) => set({ selectedIds: ids }),
  toggleSelected: (id) =>
    set((s) => {
      const newSet = new Set(s.selectedIds);
      if (newSet.has(id)) newSet.delete(id);
      else newSet.add(id);
      return { selectedIds: newSet };
    }),
  setFilterSubId: (subId) => set({ filterSubId: subId }),
  setSearchQuery: (query) => set({ searchQuery: query }),
}));
