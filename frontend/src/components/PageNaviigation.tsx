import { Link } from "react-router-dom";

type Props = {
  searchParams: URLSearchParams;
  maxPageIndex: number;
  currentPageIndex: number;
};

export function PageNavigation({
  searchParams,
  maxPageIndex,
  currentPageIndex,
}: Props) {
  const navigations: JSX.Element[] = [];

  const pageBegin: number = Math.max(1, currentPageIndex - 5);
  const pageEnd: number = Math.min(maxPageIndex, pageBegin + 9);

  for (let i = pageBegin; i <= pageEnd; i++) {
    const params = new URLSearchParams(searchParams);
    params.set("p", i.toString());
    navigations.push(
      <div
        className={`mx-2 flex aspect-square h-10 w-10 select-none items-center justify-center rounded-full text-center text-gray-900 shadow-sm shadow-gray-900 dark:text-slate-100 ${
          i == currentPageIndex ? "bg-blue-600" : "bg-gray-800"
        }`}
        key={`page-link-${i}`}
      >
        <Link to={`/search?${params.toString()}`}>{i}</Link>
      </div>
    );
  }

  return (
    <div
      key="page-navigation"
      className="flex flex-row items-center justify-center"
    >
      {navigations}
    </div>
  );
}
