import PocketBase from "pocketbase";

let pb: PocketBase;
export const getPocketBase = () => {
  if (pb) return pb;
  pb = new PocketBase("/");
  return pb;
};
