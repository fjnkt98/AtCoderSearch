import { useState, ChangeEvent, KeyboardEvent } from "react";
import { createSearchParams, useNavigate } from "react-router-dom";
import { HiPaperAirplane } from "react-icons/hi";
import { AiOutlineSearch } from "react-icons/ai";

export function SearchBar() {
  // 検索キーワード
  const [text, setText] = useState<string>("");
  // 検索用パラメータ
  // キーワード検索を実行したときの初期ページは1固定
  const [params, setParams] = useState<URLSearchParams>(
    createSearchParams({
      c: "20",
      p: "1",
    })
  );

  // 検索ボックスの入力を更新したときの処理
  const handleChange = (e: ChangeEvent<HTMLInputElement>) => {
    setText(e.target.value);

    // 入力の度に検索用パラメータを更新する
    setParams((params) => {
      params.set("q", e.target.value);
      return params;
    });
  };

  // ページ遷移のためのオブジェクト
  const navigate = useNavigate();

  // 検索結果ページへ遷移する関数
  const search = () => {
    // 空白文字だけの場合は検索を実行しない
    if (text.trim() === "") {
      return;
    }

    // 検索APIは検索結果ページに遷移して初めて実行される
    navigate(`/search?${params.toString()}`);
  };

  // エンターキーで検索を実行するためのハンドラ
  const handleKeyDown = (e: KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      e.preventDefault();
      search();
    }
  };

  return (
    <div className="my-6 flex h-12 w-1/2 min-w-min flex-row items-stretch justify-center rounded-full bg-slate-100 px-1 shadow-md shadow-gray-600">
      <div className="flex w-10 items-center bg-transparent p-1">
        <AiOutlineSearch className="w-full text-gray-900" size="1.5rem" />
      </div>
      <input
        type="text"
        className="flex-1 appearance-none bg-transparent px-2 font-notoSans text-lg text-gray-800 shadow-sm focus:border-none focus:outline-none"
        placeholder="Search Problems"
        onChange={handleChange}
        onKeyDown={handleKeyDown}
      />
      <button
        type="button"
        className="w-8 bg-transparent px-1"
        onClick={search}
      >
        <HiPaperAirplane className="rotate-90 text-blue-600" size="1.2rem" />
      </button>
    </div>
  );
}
