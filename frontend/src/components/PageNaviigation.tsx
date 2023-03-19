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
    params.set("page", i.toString());
    navigations.push(
      <Link
        key={`page-link-${i}`}
        className={`mx-2 flex aspect-square h-10 w-10 select-none items-center justify-center rounded-full text-center font-medium text-slate-200 shadow-sm shadow-gray-900 ${
          i == currentPageIndex ? "bg-blue-700" : "bg-zinc-700"
        }`}
        to={`/search?${params.toString()}`}
      >
        {i}
      </Link>
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
