import CertificateDetail from "@/components/certificate/CertificateDetail";
import { Button } from "@/components/ui/button";
import { DataTable } from "@/components/workflow/DataTable";
import { Certificate as CertificateType } from "@/domain/certificate";
import { diffDays, getLeftDays } from "@/lib/time";
import { list } from "@/repository/certificate";
import { ColumnDef } from "@tanstack/react-table";
import { useState } from "react";
import { useNavigate } from "react-router-dom";

const Certificate = () => {
  const [data, setData] = useState<CertificateType[]>([]);
  const [pageCount, setPageCount] = useState<number>(0);
  const [open, setOpen] = useState(false);
  const [selectedCertificate, setSelectedCertificate] = useState<CertificateType>();

  const fetchData = async (page: number, pageSize?: number) => {
    const resp = await list({ page: page, perPage: pageSize });
    setData(resp.items);
    setPageCount(resp.totalPages);
  };

  const navigate = useNavigate();

  const columns: ColumnDef<CertificateType>[] = [
    {
      accessorKey: "san",
      header: "域名",
      cell: ({ row }) => {
        let san: string = row.getValue("san");
        if (!san) {
          san = "";
        }

        return (
          <div>
            {san.split(";").map((item, i) => {
              return (
                <div key={i} className="max-w-[250px] truncate">
                  {item}
                </div>
              );
            })}
          </div>
        );
      },
    },
    {
      accessorKey: "expireAt",
      header: "有效期限",
      cell: ({ row }) => {
        const expireAt: string = row.getValue("expireAt");
        const data = row.original;
        const leftDays = getLeftDays(expireAt);
        const allDays = diffDays(data.expireAt, data.created);
        return (
          <div className="">
            {leftDays > 0 ? (
              <div className="text-green-500">
                {leftDays} / {allDays} 天
              </div>
            ) : (
              <div className="text-red-500">已到期</div>
            )}

            <div>{new Date(expireAt).toLocaleString().split(" ")[0]} 到期</div>
          </div>
        );
      },
    },
    {
      accessorKey: "workflow",
      header: "所属工作流",
      cell: ({ row }) => {
        const name = row.original.expand.workflow?.name;
        const workflowId: string = row.getValue("workflow");
        return (
          <div className="max-w-[200px] truncate">
            <Button
              size={"sm"}
              variant={"link"}
              onClick={() => {
                handleWorkflowClick(workflowId);
              }}
            >
              {name}
            </Button>
          </div>
        );
      },
    },
    {
      accessorKey: "created",
      header: "颁发时间",
      cell: ({ row }) => {
        const date: string = row.getValue("created");
        return new Date(date).toLocaleString();
      },
    },
    {
      id: "actions",
      cell: ({ row }) => {
        return (
          <div className="flex items-center space-x-2">
            <Button
              size="sm"
              variant={"link"}
              onClick={() => {
                handleView(row.original.id);
              }}
            >
              查看证书
            </Button>
          </div>
        );
      },
    },
  ];

  const handleWorkflowClick = (id: string) => {
    navigate(`/workflow/detail?id=${id}`);
  };

  const handleView = (id: string) => {
    setOpen(true);
    const certificate = data.find((item) => item.id === id);
    setSelectedCertificate(certificate);
  };

  return (
    <div className="flex flex-col space-y-5">
      <div className="text-muted-foreground">证书</div>

      <DataTable columns={columns} onPageChange={fetchData} data={data} pageCount={pageCount} />

      <CertificateDetail open={open} onOpenChange={setOpen} certificate={selectedCertificate} />
    </div>
  );
};

export default Certificate;
