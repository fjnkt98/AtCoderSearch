import { RangeFacetResult } from "../types/response";

type Props = {
  fieldName: string;
  facet: RangeFacetResult;
  setFilteredSearchParams: React.Dispatch<
    React.SetStateAction<URLSearchParams>
  >;
};

export function RangeFacetNavigationPart({
  fieldName,
  facet,
  setFilteredSearchParams,
}: Props) {
  return <div></div>;
}
