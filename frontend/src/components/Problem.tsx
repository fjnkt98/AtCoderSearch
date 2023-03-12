import { Item } from "../types/response";
import dayjs from "dayjs";
import timezone from "dayjs/plugin/timezone";
import utc from "dayjs/plugin/utc";
import atcoderLogo from "../assets/atcoder_logo.svg";

dayjs.extend(timezone);
dayjs.extend(utc);

type Props = {
  item: Item;
};

export function Problem({ item }: Props) {
  const startAt = dayjs(item.start_at)
    .tz("Asia/Tokyo")
    .format("YYYY/MM/DD HH:mm:ss");

  const categoryColor = new Map<string, string>([
    ["ABC", "bg-blue-600"],
    ["ABC-Like", "bg-sky-600"],
    ["AGC", "bg-yellow-600"],
    ["AGC-Like", "bg-amber-500"],
    ["AHC", "bg-green-500"],
    ["ARC", "bg-red-500"],
    ["ARC-Like", "bg-orange-700"],
    ["JAG", "bg-slate-500"],
    ["JOI", "bg-slate-600"],
    ["Marathon", "bg-slate-600"],
    ["Other Contests", "bg-slate-600"],
    ["Other Sponsored", "bg-slate-600"],
    ["PAST", "bg-slate-600"],
  ]);

  const difficultyColor = (difficulty: number): string => {
    if (difficulty < 0) {
      return "bg-black";
    } else if (difficulty < 400) {
      return "bg-slate-500";
    } else if (difficulty < 800) {
      return "bg-amber-900";
    } else if (difficulty < 1200) {
      return "bg-green-600";
    } else if (difficulty < 1600) {
      return "bg-sky-600";
    } else if (difficulty < 2000) {
      return "bg-blue-600";
    } else if (difficulty < 2400) {
      return "bg-yellow-600";
    } else if (difficulty < 2800) {
      return "bg-orange-600";
    } else if (difficulty < 3200) {
      return "bg-red-600";
    } else if (difficulty < 3600) {
      return "";
    } else {
      return "";
    }
  };

  return (
    <div className="mx-2 my-5 rounded-2xl bg-slate-200 px-2 py-2 text-gray-900 shadow-sm shadow-slate-700 dark:bg-gray-800 dark:text-slate-100">
      <div className="flex flex-row items-center">
        <img
          alt="AtCoder Logo"
          src={atcoderLogo}
          className="m-2 aspect-square h-12  rounded-full bg-white"
        />
        <div className="mx-2">
          <p className="text-xl">{item.problem_title}</p>
          <a
            href={item.problem_url}
            className="text-sm text-blue-500"
            target="_blank"
            rel="noreferrer"
          >
            {item.problem_url}
          </a>
        </div>
      </div>
      <div className="flex flex-row pt-1">
        <div className="mx-1 rounded-full bg-gray-700 px-2 py-1">{startAt}</div>
        <div
          className={`mx-1 rounded-full ${categoryColor.get(
            item.category
          )} px-2 py-1`}
        >
          {item.category}
        </div>
        <div
          className={`mx-1 rounded-full bg-gray-700 px-2 py-1 ${difficultyColor(
            item.difficulty
          )}`}
        >
          {item.difficulty}
        </div>
      </div>
    </div>
  );
}
