import dayjs from "dayjs";
import { type RecordListOptions } from "pocketbase";

import { type CertificateModel } from "@/domain/certificate";
import { getPocketBase } from "./pocketbase";

const COLLECTION_NAME = "certificate";

export type CertificateListReq = {
  page?: number;
  perPage?: number;
  state?: "expireSoon" | "expired";
};

export const list = async (req: CertificateListReq) => {
  const pb = getPocketBase();

  const page = req.page || 1;
  const perPage = req.perPage || 10;

  const options: RecordListOptions = {
    sort: "-created",
    expand: "workflow",
    requestKey: null,
  };

  if (req.state === "expireSoon") {
    options.filter = pb.filter("expireAt<{:expiredAt}", {
      expiredAt: dayjs().add(15, "d").toDate(),
    });
  } else if (req.state === "expired") {
    options.filter = pb.filter("expireAt<={:expiredAt}", {
      expiredAt: new Date(),
    });
  }

  return pb.collection(COLLECTION_NAME).getList<CertificateModel>(page, perPage, options);
};
