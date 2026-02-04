import { useEffect, useRef, useMemo, useState, useCallback } from "react";
import { useTranslation } from "react-i18next";
import {
  TrashIcon,
  CopyIcon,
  SearchIcon,
  ArrowDownIcon,
  XIcon,
  RefreshCwIcon,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { cn } from "@/lib/utils";
import { useLogStore } from "@/stores/logStore";
import { ClipboardSetText } from "@wailsjs/runtime/runtime";
import {
  GetConnections,
  CloseConnection,
  CloseAllConnections,
} from "@wailsjs/go/app/App";
import { service } from "@wailsjs/go/models";
import { toast } from "sonner";

function logLevel(line: string): string {
  const lower = line.toLowerCase();
  if (lower.includes("[error]") || lower.includes("level=error")) return "error";
  if (lower.includes("[warn") || lower.includes("level=warn")) return "warning";
  if (lower.includes("[debug]") || lower.includes("level=debug")) return "debug";
  return "info";
}

function formatBytes(bytes: number): string {
  if (bytes === 0) return "0 B";
  const k = 1024;
  const sizes = ["B", "KB", "MB", "GB"];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + " " + sizes[i];
}

function LogsTab() {
  const { t } = useTranslation();
  const {
    logs,
    filterLevel,
    searchQuery,
    autoScroll,
    loadLogs,
    clearLogs,
    setFilterLevel,
    setSearchQuery,
    setAutoScroll,
  } = useLogStore();

  const scrollRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    loadLogs();
  }, [loadLogs]);

  useEffect(() => {
    if (autoScroll && scrollRef.current) {
      scrollRef.current.scrollTop = scrollRef.current.scrollHeight;
    }
  }, [logs, autoScroll]);

  const filtered = useMemo(() => {
    let list = logs;
    if (filterLevel !== "all") {
      list = list.filter((l) => logLevel(l) === filterLevel);
    }
    if (searchQuery) {
      const q = searchQuery.toLowerCase();
      list = list.filter((l) => l.toLowerCase().includes(q));
    }
    return list;
  }, [logs, filterLevel, searchQuery]);

  const handleCopy = async () => {
    await ClipboardSetText(filtered.join("\n"));
    toast.success(t("common.success"));
  };

  const handleClear = async () => {
    await clearLogs();
  };

  return (
    <div className="flex h-full flex-col">
      <div className="flex items-center gap-2 border-b border-border px-4 py-3">
        <Select value={filterLevel} onValueChange={setFilterLevel}>
          <SelectTrigger className="w-[120px]">
            <SelectValue />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">{t("logs.level.all")}</SelectItem>
            <SelectItem value="debug">{t("logs.level.debug")}</SelectItem>
            <SelectItem value="info">{t("logs.level.info")}</SelectItem>
            <SelectItem value="warning">{t("logs.level.warning")}</SelectItem>
            <SelectItem value="error">{t("logs.level.error")}</SelectItem>
          </SelectContent>
        </Select>
        <div className="relative flex-1">
          <SearchIcon className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            className="pl-9"
            placeholder={t("logs.search")}
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
        <Button size="sm" variant="outline" onClick={handleCopy}>
          <CopyIcon className="mr-1 size-4" />
          {t("logs.copy")}
        </Button>
        <Button size="sm" variant="outline" onClick={handleClear}>
          <TrashIcon className="mr-1 size-4" />
          {t("logs.clear")}
        </Button>
        <Button
          size="sm"
          variant={autoScroll ? "default" : "outline"}
          onClick={() => setAutoScroll(!autoScroll)}
        >
          <ArrowDownIcon className="size-4" />
        </Button>
      </div>

      <div
        ref={scrollRef}
        className="flex-1 overflow-auto bg-card p-4 font-mono text-xs"
      >
        {filtered.length === 0 ? (
          <div className="flex h-full items-center justify-center text-muted-foreground">
            {t("logs.noLogs")}
          </div>
        ) : (
          <div className="space-y-px">
            {filtered.map((line, i) => {
              const level = logLevel(line);
              return (
                <div
                  key={i}
                  className={cn(
                    "whitespace-pre-wrap break-all py-px",
                    level === "error" && "text-destructive",
                    level === "warning" && "text-warning",
                    level === "debug" && "text-muted-foreground",
                  )}
                >
                  {line}
                </div>
              );
            })}
          </div>
        )}
      </div>
    </div>
  );
}

