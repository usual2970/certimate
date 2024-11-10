import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { useToast } from "@/components/ui/use-toast";
import { getErrMessage } from "@/lib/error";
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
        title: t("common.save.succeeded.message"),
        description: t("settings.notification.config.saved.message"),
      });
    } catch (e) {
      const msg = getErrMessage(e);

      toast({
        title: t("common.save.failed.message"),
        description: `${t("settings.notification.config.failed.message")}: ${msg}`,
        variant: "destructive",
      });
    }
  };

  const handlePushTestClick = async () => {
    try {
      await notifyTest("email");

      toast({
        title: t("settings.notification.push_test_message.succeeded.message"),
        description: t("settings.notification.push_test_message.succeeded.message"),
      });
    } catch (e) {
      const msg = getErrMessage(e);

      toast({
        title: t("settings.notification.push_test_message.failed.message"),
        description: `${t("settings.notification.push_test_message.failed.message")}: ${msg}`,
        variant: "destructive",
      });
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
      const msg = getErrMessage(e);

      toast({
        title: t("common.save.failed.message"),
        description: `${t("settings.notification.config.failed.message")}: ${msg}`,
        variant: "destructive",
      });
    }
  };

  return (
    <div>
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

      <Input
        className="mt-2"
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

      <Switch
        className="mt-2"
        checked={mail.data.smtpTLS}
        onCheckedChange={(e) => {
          const newData = {
            ...mail,
            data: {
              ...mail.data,
              smtpTLS: e,
            },
          };
          checkChanged(newData.data);
          setMail(newData);
        }}
      />

      <Input
        className="mt-2"
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

      <Input
        className="mt-2"
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

      <Input
        className="mt-2"
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

      <Input
        className="mt-2"
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

      <div className="mt-2">
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
                {t("common.save")}
              </Button>
            </Show>

            <Show when={!changed && mail.id != ""}>
              <Button
                variant="secondary"
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
    </div>
  );
};

export default Mail;
