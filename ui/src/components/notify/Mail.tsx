import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { useToast } from "@/components/ui/use-toast";
import { getErrMessage } from "@/lib/error";
import { NotifyChannelMail, NotifyChannels } from "@/domain/settings";
import { useNotifyContext } from "@/providers/notify";
import { update } from "@/repository/settings";
import Show from "@/components/Show";
import { notifyTest } from "@/api/notify";

type MailSetting = {
  id: string;
  name: string;
  data: NotifyChannelMail;
};

const Mail = () => {
  const { config, setChannels } = useNotifyContext();
  const { t } = useTranslation();

  const [changed, setChanged] = useState<boolean>(false);

  const [mail, setmail] = useState<MailSetting>({
    id: config.id ?? "",
    name: "notifyChannels",
    data: {
      senderAddress: "",
      receiverAddresses: "",
      smtpHostAddr: "",
      smtpHostPort: "25",
      username: "",
      password: "",
      enabled: false,
    },
  });

  const [originMail, setoriginMail] = useState<MailSetting>({
    id: config.id ?? "",
    name: "notifyChannels",
    data: {
        senderAddress: "",
        receiverAddresses: "",
        smtpHostAddr: "",
        smtpHostPort: "25",
        username: "",
        password: "",
        enabled: false,
    },
  });

  useEffect(() => {
    setChanged(false);
  }, [config]);

  useEffect(() => {
    const data = getDetailMail();
    setoriginMail({
      id: config.id ?? "",
      name: "mail",
      data,
    });
  }, [config]);

  useEffect(() => {
    const data = getDetailMail();
    setmail({
      id: config.id ?? "",
      name: "mail",
      data,
    });
  }, [config]);

  const { toast } = useToast();

  const getDetailMail = () => {
    const df: NotifyChannelMail = {
        senderAddress: "",
        receiverAddresses: "",
        smtpHostAddr: "",
        smtpHostPort: "25",
        username: "",
        password: "",
        enabled: false,
    };
    if (!config.content) {
      return df;
    }
    const chanels = config.content as NotifyChannels;
    if (!chanels.mail) {
      return df;
    }

    return chanels.mail as NotifyChannelMail;
  };

  const checkChanged = (data: NotifyChannelMail) => {
    if (data.senderAddress !== originMail.data.senderAddress || data.receiverAddresses !== originMail.data.receiverAddresses || data.smtpHostAddr !== originMail.data.smtpHostAddr || data.smtpHostPort !== originMail.data.smtpHostPort || data.username !== originMail.data.username || data.password !== originMail.data.password) {
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
          mail: {
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
      await notifyTest("mail");

      toast({
        title: t("settings.notification.config.push.test.message.success.message"),
        description: t("settings.notification.config.push.test.message.success.message"),
      });
    } catch (e) {
      const msg = getErrMessage(e);

      toast({
        title: t("settings.notification.config.push.test.message.failed.message"),
        description: `${t("settings.notification.config.push.test.message.failed.message")}: ${msg}`,
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
    setmail(newData);

    try {
      const resp = await update({
        ...config,
        name: "notifyChannels",
        content: {
          ...config.content,
          mail: {
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
        placeholder={t("settings.notification.mail.sender_address.placeholder")}
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
          setmail(newData);
        }}
      />
      <Input
        placeholder={t("settings.notification.mail.receiver_address.placeholder")}
        className="mt-2"
        value={mail.data.receiverAddresses}
        onChange={(e) => {
          const newData = {
            ...mail,
            data: {
              ...mail.data,
              receiverAddresses: e.target.value,
            },
          };
          checkChanged(newData.data);
          setmail(newData);
        }}
      />
      <Input
        placeholder={t("settings.notification.mail.smtp_host.placeholder")}
        className="mt-2"
        value={mail.data.smtpHostAddr}
        onChange={(e) => {
          const newData = {
            ...mail,
            data: {
              ...mail.data,
              smtpHostAddr: e.target.value,
            },
          };
          checkChanged(newData.data);
          setmail(newData);
        }}
        />
    <Input
        placeholder={t("settings.notification.mail.smtp_port.placeholder")}
        className="mt-2"
        value={mail.data.smtpHostPort}
        onChange={(e) => {
        const newData = {
            ...mail,
            data: {
            ...mail.data,
            smtpHostPort: e.target.value,
            },
        };
        checkChanged(newData.data);
        setmail(newData);
        }}
    />
    <Input
        placeholder={t("settings.notification.mail.username.placeholder")}
        className="mt-2"
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
        setmail(newData);
        }}
        />
    <Input
        placeholder={t("settings.notification.mail.password.placeholder")}
        className="mt-2"
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
        setmail(newData);
        }}
        />
      <div className="flex items-center space-x-1 mt-2">
        <Switch id="airplane-mode" checked={mail.data.enabled} onCheckedChange={handleSwitchChange} />
        <Label htmlFor="airplane-mode">{t("settings.notification.config.enable")}</Label>
      </div>

      <div className="flex justify-end mt-2">
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
            {t("settings.notification.config.push.test.message")}
          </Button>
        </Show>
      </div>
    </div>
  );
};

export default Mail;
