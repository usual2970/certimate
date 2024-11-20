import { WorkflowRunLog } from "@/domain/workflow";
import { logs } from "@/repository/workflow";
import { ColumnDef } from "@tanstack/react-table";
import { useState } from "react";
import { DataTable } from "./DataTable";
import { useSearchParams } from "react-router-dom";
import { Check, X } from "lucide-react";
import WorkflowLogDetail from "./WorkflowLogDetail";

const WorkflowLog = () => {
  const [data, setData] = useState<WorkflowRunLog[]>([]);
  const [pageCount, setPageCount] = useState<number>(0);

  const [searchParams] = useSearchParams();
  const id = searchParams.get("id");

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
      header: "状态",
      cell: ({ row }) => {
        const succeed: boolean = row.getValue("succeed");
        if (succeed) {
          return (
            <div className="flex items-center space-x-2">
              <div className="text-white bg-green-500 w-8 h-8 rounded-full flex items-center justify-center">
                <Check size={18} />
              </div>
              <div className="text-sone-700">通过</div>
            </div>
          );
        } else {
          return (
            <div className="flex items-center space-x-2">
              <div className="text-white bg-red-500 w-8 h-8 rounded-full flex items-center justify-center">
                <X size={18} />
              </div>
              <div className="text-stone-700">失败</div>
            </div>
          );
        }
      },
    },
    {
      accessorKey: "error",
      header: "原因",
      cell: ({ row }) => {
        let error: string = row.getValue("error");
        if (!error) {
          error = "";
        }
        return <div className="min-w-[250px] truncate">{error}</div>;
      },
    },
    {
      accessorKey: "created",
      header: "时间",
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
        <div className="text-muted-foreground mb-5">日志</div>
        <DataTable columns={columns} data={data} onPageChange={fetchData} pageCount={pageCount} onRowClick={handleRowClick} />
      </div>

      <WorkflowLogDetail open={open} onOpenChange={setOpen} log={selectedLog} />
    </div>
  );
};

export default WorkflowLog;

