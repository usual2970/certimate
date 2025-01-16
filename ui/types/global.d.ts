import { type BaseModel as PbBaseModel } from "pocketbase";

declare global {
  declare type ISO8601String = string;

  declare interface BaseModel extends PbBaseModel {
    created: ISO8601String;
    updated: ISO8601String;
    deleted?: ISO8601String;
  }

  declare type MaybeModelRecord<T extends BaseModel = BaseModel> = T | Omit<T, "id" | "created" | "updated" | "deleted">;

  declare type MaybeModelRecordWithId<T extends BaseModel = BaseModel> = T | Pick<T, "id">;

  declare interface BaseResponse<T = any> {
    code: number;
    msg: string;
    data: T;
  }
}

export {};
