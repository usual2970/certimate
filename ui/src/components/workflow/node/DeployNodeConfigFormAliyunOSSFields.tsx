import { useTranslation } from "react-i18next";
import { Form, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { validDomainName } from "@/utils/validators";

const DeployNodeFormAliyunOSSFields = () => {
  const { t } = useTranslation();

  const formSchema = z.object({
    endpoint: z
      .string({ message: t("workflow_node.deploy.form.aliyun_oss_endpoint.placeholder") })
      .url(t("common.errmsg.url_invalid"))
      .trim(),
    bucket: z
      .string({ message: t("workflow_node.deploy.form.aliyun_oss_bucket.placeholder") })
      .nonempty(t("workflow_node.deploy.form.aliyun_oss_bucket.placeholder"))
      .trim(),
    domain: z
      .string({ message: t("workflow_node.deploy.form.aliyun_oss_domain.placeholder") })
      .refine((v) => validDomainName(v, true), t("common.errmsg.domain_invalid")),
  });
  const formRule = createSchemaFieldRule(formSchema);

  return (
    <>
      <Form.Item
        name="endpoint"
        label={t("workflow_node.deploy.form.aliyun_oss_endpoint.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.aliyun_oss_endpoint.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.aliyun_oss_endpoint.placeholder")} />
      </Form.Item>

      <Form.Item
        name="bucket"
        label={t("workflow_node.deploy.form.aliyun_oss_bucket.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.aliyun_oss_bucket.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.aliyun_oss_bucket.placeholder")} />
      </Form.Item>

      <Form.Item
        name="domain"
        label={t("workflow_node.deploy.form.aliyun_oss_domain.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.aliyun_oss_domain.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.aliyun_oss_domain.placeholder")} />
      </Form.Item>
    </>
  );
};

export default DeployNodeFormAliyunOSSFields;
