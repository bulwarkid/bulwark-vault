export async function loginToVault(username, password) {
  await window.globalThis.vaultInterface.login(username, password);
}

export function getMasterSecret() {
  return window.globalThis.vaultInterface.getMasterSecret();
}
