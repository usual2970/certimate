import { Sheet, SheetContent, SheetHeader, SheetTitle } from "../ui/sheet";

import { Certificate } from "@/domain/certificate";
import { Textarea } from "../ui/textarea";
import { Button } from "../ui/button";
import { Label } from "../ui/label";
import { CustomFile, saveFiles2ZIP } from "@/lib/file";
import { useTranslation } from "react-i18next";

type WorkflowLogDetailProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  certificate?: Certificate;
};
const CertificateDetail = ({ open, onOpenChange, certificate }: WorkflowLogDetailProps) => {
  const { t } = useTranslation();
  const handleDownloadClick = async () => {
    const zipName = `${certificate?.id}-${certificate?.san}.zip`;
    const files: CustomFile[] = [
      {
        name: `${certificate?.san}.pem`,
        content: certificate?.certificate ? certificate?.certificate : "",
      },
      {
        name: `${certificate?.san}.key`,
        content: certificate?.privateKey ? certificate?.privateKey : "",
      },
    ];

    await saveFiles2ZIP(zipName, files);
  };

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent className="sm:max-w-2xl dark:text-stone-200">
        <SheetHeader>
          <SheetTitle></SheetTitle>
        </SheetHeader>

        <div className="flex flex-col space-y-5 mt-9">
          <div className="flex justify-end">
            <Button
              size={"sm"}
              onClick={() => {
                handleDownloadClick();
              }}
            >
              {t("certificate.action.download")}
            </Button>
          </div>
          <div className="flex flex-col space-y-3">
            <Label>{t("certificate.props.certificate")}</Label>
            <Textarea value={certificate?.certificate} rows={10} readOnly={true} />
          </div>
          <div className="flex flex-col space-y-3">
            <Label>{t("certificate.props.private_key")}</Label>
            <Textarea value={certificate?.privateKey} rows={10} readOnly={true} />
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
};

export default CertificateDetail;
