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
    <div className="flex flex-col justify-between items-center gap-4">
      <Link to="/problem">Problem</Link>
      <Link to="/user">User</Link>
      <Link to="/submission">Submission</Link>
    </div>
  );
}
