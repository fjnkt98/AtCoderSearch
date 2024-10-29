import type { MetaFunction } from "@remix-run/node";
import { Link } from "@remix-run/react";

export const meta: MetaFunction = () => {
  return [
    { title: "AtCoder Search" },
    { name: "description", content: "AtCoder Search" },
  ];
};

export default function Index() {
  return (
    <div className="flex flex-col justify-between items-center gap-4 py-2">
      <Link
        to="/problem"
        className="rounded-md border px-4 py-3 min-w-40 w-2/3 text-center text-lg max-w-96 shadow-sm shadow-gray-500"
      >
        問題を探す
      </Link>
      <Link
        to="/user"
        className="rounded-md border px-4 py-3 min-w-40 w-2/3 text-center text-lg max-w-96 shadow-sm shadow-gray-500"
      >
        ユーザを探す
      </Link>
      <Link
        to="/submission"
        className="rounded-md border px-4 py-3 min-w-40 w-2/3 text-center text-lg max-w-96 shadow-sm shadow-gray-500"
      >
        提出を探す
      </Link>
    </div>
  );
}
