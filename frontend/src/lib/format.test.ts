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

describe("formatBytes - boundary cases", () => {
  it("handles very small values", () => {
    expect(formatBytes(1)).toBe("1.00 B");
  });

  it("formats terabytes", () => {
    expect(formatBytes(1099511627776)).toContain("TB");
  });

  it("handles NaN gracefully", () => {
    // NaN should return "0 B" since NaN === 0 is false, but Math.log(NaN) is NaN
    const result = formatBytes(NaN);
    expect(typeof result).toBe("string");
  });

  it("handles negative values", () => {
    // Negative values produce NaN in Math.log, should not crash
    const result = formatBytes(-1);
    expect(typeof result).toBe("string");
  });

  it("handles Infinity", () => {
    const result = formatBytes(Infinity);
    expect(typeof result).toBe("string");
  });
});

describe("formatDate - boundary cases", () => {
  it("returns '-' for null/undefined-like timestamps", () => {
    expect(formatDate(0)).toBe("-");
    // @ts-expect-error testing undefined
    expect(formatDate(undefined)).toBe("-");
    // @ts-expect-error testing null
    expect(formatDate(null)).toBe("-");
  });

  it("handles very old timestamps", () => {
    const result = formatDate(1); // Jan 1 1970
    expect(result).toBeTruthy();
    expect(result).not.toBe("-");
  });

  it("handles future timestamps", () => {
    const future = Math.floor(Date.now() / 1000) + 86400 * 365;
    const result = formatDate(future);
    expect(result).toBeTruthy();
    expect(result).not.toBe("-");
  });
});

describe("protocolName - completeness", () => {
  it("covers all protocol types", () => {
    expect(protocolName(4)).toBe("SOCKS");
    expect(protocolName(10)).toBe("HTTP");
  });

  it("returns Unknown for negative values", () => {
    expect(protocolName(-1)).toBe("Unknown");
  });

  it("returns Unknown for 0", () => {
    expect(protocolName(0)).toBe("Unknown");
  });
});
