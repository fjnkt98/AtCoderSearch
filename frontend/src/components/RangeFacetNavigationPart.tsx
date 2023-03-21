import { RangeFacetResult, RangeFacetCount } from "../types/response";
import { Slider } from "@mui/material";
import { useState, useEffect } from "react";

type Props = {
  fieldName: string;
  facet: RangeFacetResult;
  setParams: React.Dispatch<React.SetStateAction<Map<string, string>>>;
};

export function RangeFacetNavigationPart({
  fieldName,
  facet,
  setParams,
}: Props) {
  const [facetCounts, setFacetCounts] = useState<RangeFacetCount[]>([]);
  useEffect(() => {
    const counts = [
      { begin: "", end: facet.start, count: Number(facet.before) ?? 0 },
      ...facet.counts,
      { begin: facet.end, end: "", count: Number(facet.after) ?? 0 },
    ];
    setFacetCounts(counts);
  }, [facet]);

  const [difficulties, setDifficulties] = useState<number[]>([0, 4000]);
  const handleChange = (event: Event, newValue: number | number[]) => {
    setDifficulties(newValue as number[]);
    const [begin, end] = newValue as number[];
    setParams((previous) => {
      previous.set(`filter.${fieldName}.from`, begin.toString());
      previous.set(`filter.${fieldName}.to`, end.toString());
      return previous;
    });
  };

  return (
    <div className="mt-3">
      <div className="flex flex-row items-center justify-between">
        <div className="p-1 text-xl">{fieldName}</div>
        <button
          className="text-lg text-blue-500"
          onClick={() => {
            setParams((previous) => {
              previous.delete(`filter.${fieldName}.from`);
              previous.delete(`filter.${fieldName}.to`);
              return previous;
            });
            setDifficulties([0, 4000]);
          }}
        >
          reset
        </button>
      </div>

      <div>
        {facetCounts.map(({ begin, end, count }, index) => (
          <div
            key={`${fieldName}-range-racet-${index}`}
            className="my-2 flex cursor-pointer flex-row items-center justify-between rounded-xl shadow-sm shadow-gray-700"
            onClick={() => {
              if (begin === "") {
                setDifficulties([Number(end), Number(end)]);
                setParams((previous) => {
                  previous.delete(`filter.${fieldName}.from`);
                  previous.set(`filter.${fieldName}.to`, end.toString());
                  return previous;
                });
              } else if (end === "") {
                setDifficulties([Number(begin), Number(begin)]);
                setParams((previous) => {
                  previous.set(`filter.${fieldName}.from`, begin.toString());
                  previous.delete(`filter.${fieldName}.to`);
                  return previous;
                });
              } else {
                setDifficulties([Number(begin), Number(end)]);
                setParams((previous) => {
                  previous.set(`filter.${fieldName}.from`, begin.toString());
                  previous.set(`filter.${fieldName}.to`, end.toString());
                  return previous;
                });
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
    </div>
  );
}
