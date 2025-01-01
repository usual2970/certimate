import { useTranslation } from "react-i18next";

const End = () => {
  const { t } = useTranslation();
  return (
    <div className="flex flex-col items-center">
      <div className="size-[20px] rounded-full bg-stone-400"></div>
      <div className="text-sm text-stone-400 mt-2">{t("workflow_node.end.title")}</div>
    </div>
  );
};

export default End;
