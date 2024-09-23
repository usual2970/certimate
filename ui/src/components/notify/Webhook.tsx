import { Input } from "../ui/input";
import { Button } from "../ui/button";
import { Switch } from "../ui/switch";
import { Label } from "../ui/label";
import { useNotify } from "@/providers/notify";
import { NotifyChannels, NotifyChannelWebhook } from "@/domain/settings";
import { useEffect, useState } from "react";
import { update } from "@/repository/settings";
import { getErrMessage } from "@/lib/error";
import { useToast } from "../ui/use-toast";

type WebhookSetting = {
  id: string;
  name: string;
  data: NotifyChannelWebhook;
};

const Webhook = () => {
  const { config, setChannels } = useNotify();

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
        title: "保存成功",
        description: "配置保存成功",
      });
    } catch (e) {
      const msg = getErrMessage(e);

      toast({
        title: "保存失败",
        description: "配置保存失败：" + msg,
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
        <Label htmlFor="airplane-mode">是否启用</Label>
      </div>

      <div className="flex justify-end mt-2">
        <Button
          onClick={() => {
            handleSaveClick();
          }}
        >
          保存
        </Button>
      </div>
    </div>
  );
};

export default Webhook;
