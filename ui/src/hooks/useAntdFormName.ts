import { useCreation } from "ahooks";
import { type FormInstance } from "antd";

export interface UseAntdFormNameOptions<T extends NonNullable<unknown> = any> {
  form: FormInstance<T>;
  name?: string;
}

const useAntdFormName = <T extends NonNullable<unknown> = any>(options: UseAntdFormNameOptions<T>) => {
  const formName = useCreation(() => `${options.name}_${Math.random().toString(36).substring(2, 10)}${new Date().getTime()}`, [options.name, options.form]);
  return formName;
};

export default useAntdFormName;
