import { Sheet, SheetContent, SheetHeader, SheetTitle } from "../ui/sheet";

import { Certificate } from "@/domain/certificate";
import { Textarea } from "../ui/textarea";
import { Button } from "../ui/button";
import { Label } from "../ui/label";
import { CustomFile, saveFiles2ZIP } from "@/lib/file";

type WorkflowLogDetailProps = {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  certificate?: Certificate;
};
const CertificateDetail = ({ open, onOpenChange, certificate }: WorkflowLogDetailProps) => {
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
      <SheetContent className="sm:max-w-2xl">
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
              下载证书
            </Button>
          </div>
          <div className="flex flex-col space-y-3">
            <Label>证书</Label>
            <Textarea value={certificate?.certificate} rows={10} readOnly={true} />
          </div>
          <div className="flex flex-col space-y-3">
            <Label>密钥</Label>
            <Textarea value={certificate?.privateKey} rows={10} readOnly={true} />
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
};

export default CertificateDetail;
