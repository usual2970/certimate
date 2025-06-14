import dayjs from "dayjs";

import { type CertificateModel } from "@/domain/certificate";
import { COLLECTION_NAME_CERTIFICATE, getPocketBase } from "./_pocketbase";

export type ListRequest = {
  keyword?: string;
  state?: "expireSoon" | "expired";
  page?: number;
  perPage?: number;
};

export const list = async (request: ListRequest) => {
  const pb = getPocketBase();

  const filters: string[] = ["deleted=null"];
  if (request.keyword) {
    filters.push(pb.filter("(subjectAltNames~{:keyword} || serialNumber={:keyword})", { keyword: request.keyword }));
  }
  if (request.state === "expireSoon") {
    filters.push(pb.filter("expireAt<{:expiredAt} && expireAt>@now", { expiredAt: dayjs().add(20, "d").toDate() }));
  } else if (request.state === "expired") {
    filters.push(pb.filter("expireAt<={:expiredAt}", { expiredAt: new Date() }));
  }

  const page = request.page || 1;
  const perPage = request.perPage || 10;
  return pb.collection(COLLECTION_NAME_CERTIFICATE).getList<CertificateModel>(page, perPage, {
    expand: "workflowId",
    filter: filters.join(" && "),
    sort: "-created",
    requestKey: null,
  });
};

export const listByWorkflowRunId = async (workflowRunId: string) => {
  const pb = getPocketBase();

  const list = await pb.collection(COLLECTION_NAME_CERTIFICATE).getFullList<CertificateModel>({
    batch: 65535,
    filter: pb.filter("workflowRunId={:workflowRunId}", { workflowRunId: workflowRunId }),
    sort: "created",
    requestKey: null,
  });

  return {
    totalItems: list.length,
    items: list,
  };
};

export const remove = async (record: MaybeModelRecordWithId<CertificateModel>) => {
  await getPocketBase()
    .collection(COLLECTION_NAME_CERTIFICATE)
    .update<CertificateModel>(record.id!, { deleted: dayjs.utc().format("YYYY-MM-DD HH:mm:ss") });
  return true;
};
