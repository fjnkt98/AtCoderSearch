import { useEffect } from "react";
import { apiHost } from "../libs/apiHost";
import { Logo } from "../components/Logo";
import { ProblemList } from "../components/ProblemList";
import { SearchBar } from "../components/SearchBar";
import { SideBar } from "../components/SideBar";
import { useSearchParams } from "react-router-dom";
import { SearchResponse } from "../types/response";
import { PageNavigation } from "../components/PageNavigation";
import { searchParamsState } from "../libs/searchParamsState";
import {
  searchResponseState,
  searchResponseTotalSelector,
  searchResponseCountSelector,
  searchResponseTimeSelector,
} from "../libs/searchResponseState";
import { useRecoilState, useRecoilValue } from "recoil";

export function SearchResult() {
  const [, setSearchParams] = useRecoilState(searchParamsState);
  const [, setSearchResponse] = useRecoilState(searchResponseState);
  const total = useRecoilValue(searchResponseTotalSelector);
  const count = useRecoilValue(searchResponseCountSelector);
  const time = useRecoilValue(searchResponseTimeSelector);

  const currentSearchParams = useSearchParams()[0];

  useEffect(() => {
    (async () => {
      setSearchParams(currentSearchParams);
      const url = new URL("/api/search", apiHost);
      const response = await fetch(`${url}?${currentSearchParams.toString()}`);
      const content: SearchResponse = await response.json();
      setSearchResponse(content);
    })();
  }, [currentSearchParams]);

  return (
    <div className="h-full w-full bg-zinc-800 text-slate-200">
      <div className="sticky top-0 mx-auto flex w-full flex-col items-center justify-center bg-zinc-800 py-2 shadow-sm shadow-black">
        <div className="flex w-3/4 flex-row items-center justify-center gap-10">
          <Logo isBig={false} />
          <SearchBar />
        </div>
        <PageNavigation />
      </div>

      <div className="flex flex-col px-6 lg:flex-row">
        <div className="mr-4 w-1/4 p-2">
          <SideBar />
        </div>
        <div className="w-3/4">
          <div className="mx-4 mt-6 text-slate-400">
            {count}件/{total}件 約{time / 1000}秒
          </div>
          <div>
            <ProblemList />
          </div>
        </div>
      </div>
    </div>
  );
}
