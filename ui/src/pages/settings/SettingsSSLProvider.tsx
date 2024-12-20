import { createContext, useContext, useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import { Button, Form, Input, message, notification, Skeleton } from "antd";
import { CheckCard } from "@ant-design/pro-components";
import { createSchemaFieldRule } from "antd-zod";
import { produce } from "immer";
import { z } from "zod";

import { SETTINGS_NAMES, SSLPROVIDERS, type SettingsModel, type SSLProviderSettingsContent, type SSLProviders } from "@/domain/settings";
import { get as getSettings, save as saveSettings } from "@/repository/settings";
import { getErrMsg } from "@/utils/error";
import { useDeepCompareEffect } from "ahooks";

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

  const [form] = Form.useForm();

  const [initialValues, setInitialValues] = useState(settings?.content?.config?.[SSLPROVIDERS.LETS_ENCRYPT]);
  const [initialChanged, setInitialChanged] = useState(false);
  useDeepCompareEffect(() => {
    setInitialValues(settings?.content?.config?.[SSLPROVIDERS.LETS_ENCRYPT]);
    setInitialChanged(settings?.content?.provider !== SSLPROVIDERS.LETS_ENCRYPT);
  }, [settings]);

  const handleFormChange = () => {
    setInitialChanged(true);
  };

  const handleFormFinish = async (fields: NonNullable<unknown>) => {
    const newSettings = produce(settings, (draft) => {
      draft.content ??= {} as SSLProviderSettingsContent;
      draft.content.provider = SSLPROVIDERS.LETS_ENCRYPT;

      draft.content.config ??= {} as SSLProviderSettingsContent["config"];
      draft.content.config[SSLPROVIDERS.LETS_ENCRYPT] = fields;
    });
    await updateSettings(newSettings);

    setInitialChanged(false);
  };

  return (
    <Form form={form} disabled={pending} layout="vertical" initialValues={initialValues} onFinish={handleFormFinish} onValuesChange={handleFormChange}>
      <Form.Item>
        <Button type="primary" htmlType="submit" disabled={!initialChanged} loading={pending}>
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
  const [form] = Form.useForm<z.infer<typeof formSchema>>();

  const [initialValues, setInitialValues] = useState(settings?.content?.config?.[SSLPROVIDERS.ZERO_SSL]);
  const [initialChanged, setInitialChanged] = useState(false);
  useDeepCompareEffect(() => {
    setInitialValues(settings?.content?.config?.[SSLPROVIDERS.ZERO_SSL]);
    setInitialChanged(settings?.content?.provider !== SSLPROVIDERS.ZERO_SSL);
  }, [settings]);

  const handleFormChange = () => {
    setInitialChanged(true);
  };

  const handleFormFinish = async (fields: z.infer<typeof formSchema>) => {
    const newSettings = produce(settings, (draft) => {
      draft.content ??= {} as SSLProviderSettingsContent;
      draft.content.provider = SSLPROVIDERS.ZERO_SSL;

      draft.content.config ??= {} as SSLProviderSettingsContent["config"];
      draft.content.config[SSLPROVIDERS.ZERO_SSL] = fields;
    });
    await updateSettings(newSettings);

    setInitialChanged(false);
  };

  return (
    <Form form={form} disabled={pending} layout="vertical" initialValues={initialValues} onFinish={handleFormFinish} onValuesChange={handleFormChange}>
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
        <Button type="primary" htmlType="submit" disabled={!initialChanged} loading={pending}>
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
  const [form] = Form.useForm<z.infer<typeof formSchema>>();

  const [initialValues, setInitialValues] = useState(settings?.content?.config?.[SSLPROVIDERS.GOOGLE_TRUST_SERVICES]);
  const [initialChanged, setInitialChanged] = useState(false);
  useDeepCompareEffect(() => {
    setInitialValues(settings?.content?.config?.[SSLPROVIDERS.GOOGLE_TRUST_SERVICES]);
    setInitialChanged(settings?.content?.provider !== SSLPROVIDERS.GOOGLE_TRUST_SERVICES);
  }, [settings]);

  const handleFormChange = () => {
    setInitialChanged(true);
  };

  const handleFormFinish = async (fields: z.infer<typeof formSchema>) => {
    const newSettings = produce(settings, (draft) => {
      draft.content ??= {} as SSLProviderSettingsContent;
      draft.content.provider = SSLPROVIDERS.GOOGLE_TRUST_SERVICES;

      draft.content.config ??= {} as SSLProviderSettingsContent["config"];
      draft.content.config[SSLPROVIDERS.GOOGLE_TRUST_SERVICES] = fields;
    });
    await updateSettings(newSettings);

    setInitialChanged(false);
  };

  return (
    <Form form={form} disabled={pending} layout="vertical" initialValues={initialValues} onFinish={handleFormFinish} onValuesChange={handleFormChange}>
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
        <Button type="primary" htmlType="submit" disabled={!initialChanged} loading={pending}>
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

  const [form] = Form.useForm();
  const [formPending, setFormPending] = useState(false);

  const [settings, setSettings] = useState<SettingsModel<SSLProviderSettingsContent>>();
  const [loading, setLoading] = useState(true);
  useEffect(() => {
    const fetchData = async () => {
      setLoading(true);

      const settings = await getSettings<SSLProviderSettingsContent>(SETTINGS_NAMES.SSL_PROVIDER);
      setSettings(settings);
      setFormProviderType(settings.content?.provider);

      setLoading(false);
    };

    fetchData();
  }, []);

  const [providerType, setFormProviderType] = useState<SSLProviders>();
  const providerFormComponent = useMemo(() => {
    switch (providerType) {
      case SSLPROVIDERS.LETS_ENCRYPT:
        return <SSLProviderEditFormLetsEncryptConfig />;
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
      setFormProviderType(resp.content?.provider);

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

      {loading ? (
        <Skeleton active />
      ) : (
        <>
          <Form form={form} disabled={formPending} layout="vertical" initialValues={{ provider: providerType }}>
            <Form.Item className="mb-2" name="provider" label={t("settings.sslprovider.form.provider.label")} initialValue={SSLPROVIDERS.LETS_ENCRYPT}>
              <CheckCard.Group className="w-full" onChange={(value) => setFormProviderType(value as SSLProviders)}>
                <CheckCard
                  avatar={<img src={"/imgs/acme/letsencrypt.svg"} className="size-8" />}
                  size="small"
                  title="Let's Encrypt"
                  value={SSLPROVIDERS.LETS_ENCRYPT}
                />
                <CheckCard avatar={<img src={"/imgs/acme/zerossl.svg"} className="size-8" />} size="small" title="ZeroSSL" value={SSLPROVIDERS.ZERO_SSL} />
                <CheckCard
                  avatar={<img src={"/imgs/acme/google.svg"} className="size-8" />}
                  size="small"
                  title="Google Trust Services"
                  value={SSLPROVIDERS.GOOGLE_TRUST_SERVICES}
                />
              </CheckCard.Group>
            </Form.Item>
          </Form>

          <div className="md:max-w-[40rem]">{providerFormComponent}</div>
        </>
      )}
    </SSLProviderContext.Provider>
  );
};

export default SettingsSSLProvider;
