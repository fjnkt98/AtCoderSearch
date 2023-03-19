import { FacetNavigation } from "./FacetNavigation";
import { SortOrder } from "./SortOrder";
import { FacetResult } from "../types/response";

type Props = {
  searchParams: URLSearchParams;
  facet: Map<string, FacetResult>;
};

export function SideBar({ searchParams, facet }: Props) {
  return (
    <div className="px-0 py-6">
      <SortOrder searchParams={searchParams} />
      <FacetNavigation searchParams={searchParams} facets={facet} />
    </div>
  );
}
