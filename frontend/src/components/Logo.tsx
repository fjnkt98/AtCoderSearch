import { Link } from "react-router-dom";
import "../index.css";

type Props = {
  isBig: boolean;
};

export function Logo({ isBig }: Props) {
  const fontSize = isBig ? "text-4xl" : "text-2xl";
  return (
    <h1 className={`text-center ${fontSize} border-2 border-solid p-2`}>
      <Link to="/">AtCoder Search</Link>
    </h1>
  );
}
