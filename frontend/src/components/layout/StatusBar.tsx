import { useTranslation } from "react-i18next";
import { ArrowUpIcon, ArrowDownIcon } from "lucide-react";
import { cn } from "@/lib/utils";
import { useAppStore } from "@/stores/appStore";
import { formatSpeed } from "@/lib/format";

export function StatusBar() {
  const { t } = useTranslation();
  const coreStatus = useAppStore((s) => s.coreStatus);
  const traffic = useAppStore((s) => s.traffic);

  return (
    <footer className="flex h-8 items-center justify-between border-t border-border bg-muted/50 px-4 text-xs text-muted-foreground">
      <div className="flex items-center gap-2">
        <span
          className={cn(
            "size-1.5 rounded-full",
            coreStatus?.running ? "bg-success" : "bg-muted-foreground/40",
          )}
        />
        <span>
          {coreStatus?.running
            ? t("status.connected")
            : t("status.disconnected")}
        </span>
        {coreStatus?.running && coreStatus.profile && (
          <>
            <span className="text-border">|</span>
            <span>{coreStatus.profile}</span>
          </>
        )}
      </div>

      {coreStatus?.running && (
        <div className="flex items-center gap-3">
          <span className="flex items-center gap-1">
            <ArrowUpIcon className="size-3" />
            {formatSpeed(traffic.upload)}
          </span>
          <span className="flex items-center gap-1">
            <ArrowDownIcon className="size-3" />
            {formatSpeed(traffic.download)}
          </span>
        </div>
      )}
    </footer>
  );
}
