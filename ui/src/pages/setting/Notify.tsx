import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import React from "react";

const Notify = () => {
  return (
    <div>
      <Accordion type={"multiple"} className="dark:text-stone-200">
        <AccordionItem value="item-1" className="dark:border-stone-200">
          <AccordionTrigger>邮箱</AccordionTrigger>
          <AccordionContent>
            Yes. It adheres to the WAI-ARIA design pattern.
          </AccordionContent>
        </AccordionItem>

        <AccordionItem value="item-2" className="dark:border-stone-200">
          <AccordionTrigger>钉钉</AccordionTrigger>
          <AccordionContent>
            Yes. It adheres to the WAI-ARIA design pattern.
          </AccordionContent>
        </AccordionItem>

        <AccordionItem value="item-3" className="dark:border-stone-200">
          <AccordionTrigger>微信公众号</AccordionTrigger>
          <AccordionContent>
            Yes. It adheres to the WAI-ARIA design pattern.
          </AccordionContent>
        </AccordionItem>

        <AccordionItem value="item-4" className="dark:border-stone-200">
          <AccordionTrigger>Telegram</AccordionTrigger>
          <AccordionContent>
            Yes. It adheres to the WAI-ARIA design pattern.
          </AccordionContent>
        </AccordionItem>

        <AccordionItem value="item-5" className="dark:border-stone-200">
          <AccordionTrigger>Webhook</AccordionTrigger>
          <AccordionContent>
            Yes. It adheres to the WAI-ARIA design pattern.
          </AccordionContent>
        </AccordionItem>
      </Accordion>
    </div>
  );
};

export default Notify;
