import dayjs from "dayjs";
import { type RecordListOptions } from "pocketbase";

import { type CertificateModel } from "@/domain/certificate";
import { getPocketBase } from "./pocketbase";

const COLLECTION_NAME = "certificate";

export type ListCertificateRequest = {
  page?: number;
  perPage?: number;
  state?: "expireSoon" | "expired";
};

export const list = async (request: ListCertificateRequest) => {
  const pb = getPocketBase();

  const page = request.page || 1;
  const perPage = request.perPage || 10;

  const options: RecordListOptions = {
    sort: "-created",
    expand: "workflow",
    requestKey: null,
  };

  if (request.state === "expireSoon") {
    options.filter = pb.filter("expireAt<{:expiredAt}", {
      expiredAt: dayjs().add(15, "d").toDate(),
    });
  } else if (request.state === "expired") {
    options.filter = pb.filter("expireAt<={:expiredAt}", {
      expiredAt: new Date(),
    });
  }

  return pb.collection(COLLECTION_NAME).getList<CertificateModel>(page, perPage, options);
};

export const remove = async (record: MaybeModelRecordWithId<CertificateModel>) => {
  record = { ...record, deleted: dayjs.utc().format("YYYY-MM-DD HH:mm:ss") };

  await getPocketBase().collection(COLLECTION_NAME).update<CertificateModel>(record.id!, record);
};
