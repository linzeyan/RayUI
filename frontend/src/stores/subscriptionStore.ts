import { create } from "zustand";
import {
  GetSubscriptions,
  AddSubscription,
  UpdateSubscription,
  DeleteSubscription,
  SyncSubscription,
  SyncAllSubscriptions,
} from "@wailsjs/go/app/App";
import { model } from "@wailsjs/go/models";

interface SubscriptionState {
  subscriptions: model.SubItem[];
  syncingIds: Set<string>;
  loading: boolean;

  loadSubscriptions: () => Promise<void>;
  addSubscription: (sub: model.SubItem) => Promise<void>;
  updateSubscription: (sub: model.SubItem) => Promise<void>;
  deleteSubscription: (id: string) => Promise<void>;
  syncSubscription: (id: string) => Promise<number>;
  syncAllSubscriptions: () => Promise<Record<string, number>>;
}

export const useSubscriptionStore = create<SubscriptionState>((set, get) => ({
  subscriptions: [],
  syncingIds: new Set(),
  loading: false,

  loadSubscriptions: async () => {
    set({ loading: true });
    try {
      const subs = await GetSubscriptions();
      set({ subscriptions: subs || [] });
    } finally {
      set({ loading: false });
    }
  },

  addSubscription: async (sub) => {
    await AddSubscription(sub);
    await get().loadSubscriptions();
  },

  updateSubscription: async (sub) => {
    await UpdateSubscription(sub);
    await get().loadSubscriptions();
  },

  deleteSubscription: async (id) => {
    await DeleteSubscription(id);
    await get().loadSubscriptions();
  },

  syncSubscription: async (id) => {
    set((s) => ({ syncingIds: new Set(s.syncingIds).add(id) }));
    try {
      const count = await SyncSubscription(id);
      await get().loadSubscriptions();
      return count;
    } finally {
      set((s) => {
        const ids = new Set(s.syncingIds);
        ids.delete(id);
        return { syncingIds: ids };
      });
    }
  },

  syncAllSubscriptions: async () => {
    const subs = get().subscriptions;
    const allIds = new Set(subs.map((s) => s.id));
    set({ syncingIds: allIds });
    try {
      const results = await SyncAllSubscriptions();
      await get().loadSubscriptions();
      return results;
    } finally {
      set({ syncingIds: new Set() });
    }
  },
}));
