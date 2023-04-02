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
    <div className="flex flex-col items-center px-8">
      <span className="mt-2 inline-block w-full px-4 text-slate-400">
        {count}件/{total}件 約{time / 1000}秒
      </span>
      {items.map((item) => {
        return <Problem key={item.problem_id} item={item} />;
      })}
    </div>
  );
}
