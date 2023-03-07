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
    <>
      <Logo isBig={false} />
      <SearchBar />
      <div className="flex flex-row justify-between">
        <SideBar searchParams={searchParams} facet={facet} />
        <div className="">
          <ProblemList items={items} />
          <PageNavigation
            searchParams={searchParams}
            maxPageIndex={pages}
            currentPageIndex={index}
          />
        </div>
      </div>
    </>
  );
}
