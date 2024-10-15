import { Access } from "@/domain/access";
import { AccessGroup } from "@/domain/access_groups";
import { getPb } from "./api";

export const list = async () => {
  const resp = await getPb().collection("access_groups").getFullList<AccessGroup>({
    sort: "-created",
    expand: "access",
  });

  return resp;
};

export const remove = async (id: string) => {
  const pb = getPb();

  // 查询有没有关联的access
  const accessGroup = await pb.collection("access").getList<Access>(1, 1, {
    filter: `group='${id}' && deleted=null`,
  });

  if (accessGroup.items.length > 0) {
    throw new Error("该分组下有授权配置，无法删除");
  }

  await pb.collection("access_groups").delete(id);
};

export const update = async (accessGroup: AccessGroup) => {
  const pb = getPb();
  if (accessGroup.id) {
    return await pb.collection("access_groups").update(accessGroup.id, accessGroup);
  }
  return await pb.collection("access_groups").create(accessGroup);
};

type UpdateByIdReq = {
  id: string;
  [key: string]: string | string[];
};

export const updateById = async (req: UpdateByIdReq) => {
  const pb = getPb();
  return await pb.collection("access_groups").update(req.id, req);
};
