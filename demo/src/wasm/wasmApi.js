async function request(method, path, data) {
    const args = { method };
    if (data && data !== "") {
        args.body = data;
    }
    const response = await fetch(path, args);
    const blob = await response.blob();
    const text = await blob.text();
    return { data: text, code: response.status };
}

export function setApi() {
    const wasmApi = { request };
    window.globalThis.wasmApi = wasmApi;
}
