import { useTranslation } from "react-i18next";
import { BookOpen } from "lucide-react";

import { cn } from "@/components/ui/utils";
import { Separator } from "@/components/ui/separator";
import { version } from "@/domain/version";

type VersionProps = {
  className?: string;
};

const Version = ({ className }: VersionProps) => {
  const { t } = useTranslation();

  return (
    <div className={cn("w-full flex pb-5 ", className)}>
      <div className="text-muted-foreground text-sm hover:text-stone-900 dark:hover:text-stone-200 flex">
        <a href="https://docs.certimate.me" target="_blank" className="flex items-center">
          <BookOpen size={16} />
          <div className="ml-1">{t("common.menu.document")}</div>
        </a>
        <Separator orientation="vertical" className="mx-2" />
        <a href="https://github.com/usual2970/certimate/releases" target="_blank">
          {version}
        </a>
      </div>
    </div>
  );
};

export default Version;
