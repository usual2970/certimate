import { type BaseModel as PbBaseModel } from "pocketbase";

declare global {
  declare interface BaseModel extends PbBaseModel {
    deleted?: string;
  }

  declare type MaybeModelRecord<T extends BaseModel = BaseModel> = T | Omit<T, "id" | "created" | "updated" | "deleted">;

  declare type MaybeModelRecordWithId<T extends BaseModel = BaseModel> = T | Pick<T, "id">;
}

export {};
