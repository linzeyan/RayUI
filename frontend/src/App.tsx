import { useEffect } from "react";
import { Layout } from "@/components/layout/Layout";
import { useAppStore } from "@/stores/appStore";
import { useStats } from "@/hooks/useStats";
import { useCoreStatus } from "@/hooks/useCoreStatus";
import { useCoreLog } from "@/hooks/useCoreLog";
import { applyTheme, watchSystemTheme } from "@/lib/theme";
import i18n from "@/i18n";
import { ProfilesPage } from "@/pages/ProfilesPage";
import { SubscriptionsPage } from "@/pages/SubscriptionsPage";
import { RoutingPage } from "@/pages/RoutingPage";
import { DNSPage } from "@/pages/DNSPage";
import { SettingsPage } from "@/pages/SettingsPage";
import { LogsPage } from "@/pages/LogsPage";
import { Toaster } from "@/components/ui/sonner";

const PAGES = {
  profiles: ProfilesPage,
  subscriptions: SubscriptionsPage,
  routing: RoutingPage,
  dns: DNSPage,
  settings: SettingsPage,
  logs: LogsPage,
} as const;

function App() {
  const currentPage = useAppStore((s) => s.currentPage);
  const config = useAppStore((s) => s.config);
  const loadConfig = useAppStore((s) => s.loadConfig);
  const loadCoreStatus = useAppStore((s) => s.loadCoreStatus);

  // Subscribe to Wails events
  useStats();
  useCoreStatus();
  useCoreLog();

  // Initial data load
  useEffect(() => {
    loadConfig();
    loadCoreStatus();
  }, [loadConfig, loadCoreStatus]);

  // Apply theme and language from config
  useEffect(() => {
    if (!config) return;
    const theme = (config.ui?.theme || "system") as "system" | "light" | "dark";
    applyTheme(theme);
    const cleanup = watchSystemTheme(theme);

    const lang = config.ui?.language || "en";
    if (i18n.language !== lang) {
      i18n.changeLanguage(lang);
    }

    return cleanup;
  }, [config]);

  const PageComponent = PAGES[currentPage];

  return (
    <Layout>
      <PageComponent />
      <Toaster position="top-right" />
    </Layout>
  );
}

export default App;
