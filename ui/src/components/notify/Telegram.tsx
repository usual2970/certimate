import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { useToast } from "@/components/ui/use-toast";
import { getErrMessage } from "@/lib/error";
import { NotifyChannels, NotifyChannelTelegram } from "@/domain/settings";
import { update } from "@/repository/settings";
import { useNotify } from "@/providers/notify";
import { notifyTest } from "@/api/notify";
import Show from "@/components/Show";

type TelegramSetting = {
  id: string;
  name: string;
  data: NotifyChannelTelegram;
};

const Telegram = () => {
  const { config, setChannels } = useNotify();
  const { t } = useTranslation();

  const [changed, setChanged] = useState<boolean>(false);

  const [telegram, setTelegram] = useState<TelegramSetting>({
    id: config.id ?? "",
    name: "notifyChannels",
    data: {
      apiToken: "",
      chatId: "",
      enabled: false,
    },
  });

  const [originTelegram, setOriginTelegram] = useState<TelegramSetting>({
    id: config.id ?? "",
    name: "notifyChannels",
    data: {
      apiToken: "",
      chatId: "",
      enabled: false,
    },
  });

  useEffect(() => {
    setChanged(false);
  }, [config]);

  useEffect(() => {
    const data = getDetailTelegram();
    setOriginTelegram({
      id: config.id ?? "",
      name: "common.provider.telegram",
      data,
    });
  }, [config]);

  useEffect(() => {
    const data = getDetailTelegram();
    setTelegram({
      id: config.id ?? "",
      name: "common.provider.telegram",
      data,
    });
  }, [config]);

  const { toast } = useToast();

  const checkChanged = (data: NotifyChannelTelegram) => {
    if (data.apiToken !== originTelegram.data.apiToken || data.chatId !== originTelegram.data.chatId) {
      setChanged(true);
    } else {
      setChanged(false);
    }
  };

  const getDetailTelegram = () => {
    const df: NotifyChannelTelegram = {
      apiToken: "",
      chatId: "",
      enabled: false,
    };
    if (!config.content) {
      return df;
    }
    const chanels = config.content as NotifyChannels;
    if (!chanels.telegram) {
      return df;
    }

    return chanels.telegram as NotifyChannelTelegram;
  };

  const handleSaveClick = async () => {
    try {
      const resp = await update({
        ...config,
        name: "notifyChannels",
        content: {
          ...config.content,
          telegram: {
            ...telegram.data,
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
      await notifyTest("telegram");

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
      ...telegram,
      data: {
        ...telegram.data,
        enabled: !telegram.data.enabled,
      },
    };
    setTelegram(newData);

    try {
      const resp = await update({
        ...config,
        name: "notifyChannels",
        content: {
          ...config.content,
          telegram: {
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
        placeholder="ApiToken"
        value={telegram.data.apiToken}
        onChange={(e) => {
          const newData = {
            ...telegram,
            data: {
              ...telegram.data,
              apiToken: e.target.value,
            },
          };

          checkChanged(newData.data);
          setTelegram(newData);
        }}
      />

      <Input
        placeholder="ChatId"
        value={telegram.data.chatId}
        onChange={(e) => {
          const newData = {
            ...telegram,
            data: {
              ...telegram.data,
              chatId: e.target.value,
            },
          };

          checkChanged(newData.data);
          setTelegram(newData);
        }}
      />

      <div className="flex items-center space-x-1 mt-2">
        <Switch id="airplane-mode" checked={telegram.data.enabled} onCheckedChange={handleSwitchChange} />
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

        <Show when={!changed && telegram.id != ""}>
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

export default Telegram;

