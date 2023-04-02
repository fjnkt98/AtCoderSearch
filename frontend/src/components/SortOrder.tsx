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
    <select
      onChange={handleChange}
      className="text-md block w-full cursor-pointer rounded-full border-none bg-zinc-900 py-2 px-2 text-slate-100 placeholder-slate-100 outline-none outline-1 outline-slate-700 hover:outline-blue-400"
    >
      <option hidden>{sortKeyMapping.get(selected ?? "-score")}</option>
      <option value="-score">{sortKeyMapping.get("-score")}</option>
      <option value="difficulty">{sortKeyMapping.get("difficulty")}</option>
      <option value="-difficulty">{sortKeyMapping.get("-difficulty")}</option>
      <option value="start_at">{sortKeyMapping.get("start_at")}</option>
      <option value="-start_at">{sortKeyMapping.get("-start_at")}</option>
    </select>
  );
}
