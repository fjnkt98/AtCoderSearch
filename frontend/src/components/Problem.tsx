import { Item } from "../types/response";
import dayjs from "dayjs";
import timezone from "dayjs/plugin/timezone";
import utc from "dayjs/plugin/utc";

dayjs.extend(timezone);
dayjs.extend(utc);

type Props = {
  item: Item;
};

export function Problem({ item }: Props) {
  const start_at = dayjs(item.start_at)
    .tz("Asia/Tokyo")
    .format("YYYY/MM/DD HH:mm:ss");

  return (
    <div className="mx-2 my-2 border-2 border-solid px-2 py-2">
      <div className="flex flex-row">
        <div className="my-auto mx-2 border-2 border-solid">Logo</div>
        <div className="border-2 border-solid">
          <div>{item.problem_title}</div>
          <div>{item.problem_url}</div>
        </div>
      </div>
      <div className="flex flex-row border-2 border-solid">
        <div className="mx-1 rounded-full border-2 border-solid px-2 py-1">
          {start_at}
        </div>
        <div className="mx-1 rounded-full border-2 border-solid px-2 py-1">
          {item.category}
        </div>
        <div className="mx-1 rounded-full border-2 border-solid px-2 py-1">
          {item.difficulty}
        </div>
      </div>
    </div>
  );
}
