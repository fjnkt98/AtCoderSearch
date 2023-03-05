import { useState, useEffect } from "react";
import { api_host } from "../libs/api_host";
import { Logo } from "../components/Logo";
import { ProblemList } from "../components/ProblemList";
import { SearchBar } from "../components/SearchBar";
import { SideBar } from "../components/SideBar";
import { useSearchParams } from "react-router-dom";
import { SearchResponse, Item } from "../types/response";

export function SearchResult() {
  const [items, setItems] = useState<Item[]>([]);

  const [searchParams] = useSearchParams();

  useEffect(() => {
    (async () => {
      const url = new URL("/api/search", api_host);
      const response = await fetch(`${url}?${searchParams.toString()}`);
      const content: SearchResponse = await response.json();
      setItems(content.items);
    })();
  }, [searchParams]);

  return (
    <>
      <Logo isBig={false} />
      <SearchBar />
      <div className="flex flex-row justify-between">
        <SideBar searchParams={searchParams} />
        <ProblemList items={items} />
      </div>
    </>
  );
}
