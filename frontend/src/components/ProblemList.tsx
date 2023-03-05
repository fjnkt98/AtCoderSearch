import { Problem } from "./Problem";
import { Item } from "../types/response";

type Props = {
  items: Item[];
};

export function ProblemList({ items }: Props) {
  return (
    <div>
      {items.map((item) => {
        return <Problem key={item.problem_id} item={item} />;
      })}
    </div>
  );
}
