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
      <Link
        key={`page-link-${i}`}
        to={`/search?${params.toString()}`}
        className={`m-1 border-2 border-solid px-4 py-2 ${
          i == currentPageIndex ? "text-red-500" : "text-black"
        }`}
      >
        {i}
      </Link>
    );
  }

  return (
    <div key="page-navigation" className="flex flex-row">
      {navigations}
    </div>
  );
}
