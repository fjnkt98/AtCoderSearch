import { Logo } from "../components/Logo";
import { ProblemList } from "../components/ProblemList";
import { SearchBar } from "../components/SearchBar";
import { SideBar } from "../components/SideBar";

export function SearchResult() {
  return (
    <>
      <Logo isBig={false} />
      <SearchBar />
      <div className="flex flex-row justify-between">
        <SideBar />
        <ProblemList />
      </div>
    </>
  );
}
