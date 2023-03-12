import { Facet } from "../types/response";

type Props = {
  searchParams: URLSearchParams;
  field: string;
  counts: Facet;
};

export function FacetCount({ searchParams, field, counts }: Props) {
  return (
    <div className="my-4 rounded-xl bg-zinc-900 py-2 px-1">
      <p className="p-1 text-xl">{field}</p>
      {counts.counts.map(({ key, count }) => (
        <div
          key={key}
          className="my-1 flex flex-row justify-between rounded-full bg-slate-800 px-4 py-1"
        >
          <div>{key}</div>
          <div>{count}</div>
        </div>
      ))}
    </div>
  );
}
