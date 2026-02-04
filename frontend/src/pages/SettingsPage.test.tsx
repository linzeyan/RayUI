import { describe, it, expect, beforeEach } from "vitest";
import { screen, waitFor } from "@testing-library/react";
import { render } from "@/test/test-utils";
import { SettingsPage } from "./SettingsPage";
import { useAppStore } from "@/stores/appStore";
import { mockAppApi, mockConfig, mockSystemInfo } from "@/test/wails-mock";

describe("SettingsPage", () => {
  beforeEach(() => {
    useAppStore.setState({
      config: mockConfig,
    });
    mockAppApi.IsAutoStartEnabled.mockResolvedValue(false);
    mockAppApi.GetSystemInfo.mockResolvedValue(mockSystemInfo);
    mockAppApi.IsCoreInstalled.mockResolvedValue(false);
    mockAppApi.CheckCoreUpdate.mockResolvedValue({
      coreType: 0,
      currentVersion: "1.0.0",
      latestVersion: "1.0.0",
      hasUpdate: false,
      downloadUrl: "",
      assetName: "",
    });
  });

  it("renders without crashing", async () => {
    render(<SettingsPage />);

    await waitFor(() => {
      expect(screen.getByText("Settings")).toBeInTheDocument();
    });
  });

  it("renders all setting sections", async () => {
    render(<SettingsPage />);

    await waitFor(() => {
      expect(screen.getByText("Proxy")).toBeInTheDocument();
    });

    expect(screen.getByText("Core")).toBeInTheDocument();
    expect(screen.getByText("Appearance")).toBeInTheDocument();
    expect(screen.getByText("General")).toBeInTheDocument();
  });

  it("displays proxy settings", async () => {
    render(<SettingsPage />);

    await waitFor(() => {
      expect(screen.getByText("Proxy Mode")).toBeInTheDocument();
    });

    expect(screen.getByText("SOCKS Port")).toBeInTheDocument();
    expect(screen.getByText("HTTP Port")).toBeInTheDocument();
    expect(screen.getByText("Allow LAN")).toBeInTheDocument();
    expect(screen.getByText("UDP Support")).toBeInTheDocument();
  });

  it("displays core settings", async () => {
    render(<SettingsPage />);

    await waitFor(() => {
      expect(screen.getByText("Log Level")).toBeInTheDocument();
    });

    expect(screen.getByText("Mux")).toBeInTheDocument();
    expect(screen.getByText("Allow Insecure TLS")).toBeInTheDocument();
    expect(screen.getByText("TLS Fingerprint")).toBeInTheDocument();
  });

  it("displays appearance settings", async () => {
    render(<SettingsPage />);

    await waitFor(() => {
      expect(screen.getByText("Theme")).toBeInTheDocument();
    });

    expect(screen.getByText("Language")).toBeInTheDocument();
  });

  it("displays general settings", async () => {
    render(<SettingsPage />);

    await waitFor(() => {
      expect(screen.getByText("Start at Login")).toBeInTheDocument();
    });

    expect(screen.getByText("Minimize on Start")).toBeInTheDocument();
    expect(screen.getByText("Close to Tray")).toBeInTheDocument();
  });

  it("displays app version info", async () => {
    render(<SettingsPage />);

    await waitFor(() => {
      expect(screen.getByText(/RayUI.*darwin\/arm64/)).toBeInTheDocument();
    });
  });

  it("renders nothing meaningful when config is not loaded", () => {
    useAppStore.setState({ config: null });
    render(<SettingsPage />);
    // SettingsPage returns null when config is null, so no settings content should appear
    expect(screen.queryByText("Settings")).not.toBeInTheDocument();
    expect(screen.queryByText("Proxy")).not.toBeInTheDocument();
  });
});
