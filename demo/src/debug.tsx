import { showAlert } from "./components/Alert";

export function init() {
  const debug = { showAlert };
  const global = window.globalThis as any;
  global.debug = debug;
}
