import { useTranslation } from "react-i18next";
import { Form, Input, Select } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import Show from "@/components/Show";
import { validPortNumber } from "@/utils/validators";

const RESOURCE_TYPE_LOADBALANCER = "loadbalancer" as const;
const RESOURCE_TYPE_LISTENER = "listener" as const;

const DeployNodeFormAliyunCLBFields = () => {
  const { t } = useTranslation();

  const formSchema = z.object({
    resourceType: z.union([z.literal(RESOURCE_TYPE_LOADBALANCER), z.literal(RESOURCE_TYPE_LISTENER)], {
      message: t("workflow_node.deploy.form.aliyun_clb_resource_type.placeholder"),
    }),
    region: z
      .string({ message: t("workflow_node.deploy.form.aliyun_clb_region.placeholder") })
      .nonempty(t("workflow_node.deploy.form.aliyun_clb_region.placeholder"))
      .trim(),
    loadbalancerId: z
      .string({ message: t("workflow_node.deploy.form.aliyun_clb_loadbalancer_id.placeholder") })
      .min(1, t("workflow_node.deploy.form.aliyun_clb_loadbalancer_id.placeholder"))
      .max(64, t("common.errmsg.string_max", { max: 64 }))
      .trim(),
    listenerPort: z
      .union([
        z
          .number()
          .refine(
            (v) => fieldResourceType === RESOURCE_TYPE_LISTENER && validPortNumber(v),
            t("workflow_node.deploy.form.aliyun_clb_listener_port.placeholder")
          ),
        z
          .string()
          .refine(
            (v) => fieldResourceType === RESOURCE_TYPE_LISTENER && validPortNumber(v),
            t("workflow_node.deploy.form.aliyun_clb_listener_port.placeholder")
          ),
      ])
      .nullish(),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const formInst = Form.useFormInstance();

  const fieldResourceType = Form.useWatch("resourceType", formInst);

  return (
    <>
      <Form.Item name="resourceType" label={t("workflow_node.deploy.form.aliyun_clb_resource_type.label")} rules={[formRule]}>
        <Select placeholder={t("workflow_node.deploy.form.aliyun_clb_resource_type.placeholder")}>
          <Select.Option key={RESOURCE_TYPE_LOADBALANCER} value={RESOURCE_TYPE_LOADBALANCER}>
            {t("workflow_node.deploy.form.aliyun_clb_resource_type.option.loadbalancer.label")}
          </Select.Option>
          <Select.Option key={RESOURCE_TYPE_LISTENER} value={RESOURCE_TYPE_LISTENER}>
            {t("workflow_node.deploy.form.aliyun_clb_resource_type.option.listener.label")}
          </Select.Option>
        </Select>
      </Form.Item>

      <Form.Item
        name="region"
        label={t("workflow_node.deploy.form.aliyun_clb_region.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.aliyun_clb_region.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.aliyun_clb_region.placeholder")} />
      </Form.Item>

      <Form.Item
        name="loadbalancerId"
        label={t("workflow_node.deploy.form.aliyun_clb_loadbalancer_id.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.aliyun_clb_loadbalancer_id.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.aliyun_clb_loadbalancer_id.placeholder")} />
      </Form.Item>

      <Show when={fieldResourceType === RESOURCE_TYPE_LISTENER}>
        <Form.Item
          name="listenerPort"
          label={t("workflow_node.deploy.form.aliyun_clb_listener_port.label")}
          rules={[formRule]}
          tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.aliyun_clb_listener_port.tooltip") }}></span>}
          initialValue={443}
        >
          <Input type="number" min={1} max={65535} placeholder={t("workflow_node.deploy.form.aliyun_clb_listener_port.placeholder")} />
        </Form.Item>
      </Show>
    </>
  );
};

export default DeployNodeFormAliyunCLBFields;
