import { useState, useEffect } from "react";
import { api_host } from "../libs/api_host";
import { Logo } from "../components/Logo";
import { ProblemList } from "../components/ProblemList";
import { SearchBar } from "../components/SearchBar";
import { SideBar } from "../components/SideBar";
import { useSearchParams } from "react-router-dom";
import { SearchResponse, Item, Facet } from "../types/response";
import { PageNavigation } from "../components/PageNaviigation";

export function SearchResult() {
  const [items, setItems] = useState<Item[]>([]);
  const [facet, setFacet] = useState<Map<string, Facet>>(new Map());
  // const [total, setTotal] = useState<number>(0);
  const [pages, setPages] = useState<number>(0);
  const [index, setIndex] = useState<number>(0);

  const [searchParams] = useSearchParams();

  useEffect(() => {
    (async () => {
      const url = new URL("/api/search", api_host);
      const response = await fetch(`${url}?${searchParams.toString()}`);
      const content: SearchResponse = await response.json();
      setItems(content.items);
      setFacet(content.stats.facet);
      // setTotal(content.stats.total);
      setPages(content.stats.pages);
      setIndex(content.stats.index);
      console.log(content.stats.time);
    })();
  }, [searchParams]);

  return (
    <div className="bg-slate-200 text-gray-900 dark:bg-gray-900 dark:text-slate-100">
      <div className="mx-auto flex w-screen flex-col items-center justify-center pt-6">
        <Logo isBig={false} />
        <SearchBar />
      </div>

      <div className="flex flex-grow flex-row px-6">
        <div className="mr-4 w-1/4 p-2">
          <SideBar searchParams={searchParams} facet={facet} />
        </div>
        <div className="w-3/4">
          <ProblemList items={items} />
          <PageNavigation
            searchParams={searchParams}
            maxPageIndex={pages}
            currentPageIndex={index}
          />
        </div>
      </div>
    </div>
  );
}
