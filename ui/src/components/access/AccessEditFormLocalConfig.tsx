import { Form, type FormInstance } from "antd";

import { type LocalAccessConfig } from "@/domain/access";
import { useAntdForm } from "@/hooks";

type AccessEditFormLocalConfigFieldValues = Partial<LocalAccessConfig>;

export type AccessEditFormLocalConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessEditFormLocalConfigFieldValues;
  onValuesChange?: (values: AccessEditFormLocalConfigFieldValues) => void;
};

const initFormModel = (): AccessEditFormLocalConfigFieldValues => {
  return {};
};

const AccessEditFormLocalConfig = ({ form, formName, disabled, initialValues, onValuesChange }: AccessEditFormLocalConfigProps) => {
  const { form: formInst, formProps } = useAntdForm({
    form: form,
    initialValues: initialValues ?? initFormModel(),
  });

  const handleFormChange = (_: unknown, values: unknown) => {
    onValuesChange?.(values as AccessEditFormLocalConfigFieldValues);
  };

  return <Form {...formProps} form={formInst} disabled={disabled} layout="vertical" name={formName} onValuesChange={handleFormChange}></Form>;
};

export default AccessEditFormLocalConfig;
