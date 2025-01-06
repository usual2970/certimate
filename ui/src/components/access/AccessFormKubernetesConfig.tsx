import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { UploadOutlined as UploadOutlinedIcon } from "@ant-design/icons";
import { Button, Form, type FormInstance, Input, Upload, type UploadFile, type UploadProps } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type AccessConfigForKubernetes } from "@/domain/access";
import { readFileContent } from "@/utils/file";

type AccessFormKubernetesConfigFieldValues = Nullish<AccessConfigForKubernetes>;

export type AccessFormKubernetesConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessFormKubernetesConfigFieldValues;
  onValuesChange?: (values: AccessFormKubernetesConfigFieldValues) => void;
};

const initFormModel = (): AccessFormKubernetesConfigFieldValues => {
  return {};
};

const AccessFormKubernetesConfig = ({ form: formInst, formName, disabled, initialValues, onValuesChange }: AccessFormKubernetesConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    kubeConfig: z
      .string()
      .trim()
      .max(20480, t("common.errmsg.string_max", { max: 20480 }))
      .nullish(),
  });
  const formRule = createSchemaFieldRule(formSchema);

  const fieldKubeConfig = Form.useWatch("kubeConfig", formInst);
  const [fieldKubeFileList, setFieldKubeFileList] = useState<UploadFile[]>([]);
  useEffect(() => {
    setFieldKubeFileList(initialValues?.kubeConfig?.trim() ? [{ uid: "-1", name: "kubeconfig", status: "done" }] : []);
  }, [initialValues?.kubeConfig]);

  const handleKubeFileChange: UploadProps["onChange"] = async ({ file }) => {
    if (file && file.status !== "removed") {
      formInst.setFieldValue("kubeConfig", await readFileContent(file.originFileObj ?? (file as unknown as File)));
      setFieldKubeFileList([file]);
    } else {
      formInst.setFieldValue("kubeConfig", "");
      setFieldKubeFileList([]);
    }

    onValuesChange?.(formInst.getFieldsValue(true));
  };

  const handleFormChange = (_: unknown, values: z.infer<typeof formSchema>) => {
    onValuesChange?.(values);
  };

  return (
    <Form
      form={formInst}
      disabled={disabled}
      initialValues={initialValues ?? initFormModel()}
      layout="vertical"
      name={formName}
      onValuesChange={handleFormChange}
    >
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

export default AccessFormKubernetesConfig;
