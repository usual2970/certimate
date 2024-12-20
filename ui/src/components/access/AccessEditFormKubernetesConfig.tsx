import { useState } from "react";
import { flushSync } from "react-dom";
import { useTranslation } from "react-i18next";
import { useDeepCompareEffect } from "ahooks";
import { Button, Form, Input, Upload, type FormInstance, type UploadFile, type UploadProps } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";
import { Upload as UploadIcon } from "lucide-react";

import { type KubernetesAccessConfig } from "@/domain/access";
import { readFileContent } from "@/utils/file";

type AccessEditFormKubernetesConfigModelType = Partial<KubernetesAccessConfig>;

export type AccessEditFormKubernetesConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  loading?: boolean;
  model?: AccessEditFormKubernetesConfigModelType;
  onModelChange?: (model: AccessEditFormKubernetesConfigModelType) => void;
};

const initModel = () => {
  return {} as AccessEditFormKubernetesConfigModelType;
};

const AccessEditFormKubernetesConfig = ({ form, formName, disabled, loading, model, onModelChange }: AccessEditFormKubernetesConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    kubeConfig: z
      .string()
      .trim()
      .min(0, t("access.form.k8s_kubeconfig.placeholder"))
      .max(20480, t("common.errmsg.string_max", { max: 20480 }))
      .nullish(),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const [initialValues, setInitialValues] = useState<Partial<z.infer<typeof formSchema>>>(model ?? initModel());
  useDeepCompareEffect(() => {
    setInitialValues(model ?? initModel());
    setKubeFileList(model?.kubeConfig?.trim() ? [{ uid: "-1", name: "kubeconfig", status: "done" }] : []);
  }, [model]);

  const [kubeFileList, setKubeFileList] = useState<UploadFile[]>([]);

  const handleFormChange = (_: unknown, fields: AccessEditFormKubernetesConfigModelType) => {
    onModelChange?.(fields);
  };

  const handleUploadChange: UploadProps["onChange"] = async ({ file }) => {
    if (file && file.status !== "removed") {
      form.setFieldValue("kubeConfig", (await readFileContent(file.originFileObj ?? (file as unknown as File))).trim());
      setKubeFileList([file]);
    } else {
      form.setFieldValue("kubeConfig", "");
      setKubeFileList([]);
    }

    flushSync(() => onModelChange?.(form.getFieldsValue()));
  };

  return (
    <Form form={form} disabled={loading || disabled} initialValues={initialValues} layout="vertical" name={formName} onValuesChange={handleFormChange}>
      <Form.Item
        name="kubeConfig"
        label={t("access.form.k8s_kubeconfig.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.k8s_kubeconfig.tooltip") }}></span>}
      >
        <Input.TextArea hidden placeholder={t("access.form.k8s_kubeconfig.placeholder")} value={form.getFieldValue("kubeConfig")} />
        <Upload beforeUpload={() => false} fileList={kubeFileList} maxCount={1} onChange={handleUploadChange}>
          <Button icon={<UploadIcon size={16} />}>{t("access.form.k8s_kubeconfig.upload")}</Button>
        </Upload>
      </Form.Item>
    </Form>
  );
};

export default AccessEditFormKubernetesConfig;
