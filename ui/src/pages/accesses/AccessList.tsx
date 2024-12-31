import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { useRequest } from "ahooks";
import { Avatar, Button, Empty, Modal, notification, Space, Table, Tooltip, Typography, type TableProps } from "antd";
import { PageHeader } from "@ant-design/pro-components";
import {
  DeleteOutlined as DeleteOutlinedIcon,
  EditOutlined as EditOutlinedIcon,
  PlusOutlined as PlusOutlinedIcon,
  SnippetsOutlined as SnippetsOutlinedIcon,
} from "@ant-design/icons";
import dayjs from "dayjs";
import { ClientResponseError } from "pocketbase";

import AccessEditModal from "@/components/access/AccessEditModal";
import { type AccessModel } from "@/domain/access";
import { accessProvidersMap } from "@/domain/provider";
import { useAccessStore } from "@/stores/access";
import { getErrMsg } from "@/utils/error";

const AccessList = () => {
  const { t } = useTranslation();

  const [modalApi, ModelContextHolder] = Modal.useModal();
  const [notificationApi, NotificationContextHolder] = notification.useNotification();

  const { accesses, loadedAtOnce, fetchAccesses, deleteAccess } = useAccessStore();

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
          <Space className="max-w-full truncate" size={4}>
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
        <Button.Group>
          <AccessEditModal
            data={record}
            preset="edit"
            trigger={
              <Tooltip title={t("access.action.edit")}>
                <Button color="primary" icon={<EditOutlinedIcon />} variant="text" />
              </Tooltip>
            }
          />

          <AccessEditModal
            data={{ ...record, id: undefined, name: `${record.name}-copy` }}
            preset="add"
            trigger={
              <Tooltip title={t("access.action.duplicate")}>
                <Button color="primary" icon={<SnippetsOutlinedIcon />} variant="text" />
              </Tooltip>
            }
          />

          <Tooltip title={t("access.action.delete")}>
            <Button
              color="danger"
              icon={<DeleteOutlinedIcon />}
              variant="text"
              onClick={() => {
                handleDeleteClick(record);
              }}
            />
          </Tooltip>
        </Button.Group>
      ),
    },
  ];
  const [tableData, setTableData] = useState<AccessModel[]>([]);
  const [tableTotal, setTableTotal] = useState<number>(0);

  const [page, setPage] = useState<number>(1);
  const [pageSize, setPageSize] = useState<number>(10);

  useEffect(() => {
    fetchAccesses().catch((err) => {
      if (err instanceof ClientResponseError && err.isAbort) {
        return;
      }

      console.error(err);
      notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
    });
  }, []);

  const { loading } = useRequest(
    () => {
      const startIndex = (page - 1) * pageSize;
      const endIndex = startIndex + pageSize;
      const items = accesses.slice(startIndex, endIndex);
      return Promise.resolve({
        items,
        totalItems: accesses.length,
      });
    },
    {
      refreshDeps: [accesses, page, pageSize],
      onSuccess: (data) => {
        setTableData(data.items);
        setTableTotal(data.totalItems);
      },
    }
  );

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
          notificationApi.error({ message: t("common.text.request_error"), description: getErrMsg(err) });
        }
      },
    });
  };

  return (
    <div className="p-4">
      {ModelContextHolder}
      {NotificationContextHolder}

      <PageHeader
        title={t("access.page.title")}
        extra={[
          <AccessEditModal
            key="create"
            preset="add"
            trigger={
              <Button type="primary" icon={<PlusOutlinedIcon />}>
                {t("access.action.add")}
              </Button>
            }
          />,
        ]}
      />

      <Table<AccessModel>
        columns={tableColumns}
        dataSource={tableData}
        loading={!loadedAtOnce || loading}
        locale={{
          emptyText: <Empty image={Empty.PRESENTED_IMAGE_SIMPLE} description={t("access.nodata")} />,
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
        rowKey={(record: AccessModel) => record.id}
        scroll={{ x: "max(100%, 960px)" }}
      />
    </div>
  );
};

export default AccessList;
