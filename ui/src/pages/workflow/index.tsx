import { Button } from "@/components/ui/button";
import { MoreHorizontal, Plus } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { ColumnDef } from "@tanstack/react-table";
import { Workflow as WorkflowType } from "@/domain/workflow";
import { DataTable } from "@/components/workflow/DataTable";
import { useState } from "react";
import { list, remove, save } from "@/repository/workflow";
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuTrigger } from "@/components/ui/dropdown-menu";
import { Switch } from "@/components/ui/switch";

import { useTranslation } from "react-i18next";
import CustomAlertDialog from "@/components/workflow/CustomAlertDialog";

const Workflow = () => {
  const navigate = useNavigate();

  const [data, setData] = useState<WorkflowType[]>([]);
  const [pageCount, setPageCount] = useState<number>(0);

  const { t } = useTranslation();

  const [alertOpen, setAlertOpen] = useState(false);

  const [alertProps, setAlertProps] = useState<{
    title: string;
    description: string;
    onConfirm: () => void;
  }>();

  const fetchData = async (page: number, pageSize?: number) => {
    const resp = await list({ page: page, perPage: pageSize });
    setData(resp.items);
    setPageCount(resp.totalPages);
  };

  const columns: ColumnDef<WorkflowType>[] = [
    {
      accessorKey: "name",
      header: "名称",
      cell: ({ row }) => {
        let name: string = row.getValue("name");
        if (!name) {
          name = "未命名工作流";
        }
        return <div className="flex items-center">{name}</div>;
      },
    },
    {
      accessorKey: "description",
      header: "描述",
      cell: ({ row }) => {
        let description: string = row.getValue("description");
        if (!description) {
          description = "-";
        }
        return description;
      },
    },
    {
      accessorKey: "executionMethod",
      header: "执行方式",
      cell: ({ row }) => {
        const method = row.getValue("executionMethod");
        if (!method) {
          return "-";
        } else if (method === "manual") {
          return "手动";
        } else if (method === "auto") {
          const crontab: string = row.getValue("crontab");
          return (
            <div className="flex flex-col">
              定时
              <div className="text-muted-foreground text-xs">{crontab}</div>
            </div>
          );
        }
      },
    },
    {
      accessorKey: "enabled",
      header: "是否启用",
      cell: ({ row }) => {
        const enabled: boolean = row.getValue("enabled");

        return (
          <>
            <Switch
              checked={enabled}
              onCheckedChange={() => {
                handleCheckedChange(row.original.id ?? "");
              }}
            />
          </>
        );
      },
    },
    {
      accessorKey: "created",
      header: "创建时间",
      cell: ({ row }) => {
        const date: string = row.getValue("created");
        return new Date(date).toLocaleString();
      },
    },
    {
      accessorKey: "updated",
      header: "更新时间",
      cell: ({ row }) => {
        const date: string = row.getValue("updated");
        return new Date(date).toLocaleString();
      },
    },

    {
      id: "actions",
      cell: ({ row }) => {
        const workflow = row.original;
        return (
          <DropdownMenu>
            <DropdownMenuTrigger asChild>
              <Button variant="ghost" className="h-8 w-8 p-0">
                <span className="sr-only">Open menu</span>
                <MoreHorizontal className="h-4 w-4" />
              </Button>
            </DropdownMenuTrigger>
            <DropdownMenuContent align="end">
              <DropdownMenuLabel>操作</DropdownMenuLabel>
              <DropdownMenuItem
                onClick={() => {
                  navigate(`/workflow/detail?id=${workflow.id}`);
                }}
              >
                编辑
              </DropdownMenuItem>
              <DropdownMenuItem
                className="text-red-500"
                onClick={() => {
                  handleDeleteClick(workflow.id ?? "");
                }}
              >
                {t("common.delete")}
              </DropdownMenuItem>
            </DropdownMenuContent>
          </DropdownMenu>
        );
      },
    },
  ];

  const handleCheckedChange = async (id: string) => {
    const resp = await save({ id, enabled: !data.find((item) => item.id === id)?.enabled });
    if (resp) {
      setData((prev) => {
        return prev.map((item) => {
          if (item.id === id) {
            return resp;
          }
          return item;
        });
      });
    }
  };

  const handleDeleteClick = (id: string) => {
    setAlertProps({
      title: "删除工作流",
      description: "确定删除工作流吗？",
      onConfirm: async () => {
        const resp = await remove(id);
        if (resp) {
          setData((prev) => {
            return prev.filter((item) => item.id !== id);
          });
        }
      },
    });
    setAlertOpen(true);
  };
  const handleCreateClick = () => {
    navigate("/workflow/detail");
  };
  return (
    <>
      <div className="flex justify-between items-center">
        <div className="text-muted-foreground">工作流</div>
        <Button onClick={handleCreateClick}>
          <Plus size={16} />
          新建工作流
        </Button>
      </div>

      <div>
        <DataTable columns={columns} data={data} onPageChange={fetchData} pageCount={pageCount} />
      </div>

      <CustomAlertDialog
        open={alertOpen}
        onOpenChange={setAlertOpen}
        title={alertProps?.title}
        description={alertProps?.description}
        confirm={alertProps?.onConfirm}
      />
    </>
  );
};

export default Workflow;
