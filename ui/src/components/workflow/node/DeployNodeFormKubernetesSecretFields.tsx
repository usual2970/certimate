import { useTranslation } from "react-i18next";
import { Form, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

const DeployNodeFormKubernetesSecretFields = () => {
  const { t } = useTranslation();

  const formSchema = z.object({
    namespace: z
      .string({ message: t("workflow_node.deploy.form.k8s_namespace.placeholder") })
      .nonempty(t("workflow_node.deploy.form.k8s_namespace.placeholder"))
      .max(256, t("common.errmsg.string_max", { max: 256 }))
      .trim(),
    secretName: z
      .string({ message: t("workflow_node.deploy.form.k8s_secret_name.placeholder") })
      .nonempty(t("workflow_node.deploy.form.k8s_secret_name.placeholder"))
      .max(256, t("common.errmsg.string_max", { max: 256 }))
      .trim(),
    secretDataKeyForCrt: z
      .string({ message: t("workflow_node.deploy.form.k8s_secret_data_key_for_crt.placeholder") })
      .nonempty(t("workflow_node.deploy.form.k8s_secret_data_key_for_crt.placeholder"))
      .max(256, t("common.errmsg.string_max", { max: 256 }))
      .trim(),
    secretDataKeyForKey: z
      .string({ message: t("workflow_node.deploy.form.k8s_secret_data_key_for_key.placeholder") })
      .nonempty(t("workflow_node.deploy.form.k8s_secret_data_key_for_key.placeholder"))
      .max(256, t("common.errmsg.string_max", { max: 256 }))
      .trim(),
  });
  const formRule = createSchemaFieldRule(formSchema);

  return (
    <>
      <Form.Item
        name="namespace"
        label={t("workflow_node.deploy.form.k8s_namespace.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.k8s_namespace.tooltip") }}></span>}
        initialValue="default"
      >
        <Input placeholder={t("workflow_node.deploy.form.k8s_namespace.placeholder")} />
      </Form.Item>

      <Form.Item
        name="secretName"
        label={t("workflow_node.deploy.form.k8s_secret_name.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.k8s_secret_name.tooltip") }}></span>}
      >
        <Input placeholder={t("workflow_node.deploy.form.k8s_secret_name.placeholder")} />
      </Form.Item>

      <Form.Item
        name="secretDataKeyForCrt"
        label={t("workflow_node.deploy.form.k8s_secret_data_key_for_crt.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.k8s_secret_data_key_for_crt.tooltip") }}></span>}
        initialValue="tls.crt"
      >
        <Input placeholder={t("workflow_node.deploy.form.k8s_secret_data_key_for_crt.placeholder")} />
      </Form.Item>

      <Form.Item
        name="secretDataKeyForKey"
        label={t("workflow_node.deploy.form.k8s_secret_data_key_for_key.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("workflow_node.deploy.form.k8s_secret_data_key_for_key.tooltip") }}></span>}
        initialValue="tls.key"
      >
        <Input placeholder={t("workflow_node.deploy.form.k8s_secret_data_key_for_key.placeholder")} />
      </Form.Item>
    </>
  );
};

export default DeployNodeFormKubernetesSecretFields;
