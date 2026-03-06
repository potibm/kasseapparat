import { StrictMode } from "react";
import { createRoot } from "react-dom/client";

import "./index.css";
import "./core/theme/theme-init.js";
import App from "./App.jsx";

const rootElement = document.getElementById("root");
if (!rootElement) throw new Error("Failed to find the root element");

createRoot(rootElement).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
