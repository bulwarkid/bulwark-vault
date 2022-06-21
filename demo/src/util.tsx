export function b64encode(data: string): string {
    let base64Data = window.btoa(data);
    return base64Data.replace("/", "_").replace("+", "-");
}

export function b64decode(base64Data: string): string {
    base64Data = base64Data.replace("_", "/").replace("-", "+");
    return window.atob(base64Data);
}

export function setImmediate(callback: () => void) {
    setTimeout(callback, 0);
}
