import { Input } from "../ui/input";
import { Button } from "../ui/button";
import { Switch } from "../ui/switch";
import { Label } from "../ui/label";
import { useNotify } from "@/providers/notify";
import { NotifyChannelLark, NotifyChannels } from "@/domain/settings";
import { useEffect, useState } from "react";
import { update } from "@/repository/settings";
import { getErrMessage } from "@/lib/error";
import { useToast } from "../ui/use-toast";
import { useTranslation } from "react-i18next";

type LarkSetting = {
  id: string;
  name: string;
  data: NotifyChannelLark;
};

const Lark = () => {
  const { config, setChannels } = useNotify();
  const { t } = useTranslation();

  const [lark, setLark] = useState<LarkSetting>({
    id: config.id ?? "",
    name: "notifyChannels",
    data: {
      webhookUrl: "",
      enabled: false,
    },
  });

  useEffect(() => {
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
    const data = getDetailLark();
    setLark({
      id: config.id ?? "",
      name: "lark",
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
        description: `${t(
          "settings.notification.config.failed.message"
        )}: ${msg}`,
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
          setLark({
            ...lark,
            data: {
              ...lark.data,
              webhookUrl: e.target.value,
            },
          });
        }}
      />
      <div className="flex items-center space-x-1 mt-2">
        <Switch
          id="airplane-mode"
          checked={lark.data.enabled}
          onCheckedChange={() => {
            setLark({
              ...lark,
              data: {
                ...lark.data,
                enabled: !lark.data.enabled,
              },
            });
          }}
        />
        <Label htmlFor="airplane-mode">
          {t("settings.notification.config.enable")}
        </Label>
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

export default Lark;
