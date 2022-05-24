export async function loginToVault(username, password) {
  await window.globalThis.vaultInterface.login(username, password);
}

export async function getMasterSecret() {
  return await window.globalThis.vaultInterface.getMasterSecret();
}

export async function getKeyDirectory() {
  return await window.globalThis.vaultInterface.getKeyDirectory();
}

export async function get(path) {
  return await window.globalThis.vaultInterface.get(path);
}

export async function put(path, data) {
  return await window.globalThis.vaultInterface.put(path, data);
}
