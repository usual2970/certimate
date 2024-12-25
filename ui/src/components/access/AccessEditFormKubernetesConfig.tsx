import { useState } from "react";
import { flushSync } from "react-dom";
import { useTranslation } from "react-i18next";
import { useDeepCompareEffect } from "ahooks";
import { Button, Form, Input, Upload, type FormInstance, type UploadFile, type UploadProps } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";
import { Upload as UploadIcon } from "lucide-react";

import { useAntdForm } from "@/hooks";
import { type KubernetesAccessConfig } from "@/domain/access";
import { readFileContent } from "@/utils/file";

type AccessEditFormKubernetesConfigModelValues = Partial<KubernetesAccessConfig>;

export type AccessEditFormKubernetesConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  model?: AccessEditFormKubernetesConfigModelValues;
  onModelChange?: (model: AccessEditFormKubernetesConfigModelValues) => void;
};

const initFormModel = (): AccessEditFormKubernetesConfigModelValues => {
  return {};
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
  const { form: formInst, formProps } = useAntdForm<z.infer<typeof formSchema>>({
    form: form,
    initialValues: model ?? initFormModel(),
  });

  const [kubeFileList, setKubeFileList] = useState<UploadFile[]>([]);
  useDeepCompareEffect(() => {
    setKubeFileList(model?.kubeConfig?.trim() ? [{ uid: "-1", name: "kubeconfig", status: "done" }] : []);
  }, [model]);

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onModelChange?.(values as AccessEditFormKubernetesConfigModelValues);
  };

  const handleKubeFileChange: UploadProps["onChange"] = async ({ file }) => {
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
    <Form {...formProps} form={formInst} disabled={disabled} layout="vertical" name={formName} onValuesChange={handleFormChange}>
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
        <Upload beforeUpload={() => false} fileList={kubeFileList} maxCount={1} onChange={handleKubeFileChange}>
          <Button icon={<UploadIcon size={16} />}>{t("access.form.k8s_kubeconfig.upload")}</Button>
        </Upload>
      </Form.Item>
    </Form>
  );
};

export default AccessEditFormKubernetesConfig;
