import { FacetCount } from "./FacetCount";
import { SortOrder } from "./SortOrder";
import { FacetResult } from "../types/response";

type Props = {
  searchParams: URLSearchParams;
  facet: Map<string, FacetResult>;
};

export function SideBar({ searchParams, facet }: Props) {
  const facets = Array.from(Object.entries(facet));

  return (
    <div className="px-0 py-6">
      <SortOrder searchParams={searchParams} />

      {facets.map(([field, counts]) => (
        <FacetCount
          key={field}
          searchParams={searchParams}
          field={field}
          counts={counts}
        />
      ))}
    </div>
  );
}
