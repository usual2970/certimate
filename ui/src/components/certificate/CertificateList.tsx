import CertificateDetail from "@/components/certificate/CertificateDetail";
import { Button } from "@/components/ui/button";
import { DataTable } from "@/components/workflow/DataTable";
import { Certificate as CertificateType } from "@/domain/certificate";
import { diffDays, getLeftDays } from "@/lib/time";
import { list } from "@/repository/certificate";
import { ColumnDef } from "@tanstack/react-table";
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useNavigate, useSearchParams } from "react-router-dom";

type CertificateListProps = {
  withPagination?: boolean;
};

const CertificateList = ({ withPagination }: CertificateListProps) => {
  const [data, setData] = useState<CertificateType[]>([]);
  const [pageCount, setPageCount] = useState<number>(0);
  const [open, setOpen] = useState(false);
  const [selectedCertificate, setSelectedCertificate] = useState<CertificateType>();

  const { t } = useTranslation();

  const [searchParams] = useSearchParams();

  const fetchData = async (page: number, pageSize?: number) => {
    const state = searchParams.get("state");
    const resp = await list({ page: page, perPage: pageSize, state: state ?? "" });
    setData(resp.items);
    setPageCount(resp.totalPages);
  };

  const navigate = useNavigate();

  const columns: ColumnDef<CertificateType>[] = [
    {
      accessorKey: "san",
      header: t("certificate.props.domain"),
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
      header: t("certificate.props.expiry"),
      cell: ({ row }) => {
        const expireAt: string = row.getValue("expireAt");
        const data = row.original;
        const leftDays = getLeftDays(expireAt);
        const allDays = diffDays(data.expireAt, data.created);
        return (
          <div className="">
            {leftDays > 0 ? (
              <div className="text-green-500">
                {leftDays} / {allDays} {t("certificate.props.expiry.days")}
              </div>
            ) : (
              <div className="text-red-500">{t("certificate.props.expiry.expired")}</div>
            )}

            <div>
              {new Date(expireAt).toLocaleString().split(" ")[0]} {t("certificate.props.expiry.text.expire")}
            </div>
          </div>
        );
      },
    },
    {
      accessorKey: "workflow",
      header: t("certificate.props.workflow"),
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
      header: t("common.text.created_at"),
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
              {t("certificate.action.view")}
            </Button>
          </div>
        );
      },
    },
  ];

  const handleWorkflowClick = (id: string) => {
    navigate(`/workflows/detail?id=${id}`);
  };

  const handleView = (id: string) => {
    setOpen(true);
    const certificate = data.find((item) => item.id === id);
    setSelectedCertificate(certificate);
  };

  return (
    <>
      <DataTable
        columns={columns}
        onPageChange={fetchData}
        data={data}
        pageCount={pageCount}
        withPagination={withPagination}
        fallback={
          <div className="flex flex-col items-center">
            <div className="text-muted-foreground">{t("certificate.nodata")}</div>
            <Button
              size={"sm"}
              className="w-[120px] mt-3"
              onClick={() => {
                navigate("/workflows/detail");
              }}
            >
              {t("workflow.action.create")}
            </Button>
          </div>
        }
      />

      <CertificateDetail open={open} onOpenChange={setOpen} certificate={selectedCertificate} />
    </>
  );
};

export default CertificateList;
