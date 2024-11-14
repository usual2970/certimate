import { useTranslation } from "react-i18next";

import { Accordion, AccordionContent, AccordionItem, AccordionTrigger } from "@/components/ui/accordion";
import DingTalk from "@/components/notify/DingTalk";
import Lark from "@/components/notify/Lark";
import NotifyTemplate from "@/components/notify/NotifyTemplate";
import Telegram from "@/components/notify/Telegram";
import Webhook from "@/components/notify/Webhook";
import ServerChan from "@/components/notify/ServerChan";
import Email from "@/components/notify/Email";
import Bark from "@/components/notify/Bark";
import { NotifyProvider } from "@/providers/notify";

const Notify = () => {
  const { t } = useTranslation();

  return (
    <>
      <NotifyProvider>
        <div className="border rounded-sm p-5 shadow-lg">
          <Accordion type={"multiple"} className="dark:text-stone-200">
            <AccordionItem value="item-1" className="dark:border-stone-200">
              <AccordionTrigger>{t("settings.notification.template.label")}</AccordionTrigger>
              <AccordionContent>
                <NotifyTemplate />
              </AccordionContent>
            </AccordionItem>
          </Accordion>
        </div>

        <div className="border rounded-md p-5 mt-7 shadow-lg">
          <Accordion type={"single"} collapsible={true} className="dark:text-stone-200">
            <AccordionItem value="item-email" className="dark:border-stone-200">
              <AccordionTrigger>{t("common.provider.email")}</AccordionTrigger>
              <AccordionContent>
                <Email />
              </AccordionContent>
            </AccordionItem>

            <AccordionItem value="item-webhook" className="dark:border-stone-200">
              <AccordionTrigger>{t("common.provider.webhook")}</AccordionTrigger>
              <AccordionContent>
                <Webhook />
              </AccordionContent>
            </AccordionItem>

            <AccordionItem value="item-dingtalk" className="dark:border-stone-200">
              <AccordionTrigger>{t("common.provider.dingtalk")}</AccordionTrigger>
              <AccordionContent>
                <DingTalk />
              </AccordionContent>
            </AccordionItem>

            <AccordionItem value="item-lark" className="dark:border-stone-200">
              <AccordionTrigger>{t("common.provider.lark")}</AccordionTrigger>
              <AccordionContent>
                <Lark />
              </AccordionContent>
            </AccordionItem>

            <AccordionItem value="item-telegram" className="dark:border-stone-200">
              <AccordionTrigger>{t("common.provider.telegram")}</AccordionTrigger>
              <AccordionContent>
                <Telegram />
              </AccordionContent>
            </AccordionItem>

            <AccordionItem value="item-serverchan" className="dark:border-stone-200">
              <AccordionTrigger>{t("common.provider.serverchan")}</AccordionTrigger>
              <AccordionContent>
                <ServerChan />
              </AccordionContent>
            </AccordionItem>

            <AccordionItem value="item-bark" className="dark:border-stone-200">
              <AccordionTrigger>{t("common.provider.bark")}</AccordionTrigger>
              <AccordionContent>
                <Bark />
              </AccordionContent>
            </AccordionItem>
          </Accordion>
        </div>
      </NotifyProvider>
    </>
  );
};

export default Notify;
