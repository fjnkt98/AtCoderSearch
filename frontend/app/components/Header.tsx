import { Link } from "@remix-run/react";

export default function Header() {
  return (
    <div className="flex flex-row w-full items-center gap-2 px-3 py-2 justify-start min-w-60">
      <img
        alt="AtCoder Search Logo"
        src="/logo.svg"
        className="aspect-square h-6 rounded-full select-none"
      />
      <Link to="/" className="text-xl font-medium select-none">
        AtCoder Search
      </Link>
    </div>
  );
}
