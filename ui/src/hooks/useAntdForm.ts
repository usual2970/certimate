import { useState } from "react";
import { Form, type FormInstance, type FormProps } from "antd";
import { useDeepCompareEffect } from "ahooks";

export interface UseAntdFormOptions<T extends NonNullable<unknown> = any> {
  form?: FormInstance<T>;
  initialValues?: Partial<T> | (() => Partial<T> | Promise<Partial<T>>);
  onSubmit?: (values: T) => any;
}

export interface UseAntdFormReturns<T extends NonNullable<unknown> = any> {
  form: FormInstance<T>;
  formProps: Omit<FormProps<T>, "children">;
  formPending: boolean;
  submit: (values?: T) => Promise<any>;
}

/**
 *
 * @param {UseAntdFormOptions} options
 * @returns {UseAntdFormReturns}
 */
const useAntdForm = <T extends NonNullable<unknown> = any>({ initialValues, form, onSubmit }: UseAntdFormOptions<T>): UseAntdFormReturns<T> => {
  const formInst = form ?? Form["useForm"]()[0];
  const [formInitialValues, setFormInitialValues] = useState<Partial<T>>();
  const [formPending, setFormPending] = useState(false);

  useDeepCompareEffect(() => {
    let unmounted = false;

    if (!initialValues) {
      return;
    }

    let temp: Promise<Partial<T>>;
    if (typeof initialValues === "function") {
      temp = Promise.resolve(initialValues());
    } else {
      temp = Promise.resolve(initialValues);
    }

    temp.then((temp) => {
      if (!unmounted) {
        type FieldName = Parameters<FormInstance<T>["getFieldValue"]>[0];
        type FieldsValue = Parameters<FormInstance<T>["setFieldsValue"]>[0];

        const obj = { ...temp };
        Object.keys(temp).forEach((key) => {
          obj[key as keyof T] = formInst!.isFieldTouched(key as FieldName) ? formInst!.getFieldValue(key as FieldName) : temp[key as keyof T];
        });

        setFormInitialValues(temp);
        formInst!.setFieldsValue(obj as FieldsValue);
      }
    });

    return () => {
      unmounted = true;
    };
  }, [formInst, initialValues]);

  const onFinish = (values: T) => {
    if (formPending) return Promise.reject(new Error("Form is pending"));

    setFormPending(true);

    return new Promise((resolve, reject) => {
      formInst
        .validateFields()
        .then(() => {
          resolve(
            Promise.resolve(onSubmit?.(values))
              .then((data) => {
                setFormPending(false);
                return data;
              })
              .catch((err) => {
                setFormPending(false);
                throw err;
              })
          );
        })
        .catch((err) => {
          setFormPending(false);
          reject(err);
        });
    });
  };

  const formProps: FormProps = {
    form: formInst,
    initialValues: formInitialValues,
    onFinish,
  };

  return {
    form: formInst,
    formProps: formProps,
    formPending: formPending,
    submit: () => {
      return onFinish(formInst.getFieldsValue(true));
    },
  };
};

export default useAntdForm;
