import { ColumnDef, flexRender, getCoreRowModel, getPaginationRowModel, PaginationState, useReactTable } from "@tanstack/react-table";

import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { Button } from "../ui/button";
import { useEffect, useState } from "react";
import Show from "../Show";
import { useTranslation } from "react-i18next";

interface DataTableProps<TData extends { id: string }, TValue> {
  columns: ColumnDef<TData, TValue>[];
  data: TData[];
  pageCount: number;
  onPageChange?: (pageIndex: number, pageSize?: number) => Promise<void>;
  onRowClick?: (id: string) => void;
  withPagination?: boolean;
  fallback?: React.ReactNode;
}

export function DataTable<TData extends { id: string }, TValue>({
  columns,
  data,
  onPageChange,
  pageCount,
  onRowClick,
  withPagination,
  fallback,
}: DataTableProps<TData, TValue>) {
  const [{ pageIndex, pageSize }, setPagination] = useState<PaginationState>({
    pageIndex: 0,
    pageSize: 10,
  });

  const { t } = useTranslation();

  const pagination = {
    pageIndex,
    pageSize,
  };

  const table = useReactTable({
    data,
    columns,
    pageCount: pageCount,
    getCoreRowModel: getCoreRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    state: {
      pagination,
    },
    onPaginationChange: setPagination,
    manualPagination: true,
  });

  useEffect(() => {
    onPageChange?.(pageIndex, pageSize);
  }, [pageIndex]);

  const handleRowClick = (id: string) => {
    onRowClick?.(id);
  };

  return (
    <div>
      <div className="rounded-md">
        <Table>
          <TableHeader>
            {table.getHeaderGroups().map((headerGroup) => (
              <TableRow key={headerGroup.id} className="dark:border-muted-foreground">
                {headerGroup.headers.map((header) => {
                  return <TableHead key={header.id}>{header.isPlaceholder ? null : flexRender(header.column.columnDef.header, header.getContext())}</TableHead>;
                })}
              </TableRow>
            ))}
          </TableHeader>
          <TableBody className="dark:text-stone-200">
            {table.getRowModel().rows?.length ? (
              table.getRowModel().rows.map((row) => (
                <TableRow
                  key={row.id}
                  data-state={row.getIsSelected() && "selected"}
                  className="dark:border-muted-foreground"
                  onClick={(e) => {
                    e.stopPropagation();
                    handleRowClick(row.original.id);
                  }}
                >
                  {row.getVisibleCells().map((cell) => (
                    <TableCell key={cell.id}>{flexRender(cell.column.columnDef.cell, cell.getContext())}</TableCell>
                  ))}
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={columns.length} className="h-24 text-center">
                  {fallback ? fallback : t("common.text.nodata")}
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </div>
      <Show when={!!withPagination}>
        <div className="flex items-center justify-end mt-5">
          <div className="flex items-center space-x-2 dark:text-stone-200">
            {table.getCanPreviousPage() && (
              <Button variant="outline" size="sm" onClick={() => table.previousPage()} disabled={!table.getCanPreviousPage()}>
                {t("common.pagination.prev")}
              </Button>
            )}

            {table.getCanNextPage && (
              <Button variant="outline" size="sm" onClick={() => table.nextPage()} disabled={!table.getCanNextPage()}>
                {t("common.pagination.next")}
              </Button>
            )}
          </div>
        </div>
      </Show>
    </div>
  );
}
