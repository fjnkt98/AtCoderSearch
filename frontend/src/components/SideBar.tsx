import { FacetCount } from "./FacetCount";
import { SortOrder } from "./SortOrder";

export function SideBar() {
  return (
    <div>
      <SortOrder />
      <FacetCount />
      <FacetCount />
    </div>
  );
}
