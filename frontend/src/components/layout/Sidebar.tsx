import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  ServerIcon,
  RssIcon,
  RouteIcon,
  GlobeIcon,
  SettingsIcon,
  ScrollTextIcon,
  PowerIcon,
  Loader2Icon,
} from "lucide-react";
import { cn } from "@/lib/utils";
import { useAppStore, type Page } from "@/stores/appStore";
import { StartCore, StopCore } from "@wailsjs/go/app/App";
import { toast } from "sonner";

const NAV_ITEMS: { key: Page; icon: typeof ServerIcon }[] = [
  { key: "profiles", icon: ServerIcon },
  { key: "subscriptions", icon: RssIcon },
  { key: "routing", icon: RouteIcon },
  { key: "dns", icon: GlobeIcon },
  { key: "settings", icon: SettingsIcon },
  { key: "logs", icon: ScrollTextIcon },
];

export function Sidebar() {
  const { t } = useTranslation();
  const currentPage = useAppStore((s) => s.currentPage);
  const setCurrentPage = useAppStore((s) => s.setCurrentPage);
  const coreStatus = useAppStore((s) => s.coreStatus);
  const loadCoreStatus = useAppStore((s) => s.loadCoreStatus);
  const [toggling, setToggling] = useState(false);

  const handleToggleCore = async () => {
    setToggling(true);
    try {
      if (coreStatus?.running) {
        await StopCore();
      } else {
        await StartCore();
      }
      await loadCoreStatus();
    } catch (err: any) {
      toast.error(err?.message || String(err));
    } finally {
      setToggling(false);
    }
  };

  return (
    <aside className="flex h-full w-[200px] flex-col border-r border-sidebar-border bg-sidebar">
      <div className="flex h-12 items-center px-4 font-semibold text-sidebar-foreground"
        style={{ WebkitAppRegion: "drag" } as React.CSSProperties}
      >
        RayUI
      </div>

      <nav className="flex-1 space-y-0.5 px-2 py-1">
        {NAV_ITEMS.map(({ key, icon: Icon }) => (
          <button
            key={key}
            onClick={() => setCurrentPage(key)}
            className={cn(
              "flex w-full items-center gap-3 rounded-lg px-3 py-2 text-sm transition-colors",
              currentPage === key
                ? "bg-sidebar-accent text-sidebar-accent-foreground font-medium"
                : "text-sidebar-foreground/70 hover:bg-sidebar-accent/50 hover:text-sidebar-foreground",
            )}
          >
            <Icon className="size-4 shrink-0" />
            {t(`nav.${key}`)}
          </button>
        ))}
      </nav>

      <div className="space-y-2 border-t border-sidebar-border px-3 py-3">
        <button
          onClick={handleToggleCore}
          disabled={toggling}
          className={cn(
            "flex w-full items-center justify-center gap-2 rounded-lg px-3 py-2 text-sm font-medium transition-colors",
            coreStatus?.running
              ? "bg-success/10 text-success hover:bg-success/20"
              : "bg-primary/10 text-primary hover:bg-primary/20",
          )}
        >
          {toggling ? (
            <Loader2Icon className="size-4 animate-spin" />
          ) : (
            <PowerIcon className="size-4" />
          )}
          {coreStatus?.running ? t("status.connected") : t("status.disconnected")}
        </button>
        {coreStatus?.running && coreStatus.profile && (
          <p className="truncate text-center text-xs text-sidebar-foreground/60">
            {coreStatus.profile}
          </p>
        )}
      </div>
    </aside>
  );
}
