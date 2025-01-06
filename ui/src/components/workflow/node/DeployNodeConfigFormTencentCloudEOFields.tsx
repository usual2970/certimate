import { useTranslation } from "react-i18next";
import { Form, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { validDomainName } from "@/utils/validators";

const DeployNodeFormTencentCloudEOFields = () => {
  const { t } = useTranslation();

  const formSchema = z.object({
    zoneId: z
      .string({ message: t("workflow_node.deploy.form.tencentcloud_eo_zone_id.placeholder") })
      .nonempty(t("workflow_node.deploy.form.tencentcloud_eo_zone_id.placeholder"))
      .trim(),
    domain: z
      .string({ message: t("workflow_node.deploy.form.tencentcloud_eo_domain.placeholder") })
      .refine((v) => validDomainName(v), t("common.errmsg.domain_invalid")),
  });
  const formRule = createSchemaFieldRule(formSchema);

  return (
    <>
      <Form.Item
        name="zoneId"
        label={t("workflow_node.deploy.form.tencentcloud_eo_zone_id.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.tencentcloud_eo_zone_id.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.tencentcloud_eo_zone_id.placeholder")} />
      </Form.Item>

      <Form.Item
        name="domain"
        label={t("workflow_node.deploy.form.tencentcloud_eo_domain.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.tencentcloud_eo_domain.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.tencentcloud_eo_domain.placeholder")} />
      </Form.Item>
    </>
  );
};

export default DeployNodeFormTencentCloudEOFields;
