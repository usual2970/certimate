import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { CheckCard } from "@ant-design/pro-components";
import { Alert, Button, Form, Input, Skeleton, message, notification } from "antd";
import { createSchemaFieldRule } from "antd-zod";
import { produce } from "immer";
import { z } from "zod";

import Show from "@/components/Show";
import { SETTINGS_NAMES, SSLPROVIDERS, type SSLProviderSettingsContent, type SSLProviders, type SettingsModel } from "@/domain/settings";
import { useAntdForm } from "@/hooks";
import { get as getSettings, save as saveSettings } from "@/repository/settings";
import { getErrMsg } from "@/utils/error";

const SSLProviderContext = createContext(
  {} as {
    pending: boolean;
    settings: SettingsModel<SSLProviderSettingsContent>;
    updateSettings: (settings: MaybeModelRecordWithId<SettingsModel<SSLProviderSettingsContent>>) => Promise<void>;
  }
);

const SSLProviderEditFormLetsEncryptConfig = () => {
  const { t } = useTranslation();

  const { pending, settings, updateSettings } = useContext(SSLProviderContext);

  const { form: formInst, formProps } = useAntdForm<NonNullable<unknown>>({
    initialValues: settings?.content?.config?.[SSLPROVIDERS.LETS_ENCRYPT],
    onSubmit: async (values) => {
      const newSettings = produce(settings, (draft) => {
        draft.content ??= {} as SSLProviderSettingsContent;
        draft.content.provider = SSLPROVIDERS.LETS_ENCRYPT;

        draft.content.config ??= {} as SSLProviderSettingsContent["config"];
        draft.content.config[SSLPROVIDERS.LETS_ENCRYPT] = values;
      });
      await updateSettings(newSettings);

      setFormChanged(false);
    },
  });

  const [formChanged, setFormChanged] = useState(false);
  useEffect(() => {
    setFormChanged(settings?.content?.provider !== SSLPROVIDERS.LETS_ENCRYPT);
  }, [settings?.content?.provider]);

  const handleFormChange = () => {
    setFormChanged(true);
  };

  return (
    <Form {...formProps} form={formInst} disabled={pending} layout="vertical" onValuesChange={handleFormChange}>
      <Form.Item>
        <Button type="primary" htmlType="submit" disabled={!formChanged} loading={pending}>
          {t("common.button.save")}
        </Button>
      </Form.Item>
    </Form>
  );
};

const SSLProviderEditFormLetsEncryptStagingConfig = () => {
  const { t } = useTranslation();

  const { pending, settings, updateSettings } = useContext(SSLProviderContext);

  const { form: formInst, formProps } = useAntdForm<NonNullable<unknown>>({
    initialValues: settings?.content?.config?.[SSLPROVIDERS.LETS_ENCRYPT_STAGING],
    onSubmit: async (values) => {
      const newSettings = produce(settings, (draft) => {
        draft.content ??= {} as SSLProviderSettingsContent;
        draft.content.provider = SSLPROVIDERS.LETS_ENCRYPT_STAGING;

        draft.content.config ??= {} as SSLProviderSettingsContent["config"];
        draft.content.config[SSLPROVIDERS.LETS_ENCRYPT_STAGING] = values;
      });
      await updateSettings(newSettings);

      setFormChanged(false);
    },
  });

  const [formChanged, setFormChanged] = useState(false);
  useEffect(() => {
    setFormChanged(settings?.content?.provider !== SSLPROVIDERS.LETS_ENCRYPT_STAGING);
  }, [settings?.content?.provider]);

  const handleFormChange = () => {
    setFormChanged(true);
  };

  return (
    <Form {...formProps} form={formInst} disabled={pending} layout="vertical" onValuesChange={handleFormChange}>
      <Form.Item>
        <Alert type="info" message={<span dangerouslySetInnerHTML={{ __html: t("settings.sslprovider.form.letsencrypt_staging_alert") }}></span>} />
      </Form.Item>

      <Form.Item>
        <Button type="primary" htmlType="submit" disabled={!formChanged} loading={pending}>
          {t("common.button.save")}
        </Button>
      </Form.Item>
    </Form>
  );
};

