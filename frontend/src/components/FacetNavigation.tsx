import { useEffect, useState, useRef } from "react";
import { HiChevronDown } from "react-icons/hi";

type Props = {
  title: string;
  children: React.ReactNode;
};

export function FacetNavigation(props: Props) {
  // 絞り込みメニューのオンオフを管理するステート
  const [showMenu, setShowMenu] = useState<boolean>(false);
  // 絞り込みメニューを開くためのHTMLエレメントへの参照
  const menuHeaderRef = useRef<HTMLDivElement>(null);
  // 絞り込みメニュー本体のHTMLエレメントへの参照
  const menuBodyRef = useRef<HTMLDivElement>(null);

  // メニュー以外をクリックしたときにメニューを閉じるためのロジック
  useEffect(() => {
    // クリックイベントを検知する関数
    const handleClickOutside = (event: MouseEvent) => {
      // タイプガード
      if (event.target instanceof Node) {
        // - ヘッダとボディがあること
        // - クリックされたのがヘッダでないこと
        // - クリックされたのがボディでないこと
        // 以上の条件をすべて満たす場合に「メニュー以外がクリックされた」と判定し、メニューを閉じる
        if (
          menuHeaderRef.current &&
          menuBodyRef.current &&
          !menuHeaderRef.current.contains(event.target) &&
          !menuBodyRef.current.contains(event.target)
        ) {
          setShowMenu(false);
        }
      }
    };

    document.addEventListener("mousedown", handleClickOutside);

    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
    };
  }, [menuHeaderRef]);

  return (
    <div className="text-md relative z-[4000]">
      <div
        className="flex cursor-pointer select-none flex-row items-center justify-between p-1.5"
        onClick={() => setShowMenu((previous) => !previous)}
        ref={menuHeaderRef}
      >
        <span className="text-md inline-block">{props.title}</span>
        <HiChevronDown className={showMenu ? "rotate-180" : "rotate-0"} />
      </div>

      <div
        className={`absolute top-7 ${
          showMenu ? "" : "hidden"
        } w-full border-2 bg-zinc-800 p-2`}
        ref={menuBodyRef}
      >
        {props.children}
      </div>
    </div>
  );
}
