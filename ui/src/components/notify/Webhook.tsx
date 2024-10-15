import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { useToast } from "@/components/ui/use-toast";
import { getErrMessage } from "@/lib/error";
import { isValidURL } from "@/lib/url";
import { NotifyChannels, NotifyChannelWebhook } from "@/domain/settings";
import { update } from "@/repository/settings";
import { useNotify } from "@/providers/notify";

type WebhookSetting = {
  id: string;
  name: string;
  data: NotifyChannelWebhook;
};

const Webhook = () => {
  const { config, setChannels } = useNotify();
  const { t } = useTranslation();

  const [webhook, setWebhook] = useState<WebhookSetting>({
    id: config.id ?? "",
    name: "notifyChannels",
    data: {
      url: "",
      enabled: false,
    },
  });

  useEffect(() => {
    const getDetailWebhook = () => {
      const df: NotifyChannelWebhook = {
        url: "",
        enabled: false,
      };
      if (!config.content) {
        return df;
      }
      const chanels = config.content as NotifyChannels;
      if (!chanels.webhook) {
        return df;
      }

      return chanels.webhook as NotifyChannelWebhook;
    };
    const data = getDetailWebhook();
    setWebhook({
      id: config.id ?? "",
      name: "webhook",
      data,
    });
  }, [config]);

  const { toast } = useToast();

  const handleSaveClick = async () => {
    try {
      webhook.data.url = webhook.data.url.trim();
      if (!isValidURL(webhook.data.url)) {
        toast({
          title: t("common.save.failed.message"),
          description: t("settings.notification.url.errmsg.invalid"),
          variant: "destructive",
        });
        return;
      }

      const resp = await update({
        ...config,
        name: "notifyChannels",
        content: {
          ...config.content,
          webhook: {
            ...webhook.data,
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

  return (
    <div>
      <Input
        placeholder="Url"
        value={webhook.data.url}
        onChange={(e) => {
          setWebhook({
            ...webhook,
            data: {
              ...webhook.data,
              url: e.target.value,
            },
          });
        }}
      />

      <div className="flex items-center space-x-1 mt-2">
        <Switch
          id="airplane-mode"
          checked={webhook.data.enabled}
          onCheckedChange={() => {
            setWebhook({
              ...webhook,
              data: {
                ...webhook.data,
                enabled: !webhook.data.enabled,
              },
            });
          }}
        />
        <Label htmlFor="airplane-mode">{t("settings.notification.config.enable")}</Label>
      </div>

      <div className="flex justify-end mt-2">
        <Button
          onClick={() => {
            handleSaveClick();
          }}
        >
          {t("common.save")}
        </Button>
      </div>
    </div>
  );
};

export default Webhook;
