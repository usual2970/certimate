import {
  Accordion,
  AccordionContent,
  AccordionItem,
  AccordionTrigger,
} from "@/components/ui/accordion";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import React from "react";

const Notify = () => {
  return (
    <>
      <div className="border rounded-sm p-5">
        <Accordion type={"multiple"} className="dark:text-stone-200">
          <AccordionItem value="item-1" className="dark:border-stone-200">
            <AccordionTrigger>模板</AccordionTrigger>
            <AccordionContent>
              <Input value="您有证书即将过期" />
              <Textarea
                className="mt-2"
                value={
                  "有{COUNT}张证书即将过期,域名分别为{DOMAINS},请保持关注！"
                }
              ></Textarea>
            </AccordionContent>
          </AccordionItem>
        </Accordion>
      </div>
      <div className="border rounded-md p-5 mt-7">
        <Accordion type={"multiple"} className="dark:text-stone-200">
          <AccordionItem value="item-2" className="dark:border-stone-200">
            <AccordionTrigger>钉钉</AccordionTrigger>
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
    </>
  );
};

export default Notify;
