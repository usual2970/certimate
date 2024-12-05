import { useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Button, Space, Table, Tooltip, Typography, type TableProps } from "antd";
import { PageHeader } from "@ant-design/pro-components";
import { Eye as EyeIcon } from "lucide-react";

import { Certificate as CertificateType } from "@/domain/certificate";
import { list as listCertificate, type CertificateListReq } from "@/repository/certificate";
import { diffDays, getLeftDays } from "@/lib/time";

const CertificateList = () => {
  const { t } = useTranslation();

  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const [loading, setLoading] = useState<boolean>(false);

  const tableColumns: TableProps<CertificateType>["columns"] = [
    {
      key: "$index",
      align: "center",
      title: "",
      width: 50,
      render: (_, __, index) => (page - 1) * pageSize + index + 1,
    },
    {
      key: "name",
      title: t("common.text.domain"),
      render: (_, record) => <Typography.Text>{record.san}</Typography.Text>,
    },
    {
      key: "expiry",
      title: t("certificate.props.expiry"),
      render: (_, record) => {
        const leftDays = getLeftDays(record.expireAt);
        const allDays = diffDays(record.expireAt, record.created);
        return (
          <Space className="max-w-full" direction="vertical" size={4}>
            {leftDays > 0 ? (
              <Typography.Text type="success">
                {leftDays} / {allDays} {t("certificate.props.expiry.days")}
              </Typography.Text>
            ) : (
              <Typography.Text type="danger">{t("certificate.props.expiry.expired")}</Typography.Text>
            )}

            <Typography.Text type="secondary">
              {new Date(record.expireAt).toLocaleString().split(" ")[0]} {t("certificate.props.expiry.text.expire")}
            </Typography.Text>
          </Space>
        );
      },
    },
    {
      key: "source",
      title: t("certificate.props.source"),
      render: (_, record) => {
        const workflowId = record.workflow;
        return workflowId ? (
          <Space className="max-w-full" direction="vertical" size={4}>
            <Typography.Text>{t("common.text.workflow")}</Typography.Text>
            <Typography.Link
              type="secondary"
              ellipsis
              onClick={() => {
                navigate(`/workflow/detail?id=${workflowId}`);
              }}
            >
              {record.expand?.workflow?.name ?? ""}
            </Typography.Link>
          </Space>
        ) : (
          <>TODO: 手动上传</>
        );
      },
    },
    {
      key: "createdAt",
      title: t("common.text.created_at"),
      ellipsis: true,
      render: (_, record) => {
        return new Date(record.created!).toLocaleString();
      },
    },
    {
      key: "updatedAt",
      title: t("common.text.updated_at"),
      ellipsis: true,
      render: (_, record) => {
        return new Date(record.updated!).toLocaleString();
      },
    },
    {
      key: "$operations",
      align: "end",
      width: 100,
      render: (_, record) => (
        <Space>
          <Tooltip title={t("common.view")}>
            <Button
              type="link"
              icon={<EyeIcon size={16} />}
              onClick={() => {
                // TODO: 查看证书详情
                alert("TODO");
              }}
            />
          </Tooltip>
        </Space>
      ),
    },
  ];
  const [tableData, setTableData] = useState<CertificateType[]>([]);
  const [tableTotal, setTableTotal] = useState<number>(0);

  const [page, setPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);

  const fetchTableData = async () => {
    if (loading) return;
    setLoading(true);

    const state = searchParams.get("state");
    const req: CertificateListReq = { page: page, perPage: pageSize };
    if (state) {
      req.state = state as CertificateListReq["state"];
    }

    try {
      const resp = await listCertificate(req);

      setTableData(resp.items);
      setTableTotal(resp.totalItems);
    } catch (err) {
      console.error(err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTableData();
  }, [page, pageSize]);

  // TODO: Empty 样式

  return (
    <>
      <PageHeader title={t("certificate.page.title")} />

      <Table<CertificateType>
        columns={tableColumns}
        dataSource={tableData}
        rowKey={(record) => record.id}
        loading={loading}
        pagination={{
          current: page,
          pageSize: pageSize,
          total: tableTotal,
          onChange: (page, pageSize) => {
            setPage(page);
            setPageSize(pageSize);
          },
          onShowSizeChange: (page, pageSize) => {
            setPage(page);
            setPageSize(pageSize);
          },
        }}
      />
    </>
  );
};

export default CertificateList;
