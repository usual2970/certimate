import { WorkflowRunLog } from "@/domain/workflow";
import { logs } from "@/repository/workflow";
import { ColumnDef } from "@tanstack/react-table";
import { useState } from "react";
import { DataTable } from "./DataTable";
import { useSearchParams } from "react-router-dom";
import { Check, X } from "lucide-react";
import WorkflowLogDetail from "./WorkflowLogDetail";
import { useTranslation } from "react-i18next";

const WorkflowLog = () => {
  const [data, setData] = useState<WorkflowRunLog[]>([]);
  const [pageCount, setPageCount] = useState<number>(0);

  const [searchParams] = useSearchParams();
  const id = searchParams.get("id");

  const { t } = useTranslation();

  const [open, setOpen] = useState(false);
  const [selectedLog, setSelectedLog] = useState<WorkflowRunLog>();

  const fetchData = async (page: number, pageSize?: number) => {
    const resp = await logs({ page: page, perPage: pageSize, id: id ?? "" });
    setData(resp.items);
    setPageCount(resp.totalPages);
  };

  const columns: ColumnDef<WorkflowRunLog>[] = [
    {
      accessorKey: "succeed",
      header: t("workflow.history.props.state"),
      cell: ({ row }) => {
        const succeed: boolean = row.getValue("succeed");
        if (succeed) {
          return (
            <div className="flex items-center space-x-2 min-w-[150px]">
              <div className="text-white bg-green-500 w-8 h-8 rounded-full flex items-center justify-center">
                <Check size={18} />
              </div>
              <div className="text-sone-700">{t("workflow.history.props.state.success")}</div>
            </div>
          );
        } else {
          return (
            <div className="flex items-center space-x-2 min-w-[150px]">
              <div className="text-white bg-red-500 w-8 h-8 rounded-full flex items-center justify-center">
                <X size={18} />
              </div>
              <div className="text-stone-700">{t("workflow.history.props.state.failed")}</div>
            </div>
          );
        }
      },
    },
    {
      accessorKey: "error",
      header: t("workflow.history.props.reason"),
      cell: ({ row }) => {
        let error: string = row.getValue("error");
        if (!error) {
          error = "";
        }
        return <div className="max-w-[300px] truncate text-red-500">{error}</div>;
      },
    },
    {
      accessorKey: "created",
      header: t("workflow.history.props.time"),
      cell: ({ row }) => {
        const date: string = row.getValue("created");
        return new Date(date).toLocaleString();
      },
    },
  ];

  const handleRowClick = (id: string) => {
    setOpen(true);
    const log = data.find((item) => item.id === id);
    setSelectedLog(log);
  };
  return (
    <div className="w-full md:w-[960px]">
      <div>
        <div className="text-muted-foreground mb-5">{t("workflow.history.page.title")}</div>
        <DataTable columns={columns} data={data} onPageChange={fetchData} pageCount={pageCount} onRowClick={handleRowClick} />
      </div>

      <WorkflowLogDetail open={open} onOpenChange={setOpen} log={selectedLog} />
    </div>
  );
};

export default WorkflowLog;
