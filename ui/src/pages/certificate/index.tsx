import CertificateList from "@/components/certificate/CertificateList";
import { useTranslation } from "react-i18next";

const Certificate = () => {
  const { t } = useTranslation();
  return (
    <div className="flex flex-col space-y-5">
      <div className="text-muted-foreground">{t("certificate.page.title")}</div>

      <CertificateList withPagination={true} />
    </div>
  );
};

export default Certificate;
