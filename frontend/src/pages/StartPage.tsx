import { Footer } from "../components/Footer";
import { Logo } from "../components/Logo";
import { SearchBar } from "../components/SearchBar";

export function StartPage() {
  return (
    <>
      <Logo isBig={true} />
      <SearchBar />
      <Footer />
    </>
  );
}
