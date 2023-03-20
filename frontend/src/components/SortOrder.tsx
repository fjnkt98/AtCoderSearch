import { ChangeEvent } from "react";
import { useNavigate } from "react-router-dom";
import { useRecoilValue } from "recoil";
import { searchParamsStateSelector } from "../libs/searchParamsState";

export function SortOrder() {
  const searchParams = useRecoilValue(searchParamsStateSelector);

  const navigate = useNavigate();

  const handleChange = (e: ChangeEvent<HTMLSelectElement>) => {
    const params = new URLSearchParams(searchParams);
    params.set("page", "1");
    params.set("sort", e.target.value);

    navigate(`/search?${params.toString()}`);
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
        <option value="start_at">開催時期(早い順)</option>
        <option value="-start_at">開催時期(遅い順))</option>
      </select>
    </div>
  );
}
