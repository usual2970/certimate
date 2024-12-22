import { create } from "zustand";
import { produce } from "immer";

import { type AccessModel } from "@/domain/access";
import { list as listAccess, save as saveAccess, remove as removeAccess } from "@/repository/access";

export interface AccessState {
  initialized: boolean;
  accesses: AccessModel[];
  createAccess: (access: MaybeModelRecord<AccessModel>) => void;
  updateAccess: (access: MaybeModelRecordWithId<AccessModel>) => void;
  deleteAccess: (access: MaybeModelRecordWithId<AccessModel>) => void;
  fetchAccesses: () => Promise<void>;
}

export const useAccessStore = create<AccessState>((set) => {
  let fetcher: Promise<AccessModel[]> | null = null; // 防止多次重复请求

  return {
    initialized: false,
    accesses: [],

    createAccess: async (access) => {
      const record = await saveAccess(access);

      set(
        produce((state: AccessState) => {
          state.accesses.unshift(record);
        })
      );
    },

    updateAccess: async (access) => {
      const record = await saveAccess(access);

      set(
        produce((state: AccessState) => {
          const index = state.accesses.findIndex((e) => e.id === record.id);
          if (index !== -1) {
            state.accesses[index] = record;
          }
        })
      );
    },

    deleteAccess: async (access) => {
      await removeAccess(access);

      set(
        produce((state: AccessState) => {
          state.accesses = state.accesses.filter((a) => a.id !== access.id);
        })
      );
    },

    fetchAccesses: async () => {
      fetcher ??= listAccess();

      try {
        const accesses = await fetcher;
        set({ accesses: accesses ?? [], initialized: true });
      } finally {
        fetcher = null;
      }
    },
  };
});
