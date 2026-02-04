import { useCallback } from "react";
import { useWailsEvent } from "./useWailsEvent";
import { useAppStore } from "@/stores/appStore";

interface TrafficEvent {
  upload: number;
  download: number;
  totalUpload: number;
  totalDownload: number;
}

export function useStats() {
  const setTraffic = useAppStore((s) => s.setTraffic);

  const handleTraffic = useCallback(
    (data: TrafficEvent) => {
      setTraffic(data);
    },
    [setTraffic],
  );

  useWailsEvent("stats:traffic", handleTraffic);
}
