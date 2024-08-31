import { Button } from "../ui/button";

type PaginationProps = {
  totalPages: number;
  currentPage: number;
  onPageChange: (page: number) => void;
};

type PageNumber = number | string;

const Pagination = ({
  totalPages,
  currentPage,
  onPageChange,
}: PaginationProps) => {
  const pageNeighbours = 2; // Number of page numbers to show on either side of the current page

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
    <div className="pagination dark:text-stone-200">
      {pages.map((page, index) => {
        if (page === "...") {
          return (
            <span key={index} className="pagination-ellipsis">
              &hellip;
            </span>
          );
        }

        return (
          <Button
            key={index}
            className={`pagination-button ${
              currentPage === page ? "active" : ""
            }`}
            variant={currentPage === page ? "default" : "outline"}
            onClick={() => onPageChange(page as number)}
          >
            {page}
          </Button>
        );
      })}
    </div>
  );
};

export default Pagination;
