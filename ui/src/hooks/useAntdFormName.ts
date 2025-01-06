import { useCreation } from "ahooks";
import { type FormInstance } from "antd";
import { nanoid } from "nanoid/non-secure";

export interface UseAntdFormNameOptions<T extends NonNullable<unknown> = any> {
  form: FormInstance<T>;
  name?: string;
}

/**
 * 生成并获取一个 antd 表单的唯一名称。
 * 通常为配合 Form 组件使用，避免页面上同时存在多个表单时若有同名的 FormItem 会产生冲突。
 * @param {UseAntdFormNameOptions} options
 * @returns {string}
 */
const useAntdFormName = <T extends NonNullable<unknown> = any>(options: UseAntdFormNameOptions<T>) => {
  const formName = useCreation(() => `${options.name}_${nanoid()}`, [options.name, options.form]);
  return formName;
};

export default useAntdFormName;