function ConnectionsTab() {
  const { t } = useTranslation();
  const [connections, setConnections] = useState<service.Connection[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchConnections = useCallback(async () => {
    setLoading(true);
    try {
      const resp = await GetConnections();
      setConnections(resp?.connections || []);
    } catch {
      // Core might not be running, just clear connections
      setConnections([]);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchConnections();
    const interval = setInterval(fetchConnections, 1000);
    return () => clearInterval(interval);
  }, [fetchConnections]);

  const handleCloseConnection = async (id: string) => {
    try {
      await CloseConnection(id);
      setConnections((prev) => prev.filter((c) => c.id !== id));
    } catch (err) {
      toast.error(String(err));
    }
  };

  const handleCloseAll = async () => {
    try {
      await CloseAllConnections();
      setConnections([]);
    } catch (err) {
      toast.error(String(err));
    }
  };

  return (
    <div className="flex h-full flex-col">
      <div className="flex items-center gap-2 border-b border-border px-4 py-3">
        <span className="text-sm text-muted-foreground">
          {connections.length} {t("logs.tabConnections").toLowerCase()}
        </span>
        <div className="flex-1" />
        <Button
          size="sm"
          variant="outline"
          onClick={fetchConnections}
          disabled={loading}
        >
          <RefreshCwIcon className={cn("mr-1 size-4", loading && "animate-spin")} />
        </Button>
        <Button
          size="sm"
          variant="outline"
          onClick={handleCloseAll}
          disabled={connections.length === 0}
        >
          <XIcon className="mr-1 size-4" />
          {t("logs.connections.closeAll")}
        </Button>
      </div>

      <div className="flex-1 overflow-auto">
        {connections.length === 0 ? (
          <div className="flex h-full items-center justify-center text-muted-foreground">
            {t("logs.connections.noConnections")}
          </div>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>{t("logs.connections.host")}</TableHead>
                <TableHead className="w-[80px]">{t("logs.connections.network")}</TableHead>
                <TableHead className="w-[100px]">{t("logs.connections.chains")}</TableHead>
                <TableHead className="w-[120px]">{t("logs.connections.rule")}</TableHead>
                <TableHead className="w-[80px] text-right">{t("logs.connections.upload")}</TableHead>
                <TableHead className="w-[80px] text-right">{t("logs.connections.download")}</TableHead>
                <TableHead className="w-[60px]" />
              </TableRow>
            </TableHeader>
            <TableBody>
              {connections.map((conn) => (
                <TableRow key={conn.id}>
                  <TableCell className="font-mono text-xs">
                    {conn.metadata?.host || conn.metadata?.destinationIP}
                    {conn.metadata?.destinationPort && `:${conn.metadata.destinationPort}`}
                  </TableCell>
                  <TableCell className="text-xs text-muted-foreground">
                    {conn.metadata?.network?.toUpperCase()}
                  </TableCell>
                  <TableCell className="text-xs text-muted-foreground">
                    {conn.chains?.join(" → ")}
                  </TableCell>
                  <TableCell className="text-xs text-muted-foreground">
                    {conn.rule}
                  </TableCell>
                  <TableCell className="text-right text-xs">
                    {formatBytes(conn.upload || 0)}
                  </TableCell>
                  <TableCell className="text-right text-xs">
                    {formatBytes(conn.download || 0)}
                  </TableCell>
                  <TableCell>
                    <Button
                      size="sm"
                      variant="ghost"
                      className="size-6 p-0"
                      onClick={() => handleCloseConnection(conn.id)}
                    >
                      <XIcon className="size-3" />
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </div>
    </div>
  );
}

export function LogsPage() {
  const { t } = useTranslation();
  const [activeTab, setActiveTab] = useState("logs");

  return (
    <Tabs value={activeTab} onValueChange={setActiveTab} className="flex h-full flex-col">
      <div className="flex items-center border-b border-border px-4">
        <TabsList className="h-auto rounded-none border-b-0 bg-transparent p-0">
          <TabsTrigger
            value="logs"
            className="rounded-none border-b-2 border-transparent px-4 py-3 data-[state=active]:border-primary data-[state=active]:bg-transparent"
          >
            {t("logs.tabLogs")}
          </TabsTrigger>
          <TabsTrigger
            value="connections"
            className="rounded-none border-b-2 border-transparent px-4 py-3 data-[state=active]:border-primary data-[state=active]:bg-transparent"
          >
            {t("logs.tabConnections")}
          </TabsTrigger>
        </TabsList>
      </div>
      <TabsContent value="logs" className="mt-0 flex-1 overflow-hidden">
        <LogsTab />
      </TabsContent>
      <TabsContent value="connections" className="mt-0 flex-1 overflow-hidden">
        <ConnectionsTab />
      </TabsContent>
    </Tabs>
  );
}
