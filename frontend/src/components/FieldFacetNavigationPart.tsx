import { FieldFacetResult } from "../types/response";
import { useState, useEffect, useRef } from "react";
import { useRecoilValue } from "recoil";
import { searchParamsStateSelector } from "../libs/searchParamsState";

type Props = {
  fieldName: string;
  facet: FieldFacetResult;
  setParams: React.Dispatch<React.SetStateAction<Map<string, string>>>;
};

export function FieldFacetNavigationPart({
  fieldName,
  facet,
  setParams,
}: Props) {
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

  // 検索パラメータをセットするためのユーティリティ関数
  const setParam = (key: string, value: string) => {
    setParams((previousParams) => {
      previousParams.set(key, value);
      return previousParams;
    });
  };

  // 検索パラメータを削除するためのユーティリティ関数
  const deleteParam = (key: string) => {
    setParams((previousParams) => {
      previousParams.delete(key);
      return previousParams;
    });
  };

  // チェックボックスが更新されたら都度検索パラメータを更新する
  useEffect(() => {
    // チェックボックスが選択されているキーをカンマで結合する
    const targetCategories = checkboxState
      .filter(([, , checked]) => checked)
      .map(([key, ,]) => key)
      .join(",");

    // 何も選択されていない場合は検索パラメータをセットしない
    if (targetCategories !== "") {
      setParam(`filter.${fieldName}`, targetCategories);
    } else {
      deleteParam(`filter.${fieldName}`);
    }
  }, [checkboxState]);

  return (
    <div>
      <div className="flex flex-row items-center justify-between">
        <div className="p-1 text-xl">{fieldName}</div>
        <button
          className="text-lg text-blue-500"
          onClick={() => {
            deleteParam(`filter.${fieldName}`);
            setCheckboxState((previous) => {
              return previous.map(([key, count]) => [key, count, false]);
            });
          }}
        >
          reset
        </button>
      </div>

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
