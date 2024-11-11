import { useTranslation } from "react-i18next";

const End = () => {
  const { t } = useTranslation();
  return (
    <div className="flex flex-col items-center">
      <div className="h-[18px] rounded-full w-[18px] bg-stone-400"></div>
      <div className="text-sm text-stone-400 mt-2">{t("workflow.node.end.title")}</div>
    </div>
  );
};

export default End;
