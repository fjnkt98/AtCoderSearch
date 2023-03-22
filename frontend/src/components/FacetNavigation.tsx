import { useState } from "react";
import { useNavigate, createSearchParams } from "react-router-dom";
import { FieldFacetNavigationPart } from "./FieldFacetNavigationPart";
import { RangeFacetNavigationPart } from "./RangeFacetNavigationPart";
import { useRecoilValue } from "recoil";
import { searchResponseFacetSelector } from "../libs/searchResponseState";
import { searchParamsStateSelector } from "../libs/searchParamsState";

export function FacetNavigation() {
  const searchParams = useRecoilValue(searchParamsStateSelector);
  const [params, setParams] = useState<Map<string, string>>(
    new Map<string, string>()
  );
  const facets = useRecoilValue(searchResponseFacetSelector);

  const navigate = useNavigate();

  return (
    <div className="my-4  flex min-w-[240px] flex-col rounded-xl bg-zinc-900 py-2 px-4">
      <FieldFacetNavigationPart
        fieldName="category"
        facet={facets.category}
        setParams={setParams}
      />
      <RangeFacetNavigationPart
        fieldName="difficulty"
        facet={facets.difficulty}
        setParams={setParams}
      />
      <button
        className="my-1 mx-2 rounded-full bg-teal-800 py-2"
        onClick={() => {
          const filteredSearchParams = createSearchParams();
          for (const [key, value] of searchParams.entries()) {
            if (!key.startsWith("filter")) {
              filteredSearchParams.set(key, value);
            }
          }
          for (const [key, value] of params.entries()) {
            filteredSearchParams.set(key, value);
          }

          navigate(`/search?${filteredSearchParams.toString()}`);
        }}
      >
        絞り込む
      </button>
    </div>
  );
}
