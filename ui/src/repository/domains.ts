import { getTimeAfter } from "@/lib/time";
import { Domain } from "@/domain/domain";
import { getPocketBase } from "./pocketbase";

type DomainListReq = {
  domain?: string;
  page?: number;
  perPage?: number;
  state?: string;
};

export const list = async (req: DomainListReq) => {
  const pb = getPocketBase();

  let page = 1;
  if (req.page) {
    page = req.page;
  }

  let perPage = 2;
  if (req.perPage) {
    perPage = req.perPage;
  }

  let filter = "";
  if (req.state === "enabled") {
    filter = "enabled=true";
  } else if (req.state === "disabled") {
    filter = "enabled=false";
  } else if (req.state === "expired") {
    filter = pb.filter("expiredAt<{:expiredAt}", {
      expiredAt: getTimeAfter(15),
    });
  }

  const response = pb.collection("domains").getList<Domain>(page, perPage, {
    sort: "-created",
    expand: "lastDeployment",
    filter: filter,
  });

  return response;
};

export const get = async (id: string) => {
  const response = await getPocketBase().collection("domains").getOne<Domain>(id);
  return response;
};

export const save = async (data: Domain) => {
  if (data.id) {
    return await getPocketBase().collection("domains").update<Domain>(data.id, data);
  }
  return await getPocketBase().collection("domains").create<Domain>(data);
};

export const remove = async (id: string) => {
  return await getPocketBase().collection("domains").delete(id);
};

type Callback = (data: Domain) => void;
export const subscribeId = (id: string, callback: Callback) => {
  return getPocketBase()
    .collection("domains")
    .subscribe<Domain>(
      id,
      (e) => {
        if (e.action === "update") {
          callback(e.record);
        }
      },
      {
        expand: "lastDeployment",
      }
    );
};

export const unsubscribeId = (id: string) => {
  getPocketBase().collection("domains").unsubscribe(id);
};
