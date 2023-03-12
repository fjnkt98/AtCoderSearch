import { ChangeEvent } from "react";
import { useNavigate } from "react-router-dom";

type Props = {
  searchParams: URLSearchParams;
};

export function SortOrder({ searchParams }: Props) {
  const navigate = useNavigate();

  const handleChange = (e: ChangeEvent<HTMLSelectElement>) => {
    searchParams.set("p", "1");
    searchParams.set("s", e.target.value);
    navigate(`/search?${searchParams.toString()}`);
  };

  return (
    <div className="w-full lg:max-w-sm">
      <select
        onChange={handleChange}
        className="w-full appearance-none rounded-md border bg-white p-2.5 text-gray-900 shadow-sm outline-none focus:border-indigo-600"
      >
        <option value="-score">関連度順</option>
        <option value="difficulty">難易度(低い順)</option>
        <option value="-difficulty">難易度(高い順)</option>
      </select>
    </div>
  );
}
