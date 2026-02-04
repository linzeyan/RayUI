import { useCallback, useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Separator } from "@/components/ui/separator";
import { useAppStore } from "@/stores/appStore";
import { applyTheme, watchSystemTheme, type Theme } from "@/lib/theme";
import {
  SetAutoStart,
  IsAutoStartEnabled,
  GetSystemInfo,
  CheckCoreUpdate,
  DownloadCoreUpdate,
  IsCoreInstalled,
} from "@wailsjs/go/app/App";
import { app, service } from "@wailsjs/go/models";
import { useWailsEvent } from "@/hooks/useWailsEvent";
import i18n from "@/i18n";
import { toast } from "sonner";

// Core type constants matching Go ECoreType
const CoreSingbox = 2;
const CoreXray = 1;

interface UpdateProgress {
  coreType: number;
  downloaded: number;
  total: number;
  status: string;
  description: string;
}

function CoreUpdateRow({
  coreType,
  coreName,
  currentVersion,
  onInstalled,
}: {
  coreType: number;
  coreName: string;
  currentVersion: string;
  onInstalled?: (version: string) => void;
}) {
  const { t } = useTranslation();
  const [checking, setChecking] = useState(false);
  const [updateInfo, setUpdateInfo] = useState<service.UpdateInfo | null>(null);
  const [downloading, setDownloading] = useState(false);
  const [progress, setProgress] = useState(0);
  const [status, setStatus] = useState("");

  const handleProgress = useCallback(
    (data: UpdateProgress) => {
      if (data.coreType !== coreType) return;
      if (data.total > 0) {
        setProgress(Math.round((data.downloaded / data.total) * 100));
      }
      setStatus(data.status);
    },
    [coreType]
  );

  useWailsEvent("update:progress", handleProgress);

  const handleCheckUpdate = async () => {
    setChecking(true);
    try {
      const info = await CheckCoreUpdate(coreType);
      setUpdateInfo(info);
      if (!info.hasUpdate) {
        toast.success(t("settings.core.upToDate"));
      }
    } catch (err) {
      toast.error(String(err));
    } finally {
      setChecking(false);
    }
  };

  const handleDownload = async () => {
    if (!updateInfo) return;
    setDownloading(true);
    setProgress(0);
    try {
      await DownloadCoreUpdate(updateInfo);
      toast.success(t("settings.core.updateDone", { version: updateInfo.latestVersion }));
      setUpdateInfo(null);
      // Notify parent to refresh displayed version
      onInstalled?.(updateInfo.latestVersion);
    } catch (err) {
      toast.error(t("settings.core.updateFailed"));
    } finally {
      setDownloading(false);
      setProgress(0);
      setStatus("");
    }
  };

  return (
    <div className="flex items-center justify-between">
      <div className="flex items-center gap-2">
        <Label>{coreName}</Label>
        <span className="text-xs text-muted-foreground">
          {currentVersion
            ? t("settings.core.version", { version: currentVersion })
            : t("settings.core.notInstalled")}
        </span>
      </div>
      <div className="flex items-center gap-2">
        {downloading ? (
          <div className="flex items-center gap-2">
            <div className="h-2 w-24 overflow-hidden rounded-full bg-muted">
              <div
                className="h-full bg-primary transition-all"
                style={{ width: `${progress}%` }}
              />
            </div>
            <span className="text-xs text-muted-foreground">
              {status === "extracting"
                ? t("settings.core.extracting")
                : t("settings.core.downloading", { percent: progress })}
            </span>
          </div>
        ) : updateInfo?.hasUpdate ? (
          <Button size="sm" variant="outline" onClick={handleDownload}>
            {t("settings.core.update", { version: updateInfo.latestVersion })}
          </Button>
        ) : (
          <Button
            size="sm"
            variant="ghost"
            onClick={handleCheckUpdate}
            disabled={checking}
          >
            {checking ? t("settings.core.checking") : t("settings.core.checkUpdate")}
          </Button>
        )}
      </div>
    </div>
  );
}

