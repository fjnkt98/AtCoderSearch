import { useState, useEffect } from "react";
import { apiHost } from "../libs/apiHost";
import { Logo } from "../components/Logo";
import { ProblemList } from "../components/ProblemList";
import { SearchBar } from "../components/SearchBar";
import { SideBar } from "../components/SideBar";
import { useSearchParams } from "react-router-dom";
import { SearchResponse } from "../types/response";
import { PageNavigation } from "../components/PageNavigation";

export function SearchResult() {
  const [response, setResponse] = useState<SearchResponse>({
    stats: {
      time: 0,
      total: 0,
      pages: 0,
      index: 0,
      count: 0,
      facet: new Map(),
    },
    items: [],
    message: null,
  });

  const searchParams = useSearchParams()[0];

  useEffect(() => {
    (async () => {
      const url = new URL("/api/search", apiHost);
      const response = await fetch(`${url}?${searchParams.toString()}`);
      const content: SearchResponse = await response.json();
      setResponse(content);
    })();
  }, [searchParams]);

  return (
    <div className="h-full w-full bg-zinc-800 text-slate-200">
      <div className="sticky top-0 mx-auto flex w-full flex-col items-center justify-center bg-zinc-800 py-2 shadow-sm shadow-black">
        <div className="flex w-3/4 flex-row items-center justify-center gap-10">
          <Logo isBig={false} />
          <SearchBar />
        </div>
        <PageNavigation
          searchParams={searchParams}
          maxPageIndex={response?.stats.pages}
          currentPageIndex={response?.stats.index}
        />
      </div>

      <div className="flex flex-col px-6 lg:flex-row">
        <div className="mr-4 w-1/4 p-2">
          <SideBar searchParams={searchParams} facet={response.stats.facet} />
        </div>
        <div className="w-3/4">
          <div className="mx-4 mt-6 text-slate-400">
            {response.stats.count}件/{response.stats.total}件 約
            {response.stats.time / 1000}秒
          </div>
          <div>
            <ProblemList items={response.items} />
          </div>
        </div>
      </div>
    </div>
  );
}
