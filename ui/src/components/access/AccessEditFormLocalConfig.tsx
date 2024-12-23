import { useState } from "react";
import { useDeepCompareEffect } from "ahooks";
import { Form, type FormInstance } from "antd";

import { type LocalAccessConfig } from "@/domain/access";

type AccessEditFormLocalConfigModelType = Partial<LocalAccessConfig>;

export type AccessEditFormLocalConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  model?: AccessEditFormLocalConfigModelType;
  onModelChange?: (model: AccessEditFormLocalConfigModelType) => void;
};

const initModel = () => {
  return {} as AccessEditFormLocalConfigModelType;
};

const AccessEditFormLocalConfig = ({ form, formName, disabled, model }: AccessEditFormLocalConfigProps) => {
  const [initialValues, setInitialValues] = useState(model ?? initModel());
  useDeepCompareEffect(() => {
    setInitialValues(model ?? initModel());
  }, [model]);

  return <Form form={form} disabled={disabled} initialValues={initialValues} layout="vertical" name={formName}></Form>;
};

export default AccessEditFormLocalConfig;