const SSLProviderEditFormZeroSSLConfig = () => {
  const { t } = useTranslation();

  const { pending, settings, updateSettings } = useContext(SSLProviderContext);

  const formSchema = z.object({
    eabKid: z
      .string({ message: t("settings.sslprovider.form.zerossl_eab_kid.placeholder") })
      .min(1, t("settings.sslprovider.form.zerossl_eab_kid.placeholder"))
      .max(256, t("common.errmsg.string_max", { max: 256 })),
    eabHmacKey: z
      .string({ message: t("settings.sslprovider.form.zerossl_eab_hmac_key.placeholder") })
      .min(1, t("settings.sslprovider.form.zerossl_eab_hmac_key.placeholder"))
      .max(256, t("common.errmsg.string_max", { max: 256 })),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm<z.infer<typeof formSchema>>({
    initialValues: settings?.content?.config?.[SSLPROVIDERS.ZERO_SSL],
    onSubmit: async (values) => {
      const newSettings = produce(settings, (draft) => {
        draft.content ??= {} as SSLProviderSettingsContent;
        draft.content.provider = SSLPROVIDERS.ZERO_SSL;

        draft.content.config ??= {} as SSLProviderSettingsContent["config"];
        draft.content.config[SSLPROVIDERS.ZERO_SSL] = values;
      });
      await updateSettings(newSettings);

      setFormChanged(false);
    },
  });

  const [formChanged, setFormChanged] = useState(false);
  useEffect(() => {
    setFormChanged(settings?.content?.provider !== SSLPROVIDERS.ZERO_SSL);
  }, [settings?.content?.provider]);

  const handleFormChange = () => {
    setFormChanged(true);
  };

  return (
    <Form {...formProps} form={formInst} disabled={pending} layout="vertical" onValuesChange={handleFormChange}>
      <Form.Item
        name="eabKid"
        label={t("settings.sslprovider.form.zerossl_eab_kid.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("settings.sslprovider.form.zerossl_eab_kid.tooltip") }}></span>}
      >
        <Input autoComplete="new-password" placeholder={t("settings.sslprovider.form.zerossl_eab_kid.placeholder")} />
      </Form.Item>

      <Form.Item
        name="eabHmacKey"
        label={t("settings.sslprovider.form.zerossl_eab_hmac_key.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("settings.sslprovider.form.zerossl_eab_hmac_key.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("settings.sslprovider.form.zerossl_eab_hmac_key.placeholder")} />
      </Form.Item>

      <Form.Item>
        <Button type="primary" htmlType="submit" disabled={!formChanged} loading={pending}>
          {t("common.button.save")}
        </Button>
      </Form.Item>
    </Form>
  );
};

const SSLProviderEditFormGoogleTrustServicesConfig = () => {
  const { t } = useTranslation();

  const { pending, settings, updateSettings } = useContext(SSLProviderContext);

  const formSchema = z.object({
    eabKid: z
      .string({ message: t("settings.sslprovider.form.gts_eab_kid.placeholder") })
      .min(1, t("settings.sslprovider.form.gts_eab_kid.placeholder"))
      .max(256, t("common.errmsg.string_max", { max: 256 })),
    eabHmacKey: z
      .string({ message: t("settings.sslprovider.form.gts_eab_hmac_key.placeholder") })
      .min(1, t("settings.sslprovider.form.gts_eab_hmac_key.placeholder"))
      .max(256, t("common.errmsg.string_max", { max: 256 })),
  });
  const formRule = createSchemaFieldRule(formSchema);
  const { form: formInst, formProps } = useAntdForm<z.infer<typeof formSchema>>({
    initialValues: settings?.content?.config?.[SSLPROVIDERS.GOOGLE_TRUST_SERVICES],
    onSubmit: async (values) => {
      const newSettings = produce(settings, (draft) => {
        draft.content ??= {} as SSLProviderSettingsContent;
        draft.content.provider = SSLPROVIDERS.GOOGLE_TRUST_SERVICES;

        draft.content.config ??= {} as SSLProviderSettingsContent["config"];
        draft.content.config[SSLPROVIDERS.GOOGLE_TRUST_SERVICES] = values;
      });
      await updateSettings(newSettings);

      setFormChanged(false);
    },
  });

  const [formChanged, setFormChanged] = useState(false);
  useEffect(() => {
    setFormChanged(settings?.content?.provider !== SSLPROVIDERS.GOOGLE_TRUST_SERVICES);
  }, [settings?.content?.provider]);

  const handleFormChange = () => {
    setFormChanged(true);
  };

  return (
    <Form {...formProps} form={formInst} disabled={pending} layout="vertical" onValuesChange={handleFormChange}>
      <Form.Item
        name="eabKid"
        label={t("settings.sslprovider.form.gts_eab_kid.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("settings.sslprovider.form.gts_eab_kid.tooltip") }}></span>}
      >
        <Input autoComplete="new-password" placeholder={t("settings.sslprovider.form.gts_eab_kid.placeholder")} />
      </Form.Item>

      <Form.Item
        name="eabHmacKey"
        label={t("settings.sslprovider.form.gts_eab_hmac_key.label")}
        rules={[formRule]}
        tooltip={<span dangerouslySetInnerHTML={{ __html: t("settings.sslprovider.form.gts_eab_hmac_key.tooltip") }}></span>}
      >
        <Input.Password autoComplete="new-password" placeholder={t("settings.sslprovider.form.gts_eab_hmac_key.placeholder")} />
      </Form.Item>

      <Form.Item>
        <Button type="primary" htmlType="submit" disabled={!formChanged} loading={pending}>
          {t("common.button.save")}
        </Button>
      </Form.Item>
    </Form>
  );
};

const SettingsSSLProvider = () => {
  const { t } = useTranslation();

  const [messageApi, MessageContextHolder] = message.useMessage();
  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const [formInst] = Form.useForm<{ provider?: string }>();
  const [formPending, setFormPending] = useState(false);

  const [settings, setSettings] = useState<SettingsModel<SSLProviderSettingsContent>>();
  const [loading, setLoading] = useState(true);
  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);

      const settings = await getSettings<SSLProviderSettingsContent>(SETTINGS_NAMES.SSL_PROVIDER);
      setSettings(settings);
      setProviderType(settings.content?.provider);

      setLoading(false);
    };

    fetchData();
  }, []);

  const [providerType, setProviderType] = useState<SSLProviders>(SSLPROVIDERS.LETS_ENCRYPT);
  const providerFormEl = useMemo(() => {
    switch (providerType) {
      case SSLPROVIDERS.LETS_ENCRYPT:
        return <SSLProviderEditFormLetsEncryptConfig />;
      case SSLPROVIDERS.LETS_ENCRYPT_STAGING:
        return <SSLProviderEditFormLetsEncryptStagingConfig />;
      case SSLPROVIDERS.ZERO_SSL:
        return <SSLProviderEditFormZeroSSLConfig />;
      case SSLPROVIDERS.GOOGLE_TRUST_SERVICES:
        return <SSLProviderEditFormGoogleTrustServicesConfig />;
    }
  }, [providerType]);

  const updateContextSettings = async (settings: MaybeModelRecordWithId<SettingsModel<SSLProviderSettingsContent>>) => {
    setFormPending(true);

    try {
      const resp = await saveSettings(settings);
      setSettings(resp);
      setProviderType(resp.content?.provider);

      messageApi.success(t("common.text.operation_succeeded"));
    } catch (err) {
      notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
    } finally {
      setFormPending(false);
    }
  };

  return (
    <SSLProviderContext.Provider
      value={{
        pending: formPending,
        settings: settings!,
        updateSettings: updateContextSettings,
      }}
    >
      {MessageContextHolder}
      {NotificationContextHolder}

      <Show when={!loading} fallback={<Skeleton active />}>
        <Form form={formInst} disabled={formPending} layout="vertical" initialValues={{ provider: providerType }}>
          <Form.Item className="mb-2" name="provider" label={t("settings.sslprovider.form.provider.label")}>
            <CheckCard.Group className="w-full" onChange={(value) => setProviderType(value as SSLProviders)}>
              <CheckCard
                avatar={<img src={"/imgs/acme/letsencrypt.svg"} className="size-8" />}
                size="small"
                title={t("settings.sslprovider.form.provider.option.letsencrypt.label")}
                description="letsencrypt.org"
                value={SSLPROVIDERS.LETS_ENCRYPT}
              />
              <CheckCard
                avatar={<img src={"/imgs/acme/letsencrypt.svg"} className="size-8" />}
                size="small"
                title={t("settings.sslprovider.form.provider.option.letsencrypt_staging.label")}
                description="letsencrypt.org"
                value={SSLPROVIDERS.LETS_ENCRYPT_STAGING}
              />
              <CheckCard
                avatar={<img src={"/imgs/acme/zerossl.svg"} className="size-8" />}
                size="small"
                title={t("settings.sslprovider.form.provider.option.zerossl.label")}
                description="zerossl.com"
                value={SSLPROVIDERS.ZERO_SSL}
              />
              <CheckCard
                avatar={<img src={"/imgs/acme/google.svg"} className="size-8" />}
                size="small"
                title={t("settings.sslprovider.form.provider.option.gts.label")}
                description="pki.goog"
                value={SSLPROVIDERS.GOOGLE_TRUST_SERVICES}
              />
            </CheckCard.Group>
          </Form.Item>

          <Form.Item>
            <Alert type="warning" message={<span dangerouslySetInnerHTML={{ __html: t("settings.sslprovider.form.provider.alert") }}></span>} />
          </Form.Item>
        </Form>

        <div className="md:max-w-[40rem]">{providerFormEl}</div>
      </Show>
    </SSLProviderContext.Provider>
  );
};

export default SettingsSSLProvider;
