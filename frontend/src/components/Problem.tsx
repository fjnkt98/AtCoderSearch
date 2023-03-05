import { Item } from "../types/response";

type Props = {
  item: Item;
};

export function Problem({ item }: Props) {
  return (
    <div>
      <p>{item.problem_id}</p>
    </div>
  );
}
