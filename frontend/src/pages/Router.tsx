import { BrowserRouter, Routes, Route } from "react-router-dom";
import { StartPage } from "./StartPage";
import { SearchResult } from "./SearchResult";

export function Router() {
  return (
    <BrowserRouter>
      <Routes>
        <Route index element={<StartPage />} />
        <Route path="search" element={<SearchResult />} />
      </Routes>
    </BrowserRouter>
  );
}
