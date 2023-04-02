import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import { Router } from "./pages/Router";
import { RecoilRoot } from "recoil";

ReactDOM.createRoot(document.getElementById("root") as HTMLElement).render(
  <RecoilRoot>
    <Router />
  </RecoilRoot>
);
