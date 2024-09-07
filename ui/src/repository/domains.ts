import { Domain, Statistic } from "@/domain/domain";
import { getPb } from "./api";
import { getTimeAfter } from "@/lib/time";

type DomainListReq = {
  domain?: string;
  page?: number;
  perPage?: number;
  state?: string;
};

export const list = async (req: DomainListReq) => {
  let page = 1;
  if (req.page) {
    page = req.page;
  }

  let perPage = 2;
  if (req.perPage) {
    perPage = req.perPage;
  }
  const pb = getPb();
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

export const statistics = async (): Promise<Statistic> => {
  const pb = getPb();
  const total = await pb.collection("domains").getList(1, 1, {});
  const expired = await pb.collection("domains").getList(1, 1, {
    filter: pb.filter("expiredAt<{:expiredAt}", {
      expiredAt: getTimeAfter(15),
    }),
  });

  const enabled = await pb.collection("domains").getList(1, 1, {
    filter: "enabled=true",
  });
  const disabled = await pb.collection("domains").getList(1, 1, {
    filter: "enabled=false",
  });

  return {
    total: total.totalItems,
    expired: expired.totalItems,
    enabled: enabled.totalItems,
    disabled: disabled.totalItems,
  };
};

export const get = async (id: string) => {
  const response = await getPb().collection("domains").getOne<Domain>(id);
  return response;
};

export const save = async (data: Domain) => {
  if (data.id) {
    return await getPb().collection("domains").update<Domain>(data.id, data);
  }
  return await getPb().collection("domains").create<Domain>(data);
};

export const remove = async (id: string) => {
  return await getPb().collection("domains").delete(id);
};

type Callback = (data: Domain) => void;
export const subscribeId = (id: string, callback: Callback) => {
  return getPb()
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
  getPb().collection("domains").unsubscribe(id);
};
