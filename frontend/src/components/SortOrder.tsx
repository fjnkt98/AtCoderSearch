import { ChangeEvent, useState, useEffect } from "react";
import { useNavigate } from "react-router-dom";
import { useRecoilValue } from "recoil";
import { searchParamsStateSelector } from "../libs/searchParamsState";

const sortKeyMapping = new Map<string, string>([
  ["-score", "関連度順"],
  ["difficulty", "難易度(低い順)"],
  ["-difficulty", "難易度(高い順)"],
  ["start_at", "開催時期(早い順)"],
  ["-start_at", "開催時期(遅い順)"],
]);

export function SortOrder() {
  const navigate = useNavigate();
  const searchParams = useRecoilValue(searchParamsStateSelector);
  const [selected, setSelected] = useState<string | null>(
    searchParams.get("sort") ?? ""
  );
  useEffect(() => {
    setSelected(searchParams.get("sort") ?? "-score");
  }, [searchParams]);

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
        <option hidden>{sortKeyMapping.get(selected ?? "-score")}</option>
        <option value="-score">{sortKeyMapping.get("-score")}</option>
        <option value="difficulty">{sortKeyMapping.get("difficulty")}</option>
        <option value="-difficulty">{sortKeyMapping.get("-difficulty")}</option>
        <option value="start_at">{sortKeyMapping.get("start_at")}</option>
        <option value="-start_at">{sortKeyMapping.get("-start_at")}</option>
      </select>
    </div>
  );
}
