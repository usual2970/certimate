import { Form, type FormInstance } from "antd";

import { useAntdForm } from "@/hooks";
import { type LocalAccessConfig } from "@/domain/access";

type AccessEditFormLocalConfigModelValues = Partial<LocalAccessConfig>;

export type AccessEditFormLocalConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  model?: AccessEditFormLocalConfigModelValues;
  onModelChange?: (model: AccessEditFormLocalConfigModelValues) => void;
};

const initFormModel = (): AccessEditFormLocalConfigModelValues => {
  return {};
};

const AccessEditFormLocalConfig = ({ form, formName, disabled, model, onModelChange }: AccessEditFormLocalConfigProps) => {
  const { form: formInst, formProps } = useAntdForm({
    form: form,
    initialValues: model ?? initFormModel(),
  });

  const handleFormChange = (_: unknown, values: unknown) => {
    onModelChange?.(values as AccessEditFormLocalConfigModelValues);
  };

  return <Form {...formProps} form={formInst} disabled={disabled} layout="vertical" name={formName} onValuesChange={handleFormChange}></Form>;
};

export default AccessEditFormLocalConfig;
