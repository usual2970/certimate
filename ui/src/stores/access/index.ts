import { create } from "zustand";
import { produce } from "immer";

import { type AccessModel } from "@/domain/access";
import { list as listAccess, save as saveAccess, remove as removeAccess } from "@/repository/access";

export interface AccessState {
  accesses: AccessModel[];
  createAccess: (access: AccessModel) => void;
  updateAccess: (access: AccessModel) => void;
  deleteAccess: (access: AccessModel) => void;
  fetchAccesses: () => Promise<void>;
}

export const useAccessStore = create<AccessState>((set) => {
  return {
    accesses: [],

    createAccess: async (access) => {
      access = await saveAccess(access);

      set(
        produce((state: AccessState) => {
          state.accesses.unshift(access);
        })
      );
    },

    updateAccess: async (access) => {
      access = await saveAccess(access);

      set(
        produce((state: AccessState) => {
          const index = state.accesses.findIndex((e) => e.id === access.id);
          state.accesses[index] = access;
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
      const accesses = await listAccess();

      set({
        accesses: accesses ?? [],
      });
    },
  };
});
