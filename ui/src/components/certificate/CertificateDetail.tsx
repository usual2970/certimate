import { useTranslation } from "react-i18next";
import { Button, Form, Input, message, Tooltip } from "antd";
import { CopyToClipboard } from "react-copy-to-clipboard";
import { Clipboard as ClipboardIcon } from "lucide-react";

import { type Certificate } from "@/domain/certificate";
import { saveFiles2Zip } from "@/lib/file";

type CertificateDetailProps = {
  data: Certificate;
};

const CertificateDetail = ({ data }: CertificateDetailProps) => {
  const { t } = useTranslation();

  const [messageApi, MessageContextHolder] = message.useMessage();

  const handleDownloadClick = async () => {
    // TODO: 支持下载多种格式
    const zipName = `${data.id}-${data.san}.zip`;
    const files = [
      {
        name: `${data.san}.pem`,
        content: data.certificate ?? "",
      },
      {
        name: `${data.san}.key`,
        content: data.privateKey ?? "",
      },
    ];

    await saveFiles2Zip(zipName, files);
  };

  return (
    <div>
      {MessageContextHolder}

      <Form layout="vertical">
        <Form.Item>
          <div className="flex items-center justify-between w-full mb-2">
            <label className="font-medium">{t("certificate.props.certificate_chain")}</label>
            <Tooltip title={t("common.copy")}>
              <CopyToClipboard
                text={data.certificate}
                onCopy={() => {
                  messageApi.success(t("common.copy.done"));
                }}
              >
                <Button type="text" icon={<ClipboardIcon size={14} />}></Button>
              </CopyToClipboard>
            </Tooltip>
          </div>
          <Input.TextArea value={data.certificate} rows={10} autoSize={{ maxRows: 10 }} readOnly />
        </Form.Item>

        <Form.Item>
          <div className="flex items-center justify-between w-full mb-2">
            <label className="font-medium">{t("certificate.props.private_key")}</label>
            <Tooltip title={t("common.copy")}>
              <CopyToClipboard
                text={data.privateKey}
                onCopy={() => {
                  messageApi.success(t("common.copy.done"));
                }}
              >
                <Button type="text" icon={<ClipboardIcon size={14} />}></Button>
              </CopyToClipboard>
            </Tooltip>
          </div>
          <Input.TextArea value={data.privateKey} rows={10} autoSize={{ maxRows: 10 }} readOnly />
        </Form.Item>
      </Form>

      <div className="flex items-center justify-end">
        <Button
          type="primary"
          onClick={() => {
            handleDownloadClick();
          }}
        >
          {t("certificate.action.download")}
        </Button>
      </div>
    </div>
  );
};

export default CertificateDetail;
