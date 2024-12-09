import { type RecordListOptions } from "pocketbase";
import moment from "moment";

import { type Certificate } from "@/domain/certificate";
import { getPocketBase } from "./pocketbase";

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
  };

  if (req.state === "expireSoon") {
    options.filter = pb.filter("expireAt<{:expiredAt}", {
      expiredAt: moment().add(15, "d").toDate(),
    });
  } else if (req.state === "expired") {
    options.filter = pb.filter("expireAt<={:expiredAt}", {
      expiredAt: new Date(),
    });
  }

  return pb.collection("certificate").getList<Certificate>(page, perPage, options);
};