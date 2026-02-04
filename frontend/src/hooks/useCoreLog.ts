import { useCallback } from "react";
import { useWailsEvent } from "./useWailsEvent";
import { useLogStore } from "@/stores/logStore";

export function useCoreLog() {
  const addLog = useLogStore((s) => s.addLog);

  const handleLog = useCallback(
    (line: string) => {
      addLog(line);
    },
    [addLog],
  );

  useWailsEvent("core:log", handleLog);
}
