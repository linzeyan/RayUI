import { describe, it, expect, beforeEach } from "vitest";
import { screen, waitFor } from "@testing-library/react";
import { render, userEvent } from "@/test/test-utils";
import { DNSPage } from "./DNSPage";
import { mockAppApi } from "@/test/wails-mock";

describe("DNSPage", () => {
  const mockDNS = {
    remoteDns: "https://1.1.1.1/dns-query",
    directDns: "https://223.5.5.5/dns-query",
    bootstrapDns: "8.8.8.8",
    useSystemHosts: true,
    fakeIP: false,
    hosts: "",
    domainStrategy: "prefer_ipv4",
  };

  beforeEach(() => {
    mockAppApi.GetDNSConfig.mockResolvedValue(mockDNS);
  });

  it("renders without crashing", async () => {
    render(<DNSPage />);

    await waitFor(() => {
      expect(screen.getByText("DNS")).toBeInTheDocument();
    });
  });

  it("displays DNS configuration fields", async () => {
    render(<DNSPage />);

    await waitFor(() => {
      expect(screen.getByText("Remote DNS")).toBeInTheDocument();
    });

    expect(screen.getByText("Direct DNS")).toBeInTheDocument();
    expect(screen.getByText("Bootstrap DNS")).toBeInTheDocument();
    expect(screen.getByText("FakeIP (sing-box only)")).toBeInTheDocument();
    expect(screen.getByText("Use System Hosts")).toBeInTheDocument();
    expect(screen.getByText("Custom Hosts")).toBeInTheDocument();
  });

  it("loads DNS config from backend", async () => {
    render(<DNSPage />);

    await waitFor(() => {
      expect(mockAppApi.GetDNSConfig).toHaveBeenCalled();
    });

    // Check values are loaded
    const remoteDnsInput = screen.getByDisplayValue("https://1.1.1.1/dns-query");
    expect(remoteDnsInput).toBeInTheDocument();
  });

  it("has a save button", async () => {
    render(<DNSPage />);

    await waitFor(() => {
      expect(screen.getByText("Save")).toBeInTheDocument();
    });
  });

  it("calls UpdateDNSConfig on save", async () => {
    const user = userEvent.setup();
    mockAppApi.UpdateDNSConfig.mockResolvedValue(undefined);

    render(<DNSPage />);

    await waitFor(() => {
      expect(screen.getByText("Save")).toBeInTheDocument();
    });

    await user.click(screen.getByText("Save"));

    await waitFor(() => {
      expect(mockAppApi.UpdateDNSConfig).toHaveBeenCalled();
    });
  });
});
