import { useCallback } from "react";
import { useWailsEvent } from "./useWailsEvent";
import { useAppStore } from "@/stores/appStore";
import { model } from "@wailsjs/go/models";

export function useCoreStatus() {
  const setCoreStatus = useAppStore((s) => s.setCoreStatus);

  const handleStatus = useCallback(
    (data: model.CoreStatus) => {
      setCoreStatus(data);
    },
    [setCoreStatus],
  );

  useWailsEvent("core:status", handleStatus);
}
