import { FacetCount } from "./FacetCount";
import { SortOrder } from "./SortOrder";

type Props = {
  searchParams: URLSearchParams;
};

export function SideBar({ searchParams }: Props) {
  return (
    <div>
      <SortOrder searchParams={searchParams} />
      <FacetCount />
      <FacetCount />
    </div>
  );
}
