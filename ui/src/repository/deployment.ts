import { Deployment, DeploymentListReq } from "@/domain/deployment";
import { getPb } from "./api";

export const list = async (req: DeploymentListReq) => {
  let page = 1;
  if (req.page) {
    page = req.page;
  }

  let perPage = 10;
  if (req.perPage) {
    perPage = req.perPage;
  }
  let filter = "domain!=null";
  if (req.domain) {
    filter = `domain="${req.domain}"`;
  }
  return await getPb()
    .collection("deployments")
    .getList<Deployment>(page, perPage, {
      filter: filter,
      sort: "-id",
      expand: "domain",
    });
};
