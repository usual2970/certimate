import { useCallback, useEffect, useState } from "react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Button, Divider, Empty, Menu, notification, Radio, Space, Table, theme, Tooltip, Typography, type MenuProps, type TableProps } from "antd";
import { PageHeader } from "@ant-design/pro-components";
import { Eye as EyeIcon, Filter as FilterIcon } from "lucide-react";
import dayjs from "dayjs";
import { ClientResponseError } from "pocketbase";

import CertificateDetailDrawer from "@/components/certificate/CertificateDetailDrawer";
import { CertificateModel } from "@/domain/certificate";
import { list as listCertificate, type CertificateListReq } from "@/repository/certificate";
import { getErrMsg } from "@/utils/error";

const CertificateList = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const { t } = useTranslation();

  const { token: themeToken } = theme.useToken();

  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const [loading, setLoading] = useState<boolean>(false);

  const tableColumns: TableProps<CertificateModel>["columns"] = [
    {
      key: "$index",
      align: "center",
      fixed: "left",
      width: 50,
      render: (_, __, index) => (page - 1) * pageSize + index + 1,
    },
    {
      key: "name",
      title: t("certificate.props.san"),
      render: (_, record) => <Typography.Text>{record.san}</Typography.Text>,
    },
    {
      key: "expiry",
      title: t("certificate.props.expiry"),
      ellipsis: true,
      defaultFilteredValue: searchParams.has("state") ? [searchParams.get("state") as string] : undefined,
      filterDropdown: ({ setSelectedKeys, confirm, clearFilters }) => {
        const items: Required<MenuProps>["items"] = [
          ["expireSoon", "certificate.props.expiry.filter.expire_soon"],
          ["expired", "certificate.props.expiry.filter.expired"],
        ].map(([key, label]) => {
          return {
            key,
            label: <Radio checked={filters["state"] === key}>{t(label)}</Radio>,
            onClick: () => {
              if (filters["state"] !== key) {
                setFilters((prev) => ({ ...prev, state: key }));
                setSelectedKeys([key]);
              }

              confirm({ closeDropdown: true });
            },
          };
        });

        const handleResetClick = () => {
          setFilters((prev) => ({ ...prev, state: undefined }));
          setSelectedKeys([]);
          clearFilters?.();
          confirm();
        };

        const handleConfirmClick = () => {
          confirm();
        };

        return (
          <div style={{ padding: 0 }}>
            <Menu items={items} selectable={false} />
            <Divider style={{ margin: 0 }} />
            <Space className="justify-end w-full" style={{ padding: themeToken.paddingSM }}>
              <Button size="small" disabled={!filters.state} onClick={handleResetClick}>
                {t("common.button.reset")}
              </Button>
              <Button type="primary" size="small" onClick={handleConfirmClick}>
                {t("common.button.ok")}
              </Button>
            </Space>
          </div>
        );
      },
      filterIcon: () => <FilterIcon size={14} />,
      render: (_, record) => {
        const total = dayjs(record.expireAt).diff(dayjs(record.created), "d") + 1;
        const left = dayjs(record.expireAt).diff(dayjs(), "d");
        return (
          <Space className="max-w-full" direction="vertical" size={4}>
            {left > 0 ? (
              <Typography.Text type="success">{t("certificate.props.expiry.left_days", { left, total })}</Typography.Text>
            ) : (
              <Typography.Text type="danger">{t("certificate.props.expiry.expired")}</Typography.Text>
            )}

            <Typography.Text type="secondary">
              {t("certificate.props.expiry.expiration", { date: dayjs(record.expireAt).format("YYYY-MM-DD") })}
            </Typography.Text>
          </Space>
        );
      },
    },
    {
      key: "source",
      title: t("certificate.props.source"),
      ellipsis: true,
      render: (_, record) => {
        const workflowId = record.workflow;
        return workflowId ? (
          <Space className="max-w-full" direction="vertical" size={4}>
            <Typography.Text>{t("certificate.props.source.workflow")}</Typography.Text>
            <Typography.Link
              type="secondary"
              ellipsis
              onClick={() => {
                navigate(`/workflows/detail?id=${workflowId}`);
              }}
            >
              {record.expand?.workflow?.name ?? ""}
            </Typography.Link>
          </Space>
        ) : (
          <>TODO: 支持手动上传</>
        );
      },
    },
    {
      key: "createdAt",
      title: t("certificate.props.created_at"),
      ellipsis: true,
      render: (_, record) => {
        return dayjs(record.created!).format("YYYY-MM-DD HH:mm:ss");
      },
    },
    {
      key: "updatedAt",
      title: t("certificate.props.updated_at"),
      ellipsis: true,
      render: (_, record) => {
        return dayjs(record.updated!).format("YYYY-MM-DD HH:mm:ss");
      },
    },
    {
      key: "$action",
      align: "end",
      fixed: "right",
      width: 120,
      render: (_, record) => (
        <Space size={0}>
          <CertificateDetailDrawer
            data={record}
            trigger={
              <Tooltip title={t("certificate.action.view")}>
                <Button type="link" icon={<EyeIcon size={16} />} />
              </Tooltip>
            }
          />
        </Space>
      ),
    },
  ];
  const [tableData, setTableData] = useState<CertificateModel[]>([]);
  const [tableTotal, setTableTotal] = useState<number>(0);

  const [filters, setFilters] = useState<Record<string, unknown>>(() => {
    return {
      state: searchParams.get("state"),
    };
  });

  const [page, setPage] = useState<number>(() => parseInt(+searchParams.get("page")! + "") || 1);
  const [pageSize, setPageSize] = useState<number>(() => parseInt(+searchParams.get("perPage")! + "") || 10);

  const fetchTableData = useCallback(async () => {
    if (loading) return;
    setLoading(true);

    try {
      const resp = await listCertificate({
        page: page,
        perPage: pageSize,
        state: filters["state"] as CertificateListReq["state"],
      });

      setTableData(resp.items);
      setTableTotal(resp.totalItems);
    } catch (err) {
      if (err instanceof ClientResponseError && err.isAbort) {
        return;
      }

      console.error(err);
      notificationApi.error({ message: t("common.text.request_error"), description: <>{getErrMsg(err)}</> });
    } finally {
      setLoading(false);
    }
  }, [filters, page, pageSize]);

  useEffect(() => {
    fetchTableData();
  }, [fetchTableData]);

  return (
    <>
      {NotificationContextHolder}

      <PageHeader title={t("certificate.page.title")} />

      <Table<CertificateModel>
        columns={tableColumns}
        dataSource={tableData}
        loading={loading}
        locale={{
          emptyText: <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={t("certificate.nodata")} />,
        }}
        pagination={{
          current: page,
          pageSize: pageSize,
          total: tableTotal,
          showSizeChanger: true,
          onChange: (page: number, pageSize: number) => {
            setPage(page);
            setPageSize(pageSize);
          },
          onShowSizeChange: (page: number, pageSize: number) => {
            setPage(page);
            setPageSize(pageSize);
          },
        }}
        rowKey={(record: CertificateModel) => record.id}
        scroll={{ x: "max(100%, 960px)" }}
      />
    </>
  );
};

export default CertificateList;
