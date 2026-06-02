import React from "react";
import ReactDOM from "react-dom/client";
import { App } from "./App";
import { GlobalFilterProvider } from "./context/GlobalFilterContext";
import "./styles/globals.css";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <GlobalFilterProvider>
      <App />
    </GlobalFilterProvider>
  </React.StrictMode>
);
