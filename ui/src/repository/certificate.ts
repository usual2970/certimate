import { Certificate } from "@/domain/certificate";
import { getPb } from "./api";

type CertificateListReq = {
  page?: number;
  perPage?: number;
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

  const response = pb.collection("certificate").getList<Certificate>(page, perPage, {
    sort: "-created",
    expand: "workflow",
  });

  return response;
};
