import { FieldFacetResult } from "../types/response";
import { useState, useEffect } from "react";

type Props = {
  fieldName: string;
  facet: FieldFacetResult;
  setFilteredSearchParams: React.Dispatch<
    React.SetStateAction<URLSearchParams>
  >;
};

export function FieldFacetNavigationPart({
  fieldName,
  facet,
  setFilteredSearchParams,
}: Props) {
  const [checkboxState, setCheckboxState] = useState<
    [string, number, boolean][]
  >(facet.counts.map(({ key, count }) => [key, count, false]));

  // setCheckboxState(facet.counts.map(({ key, count }) => [key, count, false]));
  // useEffect(() => {
  //   setCheckboxState(facet.counts.map(({ key, count }) => [key, count, false]));
  // }, [facet]);

  const handleCheckboxChange = (index: number, isChecked: boolean) => {
    setCheckboxState((previousState) => {
      previousState[index][2] = isChecked;
      return previousState;
    });

    setFilteredSearchParams((previousParams) => {
      for (const [index, [key, , isChecked]] of checkboxState.entries()) {
        console.log(isChecked);
        if (isChecked) {
          previousParams.set(`filter[${fieldName}][${index}]`, key);
        } else {
          previousParams.delete(`filter[${fieldName}][${index}]`);
        }
      }
      console.log(previousParams.toString());
      return previousParams;
    });
  };

  return (
    <div>
      <div className="flex flex-row items-center justify-between">
        <div className="p-1 text-xl">{fieldName}</div>
        <button
          className="text-lg text-blue-500"
          onClick={() => {
            setFilteredSearchParams((previousParams) => {
              const target = Array.from(previousParams.keys()).filter((key) =>
                key.startsWith(`filter[${fieldName}]`)
              );
              for (const key of target) {
                previousParams.delete(key);
              }
              return previousParams;
            });
            // setCheckboxState((previous) => {
            //   return previous.map(([key, count]) => [key, count, false]);
            // });
          }}
        >
          reset
        </button>
      </div>

      <div>
        {checkboxState.map(([key, count, checked], index) => (
          <FilteringCheckBox
            key={`${fieldName}-filtering-checkbox-${index}`}
            facetKey={key}
            count={count}
            index={index}
            onCheckboxChange={handleCheckboxChange}
          />
        ))}
      </div>
    </div>
  );
}

type CheckboxProps = {
  facetKey: string;
  count: number;
  index: number;
  onCheckboxChange: (index: number, isChecked: boolean) => void;
};

function FilteringCheckBox({
  facetKey,
  count,
  index,
  onCheckboxChange,
}: CheckboxProps) {
  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    onCheckboxChange(index, event.target.checked);
  };
  return (
    <div className="flex flex-row items-center justify-between">
      <input
        id={`${facetKey}-filtering-${index}`}
        type="checkbox"
        className="inline-block h-4 w-4 rounded-xl focus:border-blue-600 focus:ring-2 focus:ring-blue-600 focus:ring-opacity-50"
        onChange={handleChange}
      />
      <label
        htmlFor={`${facetKey}-filtering-${index}`}
        className="inline-block flex-auto cursor-pointer select-none break-all"
      >
        {facetKey}
      </label>
      <label
        htmlFor={`${facetKey}-filtering-${index}`}
        className="inline-block flex-auto cursor-pointer select-none break-all"
      >
        {count}
      </label>
    </div>
  );
}
