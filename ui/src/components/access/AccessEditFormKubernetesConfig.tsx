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
  model?: AccessEditFormKubernetesConfigModelType;
  onModelChange?: (model: AccessEditFormKubernetesConfigModelType) => void;
};

const initModel = () => {
  return {} as AccessEditFormKubernetesConfigModelType;
};

const AccessEditFormKubernetesConfig = ({ form, formName, disabled, model, onModelChange }: AccessEditFormKubernetesConfigProps) => {
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
  const formInst = form as FormInstance<z.infer<typeof formSchema>>;

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
      formInst.setFieldValue("kubeConfig", (await readFileContent(file.originFileObj ?? (file as unknown as File))).trim());
      setKubeFileList([file]);
    } else {
      formInst.setFieldValue("kubeConfig", "");
      setKubeFileList([]);
    }

    flushSync(() => onModelChange?.(formInst.getFieldsValue(true)));
  };

  return (
    <Form form={form} disabled={disabled} initialValues={initialValues} layout="vertical" name={formName} onValuesChange={handleFormChange}>
      <Form.Item name="kubeConfig" noStyle rules={[formRule]}>
        <Input.TextArea
          autoComplete="new-password"
          hidden
          placeholder={t("access.form.k8s_kubeconfig.placeholder")}
          value={formInst.getFieldValue("kubeConfig")}
        />
      </Form.Item>
      <Form.Item
        label={t("access.form.k8s_kubeconfig.label")}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.k8s_kubeconfig.tooltip") }}></span>}
      >
        <Upload beforeUpload={() => false} fileList={kubeFileList} maxCount={1} onChange={handleUploadChange}>
          <Button icon={<UploadIcon size={16} />}>{t("access.form.k8s_kubeconfig.upload")}</Button>
        </Upload>
      </Form.Item>
    </Form>
  );
};

export default AccessEditFormKubernetesConfig;
