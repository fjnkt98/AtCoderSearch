import { useState } from "react";
import { FieldFacetResult, RangeFacetResult } from "../types/response";
import { useNavigate, useSearchParams } from "react-router-dom";
import { FieldFacetNavigationPart } from "./FieldFacetNavigationPart";

type Props = {
  searchParams: URLSearchParams;
  facets: {
    category: FieldFacetResult;
    difficulty: RangeFacetResult;
  };
};

export function FacetNavigation({ searchParams, facets }: Props) {
  const [filteredSearchParams, setFilteredSearchParams] =
    useState<URLSearchParams>(searchParams);

  const navigate = useNavigate();

  return (
    <div className="my-4 flex min-w-[120px] flex-col rounded-xl bg-zinc-900 py-2 px-2">
      <FieldFacetNavigationPart
        fieldName="category"
        facet={facets.category}
        setFilteredSearchParams={setFilteredSearchParams}
      />

      <button
        className=""
        onClick={() => {
          navigate(`/search?${filteredSearchParams.toString()}`);
        }}
      >
        絞り込む
      </button>
    </div>
  );
}
