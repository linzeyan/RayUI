export type Theme = "system" | "light" | "dark";

function getSystemTheme(): "light" | "dark" {
  return window.matchMedia("(prefers-color-scheme: dark)").matches
    ? "dark"
    : "light";
}

export function applyTheme(theme: Theme) {
  const resolved = theme === "system" ? getSystemTheme() : theme;
  document.documentElement.classList.toggle("dark", resolved === "dark");
}

export function watchSystemTheme(theme: Theme) {
  if (theme !== "system") return () => {};
  const mq = window.matchMedia("(prefers-color-scheme: dark)");
  const handler = () => applyTheme("system");
  mq.addEventListener("change", handler);
  return () => mq.removeEventListener("change", handler);
}
