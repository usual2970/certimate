import { useTranslation } from "react-i18next";
import { Form, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { validDomainName } from "@/utils/validators";

const DeployNodeFormHuaweiCloudCDNFields = () => {
  const { t } = useTranslation();

  const formSchema = z.object({
    region: z
      .string({ message: t("workflow_node.deploy.form.huaweicloud_cdn_region.placeholder") })
      .nonempty(t("workflow_node.deploy.form.huaweicloud_cdn_region.placeholder"))
      .trim(),
    domain: z
      .string({ message: t("workflow_node.deploy.form.huaweicloud_cdn_domain.placeholder") })
      .refine((v) => validDomainName(v), t("common.errmsg.domain_invalid")),
  });
  const formRule = createSchemaFieldRule(formSchema);

  return (
    <>
      <Form.Item
        name="region"
        label={t("workflow_node.deploy.form.huaweicloud_cdn_region.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.huaweicloud_cdn_region.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.huaweicloud_cdn_region.placeholder")} />
      </Form.Item>

      <Form.Item
        name="domain"
        label={t("workflow_node.deploy.form.huaweicloud_cdn_domain.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.huaweicloud_cdn_domain.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.huaweicloud_cdn_domain.placeholder")} />
      </Form.Item>
    </>
  );
};

export default DeployNodeFormHuaweiCloudCDNFields;
