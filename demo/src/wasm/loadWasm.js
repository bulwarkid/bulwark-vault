import { onWasmLoaded } from "./vault";

// This is a polyfill for FireFox and Safari
if (!WebAssembly.instantiateStreaming) {
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}

// Promise to load the wasm file
function loadWasm(path) {
    const go = new window.Go();

    return new Promise((resolve, reject) => {
        const response = fetch(path);
        WebAssembly.instantiateStreaming(response, go.importObject)
            .then((result) => {
                go.run(result.instance);
                resolve(result.instance);
            })
            .catch((error) => {
                reject(error);
            });
    });
}

export function load() {
    console.log("Loading wasm");
    loadWasm("main.wasm")
        .then((wasm) => {
            console.log("Loaded WASM file");
            onWasmLoaded();
        })
        .catch((error) => {
            console.log("Could not load WASM:", error);
        });
}
