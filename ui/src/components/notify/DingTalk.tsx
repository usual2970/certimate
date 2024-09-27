import { Input } from "../ui/input";
import { Button } from "../ui/button";
import { Switch } from "../ui/switch";
import { Label } from "../ui/label";
import { useNotify } from "@/providers/notify";
import { NotifyChannelDingTalk, NotifyChannels } from "@/domain/settings";
import { useEffect, useState } from "react";
import { update } from "@/repository/settings";
import { getErrMessage } from "@/lib/error";
import { useToast } from "../ui/use-toast";
import { useTranslation } from 'react-i18next'

type DingTalkSetting = {
  id: string;
  name: string;
  data: NotifyChannelDingTalk;
};

const DingTalk = () => {
  const { config, setChannels } = useNotify();
  const { t } = useTranslation();

  const [dingtalk, setDingtalk] = useState<DingTalkSetting>({
    id: config.id ?? "",
    name: "notifyChannels",
    data: {
      accessToken: "",
      secret: "",
      enabled: false,
    },
  });

  useEffect(() => {
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
    const data = getDetailDingTalk();
    setDingtalk({
      id: config.id ?? "",
      name: "dingtalk",
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
          dingtalk: {
            ...dingtalk.data,
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
        placeholder="AccessToken"
        value={dingtalk.data.accessToken}
        onChange={(e) => {
          setDingtalk({
            ...dingtalk,
            data: {
              ...dingtalk.data,
              accessToken: e.target.value,
            },
          });
        }}
      />
      <Input
        placeholder={t('access.form.ding.access.token.placeholder')}
        className="mt-2"
        value={dingtalk.data.secret}
        onChange={(e) => {
          setDingtalk({
            ...dingtalk,
            data: {
              ...dingtalk.data,
              secret: e.target.value,
            },
          });
        }}
      />
      <div className="flex items-center space-x-1 mt-2">
        <Switch
          id="airplane-mode"
          checked={dingtalk.data.enabled}
          onCheckedChange={() => {
            setDingtalk({
              ...dingtalk,
              data: {
                ...dingtalk.data,
                enabled: !dingtalk.data.enabled,
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

export default DingTalk;
