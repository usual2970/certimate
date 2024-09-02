import {
  Pagination,
  PaginationContent,
  PaginationEllipsis,
  PaginationItem,
  PaginationLink,
} from "../ui/pagination";

type PaginationProps = {
  totalPages: number;
  currentPage: number;
  onPageChange: (page: number) => void;
};

type PageNumber = number | string;

const XPagination = ({
  totalPages,
  currentPage,
  onPageChange,
}: PaginationProps) => {
  const pageNeighbours = 1; // Number of page numbers to show on either side of the current page

  const getPageNumbers = () => {
    const totalNumbers = pageNeighbours * 2 + 3; // total pages to display (left + right neighbours + current + 2 for start and end)
    const totalBlocks = totalNumbers + 2; // adding 2 for the start and end page numbers

    if (totalPages > totalBlocks) {
      let pages: PageNumber[] = [];

      const leftBound = Math.max(2, currentPage - pageNeighbours);
      const rightBound = Math.min(totalPages - 1, currentPage + pageNeighbours);

      const beforeLastPage = totalPages - 1;

      pages = range(leftBound, rightBound);

      if (currentPage > pageNeighbours + 2) {
        pages.unshift("...");
      }
      if (currentPage < beforeLastPage - pageNeighbours) {
        pages.push("...");
      }

      pages.unshift(1);
      pages.push(totalPages);

      return pages;
    }

    return range(1, totalPages);
  };

  const range = (from: number, to: number, step = 1) => {
    let i = from;
    const range = [];

    while (i <= to) {
      range.push(i);
      i += step;
    }

    return range;
  };

  const pages = getPageNumbers();

  return (
    <>
      <Pagination className="dark:text-stone-200 justify-end mt-3">
        <PaginationContent>
          {pages.map((page, index) => {
            if (page === "...") {
              return (
                <PaginationItem key={index}>
                  <PaginationEllipsis />
                </PaginationItem>
              );
            }

            return (
              <PaginationItem key={index}>
                <PaginationLink
                  href="#"
                  isActive={currentPage == page}
                  onClick={(e) => {
                    e.preventDefault();
                    onPageChange(page as number);
                  }}
                >
                  {page}
                </PaginationLink>
              </PaginationItem>
            );
          })}
        </PaginationContent>
      </Pagination>
    </>
  );
};

export default XPagination;
