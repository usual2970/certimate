import { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Avatar, Button, Empty, Modal, notification, Space, Table, Tooltip, Typography, type TableProps } from "antd";
import { PageHeader } from "@ant-design/pro-components";
import { Copy as CopyIcon, Pencil as PencilIcon, Plus as PlusIcon, Trash2 as Trash2Icon } from "lucide-react";
import dayjs from "dayjs";
import { ClientResponseError } from "pocketbase";

import AccessEditDialog from "@/components/certimate/AccessEditDialog";
import { accessProvidersMap, type AccessModel } from "@/domain/access";
import { useAccessStore } from "@/stores/access";

const AccessList = () => {
  const { t } = useTranslation();

  const [modalApi, ModelContextHolder] = Modal.useModal();
  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const { accesses, fetchAccesses, deleteAccess } = useAccessStore();

  const [loading, setLoading] = useState<boolean>(false);

  const tableColumns: TableProps<AccessModel>["columns"] = [
    {
      key: "$index",
      align: "center",
      fixed: "left",
      width: 50,
      render: (_, __, index) => (page - 1) * pageSize + index + 1,
    },
    {
      key: "name",
      title: t("access.props.name"),
      ellipsis: true,
      render: (_, record) => <>{record.name}</>,
    },
    {
      key: "provider",
      title: t("access.props.provider"),
      ellipsis: true,
      render: (_, record) => {
        return (
          <Space className="max-w-full truncate" align="center" size={4}>
            <Avatar src={accessProvidersMap.get(record.configType)?.icon} size="small" />
            <Typography.Text ellipsis>{t(accessProvidersMap.get(record.configType)?.name ?? "")}</Typography.Text>
          </Space>
        );
      },
    },
    {
      key: "createdAt",
      title: t("access.props.created_at"),
      ellipsis: true,
      render: (_, record) => {
        return dayjs(record.created!).format("YYYY-MM-DD HH:mm:ss");
      },
    },
    {
      key: "updatedAt",
      title: t("access.props.updated_at"),
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
        <>
          <Space size={0}>
            <AccessEditDialog
              trigger={
                <Tooltip title={t("access.action.edit")}>
                  <Button type="link" icon={<PencilIcon size={16} />} />
                </Tooltip>
              }
              op="edit"
              data={record}
            />

            <AccessEditDialog
              trigger={
                <Tooltip title={t("access.action.copy")}>
                  <Button type="link" icon={<CopyIcon size={16} />} />
                </Tooltip>
              }
              op="copy"
              data={record}
            />

            <Tooltip title={t("access.action.delete")}>
              <Button
                type="link"
                danger={true}
                icon={<Trash2Icon size={16} />}
                onClick={() => {
                  handleDeleteClick(record);
                }}
              />
            </Tooltip>
          </Space>
        </>
      ),
    },
  ];
  const [tableData, setTableData] = useState<AccessModel[]>([]);
  const [tableTotal, setTableTotal] = useState<number>(0);

  const [page, setPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);

  useEffect(() => {
    fetchAccesses();
  }, []);

  const fetchTableData = useCallback(async () => {
    if (loading) return;
    setLoading(true);

    try {
      const startIndex = (page - 1) * pageSize;
      const endIndex = startIndex + pageSize;
      const items = accesses.slice(startIndex, endIndex);

      setTableData(items);
      setTableTotal(accesses.length);
    } catch (err) {
      if (err instanceof ClientResponseError && err.isAbort) {
        return;
      }

      console.error(err);
      notificationApi.error({ message: t("common.text.request_error"), description: <>{String(err)}</> });
    } finally {
      setLoading(false);
    }
  }, [page, pageSize, accesses]);

  useEffect(() => {
    fetchTableData();
  }, [fetchTableData]);

  const handleDeleteClick = async (data: AccessModel) => {
    modalApi.confirm({
      title: t("access.action.delete"),
      content: t("access.action.delete.confirm"),
      onOk: async () => {
        // TODO: 有关联数据的不允许被删除
        try {
          await deleteAccess(data);
        } catch (err) {
          console.error(err);
          notificationApi.error({ message: t("common.text.request_error"), description: <>{String(err)}</> });
        }
      },
    });
  };

  return (
    <>
      {ModelContextHolder}
      {NotificationContextHolder}

      <PageHeader
        title={t("access.page.title")}
        extra={[
          <AccessEditDialog
            key="create"
            trigger={
              <Button key="create" type="primary" icon={<PlusIcon size={16} />}>
                {t("access.action.add")}
              </Button>
            }
            op="add"
          />,
        ]}
      />

      <Table<AccessModel>
        columns={tableColumns}
        dataSource={tableData}
        loading={loading}
        locale={{
          emptyText: <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={t("access.nodata")} />,
        }}
        pagination={{
          current: page,
          pageSize: pageSize,
          total: tableTotal,
          onChange: (page: number, pageSize: number) => {
            setPage(page);
            setPageSize(pageSize);
          },
          onShowSizeChange: (page: number, pageSize: number) => {
            setPage(page);
            setPageSize(pageSize);
          },
        }}
        rowKey={(record: AccessModel) => record.id}
        scroll={{ x: "max(100%, 960px)" }}
      />
    </>
  );
};

export default AccessList;
