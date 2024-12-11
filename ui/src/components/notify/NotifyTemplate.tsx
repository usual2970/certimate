import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";

import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { useToast } from "@/components/ui/use-toast";
import { defaultNotifyTemplate, NotifyTemplates, NotifyTemplate as NotifyTemplateT } from "@/domain/settings";
import { get, save } from "@/repository/settings";

const NotifyTemplate = () => {
  const [id, setId] = useState("");
  const [templates, setTemplates] = useState<NotifyTemplateT[]>([defaultNotifyTemplate]);

  const { toast } = useToast();
  const { t } = useTranslation();

  useEffect(() => {
    const featchData = async () => {
      const resp = await get("templates");

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
    const resp = await save({
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
      title: t("common.text.operation_succeeded"),
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

      <div className="text-muted-foreground text-sm mt-1">{t("settings.notification.template.variables.tips.title")}</div>

      <Textarea
        className="mt-2"
        value={templates[0].content}
        onChange={(e) => {
          handleContentChange(e.target.value);
        }}
      ></Textarea>
      <div className="text-muted-foreground text-sm mt-1">{t("settings.notification.template.variables.tips.content")}</div>
      <div className="flex justify-end mt-2">
        <Button onClick={handleSaveClick}>{t("common.button.save")}</Button>
      </div>
    </div>
  );
};

export default NotifyTemplate;
