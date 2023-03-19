import { FacetResult } from "../types/response";
import { Link, useNavigate } from "react-router-dom";

type Props = {
  searchParams: URLSearchParams;
  field: string;
  counts: FacetResult;
};

export function FacetCount({ searchParams, field, counts }: Props) {
  const navigate = useNavigate();
  const handleClick = () => {
    searchParams.delete(`filter[${field}][]`);
    searchParams.delete(`filter[${field}][from]`);
    searchParams.delete(`filter[${field}][to]`);

    navigate(`/search?${searchParams.toString()}`);
  };

  return (
    <div className="my-4 rounded-xl bg-zinc-900 py-2 px-1">
      <div className="flex flex-row items-center justify-between">
        <p className="p-1 text-xl">{field}</p>
        <button className="text-lg text-blue-500" onClick={handleClick}>
          reset
        </button>
      </div>
      {counts.counts.map(({ key, count }) => {
        const params = new URLSearchParams(searchParams);
        if (counts.range_info == null) {
          params.set(`filter[${field}][]`, key);
        } else {
          const start = Number(key);
          const end = Number(key) + Number(counts.range_info.gap);
          params.set(`filter[${field}][from]`, start.toString());
          params.set(`filter[${field}][to]`, end.toString());
        }
        const linkTo = `/search?${params.toString()}`;

        return (
          <Link
            key={key}
            to={linkTo}
            className="my-1 flex flex-row justify-between rounded-full bg-slate-800 px-4 py-2"
          >
            <div>{key}</div>
            <div>{count}</div>
          </Link>
        );
      })}
    </div>
  );
}
