import { Footer } from "../components/Footer";
import { Logo } from "../components/Logo";
import { SearchBar } from "../components/SearchBar";

export function StartPage() {
  return (
    <>
      <div className="container mx-auto flex w-full justify-center border-2 border-solid border-sky-500 py-10">
        <Logo isBig={true} />
      </div>
      <SearchBar />
      <Footer />
    </>
  );
}
