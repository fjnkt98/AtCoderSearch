import { Link } from "react-router-dom";
import "../index.css";

type Props = {
  isBig: boolean;
};

export function Logo({ isBig }: Props) {
  const fontSize = isBig ? "text-6xl" : "text-3xl";
  return (
    <h1 className={`text-center ${fontSize} p-1`}>
      <Link className="font-roboto" to="/">
        AtCoder Search
      </Link>
    </h1>
  );
}
