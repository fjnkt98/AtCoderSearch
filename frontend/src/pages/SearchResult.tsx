import { useEffect } from "react";
import { apiHost } from "../libs/apiHost";
import { Logo } from "../components/Logo";
import { ProblemList } from "../components/ProblemList";
import { SearchBar } from "../components/SearchBar";
import { useSearchParams } from "react-router-dom";
import { SearchResponse } from "../types/response";
import { PageNavigation } from "../components/PageNavigation";
import { FieldFacetNavigationPart } from "../components/FieldFacetNavigationPart";
import { RangeFacetNavigationPart } from "../components/RangeFacetNavigationPart";
import { searchParamsState } from "../libs/searchParamsState";
import { searchResponseState } from "../libs/searchResponseState";
import { useRecoilState, useRecoilValue } from "recoil";
import { SortOrder } from "../components/SortOrder";
import { searchResponseFacetSelector } from "../libs/searchResponseState";
import { FacetNavigation } from "../components/FacetNavigation";

export function SearchResult() {
  const [, setSearchParams] = useRecoilState(searchParamsState);
  const [, setSearchResponse] = useRecoilState(searchResponseState);

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

  const facets = useRecoilValue(searchResponseFacetSelector);

  return (
    <div className="h-full w-full bg-zinc-800 text-slate-100">
      <div className="sticky top-0 z-[5000] mx-auto flex w-full flex-col items-center justify-center bg-zinc-800 py-2 shadow-sm shadow-gray-900">
        <div className="flex w-3/4 flex-row items-center justify-center gap-10">
          <Logo isBig={false} />
          <SearchBar />
        </div>
        <PageNavigation />
      </div>

      <div className="mx-auto my-2 flex w-3/4 min-w-[600px] flex-col items-center justify-center text-slate-100">
        <div className="my-2 flex w-full max-w-[800px] flex-row items-center justify-between">
          <div className="mx-2 flex-1">
            <SortOrder />
          </div>
          <div className="mx-2 flex-1">
            <FacetNavigation title="Category">
              <FieldFacetNavigationPart
                fieldName="category"
                facet={facets.category}
              />
            </FacetNavigation>
          </div>
          <div className="mx-2 flex-1">
            <FacetNavigation title="Difficulty">
              <RangeFacetNavigationPart
                fieldName="difficulty"
                facet={facets.difficulty}
              />
            </FacetNavigation>
          </div>
        </div>
        <ProblemList />
      </div>
    </div>
  );
}
