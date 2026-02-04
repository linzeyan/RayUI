import { describe, it, expect } from "vitest";
import { formatBytes, formatSpeed, formatDate, protocolName, coreName } from "./format";

describe("formatBytes", () => {
  it("returns '0 B' for zero", () => {
    expect(formatBytes(0)).toBe("0 B");
  });

  it("formats bytes", () => {
    expect(formatBytes(500)).toBe("500 B");
  });

  it("formats kilobytes", () => {
    expect(formatBytes(1024)).toBe("1.00 KB");
    expect(formatBytes(1536)).toBe("1.50 KB");
  });

  it("formats megabytes", () => {
    expect(formatBytes(1048576)).toBe("1.00 MB");
    expect(formatBytes(10485760)).toBe("10.0 MB");
  });

  it("formats gigabytes", () => {
    expect(formatBytes(1073741824)).toBe("1.00 GB");
    expect(formatBytes(107374182400)).toBe("100 GB");
  });
});

describe("formatSpeed", () => {
  it("appends /s to formatted bytes", () => {
    expect(formatSpeed(0)).toBe("0 B/s");
    expect(formatSpeed(1024)).toBe("1.00 KB/s");
    expect(formatSpeed(1048576)).toBe("1.00 MB/s");
  });
});

describe("formatDate", () => {
  it("returns '-' for zero timestamp", () => {
    expect(formatDate(0)).toBe("-");
  });

  it("formats unix timestamp to locale string", () => {
    const result = formatDate(1700000000);
    // Should be a non-empty string
    expect(result).toBeTruthy();
    expect(result).not.toBe("-");
  });
});

describe("protocolName", () => {
  it("returns correct protocol names", () => {
    expect(protocolName(1)).toBe("VMess");
    expect(protocolName(3)).toBe("Shadowsocks");
    expect(protocolName(5)).toBe("VLESS");
    expect(protocolName(6)).toBe("Trojan");
    expect(protocolName(7)).toBe("Hysteria2");
    expect(protocolName(8)).toBe("TUIC");
    expect(protocolName(9)).toBe("WireGuard");
  });

  it("returns 'Unknown' for unrecognized types", () => {
    expect(protocolName(99)).toBe("Unknown");
  });
});

describe("coreName", () => {
  it("returns correct core names", () => {
    expect(coreName(0)).toBe("Auto");
    expect(coreName(1)).toBe("xray");
    expect(coreName(2)).toBe("sing-box");
  });

  it("returns 'Unknown' for unrecognized types", () => {
    expect(coreName(99)).toBe("Unknown");
  });
});
