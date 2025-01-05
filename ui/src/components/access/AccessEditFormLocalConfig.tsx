import { Form, type FormInstance } from "antd";

import { type AccessConfigForLocal } from "@/domain/access";

type AccessEditFormLocalConfigFieldValues = Partial<AccessConfigForLocal>;

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
  const handleFormChange = (_: unknown, values: unknown) => {
    onValuesChange?.(values as AccessEditFormLocalConfigFieldValues);
  };

  return (
    <Form
      form={form}
      disabled={disabled}
      initialValues={initialValues ?? initFormModel()}
      layout="vertical"
      name={formName}
      onValuesChange={handleFormChange}
    ></Form>
  );
};

export default AccessEditFormLocalConfig;