export function SettingsPage() {
  const { t } = useTranslation();
  const config = useAppStore((s) => s.config);
  const updateConfig = useAppStore((s) => s.updateConfig);
  const setProxyMode = useAppStore((s) => s.setProxyMode);

  const [autoStartEnabled, setAutoStartEnabled] = useState(false);
  const [systemInfo, setSystemInfo] = useState<app.SystemInfo | null>(null);
  const [saving, setSaving] = useState(false);
  const [coreVersions, setCoreVersions] = useState<{ singbox: string; xray: string }>({
    singbox: "",
    xray: "",
  });
  const [singboxInstalled, setSingboxInstalled] = useState(false);

  useEffect(() => {
    IsAutoStartEnabled().then(setAutoStartEnabled);
    GetSystemInfo().then(setSystemInfo);
    // Check if sing-box is installed (required for TUN mode)
    IsCoreInstalled(CoreSingbox).then(setSingboxInstalled);
    // Check for stored core versions (from version files)
    CheckCoreUpdate(CoreSingbox)
      .then((info) => setCoreVersions((v) => ({ ...v, singbox: info.currentVersion })))
      .catch(() => {});
    CheckCoreUpdate(CoreXray)
      .then((info) => setCoreVersions((v) => ({ ...v, xray: info.currentVersion })))
      .catch(() => {});
  }, []);

  if (!config) return null;

  const handleAutoStart = async (enabled: boolean) => {
    await SetAutoStart(enabled);
    setAutoStartEnabled(enabled);
  };

  const handleThemeChange = (theme: string) => {
    const t = theme as Theme;
    applyTheme(t);
    watchSystemTheme(t);
    handleSave({ ui: { ...config.ui, theme } });
  };

  const handleLanguageChange = (lang: string) => {
    i18n.changeLanguage(lang);
    handleSave({ ui: { ...config.ui, language: lang } });
  };

  const handleSave = async (partial: Record<string, any>) => {
    setSaving(true);
    try {
      const newConfig = { ...config, ...partial };
      await updateConfig(newConfig as any);
      toast.success(t("settings.saved"));
    } finally {
      setSaving(false);
    }
  };

  const handleProxyModeChange = async (mode: string) => {
    await setProxyMode(Number(mode));
  };

  const updateInbound = (index: number, key: string, value: any) => {
    const inbounds = [...(config.inbounds || [])];
    if (inbounds[index]) {
      inbounds[index] = { ...inbounds[index], [key]: value } as any;
      handleSave({ inbounds });
    }
  };

  const socksInbound = config.inbounds?.[0];
  const httpInbound = config.inbounds?.[1];

  return (
    <div className="flex h-full flex-col">
      <div className="border-b border-border px-6 py-3">
        <h1 className="text-xl font-semibold">{t("settings.title")}</h1>
      </div>

      <div className="flex-1 space-y-6 overflow-auto p-6">
        {/* Proxy */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base">{t("settings.proxy.title")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <Label>{t("settings.proxy.mode")}</Label>
              <Select
                value={String(config.proxyMode)}
                onValueChange={handleProxyModeChange}
              >
                <SelectTrigger className="mt-1 w-[200px]">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="0">{t("settings.proxy.modeManual")}</SelectItem>
                  <SelectItem value="1">{t("settings.proxy.modeSystem")}</SelectItem>
                  <SelectItem value="2" disabled={!singboxInstalled}>
                    {t("settings.proxy.modeTun")}
                    {!singboxInstalled && " (sing-box required)"}
                  </SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="grid grid-cols-2 gap-4">
              <div>
                <Label>{t("settings.proxy.socksPort")}</Label>
                <Input
                  className="mt-1"
                  type="number"
                  value={socksInbound?.port ?? 10808}
                  onChange={(e) => updateInbound(0, "port", Number(e.target.value))}
                />
              </div>
              <div>
                <Label>{t("settings.proxy.httpPort")}</Label>
                <Input
                  className="mt-1"
                  type="number"
                  value={httpInbound?.port ?? 10809}
                  onChange={(e) => updateInbound(1, "port", Number(e.target.value))}
                />
              </div>
            </div>
            <div className="flex items-center justify-between">
              <Label>{t("settings.proxy.allowLAN")}</Label>
              <Switch
                checked={socksInbound?.allowLAN ?? false}
                onCheckedChange={(v) => updateInbound(0, "allowLAN", v)}
              />
            </div>
            <div className="flex items-center justify-between">
              <Label>{t("settings.proxy.udp")}</Label>
              <Switch
                checked={socksInbound?.udpEnabled ?? true}
                onCheckedChange={(v) => updateInbound(0, "udpEnabled", v)}
              />
            </div>
            <div className="flex items-center justify-between">
              <Label>{t("settings.proxy.sniffing")}</Label>
              <Switch
                checked={socksInbound?.sniffingEnabled ?? true}
                onCheckedChange={(v) => updateInbound(0, "sniffingEnabled", v)}
              />
            </div>
          </CardContent>
        </Card>

        {/* Core */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base">{t("settings.core.title")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <CoreUpdateRow
              coreType={CoreSingbox}
              coreName={t("settings.core.singbox")}
              currentVersion={coreVersions.singbox}
              onInstalled={(v) => {
                setCoreVersions((prev) => ({ ...prev, singbox: v }));
                setSingboxInstalled(true);
              }}
            />
            <CoreUpdateRow
              coreType={CoreXray}
              coreName={t("settings.core.xray")}
              currentVersion={coreVersions.xray}
              onInstalled={(v) => setCoreVersions((prev) => ({ ...prev, xray: v }))}
            />
            <Separator />
            <div>
              <Label>{t("settings.core.logLevel")}</Label>
              <Select
                value={config.coreBasic?.logLevel ?? "info"}
                onValueChange={(v) =>
                  handleSave({ coreBasic: { ...config.coreBasic, logLevel: v } })
                }
              >
                <SelectTrigger className="mt-1 w-[200px]">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="debug">Debug</SelectItem>
                  <SelectItem value="info">Info</SelectItem>
                  <SelectItem value="warning">Warning</SelectItem>
                  <SelectItem value="error">Error</SelectItem>
                  <SelectItem value="none">None</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="flex items-center justify-between">
              <Label>{t("settings.core.mux")}</Label>
              <Switch
                checked={config.coreBasic?.muxEnabled ?? false}
                onCheckedChange={(v) =>
                  handleSave({ coreBasic: { ...config.coreBasic, muxEnabled: v } })
                }
              />
            </div>
            <div className="flex items-center justify-between">
              <Label>{t("settings.core.allowInsecure")}</Label>
              <Switch
                checked={config.coreBasic?.allowInsecure ?? false}
                onCheckedChange={(v) =>
                  handleSave({ coreBasic: { ...config.coreBasic, allowInsecure: v } })
                }
              />
            </div>
            <div>
              <Label>{t("settings.core.fingerprint")}</Label>
              <Select
                value={config.coreBasic?.fingerprint ?? "chrome"}
                onValueChange={(v) =>
                  handleSave({ coreBasic: { ...config.coreBasic, fingerprint: v } })
                }
              >
                <SelectTrigger className="mt-1 w-[200px]">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="chrome">Chrome</SelectItem>
                  <SelectItem value="firefox">Firefox</SelectItem>
                  <SelectItem value="safari">Safari</SelectItem>
                  <SelectItem value="edge">Edge</SelectItem>
                  <SelectItem value="random">Random</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </CardContent>
        </Card>

        {/* Appearance */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base">{t("settings.appearance.title")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <Label>{t("settings.appearance.theme")}</Label>
              <Select
                value={config.ui?.theme ?? "system"}
                onValueChange={handleThemeChange}
              >
                <SelectTrigger className="mt-1 w-[200px]">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="system">{t("settings.appearance.themeSystem")}</SelectItem>
                  <SelectItem value="light">{t("settings.appearance.themeLight")}</SelectItem>
                  <SelectItem value="dark">{t("settings.appearance.themeDark")}</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div>
              <Label>{t("settings.appearance.language")}</Label>
              <Select
                value={config.ui?.language ?? "en"}
                onValueChange={handleLanguageChange}
              >
                <SelectTrigger className="mt-1 w-[200px]">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="en">English</SelectItem>
                  <SelectItem value="zh-TW">正體中文</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </CardContent>
        </Card>

        {/* General */}
        <Card>
          <CardHeader>
            <CardTitle className="text-base">{t("settings.general.title")}</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="flex items-center justify-between">
              <Label>{t("settings.general.autoStart")}</Label>
              <Switch
                checked={autoStartEnabled}
                onCheckedChange={handleAutoStart}
              />
            </div>
            <div className="flex items-center justify-between">
              <Label>{t("settings.general.autoHide")}</Label>
              <Switch
                checked={config.ui?.autoHideOnStart ?? false}
                onCheckedChange={(v) =>
                  handleSave({ ui: { ...config.ui, autoHideOnStart: v } })
                }
              />
            </div>
            <div className="flex items-center justify-between">
              <Label>{t("settings.general.closeToTray")}</Label>
              <Switch
                checked={config.ui?.closeToTray ?? true}
                onCheckedChange={(v) =>
                  handleSave({ ui: { ...config.ui, closeToTray: v } })
                }
              />
            </div>
            {systemInfo?.os === "darwin" && (
              <div className="flex items-center justify-between">
                <Label>{t("settings.general.showInDock")}</Label>
                <Switch
                  checked={config.ui?.showInDock ?? true}
                  onCheckedChange={(v) =>
                    handleSave({ ui: { ...config.ui, showInDock: v } })
                  }
                />
              </div>
            )}
          </CardContent>
        </Card>

        {/* App info */}
        {systemInfo && (
          <p className="pb-4 text-center text-xs text-muted-foreground">
            RayUI {systemInfo.appVersion} · {systemInfo.os}/{systemInfo.arch}
          </p>
        )}
      </div>
    </div>
  );
}
