import { getPocketBase } from "./_pocketbase";

const COLLECTION_NAME = "_superusers";

export const authWithPassword = (username: string, password: string) => {
  return getPocketBase().collection(COLLECTION_NAME).authWithPassword(username, password);
};

export const getAuthStore = () => {
  return getPocketBase().authStore;
};

export const save = (data: { email: string } | { password: string }) => {
  return getPocketBase()
    .collection(COLLECTION_NAME)
    .update(getAuthStore().record?.id || "", data);
};
