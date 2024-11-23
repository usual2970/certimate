import { Button } from "@/components/ui/button";
import { MoreHorizontal, Plus } from "lucide-react";
import { useNavigate, useSearchParams } from "react-router-dom";
import { ColumnDef } from "@tanstack/react-table";
import { Workflow as WorkflowType } from "@/domain/workflow";
import { DataTable } from "@/components/workflow/DataTable";
import { useState } from "react";
import { list, remove, save, WorkflowListReq } from "@/repository/workflow";
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

  const [searchParams] = useSearchParams();

  const fetchData = async (page: number, pageSize?: number) => {
    const state = searchParams.get("state");
    const req: WorkflowListReq = { page: page, perPage: pageSize };
    if (state && state == "enabled") {
      req.enabled = true;
    }
    const resp = await list(req);
    setData(resp.items);
    setPageCount(resp.totalPages);
  };

  const columns: ColumnDef<WorkflowType>[] = [
    {
      accessorKey: "name",
      header: t("workflow.props.name"),
      cell: ({ row }) => {
        let name: string = row.getValue("name");
        if (!name) {
          name = t("workflow.props.name.default");
        }
        return <div className="max-w-[150px] truncate">{name}</div>;
      },
    },
    {
      accessorKey: "description",
      header: t("workflow.props.description"),
      cell: ({ row }) => {
        let description: string = row.getValue("description");
        if (!description) {
          description = "-";
        }
        return <div className="max-w-[200px] truncate">{description}</div>;
      },
    },
    {
      accessorKey: "type",
      header: t("workflow.props.executionMethod"),
      cell: ({ row }) => {
        const method = row.getValue("type");
        if (!method) {
          return "-";
        } else if (method === "manual") {
          return t("workflow.node.start.form.executionMethod.options.manual");
        } else if (method === "auto") {
          const crontab: string = row.original.crontab ?? "";
          return (
            <div className="flex flex-col space-y-1">
              <div>{t("workflow.node.start.form.executionMethod.options.auto")}</div>
              <div className="text-muted-foreground text-xs">{crontab}</div>
            </div>
          );
        }
      },
    },
    {
      accessorKey: "enabled",
      header: t("workflow.props.enabled"),
      cell: ({ row }) => {
        const enabled: boolean = row.getValue("enabled");

        return (
          <>
            <Switch
              checked={enabled}
              onCheckedChange={() => {
                handleCheckedChange(row.original.id ?? "");
              }}
              onClick={(e) => {
                e.stopPropagation();
              }}
            />
          </>
        );
      },
    },
    {
      accessorKey: "created",
      header: t("workflow.props.created"),
      cell: ({ row }) => {
        const date: string = row.getValue("created");
        return new Date(date).toLocaleString();
      },
    },
    {
      accessorKey: "updated",
      header: t("workflow.props.updated"),
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
              <DropdownMenuLabel>{t("workflow.action")}</DropdownMenuLabel>
              <DropdownMenuItem
                onClick={(e) => {
                  e.stopPropagation();
                  navigate(`/workflow/detail?id=${workflow.id}`);
                }}
              >
                {t("workflow.action.edit")}
              </DropdownMenuItem>
              <DropdownMenuItem
                className="text-red-500"
                onClick={(e) => {
                  e.stopPropagation();
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
      title: t("workflow.action.delete.alert.title"),
      description: t("workflow.action.delete.alert.description"),
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

  const handleRowClick = (id: string) => {
    navigate(`/workflow/detail?id=${id}`);
  };
  return (
    <>
      <div className="flex justify-between items-center">
        <div className="text-muted-foreground">{t("workflow.page.title")}</div>
        <Button onClick={handleCreateClick}>
          <Plus size={16} />
          {t("workflow.action.create")}
        </Button>
      </div>

      <div>
        <DataTable columns={columns} data={data} onPageChange={fetchData} pageCount={pageCount} onRowClick={handleRowClick} withPagination />
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
