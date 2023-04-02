import { Footer } from "../components/Footer";
import { Logo } from "../components/Logo";
import { SearchBar } from "../components/SearchBar";

export function StartPage() {
  return (
    <>
      <div className="flex min-h-screen flex-col items-center justify-center bg-zinc-800 dark:text-slate-200">
        <div className="mx-auto flex w-screen flex-grow flex-col items-center justify-center">
          <Logo isBig={true} />
          <SearchBar />
        </div>
        <div className="min-h-2 w-screen flex-shrink-0 text-center">
          <Footer />
        </div>
      </div>
    </>
  );
}
