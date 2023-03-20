import { Problem } from "./Problem";
import { useRecoilValue } from "recoil";
import { searchResponseItemsSelector } from "../libs/searchResponseState";

export function ProblemList() {
  const items = useRecoilValue(searchResponseItemsSelector);

  return (
    <div>
      {items.map((item) => {
        return <Problem key={item.problem_id} item={item} />;
      })}
    </div>
  );
}
