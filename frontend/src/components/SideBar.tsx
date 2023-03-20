import { FacetNavigation } from "./FacetNavigation";
import { SortOrder } from "./SortOrder";
import { FieldFacetResult, RangeFacetResult } from "../types/response";

type Props = {
  searchParams: URLSearchParams;
  facets: {
    category: FieldFacetResult;
    difficulty: RangeFacetResult;
  };
};

export function SideBar({ searchParams, facets }: Props) {
  return (
    <div className="px-0 py-6">
      <SortOrder searchParams={searchParams} />
      <FacetNavigation searchParams={searchParams} facets={facets} />
    </div>
  );
}
