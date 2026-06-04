import React from "react";
import ReactDOM from "react-dom/client";
import { App } from "./App";
import { ErrorBoundary } from "./components/ErrorBoundary";
import { GlobalFilterProvider } from "./context/GlobalFilterContext";
import "./styles/globals.css";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ErrorBoundary>
      <GlobalFilterProvider>
        <App />
      </GlobalFilterProvider>
    </ErrorBoundary>
  </React.StrictMode>
);
