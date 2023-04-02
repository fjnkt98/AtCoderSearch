import { FieldFacetResult } from "../types/response";
import { useState, useEffect, useRef } from "react";
import { useRecoilValue } from "recoil";
import { useNavigate, createSearchParams } from "react-router-dom";
import { searchParamsStateSelector } from "../libs/searchParamsState";

type Props = {
  fieldName: string;
  facet: FieldFacetResult;
};

export function FieldFacetNavigationPart({ fieldName, facet }: Props) {
  // 子要素のチェックボックスの情報を保持・管理するステート
  // ステートの更新は子要素から行う
  const [checkboxState, setCheckboxState] = useState<
    [string, number, boolean][]
  >(facet.counts.map(({ key, count }) => [key, count, false]));

  const searchParams = useRecoilValue(searchParamsStateSelector);
  useEffect(() => {
    // 検索パラメータからチェックボックスの選択状態を更新する
    const categories = (searchParams.get(`filter.${fieldName}`) ?? "").split(
      ","
    );
    setCheckboxState(
      facet.counts.map(({ key, count }) => {
        if (categories.includes(key)) {
          return [key, count, true];
        } else {
          return [key, count, false];
        }
      })
    );
  }, [facet]);

  const navigate = useNavigate();

  return (
    <div className="my-2 flex-1">
      <div>
        {checkboxState.map(([key, count, checked], index) => (
          <FilteringCheckbox
            key={`${fieldName}-filtering-checkbox-${index}`}
            facetKey={key}
            count={count}
            index={index}
            checked={checked}
            setCheckboxState={setCheckboxState}
          />
        ))}
      </div>

      <div className="flex flex-row items-center justify-between pt-3">
        <button
          className="mx-2 rounded-full bg-blue-700 py-1 px-2 text-sm text-slate-100"
          onClick={() => {
            const filteredSearchParams = createSearchParams(searchParams);
            filteredSearchParams.delete(`filter.${fieldName}`);

            const targetCategories = checkboxState
              .filter(([, , checked]) => checked)
              .map(([key, ,]) => key)
              .join(",");

            if (targetCategories !== "") {
              filteredSearchParams.set(`filter.${fieldName}`, targetCategories);
            }
            navigate(`/search?${filteredSearchParams.toString()}`);
          }}
        >
          絞り込む
        </button>
        <button
          className="text-lg text-blue-500"
          onClick={() => {
            setCheckboxState((previous) => {
              return previous.map(([key, count]) => [key, count, false]);
            });

            const filteredSearchParams = createSearchParams(searchParams);
            filteredSearchParams.delete(`filter.${fieldName}`);
            navigate(`/search?${filteredSearchParams.toString()}`);
          }}
        >
          Reset
        </button>
      </div>
    </div>
  );
}

type CheckboxProps = {
  facetKey: string;
  count: number;
  index: number;
  checked: boolean;
  setCheckboxState: React.Dispatch<
    React.SetStateAction<[string, number, boolean][]>
  >;
};

function FilteringCheckbox({
  facetKey,
  count,
  index,
  checked,
  setCheckboxState,
}: CheckboxProps) {
  const checkbox = useRef<HTMLInputElement>(null);
  if (checkbox != null && checkbox.current != null) {
    checkbox.current.checked = checked;
  }

  return (
    <div className="my-2 flex flex-row items-center justify-between rounded-xl shadow-sm shadow-gray-700">
      <input
        id={`${facetKey}-filtering-${index}`}
        type="checkbox"
        className="mx-1 inline-block h-4 w-4 rounded-xl focus:border-blue-600 focus:ring-2 focus:ring-blue-600 focus:ring-opacity-50"
        onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
          setCheckboxState((previousState) => {
            const newState: [string, number, boolean][] = previousState.map(
              ([key, count, checked]) => [key, count, checked]
            );
            newState[index][2] = e.target.checked;
            return newState;
          });
        }}
        defaultChecked={checked}
        ref={checkbox}
      />
      <label
        htmlFor={`${facetKey}-filtering-${index}`}
        className="mx-1 inline-block flex-grow cursor-pointer select-none break-all"
      >
        {facetKey}
      </label>
      <label
        htmlFor={`${facetKey}-filtering-${index}`}
        className="mx-1 inline-block cursor-pointer select-none break-all"
      >
        {count}
      </label>
    </div>
  );
}
