import { useState, useEffect } from "react";
import { api_host } from "../libs/api_host";
import { Logo } from "../components/Logo";
import { ProblemList } from "../components/ProblemList";
import { SearchBar } from "../components/SearchBar";
import { SideBar } from "../components/SideBar";
import { useSearchParams } from "react-router-dom";
import { SearchResponse, Item, Facet } from "../types/response";

export function SearchResult() {
  const [items, setItems] = useState<Item[]>([]);
  const [facet, setFacet] = useState<Map<string, Facet>>(new Map());

  const [searchParams] = useSearchParams();

  useEffect(() => {
    (async () => {
      const url = new URL("/api/search", api_host);
      const response = await fetch(`${url}?${searchParams.toString()}`);
      const content: SearchResponse = await response.json();
      setItems(content.items);
      setFacet(content.stats.facet);
    })();
  }, [searchParams]);

  return (
    <>
      <Logo isBig={false} />
      <SearchBar />
      <div className="flex flex-row justify-between">
        <SideBar searchParams={searchParams} facet={facet} />
        <ProblemList items={items} />
      </div>
    </>
  );
}
