import { COLLECTION_NAME_ADMIN, getPocketBase } from "./_pocketbase";

export const authWithPassword = (username: string, password: string) => {
  return getPocketBase().collection(COLLECTION_NAME_ADMIN).authWithPassword(username, password);
};

export const getAuthStore = () => {
  return getPocketBase().authStore;
};

export const save = (data: { email: string } | { password: string; passwordConfirm: string }) => {
  return getPocketBase()
    .collection(COLLECTION_NAME_ADMIN)
    .update(getAuthStore().record?.id || "", data);
};
