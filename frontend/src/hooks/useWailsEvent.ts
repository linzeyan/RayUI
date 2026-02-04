import { useEffect } from "react";
import { EventsOn } from "@wailsjs/runtime/runtime";

export function useWailsEvent(
  eventName: string,
  callback: (...data: any[]) => void,
) {
  useEffect(() => {
    const cancel = EventsOn(eventName, callback);
    return cancel;
  }, [eventName, callback]);
}
