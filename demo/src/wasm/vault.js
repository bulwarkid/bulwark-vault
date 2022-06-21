import { setImmediate } from "../util";

const listeners = [];

async function waitForWasmLoad() {
    if (window.globalThis.vaultInterface) {
        return Promise.resolve();
    }
    return new Promise((resolve) => {
        listeners.push(() => {
            resolve();
        });
    });
}

export function onWasmLoaded() {
    for (const listener of listeners) {
        setImmediate(() => {
            listener();
        });
    }
}

export async function loginToVault(username, password) {
    await waitForWasmLoad();
    await window.globalThis.vaultInterface.login(username, password);
}

export async function getMasterSecret() {
    await waitForWasmLoad();
    return await window.globalThis.vaultInterface.getMasterSecret();
}

export async function getKeyDirectory() {
    await waitForWasmLoad();
    return await window.globalThis.vaultInterface.getKeyDirectory();
}

export async function get(path) {
    await waitForWasmLoad();
    return await window.globalThis.vaultInterface.get(path);
}

export async function put(path, data) {
    await waitForWasmLoad();
    return await window.globalThis.vaultInterface.put(path, data);
}

export async function getAuthData(publicKeyBase64, encryptionKeyBase64) {
    await waitForWasmLoad();
    return await window.globalThis.vaultInterface.getAuthData(
        publicKeyBase64,
        encryptionKeyBase64
    );
}

export async function createAuthData(data) {
    await waitForWasmLoad();
    return await window.globalThis.vaultInterface.createAuthData(data);
}
