import DingTalk from "@/components/notify/DingTalk";
import NotifyTemplate from "@/components/notify/NotifyTemplate";
import Telegram from "@/components/notify/Telegram";
import Webhook from "@/components/notify/Webhook";
import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { NotifyProvider } from "@/providers/notify";

const Notify = () => {
  return (
    <>
      <NotifyProvider>
        <div className="border rounded-sm p-5 shadow-lg">
          <Accordion type={"multiple"} className="dark:text-stone-200">
            <AccordionItem value="item-1" className="dark:border-stone-200">
              <AccordionTrigger>模板</AccordionTrigger>
              <AccordionContent>
                <NotifyTemplate />
              </AccordionContent>
            </AccordionItem>
          </Accordion>
        </div>
        <div className="border rounded-md p-5 mt-7 shadow-lg">
          <Accordion type={"multiple"} className="dark:text-stone-200">
            <AccordionItem value="item-2" className="dark:border-stone-200">
              <AccordionTrigger>钉钉</AccordionTrigger>
              <AccordionContent>
                <DingTalk />
              </AccordionContent>
            </AccordionItem>

            <AccordionItem value="item-4" className="dark:border-stone-200">
              <AccordionTrigger>Telegram</AccordionTrigger>
              <AccordionContent>
                <Telegram />
              </AccordionContent>
            </AccordionItem>

            <AccordionItem value="item-5" className="dark:border-stone-200">
              <AccordionTrigger>Webhook</AccordionTrigger>
              <AccordionContent>
                <Webhook />
              </AccordionContent>
            </AccordionItem>
          </Accordion>
        </div>
      </NotifyProvider>
    </>
  );
};

export default Notify;
