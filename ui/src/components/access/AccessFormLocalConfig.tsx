import { Form, type FormInstance } from "antd";

import { type AccessConfigForLocal } from "@/domain/access";

type AccessFormLocalConfigFieldValues = Nullish<AccessConfigForLocal>;

export type AccessFormLocalConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessFormLocalConfigFieldValues;
  onValuesChange?: (values: AccessFormLocalConfigFieldValues) => void;
};

const initFormModel = (): AccessFormLocalConfigFieldValues => {
  return {};
};

const AccessFormLocalConfig = ({ form: formInst, formName, disabled, initialValues, onValuesChange }: AccessFormLocalConfigProps) => {
  const handleFormChange = (_: unknown, values: any) => {
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
    ></Form>
  );
};

export default AccessFormLocalConfig;
