import { useState } from "react";
import { flushSync } from "react-dom";
import { useTranslation } from "react-i18next";
import { UploadOutlined as UploadOutlinedIcon } from "@ant-design/icons";
import { useDeepCompareEffect } from "ahooks";
import { Button, Form, type FormInstance, Input, Upload, type UploadFile, type UploadProps } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type AccessConfigForKubernetes } from "@/domain/access";
import { readFileContent } from "@/utils/file";

type AccessEditFormKubernetesConfigFieldValues = Partial<AccessConfigForKubernetes>;

export type AccessEditFormKubernetesConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessEditFormKubernetesConfigFieldValues;
  onValuesChange?: (values: AccessEditFormKubernetesConfigFieldValues) => void;
};

const initFormModel = (): AccessEditFormKubernetesConfigFieldValues => {
  return {};
};

const AccessEditFormKubernetesConfig = ({ form, formName, disabled, initialValues, onValuesChange }: AccessEditFormKubernetesConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    kubeConfig: z
      .string()
      .trim()
      .max(20480, t("common.errmsg.string_max", { max: 20480 }))
      .nullish(),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const fieldKubeConfig = Form.useWatch("kubeConfig", form);
  const [fieldKubeFileList, setFieldKubeFileList] = useState<UploadFile[]>([]);
  useDeepCompareEffect(() => {
    setFieldKubeFileList(initialValues?.kubeConfig?.trim() ? [{ uid: "-1", name: "kubeconfig", status: "done" }] : []);
  }, [initialValues]);

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onValuesChange?.(values as AccessEditFormKubernetesConfigFieldValues);
  };

  const handleKubeFileChange: UploadProps["onChange"] = async ({ file }) => {
    if (file && file.status !== "removed") {
      form.setFieldValue("kubeConfig", await readFileContent(file.originFileObj ?? (file as unknown as File)));
      setFieldKubeFileList([file]);
    } else {
      form.setFieldValue("kubeConfig", "");
      setFieldKubeFileList([]);
    }

    flushSync(() => onValuesChange?.(form.getFieldsValue(true)));
  };

  return (
    <Form form={form} disabled={disabled} initialValues={initialValues ?? initFormModel()} layout="vertical" name={formName} onValuesChange={handleFormChange}>
      <Form.Item name="kubeConfig" noStyle rules={[formRule]}>
        <Input.TextArea autoComplete="new-password" hidden placeholder={t("access.form.k8s_kubeconfig.placeholder")} value={fieldKubeConfig} />
      </Form.Item>
      <Form.Item
        label={t("access.form.k8s_kubeconfig.label")}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.k8s_kubeconfig.tooltip") }}></span>}
      >
        <Upload beforeUpload={() => false} fileList={fieldKubeFileList} maxCount={1} onChange={handleKubeFileChange}>
          <Button icon={<UploadOutlinedIcon />}>{t("access.form.k8s_kubeconfig.upload")}</Button>
        </Upload>
      </Form.Item>
    </Form>
  );
};

export default AccessEditFormKubernetesConfig;
