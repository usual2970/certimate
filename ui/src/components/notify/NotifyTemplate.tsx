import { Input } from "../ui/input";
import { Textarea } from "../ui/textarea";
import { Button } from "../ui/button";
import { useEffect, useState } from "react";
import {
  defaultNotifyTemplate,
  NotifyTemplates,
  NotifyTemplate as NotifyTemplateT,
} from "@/domain/settings";
import { getSetting, update } from "@/repository/settings";
import { useToast } from "../ui/use-toast";
import { useTranslation } from "react-i18next";

const NotifyTemplate = () => {
  const [id, setId] = useState("");
  const [templates, setTemplates] = useState<NotifyTemplateT[]>([
    defaultNotifyTemplate,
  ]);

  const { toast } = useToast();
  const { t } = useTranslation();

  useEffect(() => {
    const featchData = async () => {
      const resp = await getSetting("templates");

      if (resp.content) {
        setTemplates((resp.content as NotifyTemplates).notifyTemplates);
        setId(resp.id ? resp.id : "");
      }
    };
    featchData();
  }, []);

  const handleTitleChange = (val: string) => {
    const template = templates[0];

    setTemplates([
      {
        ...template,
        title: val,
      },
    ]);
  };

  const handleContentChange = (val: string) => {
    const template = templates[0];

    setTemplates([
      {
        ...template,
        content: val,
      },
    ]);
  };

  const handleSaveClick = async () => {
    const resp = await update({
      id: id,
      content: {
        notifyTemplates: templates,
      },
      name: "templates",
    });

    if (resp.id) {
      setId(resp.id);
    }

    toast({
      title: t("common.save.succeeded.message"),
      description: t("settings.notification.template.saved.message"),
    });
  };

  return (
    <div>
      <Input
        value={templates[0].title}
        onChange={(e) => {
          handleTitleChange(e.target.value);
        }}
      />

      <div className="text-muted-foreground text-sm mt-1">
        {t("settings.notification.template.variables.tips.title")}
      </div>

      <Textarea
        className="mt-2"
        value={templates[0].content}
        onChange={(e) => {
          handleContentChange(e.target.value);
        }}
      ></Textarea>
      <div className="text-muted-foreground text-sm mt-1">
        {t("settings.notification.template.variables.tips.content")}
      </div>
      <div className="flex justify-end mt-2">
        <Button onClick={handleSaveClick}>{t("common.save")}</Button>
      </div>
    </div>
  );
};

export default NotifyTemplate;
