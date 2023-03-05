import { ChangeEvent } from "react";
import { useNavigate } from "react-router-dom";

export function SortOrder() {
  const navigate = useNavigate();

  const handleChange = (e: ChangeEvent<HTMLSelectElement>) => {
    navigate(`/search?q=atcoder&s=${e.target.value}`);
  };

  return (
    <div className="relative w-full lg:max-w-sm">
      <select
        onChange={handleChange}
        className="w-full appearance-none rounded-md border bg-white p-2.5 text-gray-500 shadow-sm outline-none focus:border-indigo-600"
      >
        <option value="-score">デフォルト</option>
        <option value="difficulty">難易度(低い順)</option>
        <option value="-difficulty">難易度(高い順)</option>
      </select>
    </div>
  );
}
