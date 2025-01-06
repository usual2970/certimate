import { useTranslation } from "react-i18next";
import { Form, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { validDomainName } from "@/utils/validators";

const DeployNodeFormTencentCloudCOSFields = () => {
  const { t } = useTranslation();

  const formSchema = z.object({
    region: z
      .string({ message: t("workflow_node.deploy.form.tencentcloud_cos_region.placeholder") })
      .nonempty(t("workflow_node.deploy.form.tencentcloud_cos_region.placeholder"))
      .trim(),
    bucket: z
      .string({ message: t("workflow_node.deploy.form.tencentcloud_cos_bucket.placeholder") })
      .nonempty(t("workflow_node.deploy.form.tencentcloud_cos_bucket.placeholder"))
      .trim(),
    domain: z
      .string({ message: t("workflow_node.deploy.form.tencentcloud_cos_domain.placeholder") })
      .refine((v) => validDomainName(v), t("common.errmsg.domain_invalid")),
  });
  const formRule = createSchemaFieldRule(formSchema);

  return (
    <>
      <Form.Item
        name="region"
        label={t("workflow_node.deploy.form.tencentcloud_cos_region.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.tencentcloud_cos_region.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.tencentcloud_cos_region.placeholder")} />
      </Form.Item>

      <Form.Item
        name="bucket"
        label={t("workflow_node.deploy.form.tencentcloud_cos_bucket.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.tencentcloud_cos_bucket.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.tencentcloud_cos_bucket.placeholder")} />
      </Form.Item>

      <Form.Item
        name="domain"
        label={t("workflow_node.deploy.form.tencentcloud_cos_domain.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.tencentcloud_cos_domain.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.tencentcloud_cos_domain.placeholder")} />
      </Form.Item>
    </>
  );
};

export default DeployNodeFormTencentCloudCOSFields;
