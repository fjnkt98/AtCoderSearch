import { RangeFacetResult, RangeFacetCount } from "../types/response";
import { Slider } from "@mui/material";
import { useState, useEffect } from "react";
import { createSearchParams, useNavigate } from "react-router-dom";
import { useRecoilValue } from "recoil";
import { searchParamsStateSelector } from "../libs/searchParamsState";

type Props = {
  fieldName: string;
  facet: RangeFacetResult;
};

export function RangeFacetNavigationPart({ fieldName, facet }: Props) {
  const [facetCounts, setFacetCounts] = useState<RangeFacetCount[]>([]);
  useEffect(() => {
    const counts = [
      { begin: "", end: facet.start, count: Number(facet.before) ?? 0 },
      ...facet.counts,
      { begin: facet.end, end: "", count: Number(facet.after) ?? 0 },
    ];
    setFacetCounts(counts);
  }, [facet]);

  const searchParams = useRecoilValue(searchParamsStateSelector);

  const [difficulties, setDifficulties] = useState<number[]>([0, 4000]);
  const handleChange = (event: Event, newValue: number | number[]) => {
    setDifficulties(newValue as number[]);
  };

  const navigate = useNavigate();

  return (
    <div className="mt-3">
      <div>
        {facetCounts.map(({ begin, end, count }, index) => (
          <div
            key={`${fieldName}-range-racet-${index}`}
            className="my-2 flex cursor-pointer flex-row items-center justify-between rounded-xl shadow-sm shadow-gray-700"
            onClick={() => {
              if (begin === "") {
                setDifficulties([-999999, Number(end)]);
              } else if (end === "") {
                setDifficulties([Number(begin), 999999]);
              } else {
                setDifficulties([Number(begin), Number(end)]);
              }
            }}
          >
            <span className="mx-1 inline-block flex-grow select-none break-all">{`${begin} ~ ${end}`}</span>
            <span className="mx-1 inline-block select-none break-all">
              {count}
            </span>
          </div>
        ))}
      </div>

      <div className="px-4">
        <Slider
          value={difficulties}
          onChange={handleChange}
          min={0}
          max={4000}
          step={400}
          marks
          valueLabelDisplay="auto"
        />
      </div>

      <div className="flex flex-row items-center justify-between py-2">
        <button
          className="mx-2 rounded-full bg-blue-700 py-1 px-2 text-sm text-slate-100"
          onClick={() => {
            const filteredSearchParams = createSearchParams(searchParams);
            filteredSearchParams.delete(`filter.${fieldName}.to`);
            filteredSearchParams.delete(`filter.${fieldName}.from`);

            const [begin, end] = difficulties;
            if (begin !== -999999) {
              filteredSearchParams.set(
                `filter.${fieldName}.from`,
                begin.toString()
              );
            }
            if (end !== 999999) {
              filteredSearchParams.set(
                `filter.${fieldName}.to`,
                end.toString()
              );
            }

            navigate(`/search?${filteredSearchParams.toString()}`);
          }}
        >
          絞り込む
        </button>

        <button
          className="text-lg text-blue-500"
          onClick={() => {
            setDifficulties([0, 4000]);

            const filteredSearchParams = createSearchParams(searchParams);
            filteredSearchParams.delete(`filter.${fieldName}.to`);
            filteredSearchParams.delete(`filter.${fieldName}.from`);
            navigate(`/search?${filteredSearchParams.toString()}`);
          }}
        >
          Reset
        </button>
      </div>
    </div>
  );
}
