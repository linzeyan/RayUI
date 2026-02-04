export function formatBytes(bytes: number): string {
  if (bytes === 0) return "0 B";
  const units = ["B", "KB", "MB", "GB", "TB"];
  const i = Math.floor(Math.log(bytes) / Math.log(1024));
  const value = bytes / Math.pow(1024, i);
  return `${value.toFixed(value >= 100 ? 0 : value >= 10 ? 1 : 2)} ${units[i]}`;
}

export function formatSpeed(bytesPerSec: number): string {
  return `${formatBytes(bytesPerSec)}/s`;
}

export function formatDate(timestamp: number): string {
  if (!timestamp) return "-";
  return new Date(timestamp * 1000).toLocaleString();
}

const PROTOCOL_NAMES: Record<number, string> = {
  1: "VMess",
  3: "Shadowsocks",
  4: "SOCKS",
  5: "VLESS",
  6: "Trojan",
  7: "Hysteria2",
  8: "TUIC",
  9: "WireGuard",
  10: "HTTP",
};

export function protocolName(configType: number): string {
  return PROTOCOL_NAMES[configType] || "Unknown";
}

const CORE_NAMES: Record<number, string> = {
  0: "Auto",
  1: "xray",
  2: "sing-box",
};

export function coreName(coreType: number): string {
  return CORE_NAMES[coreType] || "Unknown";
}
