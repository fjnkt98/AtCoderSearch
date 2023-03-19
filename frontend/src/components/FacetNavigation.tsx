import { FacetResult } from "../types/response";
import { useNavigate } from "react-router-dom";

type Props = {
  searchParams: URLSearchParams;
  facets: Map<string, FacetResult>;
};

export function FacetNavigation({ searchParams, facets }: Props) {
  const navigate = useNavigate();
  const facetEntries = (
    facets: Map<string, FacetResult>
  ): [string, FacetResult][] => {
    return Object.entries(facets);
  };

  return (
    <div className="my-4 flex min-w-[120px] flex-col rounded-xl bg-zinc-900 py-2 px-2">
      {facetEntries(facets).map(([field, facet]) => {
        return (
          <div key={`facet-${field}`}>
            <div className="flex flex-row items-center justify-between">
              <div className="p-1 text-xl">{field}</div>
              <button
                className="text-lg text-blue-500"
                onClick={() => {
                  const target = Array.from(searchParams.keys()).filter((key) =>
                    key.startsWith(`filter[${field}]`)
                  );
                  for (const key of target) {
                    searchParams.delete(key);
                  }
                }}
              >
                reset
              </button>
            </div>

            <div>
              {facet.counts.map(({ key, count }, i) => {
                return (
                  <div
                    key={`facet-item-${i}`}
                    className="flex flex-row items-center justify-between"
                  >
                    <input
                      id={`${field}-filtering-${i}`}
                      type="checkbox"
                      className="inline-block h-4 w-4 rounded-xl focus:border-blue-600 focus:ring-2 focus:ring-blue-600 focus:ring-opacity-50"
                      onChange={(e) => {
                        if (facet.range_info == null) {
                          if (e.target.checked) {
                            searchParams.set(`filter[${field}][${i}]`, key);
                          } else {
                            searchParams.delete(`filter[${field}][${i}]`);
                          }
                        }
                      }}
                    />
                    <label
                      htmlFor={`${field}-filtering-${i}`}
                      className="inline-block flex-auto cursor-pointer select-none break-all"
                    >
                      {key}
                    </label>
                    <span className="inline-block">{count}</span>
                  </div>
                );
              })}
            </div>
          </div>
        );
      })}

      <button
        className=""
        onClick={() => {
          navigate(`/search?${searchParams.toString()}`);
        }}
      >
        絞り込む
      </button>
    </div>
  );
}
