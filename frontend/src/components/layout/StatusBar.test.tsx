import { describe, it, expect, beforeEach } from "vitest";
import { screen } from "@testing-library/react";
import { render } from "@/test/test-utils";
import { StatusBar } from "./StatusBar";
import { useAppStore } from "@/stores/appStore";

describe("StatusBar", () => {
  beforeEach(() => {
    useAppStore.setState({
      coreStatus: null,
      traffic: { upload: 0, download: 0, totalUpload: 0, totalDownload: 0 },
    });
  });

  it("renders without crashing", () => {
    render(<StatusBar />);
    expect(screen.getByText("Disconnected")).toBeInTheDocument();
  });

  it("shows disconnected status when core is not running", () => {
    useAppStore.setState({
      coreStatus: { running: false, coreType: 0, version: "1.0.0" },
    });
    render(<StatusBar />);

    expect(screen.getByText("Disconnected")).toBeInTheDocument();
  });

  it("shows connected status when core is running", () => {
    useAppStore.setState({
      coreStatus: { running: true, coreType: 0, version: "1.0.0" },
    });
    render(<StatusBar />);

    expect(screen.getByText("Connected")).toBeInTheDocument();
  });

  it("shows profile name when connected", () => {
    useAppStore.setState({
      coreStatus: {
        running: true,
        coreType: 0,
        version: "1.0.0",
        profile: "My Server",
      },
    });
    render(<StatusBar />);

    expect(screen.getByText("My Server")).toBeInTheDocument();
  });

  it("shows traffic stats when connected", () => {
    useAppStore.setState({
      coreStatus: { running: true, coreType: 0, version: "1.0.0" },
      traffic: {
        upload: 1024,
        download: 2048,
        totalUpload: 0,
        totalDownload: 0,
      },
    });
    render(<StatusBar />);

    // Traffic should be displayed (formatted)
    expect(screen.getByText("1.00 KB/s")).toBeInTheDocument();
    expect(screen.getByText("2.00 KB/s")).toBeInTheDocument();
  });

  it("does not show traffic when disconnected", () => {
    useAppStore.setState({
      coreStatus: { running: false, coreType: 0, version: "1.0.0" },
      traffic: {
        upload: 1024,
        download: 2048,
        totalUpload: 0,
        totalDownload: 0,
      },
    });
    render(<StatusBar />);

    expect(screen.queryByText("1.00 KB/s")).not.toBeInTheDocument();
    expect(screen.queryByText("2.00 KB/s")).not.toBeInTheDocument();
  });
});
