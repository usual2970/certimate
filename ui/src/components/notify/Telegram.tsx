import { Input } from "../ui/input";
import { Button } from "../ui/button";
import { Switch } from "../ui/switch";
import { Label } from "../ui/label";
import { useNotify } from "@/providers/notify";
import { NotifyChannels, NotifyChannelTelegram } from "@/domain/settings";
import { useEffect, useState } from "react";
import { update } from "@/repository/settings";
import { getErrMessage } from "@/lib/error";
import { useToast } from "../ui/use-toast";
import { useTranslation } from "react-i18next";

type TelegramSetting = {
  id: string;
  name: string;
  data: NotifyChannelTelegram;
};

const Telegram = () => {
  const { config, setChannels } = useNotify();
  const { t } = useTranslation();

  const [telegram, setTelegram] = useState<TelegramSetting>({
    id: config.id ?? "",
    name: "notifyChannels",
    data: {
      apiToken: "",
      chatId: "",
      enabled: false,
    },
  });

  useEffect(() => {
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
    const data = getDetailTelegram();
    setTelegram({
      id: config.id ?? "",
      name: "telegram",
      data,
    });
  }, [config]);

  const { toast } = useToast();

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
        title: t('save.succeed'),
        description: t('setting.notify.config.save.succeed'),
      });
    } catch (e) {
      const msg = getErrMessage(e);

      toast({
        title: t('save.failed'),
        description: `${t('setting.notify.config.save.failed')}: ${msg}`,
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
          setTelegram({
            ...telegram,
            data: {
              ...telegram.data,
              apiToken: e.target.value,
            },
          });
        }}
      />

      <Input
        placeholder="ChatId"
        value={telegram.data.chatId}
        onChange={(e) => {
          setTelegram({
            ...telegram,
            data: {
              ...telegram.data,
              chatId: e.target.value,
            },
          });
        }}
      />

      <div className="flex items-center space-x-1 mt-2">
        <Switch
          id="airplane-mode"
          checked={telegram.data.enabled}
          onCheckedChange={() => {
            setTelegram({
              ...telegram,
              data: {
                ...telegram.data,
                enabled: !telegram.data.enabled,
              },
            });
          }}
        />
        <Label htmlFor="airplane-mode">{t('setting.notify.config.enable')}</Label>
      </div>

      <div className="flex justify-end mt-2">
        <Button
          onClick={() => {
            handleSaveClick();
          }}
        >
          {t('save')}
        </Button>
      </div>
    </div>
  );
};

export default Telegram;
