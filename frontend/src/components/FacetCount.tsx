import { Facet } from "../types/response";

type Props = {
  searchParams: URLSearchParams;
  field: string;
  counts: Facet;
};

export function FacetCount({ searchParams, field, counts }: Props) {
  return (
    <div>
      <p>{field}</p>
      {counts.counts.map(({ key, count }) => (
        <p key={key}>
          {key}: {count}
        </p>
      ))}
    </div>
  );
}
