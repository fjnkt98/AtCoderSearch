import { Link } from "@remix-run/react";

interface Props {
  to: string;
  params: URLSearchParams;
  current: number;
  end?: number;
  width?: number;
}

const defaultWidth = 5;

export default function Pagination({ to, params, current, end, width }: Props) {
  let labels: string[] = [];
  let containsBegin = false;
  let containsEnd = false;

  const left = Math.max(current - Math.floor((width ?? defaultWidth) / 2), 1);
  const right = Math.min(left + (width ?? defaultWidth), end ?? 999999);
  for (let i = left; i <= right; i++) {
    if (i === 1) {
      containsBegin = true;
    }
    if (i === end) {
      containsEnd = true;
    }
    labels.push(i.toString());
  }

  if (!containsBegin) {
    if (labels[0] === "2") {
      labels = ["1", ...labels];
    } else {
      labels = ["1", "...", ...labels];
    }
  }
  if (!containsEnd) {
    if (end == null) {
      labels = [...labels, "..."];
    } else {
      labels = [...labels, "...", end.toString()];
    }
  }

  return (
    <div className="flex flex-row gap-1 items-center justify-center">
      {labels.map((label) => {
        if (label === "...") {
          return (
            <span
              key={label}
              className="block w-10 h-10 border border-gray-500 dark:border-gray-400 text-center rounded-lg p-2"
            >
              {label}
            </span>
          );
        }

        const p = new URLSearchParams(params);
        p.set("page", label);
        return (
          <Link
            key={label}
            to={`${to}?${p.toString()}`}
            className={`w-10 h-10 border border-gray-500 dark:border-gray-400 text-center rounded-lg p-2 ${
              current.toString() === label
                ? "bg-gray-950 text-gray-50 dark:bg-gray-50 dark:text-gray-950"
                : ""
            }`}
          >
            {label}
          </Link>
        );
      })}
    </div>
  );
}
