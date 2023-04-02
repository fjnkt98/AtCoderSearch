import { Problem } from "./Problem";
import { useRecoilValue } from "recoil";
import {
  searchResponseItemsSelector,
  searchResponseTotalSelector,
  searchResponseCountSelector,
  searchResponseTimeSelector,
} from "../libs/searchResponseState";

export function ProblemList() {
  const items = useRecoilValue(searchResponseItemsSelector);
  const total = useRecoilValue(searchResponseTotalSelector);
  const count = useRecoilValue(searchResponseCountSelector);
  const time = useRecoilValue(searchResponseTimeSelector);

  return (
    <div className="">
      <span className="inline-block px-4 text-slate-400">
        {count}件/{total}件 約{time / 1000}秒
      </span>
      {items.map((item) => {
        return <Problem key={item.problem_id} item={item} />;
      })}
    </div>
  );
}
