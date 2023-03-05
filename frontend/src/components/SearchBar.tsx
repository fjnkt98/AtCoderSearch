import { useState, ChangeEvent } from "react";
import { Link, createSearchParams } from "react-router-dom";

export function SearchBar() {
  const [text, setText] = useState<string>("");
  const [url, setUrl] = useState<string>("/search");

  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    setText(() => {
      return e.target.value;
    });
    const params: string = createSearchParams({
      q: e.target.value,
    }).toString();
    setUrl(`/search?${params}`);
  };

  return (
    <div className="my-2 flex flex-row items-center">
      <input
        type="text"
        placeholder="Search Problems"
        onChange={handleChange}
      ></input>
      <button type="button">
        <Link to={url}>Search</Link>
      </button>
    </div>
  );
}
