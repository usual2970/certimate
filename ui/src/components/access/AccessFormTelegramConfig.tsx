import { useTranslation } from "react-i18next";
import { Form, type FormInstance, Input } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { z } from "zod";

import { type AccessConfigForTelegram } from "@/domain/access";

type AccessFormTelegramConfigFieldValues = Nullish<AccessConfigForTelegram>;

export type AccessFormTelegramConfigProps = {
  form: FormInstance;
  formName: string;
  disabled?: boolean;
  initialValues?: AccessFormTelegramConfigFieldValues;
  onValuesChange?: (values: AccessFormTelegramConfigFieldValues) => void;
};

const initFormModel = (): AccessFormTelegramConfigFieldValues => {
  return {
    botToken: "",
  };
};

const AccessFormTelegramConfig = ({ form: formInst, formName, disabled, initialValues, onValuesChange }: AccessFormTelegramConfigProps) => {
  const { t } = useTranslation();

  const formSchema = z.object({
    botToken: z
      .string({ message: t("access.form.telegram_bot_token.placeholder") })
      .min(1, t("access.form.telegram_bot_token.placeholder"))
      .max(256, t("common.errmsg.string_max", { max: 256 })),
    defaultChatId: z
      .preprocess(
        (v) => (v == null || v === "" ? undefined : Number(v)),
        z
          .number()
          .nullish()
          .refine((v) => {
            if (v == null || v + "" === "") return true;
            return /^\d+$/.test(v + "") && +v! > 0;
          }, t("access.form.telegram_default_chat_id.placeholder"))
      )
      .nullish(),
  });
  const formRule = createSchemaFieldRule(formSchema);

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
      <Form.Item
        name="botToken"
        label={t("access.form.telegram_bot_token.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.telegram_bot_token.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("access.form.telegram_bot_token.placeholder")} />
      </Form.Item>

      <Form.Item
        name="defaultChatId"
        label={t("access.form.telegram_default_chat_id.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("access.form.telegram_default_chat_id.tooltip") }}></span>}
      >
        <Input type="number" allowClear placeholder={t("access.form.telegram_default_chat_id.placeholder")} />
      </Form.Item>
    </Form>
  );
};

export default AccessFormTelegramConfig;
