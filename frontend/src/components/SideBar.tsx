import { FacetNavigation } from "./FacetNavigation";
import { SortOrder } from "./SortOrder";

export function SideBar() {
  return (
    <div className="min-w-[240px] px-0 py-6">
      <SortOrder />
      <FacetNavigation />
    </div>
  );
}
