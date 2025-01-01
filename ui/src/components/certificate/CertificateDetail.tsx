import { CopyToClipboard } from "react-copy-to-clipboard";
import { useTranslation } from "react-i18next";
import { CopyOutlined as CopyOutlinedIcon, DownOutlined as DownOutlinedIcon, LikeOutlined as LikeOutlinedIcon } from "@ant-design/icons";
import { Button, Dropdown, Form, Input, message, Space, Tooltip } from "antd";
import dayjs from "dayjs";

import { type CertificateModel } from "@/domain/certificate";
import { saveFiles2Zip } from "@/utils/file";

export type CertificateDetailProps = {
  className?: string;
  style?: React.CSSProperties;
  data: CertificateModel;
};

const CertificateDetail = ({ data, ...props }: CertificateDetailProps) => {
  const { t } = useTranslation();

  const [messageApi, MessageContextHolder] = message.useMessage();

  const handleDownloadPEMClick = async () => {
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
    <div {...props}>
      {MessageContextHolder}

      <Form layout="vertical">
        <Form.Item label={t("certificate.props.san")}>
          <Input value={data.san} placeholder="" />
        </Form.Item>

        <Form.Item label={t("certificate.props.expiry")}>
          <Input value={dayjs(data.expireAt).format("YYYY-MM-DD HH:mm:ss")} placeholder="" />
        </Form.Item>

        <Form.Item>
          <div className="flex items-center justify-between w-full mb-2">
            <label>{t("certificate.props.certificate_chain")}</label>
            <Tooltip title={t("common.button.copy")}>
              <CopyToClipboard
                text={data.certificate}
                onCopy={() => {
                  messageApi.success(t("common.text.copied"));
                }}
              >
                <Button size="small" type="text" icon={<CopyOutlinedIcon />}></Button>
              </CopyToClipboard>
            </Tooltip>
          </div>
          <Input.TextArea value={data.certificate} rows={10} autoSize={{ maxRows: 10 }} readOnly />
        </Form.Item>

        <Form.Item>
          <div className="flex items-center justify-between w-full mb-2">
            <label>{t("certificate.props.private_key")}</label>
            <Tooltip title={t("common.button.copy")}>
              <CopyToClipboard
                text={data.privateKey}
                onCopy={() => {
                  messageApi.success(t("common.text.copied"));
                }}
              >
                <Button size="small" type="text" icon={<CopyOutlinedIcon />}></Button>
              </CopyToClipboard>
            </Tooltip>
          </div>
          <Input.TextArea value={data.privateKey} rows={10} autoSize={{ maxRows: 10 }} readOnly />
        </Form.Item>
      </Form>

      <div className="flex items-center justify-end">
        <Dropdown
          menu={{
            items: [
              {
                key: "PEM",
                label: "PEM",
                extra: <LikeOutlinedIcon />,
                onClick: () => handleDownloadPEMClick(),
              },
              {
                key: "PFX",
                label: "PFX",
                onClick: () => {
                  // TODO: 下载 PFX 格式证书
                  alert("TODO");
                },
              },
              {
                key: "JKS",
                label: "JKS",
                onClick: () => {
                  // TODO: 下载 JKS 格式证书
                  alert("TODO");
                },
              },
            ],
          }}
        >
          <Button type="primary">
            <Space>
              <span>{t("certificate.action.download")}</span>
              <DownOutlinedIcon />
            </Space>
          </Button>
        </Dropdown>
      </div>
    </div>
  );
};

export default CertificateDetail;
