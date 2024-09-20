import { Access } from "./access";

export type AccessGroup = {
  id?: string;
  name?: string;
  access?: string[];
  expand?: {
    access: Access[];
  };
};
