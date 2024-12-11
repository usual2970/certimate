import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Switch } from "@/components/ui/switch";
import { Label } from "@/components/ui/label";
import { useNotifyContext } from "@/providers/notify";
import { NotifyChannelLark, NotifyChannels } from "@/domain/settings";
import { useEffect, useState } from "react";
import { save } from "@/repository/settings";
import { getErrMsg } from "@/utils/error";
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
      const resp = await save({
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
    } finally {
      setTesting(false);
    }
  };

  const [testing, setTesting] = useState<boolean>(false);
  const handlePushTestClick = async () => {
    if (testing) return;

    try {
      setTesting(true);

      await notifyTest("lark");

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
      const resp = await save({
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
      <div>
        <Label>{t("settings.notification.lark.webhook_url.label")}</Label>
        <Input
          placeholder={t("settings.notification.lark.webhook_url.placeholder")}
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
      </div>

      <div className="flex justify-between gap-4">
        <div className="flex items-center space-x-1">
          <Switch id="airplane-mode" checked={lark.data.enabled} onCheckedChange={handleSwitchChange} />
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

          <Show when={!changed && lark.id != ""}>
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

export default Lark;
