import React from "react";
import ReactDOM from "react-dom/client";
import "./style/tailwind.css";
import "./style/index.css";
import { App } from "./app/App";
import reportWebVitals from "./reportWebVitals";
import * as debug from "./debug";
import * as loadWasm from "./wasm/loadWasm";
import * as wasmApi from "./wasm/wasmApi";

const root = ReactDOM.createRoot(
    document.getElementById("root") as HTMLElement
);
root.render(
    <React.StrictMode>
        <App />
    </React.StrictMode>
);

// If you want to start measuring performance in your app, pass a function
// to log results (for example: reportWebVitals(console.log))
// or send to an analytics endpoint. Learn more: https://bit.ly/CRA-vitals
reportWebVitals();

wasmApi.setApi();
debug.init();

setTimeout(() => {
    loadWasm.load();
}, 0);
