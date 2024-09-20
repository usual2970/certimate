import { BookOpen } from "lucide-react";

import { Separator } from "../ui/separator";
import { version } from "@/domain/version";

const Version = () => {
  return (
    <div className="fixed right-0 bottom-0 w-full flex justify-between p-5">
      <div className=""></div>
      <div className="text-muted-foreground text-sm hover:text-stone-900 dark:hover:text-stone-200 flex">
        <a
          href="https://docs.certimate.me"
          target="_blank"
          className="flex items-center"
        >
          <BookOpen size={16} />
          <div className="ml-1">文档</div>
        </a>
        <Separator orientation="vertical" className="mx-2" />
        <a
          href="https://github.com/usual2970/certimate/releases"
          target="_blank"
        >
          {version}
        </a>
      </div>
    </div>
  );
};

export default Version;
