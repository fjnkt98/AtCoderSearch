import { Link } from "react-router-dom";
import { searchParamsStateSelector } from "../libs/searchParamsState";
import {
  searchResponseIndexSelector,
  searchResponsePagesSelector,
} from "../libs/searchResponseState";
import { useRecoilValue } from "recoil";

export function PageNavigation() {
  const navigations: JSX.Element[] = [];
  const index = useRecoilValue(searchResponseIndexSelector);
  const pages = useRecoilValue(searchResponsePagesSelector);
  const searchParams = useRecoilValue(searchParamsStateSelector);

  const pageBegin: number = Math.max(1, index - 5);
  const pageEnd: number = Math.min(pages, pageBegin + 9);

  for (let i = pageBegin; i <= pageEnd; i++) {
    const params = new URLSearchParams(searchParams);
    params.set("page", i.toString());
    navigations.push(
      <Link
        key={`page-link-${i}`}
        className={`mx-2 flex aspect-square h-10 w-10 select-none items-center justify-center rounded-full text-center font-medium text-slate-200 shadow-sm shadow-gray-900 ${
          i == index ? "bg-blue-700" : "bg-zinc-700"
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
