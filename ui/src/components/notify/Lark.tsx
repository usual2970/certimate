import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import { useNotifyContext } from "@/providers/notify";
import { NotifyChannelLark, NotifyChannels } from "@/domain/settings";
import { useEffect, useState } from "react";
import { update } from "@/repository/settings";
import { getErrMessage } from "@/lib/error";
import { useToast } from "@/components/ui/use-toast";
import { useTranslation } from "react-i18next";
import { notifyTest } from "@/api/notify";
import Show from "@/components/Show";

type LarkSetting = {
  id: string;
  name: string;
  data: NotifyChannelLark;
};

const Lark = () => {
  const { config, setChannels } = useNotifyContext();
  const { t } = useTranslation();

  const [changed, setChanged] = useState<boolean>(false);

  const [lark, setLark] = useState<LarkSetting>({
    id: config.id ?? "",
    name: "notifyChannels",
    data: {
      webhookUrl: "",
      enabled: false,
    },
  });

  const [originLark, setOriginLark] = useState<LarkSetting>({
    id: config.id ?? "",
    name: "notifyChannels",
    data: {
      webhookUrl: "",
      enabled: false,
    },
  });

  useEffect(() => {
    setChanged(false);
  }, [config]);

  useEffect(() => {
    const data = getDetailLark();
    setOriginLark({
      id: config.id ?? "",
      name: "lark",
      data,
    });
  }, [config]);

  useEffect(() => {
    const data = getDetailLark();
    setLark({
      id: config.id ?? "",
      name: "lark",
      data,
    });
  }, [config]);

  const { toast } = useToast();

  const checkChanged = (data: NotifyChannelLark) => {
    if (data.webhookUrl !== originLark.data.webhookUrl) {
      setChanged(true);
    } else {
      setChanged(false);
    }
  };

  const getDetailLark = () => {
    const df: NotifyChannelLark = {
      webhookUrl: "",
      enabled: false,
    };
    if (!config.content) {
      return df;
    }
    const chanels = config.content as NotifyChannels;
    if (!chanels.lark) {
      return df;
    }

    return chanels.lark as NotifyChannelLark;
  };

  const handleSaveClick = async () => {
    try {
      const resp = await update({
        ...config,
        name: "notifyChannels",
        content: {
          ...config.content,
          lark: {
            ...lark.data,
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
      await notifyTest("lark");

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
      ...lark,
      data: {
        ...lark.data,
        enabled: !lark.data.enabled,
      },
    };
    setLark(newData);

    try {
      const resp = await update({
        ...config,
        name: "notifyChannels",
        content: {
          ...config.content,
          lark: {
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
        placeholder="Webhook Url"
        value={lark.data.webhookUrl}
        onChange={(e) => {
          const newData = {
            ...lark,
            data: {
              ...lark.data,
              webhookUrl: e.target.value,
            },
          };

          checkChanged(newData.data);
          setLark(newData);
        }}
      />
      <div className="flex items-center space-x-1 mt-2">
        <Switch id="airplane-mode" checked={lark.data.enabled} onCheckedChange={handleSwitchChange} />
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

        <Show when={!changed && lark.id != ""}>
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

export default Lark;
