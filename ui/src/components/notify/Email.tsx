import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { useToast } from "@/components/ui/use-toast";
import { getErrMsg } from "@/utils/error";
import { NotifyChannelEmail, NotifyChannels } from "@/domain/settings";
import { useNotifyContext } from "@/providers/notify";
import { update } from "@/repository/settings";
import Show from "@/components/Show";
import { notifyTest } from "@/api/notify";

type EmailSetting = {
  id: string;
  name: string;
  data: NotifyChannelEmail;
};

const Mail = () => {
  const { config, setChannels } = useNotifyContext();
  const { t } = useTranslation();

  const [changed, setChanged] = useState<boolean>(false);

  const [mail, setMail] = useState<EmailSetting>({
    id: config.id ?? "",
    name: "notifyChannels",
    data: {
      smtpHost: "",
      smtpPort: 465,
      smtpTLS: true,
      username: "",
      password: "",
      senderAddress: "",
      receiverAddress: "",
      enabled: false,
    },
  });

  const [originMail, setOriginMail] = useState<EmailSetting>({
    id: config.id ?? "",
    name: "notifyChannels",
    data: {
      smtpHost: "",
      smtpPort: 465,
      smtpTLS: true,
      username: "",
      password: "",
      senderAddress: "",
      receiverAddress: "",
      enabled: false,
    },
  });

  useEffect(() => {
    setChanged(false);
  }, [config]);

  useEffect(() => {
    const data = getDetailMail();
    setOriginMail({
      id: config.id ?? "",
      name: "email",
      data,
    });
  }, [config]);

  useEffect(() => {
    const data = getDetailMail();
    setMail({
      id: config.id ?? "",
      name: "email",
      data,
    });
  }, [config]);

  const { toast } = useToast();

  const getDetailMail = () => {
    const df: NotifyChannelEmail = {
      smtpHost: "smtp.example.com",
      smtpPort: 465,
      smtpTLS: true,
      username: "",
      password: "",
      senderAddress: "",
      receiverAddress: "",
      enabled: false,
    };
    if (!config.content) {
      return df;
    }
    const chanels = config.content as NotifyChannels;
    if (!chanels.email) {
      return df;
    }

    return chanels.email as NotifyChannelEmail;
  };

  const checkChanged = (data: NotifyChannelEmail) => {
    if (
      data.smtpHost !== originMail.data.smtpHost ||
      data.smtpPort !== originMail.data.smtpPort ||
      data.smtpTLS !== originMail.data.smtpTLS ||
      data.username !== originMail.data.username ||
      data.password !== originMail.data.password ||
      data.senderAddress !== originMail.data.senderAddress ||
      data.receiverAddress !== originMail.data.receiverAddress
    ) {
      setChanged(true);
    } else {
      setChanged(false);
    }
  };

  const handleSaveClick = async () => {
    try {
      const resp = await update({
        ...config,
        name: "notifyChannels",
        content: {
          ...config.content,
          email: {
            ...mail.data,
          },
        },
      });

      setChannels(resp);
      toast({
        title: t("common.text.operation_succeeded"),
        description: t("settings.notification.config.saved.message"),
      });
    } catch (e) {
      const msg = getErrMsg(e);

      toast({
        title: t("common.text.operation_failed"),
        description: `${t("settings.notification.config.failed.message")}: ${msg}`,
        variant: "destructive",
      });
    }
  };

  const [testing, setTesting] = useState<boolean>(false);
  const handlePushTestClick = async () => {
    if (testing) return;

    try {
      setTesting(true);

      await notifyTest("email");

      toast({
        title: t("settings.notification.push_test_message.succeeded.message"),
        description: t("settings.notification.push_test_message.succeeded.message"),
      });
    } catch (e) {
      const msg = getErrMsg(e);

      toast({
        title: t("settings.notification.push_test_message.failed.message"),
        description: `${t("settings.notification.push_test_message.failed.message")}: ${msg}`,
        variant: "destructive",
      });
    } finally {
      setTesting(false);
    }
  };

  const handleSwitchChange = async () => {
    const newData = {
      ...mail,
      data: {
        ...mail.data,
        enabled: !mail.data.enabled,
      },
    };
    setMail(newData);

    try {
      const resp = await update({
        ...config,
        name: "notifyChannels",
        content: {
          ...config.content,
          email: {
            ...newData.data,
          },
        },
      });

      setChannels(resp);
    } catch (e) {
      const msg = getErrMsg(e);

      toast({
        title: t("common.text.operation_failed"),
        description: `${t("settings.notification.config.failed.message")}: ${msg}`,
        variant: "destructive",
      });
    }
  };

  return (
    <div className="flex flex-col space-y-4">
      <div className="flex space-x-4">
        <div className="w-2/5">
          <Label>{t("settings.notification.email.smtp_host.label")}</Label>
          <Input
            placeholder={t("settings.notification.email.smtp_host.placeholder")}
            value={mail.data.smtpHost}
            onChange={(e) => {
              const newData = {
                ...mail,
                data: {
                  ...mail.data,
                  smtpHost: e.target.value,
                },
              };
              checkChanged(newData.data);
              setMail(newData);
            }}
          />
        </div>

        <div className="w-2/5">
          <Label>{t("settings.notification.email.smtp_port.label")}</Label>
          <Input
            type="number"
            placeholder={t("settings.notification.email.smtp_port.placeholder")}
            value={mail.data.smtpPort}
            onChange={(e) => {
              const newData = {
                ...mail,
                data: {
                  ...mail.data,
                  smtpPort: +e.target.value || 0,
                },
              };
              checkChanged(newData.data);
              setMail(newData);
            }}
          />
        </div>

        <div className="w-1/5">
          <Label>{t("settings.notification.email.smtp_tls.label")}</Label>
          <Switch
            className="block mt-2"
            checked={mail.data.smtpTLS}
            onCheckedChange={(e) => {
              const newData = {
                ...mail,
                data: {
                  ...mail.data,
                  smtpPort: e && mail.data.smtpPort === 25 ? 465 : !e && mail.data.smtpPort === 465 ? 25 : mail.data.smtpPort,
                  smtpTLS: e,
                },
              };
              checkChanged(newData.data);
              setMail(newData);
            }}
          />
        </div>
      </div>

      <div className="flex space-x-4">
        <div className="w-1/2">
          <Label>{t("settings.notification.email.username.label")}</Label>
          <Input
            placeholder={t("settings.notification.email.username.placeholder")}
            value={mail.data.username}
            onChange={(e) => {
              const newData = {
                ...mail,
                data: {
                  ...mail.data,
                  username: e.target.value,
                },
              };
              checkChanged(newData.data);
              setMail(newData);
            }}
          />
        </div>

        <div className="w-1/2">
          <Label>{t("settings.notification.email.password.label")}</Label>
          <Input
            placeholder={t("settings.notification.email.password.placeholder")}
            value={mail.data.password}
            onChange={(e) => {
              const newData = {
                ...mail,
                data: {
                  ...mail.data,
                  password: e.target.value,
                },
              };
              checkChanged(newData.data);
              setMail(newData);
            }}
          />
        </div>
      </div>

      <div>
        <Label>{t("settings.notification.email.sender_address.label")}</Label>
        <Input
          placeholder={t("settings.notification.email.sender_address.placeholder")}
          value={mail.data.senderAddress}
          onChange={(e) => {
            const newData = {
              ...mail,
              data: {
                ...mail.data,
                senderAddress: e.target.value,
              },
            };
            checkChanged(newData.data);
            setMail(newData);
          }}
        />
      </div>

      <div>
        <Label>{t("settings.notification.email.receiver_address.label")}</Label>
        <Input
          placeholder={t("settings.notification.email.receiver_address.placeholder")}
          value={mail.data.receiverAddress}
          onChange={(e) => {
            const newData = {
              ...mail,
              data: {
                ...mail.data,
                receiverAddress: e.target.value,
              },
            };
            checkChanged(newData.data);
            setMail(newData);
          }}
        />
      </div>

      <div className="flex justify-between gap-4">
        <div className="flex items-center space-x-1">
          <Switch id="airplane-mode" checked={mail.data.enabled} onCheckedChange={handleSwitchChange} />
          <Label htmlFor="airplane-mode">{t("settings.notification.config.enable")}</Label>
        </div>

        <div className="flex items-center space-x-1">
          <Show when={changed}>
            <Button
              onClick={() => {
                handleSaveClick();
              }}
            >
              {t("common.button.save")}
            </Button>
          </Show>

          <Show when={!changed && mail.id != ""}>
            <Button
              variant="secondary"
              loading={testing}
              onClick={() => {
                handlePushTestClick();
              }}
            >
              {t("settings.notification.push_test_message")}
            </Button>
          </Show>
        </div>
      </div>
    </div>
  );
};

export default Mail;
