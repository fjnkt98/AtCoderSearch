import { useRef } from "react";
import { Link } from "react-router-dom";

export function SearchBar() {
  const inputRef = useRef<HTMLInputElement>(null);

  return (
    <div className="my-2 flex flex-row items-center">
      <input type="text" placeholder="Search Problems" ref={inputRef}></input>
      <button type="button">
        <Link to="/search">Search</Link>
      </button>
    </div>
  );
}
