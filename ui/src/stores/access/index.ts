import { produce } from "immer";
import { create } from "zustand";

import { type AccessModel } from "@/domain/access";
import { list as listAccesses, remove as removeAccess, save as saveAccess } from "@/repository/access";

export interface AccessesState {
  accesses: AccessModel[];
  loading: boolean;
  loadedAtOnce: boolean;

  fetchAccesses: () => Promise<void>;
  createAccess: (access: MaybeModelRecord<AccessModel>) => Promise<AccessModel>;
  updateAccess: (access: MaybeModelRecordWithId<AccessModel>) => Promise<AccessModel>;
  deleteAccess: (access: MaybeModelRecordWithId<AccessModel>) => Promise<AccessModel>;
}

export const useAccessesStore = create<AccessesState>((set) => {
  let fetcher: Promise<AccessModel[]> | null = null; // 防止多次重复请求

  return {
    accesses: [],
    loading: false,
    loadedAtOnce: false,

    fetchAccesses: async () => {
      fetcher ??= listAccesses().then((res) => res.items);

      try {
        set({ loading: true });
        const accesses = await fetcher;
        set({ accesses: accesses ?? [], loadedAtOnce: true });
      } finally {
        fetcher = null;
        set({ loading: false });
      }
    },

    createAccess: async (access) => {
      const record = await saveAccess(access);
      set(
        produce((state: AccessesState) => {
          state.accesses.unshift(record);
        })
      );

      return record as AccessModel;
    },

    updateAccess: async (access) => {
      const record = await saveAccess(access);
      set(
        produce((state: AccessesState) => {
          const index = state.accesses.findIndex((e) => e.id === record.id);
          if (index !== -1) {
            state.accesses[index] = record;
          }
        })
      );

      return record as AccessModel;
    },

    deleteAccess: async (access) => {
      await removeAccess(access);
      set(
        produce((state: AccessesState) => {
          state.accesses = state.accesses.filter((a) => a.id !== access.id);
        })
      );

      return access as AccessModel;
    },
  };
});
