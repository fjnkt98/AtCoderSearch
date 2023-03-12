import { FacetCount } from "./FacetCount";
import { SortOrder } from "./SortOrder";
import { Facet } from "../types/response";

type Props = {
  searchParams: URLSearchParams;
  facet: Map<string, Facet>;
};

export function SideBar({ searchParams, facet }: Props) {
  const facets = [];
  for (const [field, counts] of Object.entries(facet)) {
    facets.push({
      field,
      counts,
    });
  }
  facets.sort((a, b) => (a.field > b.field ? 1 : -1));

  return (
    <div className="px-0">
      <SortOrder searchParams={searchParams} />

      {facets.map(({ field, counts }) => (
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
