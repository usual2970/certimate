import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { useToast } from "@/components/ui/use-toast";
import { getErrMessage } from "@/lib/error";
import { NotifyChannels, NotifyChannelBark } from "@/domain/settings";
import { update } from "@/repository/settings";
import { useNotifyContext } from "@/providers/notify";
import { notifyTest } from "@/api/notify";
import Show from "@/components/Show";

type BarkSetting = {
  id: string;
  name: string;
  data: NotifyChannelBark;
};

const Bark = () => {
  const { config, setChannels } = useNotifyContext();
  const { t } = useTranslation();

  const [changed, setChanged] = useState<boolean>(false);

  const [bark, setBark] = useState<BarkSetting>({
    id: config.id ?? "",
    name: "notifyChannels",
    data: {
      serverUrl: "",
      deviceKey: "",
      enabled: false,
    },
  });

  const [originBark, setOriginBark] = useState<BarkSetting>({
    id: config.id ?? "",
    name: "notifyChannels",
    data: {
      serverUrl: "",
      deviceKey: "",
      enabled: false,
    },
  });

  useEffect(() => {
    setChanged(false);
  }, [config]);

  useEffect(() => {
    const data = getDetailBark();
    setOriginBark({
      id: config.id ?? "",
      name: "common.provider.bark",
      data,
    });
  }, [config]);

  useEffect(() => {
    const data = getDetailBark();
    setBark({
      id: config.id ?? "",
      name: "common.provider.bark",
      data,
    });
  }, [config]);

  const { toast } = useToast();

  const checkChanged = (data: NotifyChannelBark) => {
    if (data.serverUrl !== originBark.data.serverUrl || data.deviceKey !== originBark.data.deviceKey) {
      setChanged(true);
    } else {
      setChanged(false);
    }
  };

  const getDetailBark = () => {
    const df: NotifyChannelBark = {
      serverUrl: "",
      deviceKey: "",
      enabled: false,
    };
    if (!config.content) {
      return df;
    }
    const chanels = config.content as NotifyChannels;
    if (!chanels.bark) {
      return df;
    }

    return chanels.bark as NotifyChannelBark;
  };

  const handleSaveClick = async () => {
    try {
      const resp = await update({
        ...config,
        name: "notifyChannels",
        content: {
          ...config.content,
          bark: {
            ...bark.data,
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
      await notifyTest("bark");

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
      ...bark,
      data: {
        ...bark.data,
        enabled: !bark.data.enabled,
      },
    };
    setBark(newData);

    try {
      const resp = await update({
        ...config,
        name: "notifyChannels",
        content: {
          ...config.content,
          bark: {
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
    <div className="flex flex-col space-y-4">
      <div>
        <Label>{t("settings.notification.bark.server_url.label")}</Label>
        <Input
          placeholder={t("settings.notification.bark.server_url.placeholder")}
          value={bark.data.serverUrl}
          onChange={(e) => {
            const newData = {
              ...bark,
              data: {
                ...bark.data,
                serverUrl: e.target.value,
              },
            };

            checkChanged(newData.data);
            setBark(newData);
          }}
        />
      </div>

      <div>
        <Label>{t("settings.notification.bark.device_key.label")}</Label>
        <Input
          placeholder={t("settings.notification.bark.device_key.placeholder")}
          value={bark.data.deviceKey}
          onChange={(e) => {
            const newData = {
              ...bark,
              data: {
                ...bark.data,
                deviceKey: e.target.value,
              },
            };

            checkChanged(newData.data);
            setBark(newData);
          }}
        />
      </div>

      <div className="flex justify-between gap-4">
        <div className="flex items-center space-x-1">
          <Switch id="airplane-mode" checked={bark.data.enabled} onCheckedChange={handleSwitchChange} />
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

          <Show when={!changed && bark.id != ""}>
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
  );
};

export default Bark;
