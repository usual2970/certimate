import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate, useSearchParams } from "react-router-dom";
import { DeleteOutlined as DeleteOutlinedIcon, ReloadOutlined as ReloadOutlinedIcon, SelectOutlined as SelectOutlinedIcon } from "@ant-design/icons";
import { PageHeader } from "@ant-design/pro-components";
import { useRequest } from "ahooks";
import {
  Button,
  Card,
  Divider,
  Empty,
  Flex,
  Input,
  Menu,
  type MenuProps,
  Modal,
  Radio,
  Space,
  Table,
  type TableProps,
  Tooltip,
  Typography,
  notification,
  theme,
} from "antd";
import dayjs from "dayjs";
import { ClientResponseError } from "pocketbase";

import CertificateDetailDrawer from "@/components/certificate/CertificateDetailDrawer";
import { CERTIFICATE_SOURCES, type CertificateModel } from "@/domain/certificate";
import { list as listCertificates, type ListRequest as listCertificatesRequest, remove as removeCertificate } from "@/repository/certificate";
import { getErrMsg } from "@/utils/error";

const CertificateList = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();

  const { t } = useTranslation();

  const { token: themeToken } = theme.useToken();

  const [modalApi, ModalContextHolder] = Modal.useModal();
  const [notificationApi, NotificationContextHolder] = notification.useNotification();

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
      title: t("certificate.props.subject_alt_names"),
      render: (_, record) => <Typography.Text>{record.subjectAltNames}</Typography.Text>,
    },
    {
      key: "expiry",
      title: t("certificate.props.validity"),
      ellipsis: true,
      defaultFilteredValue: searchParams.has("state") ? [searchParams.get("state") as string] : undefined,
      filterDropdown: ({ setSelectedKeys, confirm, clearFilters }) => {
        const items: Required<MenuProps>["items"] = [
          ["expireSoon", "certificate.props.validity.filter.expire_soon"],
          ["expired", "certificate.props.validity.filter.expired"],
        ].map(([key, label]) => {
          return {
            key,
            label: <Radio checked={filters["state"] === key}>{t(label)}</Radio>,
            onClick: () => {
              if (filters["state"] !== key) {
                setPage(1);
                setFilters((prev) => ({ ...prev, state: key }));
                setSelectedKeys([key]);
              }

              confirm({ closeDropdown: true });
            },
          };
        });

        const handleResetClick = () => {
          setPage(1);
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
            <Space className="w-full justify-end" style={{ padding: themeToken.paddingSM }}>
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
      render: (_, record) => {
        const total = dayjs(record.expireAt).diff(dayjs(record.created), "d") + 1;
        const left = dayjs(record.expireAt).diff(dayjs(), "d");
        return (
          <Space className="max-w-full" direction="vertical" size={4}>
            {left > 0 ? (
              <Typography.Text type="success">{t("certificate.props.validity.left_days", { left, total })}</Typography.Text>
            ) : (
              <Typography.Text type="danger">{t("certificate.props.validity.expired")}</Typography.Text>
            )}

            <Typography.Text type="secondary">
              {t("certificate.props.validity.expiration", { date: dayjs(record.expireAt).format("YYYY-MM-DD") })}
            </Typography.Text>
          </Space>
        );
      },
    },
    {
      key: "issuer",
      title: t("certificate.props.brand"),
      render: (_, record) => (
        <Space className="max-w-full" direction="vertical" size={4}>
          <Typography.Text>{record.issuer}</Typography.Text>
          <Typography.Text>{record.keyAlgorithm}</Typography.Text>
        </Space>
      ),
    },
    {
      key: "source",
      title: t("certificate.props.source"),
      ellipsis: true,
      render: (_, record) => {
        if (record.source === CERTIFICATE_SOURCES.WORKFLOW) {
          const workflowId = record.workflowId;
          return (
            <Space className="max-w-full" direction="vertical" size={4}>
              <Typography.Text>{t("certificate.props.source.workflow")}</Typography.Text>
              <Typography.Link
                type="secondary"
                ellipsis
                onClick={() => {
                  if (workflowId) {
                    navigate(`/workflows/${workflowId}`);
                  }
                }}
              >
                {record.expand?.workflowId?.name ?? <span className="font-mono">{t(`#${workflowId}`)}</span>}
              </Typography.Link>
            </Space>
          );
        } else if (record.source === CERTIFICATE_SOURCES.UPLOAD) {
          return <Typography.Text>{t("certificate.props.source.upload")}</Typography.Text>;
        }

        return <></>;
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
        <Space.Compact>
          <CertificateDetailDrawer
            data={record}
            trigger={
              <Tooltip title={t("certificate.action.view")}>
                <Button color="primary" icon={<SelectOutlinedIcon />} variant="text" />
              </Tooltip>
            }
          />

          <Tooltip title={t("certificate.action.delete")}>
            <Button color="danger" icon={<DeleteOutlinedIcon />} variant="text" onClick={() => handleDeleteClick(record)} />
          </Tooltip>
        </Space.Compact>
      ),
    },
  ];
  const [tableData, setTableData] = useState<CertificateModel[]>([]);
  const [tableTotal, setTableTotal] = useState<number>(0);

  const [filters, setFilters] = useState<Record<string, unknown>>(() => {
    return {
      keyword: searchParams.get("keyword"),
      state: searchParams.get("state"),
    };
  });

  const [page, setPage] = useState<number>(() => parseInt(+searchParams.get("page")! + "") || 1);
  const [pageSize, setPageSize] = useState<number>(() => parseInt(+searchParams.get("perPage")! + "") || 10);

  const {
    loading,
    error: loadedError,
    run: refreshData,
  } = useRequest(
    () => {
      return listCertificates({
        keyword: filters["keyword"] as string,
        state: filters["state"] as listCertificatesRequest["state"],
        page: page,
        perPage: pageSize,
      });
    },
    {
      refreshDeps: [filters, page, pageSize],
      onSuccess: (res) => {
        setTableData(res.items);
        setTableTotal(res.totalItems);
      },
      onError: (err) => {
        if (err instanceof ClientResponseError && err.isAbort) {
          return;
        }

        console.error(err);
        notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });

        throw err;
      },
    }
  );

  const handleSearch = (value: string) => {
    setFilters((prev) => ({ ...prev, keyword: value.trim() }));
    setPage(1);
  };

  const handleReloadClick = () => {
    if (loading) return;

    refreshData();
  };

  const handleDeleteClick = (certificate: CertificateModel) => {
    modalApi.confirm({
      title: t("certificate.action.delete"),
      content: t("certificate.action.delete.confirm"),
      onOk: async () => {
        try {
          const resp = await removeCertificate(certificate);
          if (resp) {
            setTableData((prev) => prev.filter((item) => item.id !== certificate.id));
            refreshData();
          }
        } catch (err) {
          console.error(err);
          notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
        }
      },
    });
  };

  return (
    <div className="p-4">
      {ModalContextHolder}
      {NotificationContextHolder}

      <PageHeader title={t("certificate.page.title")} />

      <Card size="small">
        <div className="mb-4">
          <Flex gap="small">
            <div className="flex-1">
              <Input.Search allowClear defaultValue={filters["keyword"] as string} placeholder={t("certificate.search.placeholder")} onSearch={handleSearch} />
            </div>
            <div>
              <Button icon={<ReloadOutlinedIcon spin={loading} />} onClick={handleReloadClick} />
            </div>
          </Flex>
        </div>

        <Table<CertificateModel>
          columns={tableColumns}
          dataSource={tableData}
          loading={loading}
          locale={{
            emptyText: <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={getErrMsg(loadedError ?? t("certificate.nodata"))} />,
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
          rowKey={(record) => record.id}
          scroll={{ x: "max(100%, 960px)" }}
        />
      </Card>
    </div>
  );
};

export default CertificateList;
