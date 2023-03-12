import { Facet } from "../types/response";
import { Link, useNavigate } from "react-router-dom";

type Props = {
  searchParams: URLSearchParams;
  field: string;
  counts: Facet;
};

export function FacetCount({ searchParams, field, counts }: Props) {
  const navigate = useNavigate();
  const handleClick = () => {
    searchParams.delete(`f[${field}][]`);
    searchParams.delete(`f[${field}][from]`);
    searchParams.delete(`f[${field}][to]`);

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
        if (counts.gap == null) {
          params.set(`f[${field}][]`, key);
        } else {
          const start = Number(key);
          const end = Number(key) + Number(counts.gap);
          params.set(`f[${field}][from]`, start.toString());
          params.set(`f[${field}][to]`, end.toString());
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
