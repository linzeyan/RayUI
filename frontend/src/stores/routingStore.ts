import { create } from "zustand";
import {
  GetRoutings,
  AddRouting,
  UpdateRouting,
  DeleteRouting,
  SetActiveRouting,
} from "@wailsjs/go/app/App";
import { model } from "@wailsjs/go/models";

interface RoutingState {
  routings: model.RoutingItem[];
  loading: boolean;

  loadRoutings: () => Promise<void>;
  addRouting: (routing: model.RoutingItem) => Promise<void>;
  updateRouting: (routing: model.RoutingItem) => Promise<void>;
  deleteRouting: (id: string) => Promise<void>;
  setActiveRouting: (id: string) => Promise<void>;
}

export const useRoutingStore = create<RoutingState>((set, get) => ({
  routings: [],
  loading: false,

  loadRoutings: async () => {
    set({ loading: true });
    try {
      const routings = await GetRoutings();
      set({ routings: routings || [] });
    } finally {
      set({ loading: false });
    }
  },

  addRouting: async (routing) => {
    await AddRouting(routing);
    await get().loadRoutings();
  },

  updateRouting: async (routing) => {
    await UpdateRouting(routing);
    await get().loadRoutings();
  },

  deleteRouting: async (id) => {
    await DeleteRouting(id);
    await get().loadRoutings();
  },

  setActiveRouting: async (id) => {
    await SetActiveRouting(id);
    await get().loadRoutings();
  },
}));
