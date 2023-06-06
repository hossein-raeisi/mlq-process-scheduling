import React from "react";
import ReactDOM from "react-dom/client";
import "./index.css";
import App from "./App";

console.log(typeof document.getElementById("root"));
ReactDOM.createRoot(document.getElementById("root") as HTMLElement)
    .render(<React.StrictMode><App /></React.StrictMode>);
