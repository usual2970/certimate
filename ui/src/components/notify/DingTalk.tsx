import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { useToast } from "@/components/ui/use-toast";
import { getErrMsg } from "@/utils/error";
import { NotifyChannelDingTalk, NotifyChannels } from "@/domain/settings";
import { useNotifyContext } from "@/providers/notify";
import { save } from "@/repository/settings";
import Show from "@/components/Show";
import { notifyTest } from "@/api/notify";

type DingTalkSetting = {
  id: string;
  name: string;
  data: NotifyChannelDingTalk;
};

const DingTalk = () => {
  const { config, setChannels } = useNotifyContext();
  const { t } = useTranslation();

  const [changed, setChanged] = useState<boolean>(false);

  const [dingtalk, setDingtalk] = useState<DingTalkSetting>({
    id: config.id ?? "",
    name: "notifyChannels",
    data: {
      accessToken: "",
      secret: "",
      enabled: false,
    },
  });

  const [originDingtalk, setOriginDingtalk] = useState<DingTalkSetting>({
    id: config.id ?? "",
    name: "notifyChannels",
    data: {
      accessToken: "",
      secret: "",
      enabled: false,
    },
  });

  useEffect(() => {
    setChanged(false);
  }, [config]);

  useEffect(() => {
    const data = getDetailDingTalk();
    setOriginDingtalk({
      id: config.id ?? "",
      name: "dingtalk",
      data,
    });
  }, [config]);

  useEffect(() => {
    const data = getDetailDingTalk();
    setDingtalk({
      id: config.id ?? "",
      name: "dingtalk",
      data,
    });
  }, [config]);

  const { toast } = useToast();

  const getDetailDingTalk = () => {
    const df: NotifyChannelDingTalk = {
      accessToken: "",
      secret: "",
      enabled: false,
    };
    if (!config.content) {
      return df;
    }
    const chanels = config.content as NotifyChannels;
    if (!chanels.dingtalk) {
      return df;
    }

    return chanels.dingtalk as NotifyChannelDingTalk;
  };

  const checkChanged = (data: NotifyChannelDingTalk) => {
    if (data.accessToken !== originDingtalk.data.accessToken || data.secret !== originDingtalk.data.secret) {
      setChanged(true);
    } else {
      setChanged(false);
    }
  };

  const handleSaveClick = async () => {
    try {
      const resp = await save({
        ...config,
        name: "notifyChannels",
        content: {
          ...config.content,
          dingtalk: {
            ...dingtalk.data,
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

      await notifyTest("dingtalk");

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
      ...dingtalk,
      data: {
        ...dingtalk.data,
        enabled: !dingtalk.data.enabled,
      },
    };
    setDingtalk(newData);

    try {
      const resp = await save({
        ...config,
        name: "notifyChannels",
        content: {
          ...config.content,
          dingtalk: {
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
        <Label>{t("settings.notification.dingtalk.access_token.label")}</Label>
        <Input
          placeholder={t("settings.notification.dingtalk.access_token.placeholder")}
          value={dingtalk.data.accessToken}
          onChange={(e) => {
            const newData = {
              ...dingtalk,
              data: {
                ...dingtalk.data,
                accessToken: e.target.value,
              },
            };
            checkChanged(newData.data);
            setDingtalk(newData);
          }}
        />
      </div>

      <div>
        <Label>{t("settings.notification.dingtalk.secret.label")}</Label>
        <Input
          placeholder={t("settings.notification.dingtalk.secret.placeholder")}
          value={dingtalk.data.secret}
          onChange={(e) => {
            const newData = {
              ...dingtalk,
              data: {
                ...dingtalk.data,
                secret: e.target.value,
              },
            };
            checkChanged(newData.data);
            setDingtalk(newData);
          }}
        />
      </div>

      <div className="flex justify-between gap-4">
        <div className="flex items-center space-x-1">
          <Switch id="airplane-mode" checked={dingtalk.data.enabled} onCheckedChange={handleSwitchChange} />
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

          <Show when={!changed && dingtalk.id != ""}>
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

export default DingTalk;
