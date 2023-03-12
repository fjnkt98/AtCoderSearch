import { Link } from "react-router-dom";
import "../index.css";

type Props = {
  isBig: boolean;
};

export function Logo({ isBig }: Props) {
  const fontSize = isBig ? "text-6xl" : "text-2xl";
  return (
    <h1
      className={`text-center ${fontSize} border-2 border-solid border-red-500 p-1`}
    >
      <Link className="font-roboto shadow-lg" to="/">
        AtCoder Search
      </Link>
    </h1>
  );
}
