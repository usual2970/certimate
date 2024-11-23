import { Certificate } from "@/domain/certificate";
import { getPb } from "./api";
import { RecordListOptions } from "pocketbase";
import { getTimeAfter } from "@/lib/time";

type CertificateListReq = {
  page?: number;
  perPage?: number;
  state?: string;
};

export const list = async (req: CertificateListReq) => {
  const pb = getPb();

  let page = 1;
  if (req.page) {
    page = req.page;
  }

  let perPage = 2;
  if (req.perPage) {
    perPage = req.perPage;
  }

  const options: RecordListOptions = {
    sort: "-created",
    expand: "workflow",
  };

  if (req.state === "expireSoon") {
    options.filter = pb.filter("expireAt<{:expiredAt}", {
      expiredAt: getTimeAfter(15),
    });
  } else if (req.state === "expired") {
    options.filter = pb.filter("expireAt<={:expiredAt}", {
      expiredAt: new Date(),
    });
  }

  const response = pb.collection("certificate").getList<Certificate>(page, perPage, options);

  return response;
};
