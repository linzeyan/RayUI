import { useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import {
  PlusIcon,
  ImportIcon,
  ZapIcon,
  SearchIcon,
  CopyIcon,
  TrashIcon,
  PencilIcon,
  PlayIcon,
  CheckCircleIcon,
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
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import { Checkbox } from "@/components/ui/checkbox";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { Badge } from "@/components/ui/badge";
import { cn } from "@/lib/utils";
import { protocolName } from "@/lib/format";
import { useProfileStore } from "@/stores/profileStore";
import { useSubscriptionStore } from "@/stores/subscriptionStore";
import { useAppStore } from "@/stores/appStore";
import { ProfileDialog } from "@/components/common/ProfileDialog";
import { ImportDialog } from "@/components/common/ImportDialog";
import { model } from "@wailsjs/go/models";
import { toast } from "sonner";
import { ClipboardSetText } from "@wailsjs/runtime/runtime";

type SortKey = "remarks" | "address" | "port" | "configType" | "network";
type SortDir = "asc" | "desc";

// Format bytes/sec to Mbps
function formatSpeed(bytesPerSec: number): string {
  if (bytesPerSec <= 0) return "-";
  const mbps = (bytesPerSec * 8) / (1024 * 1024);
  if (mbps < 1) return `${(mbps * 1000).toFixed(0)} Kbps`;
  return `${mbps.toFixed(1)} Mbps`;
}

export function ProfilesPage() {
  const { t } = useTranslation();
  const {
    profiles,
    speedResults,
    selectedIds,
    filterSubId,
    searchQuery,
    loading,
    loadProfiles,
    deleteProfiles,
    setActiveProfile,
    importFromText,
    exportShareLink,
    testProfiles,
    testAllProfiles,
    setSelectedIds,
    toggleSelected,
    setFilterSubId,
    setSearchQuery,
  } = useProfileStore();

  const { subscriptions, loadSubscriptions } = useSubscriptionStore();
  const config = useAppStore((s) => s.config);

  const [profileDialogOpen, setProfileDialogOpen] = useState(false);
  const [importDialogOpen, setImportDialogOpen] = useState(false);
  const [editingProfile, setEditingProfile] = useState<model.ProfileItem | null>(null);
  const [sortKey, setSortKey] = useState<SortKey>("remarks");
  const [sortDir, setSortDir] = useState<SortDir>("asc");
  const [testing, setTesting] = useState(false);

  useEffect(() => {
    loadProfiles(filterSubId);
    loadSubscriptions();
  }, [loadProfiles, loadSubscriptions, filterSubId]);

  const subMap = useMemo(() => {
    const m = new Map<string, string>();
    subscriptions.forEach((s) => m.set(s.id, s.remarks));
    return m;
  }, [subscriptions]);

  const filtered = useMemo(() => {
    let list = profiles;
    if (searchQuery) {
      const q = searchQuery.toLowerCase();
      list = list.filter(
        (p) =>
          p.remarks.toLowerCase().includes(q) ||
          p.address.toLowerCase().includes(q),
      );
    }
    list = [...list].sort((a, b) => {
      const av = a[sortKey] ?? "";
      const bv = b[sortKey] ?? "";
      const cmp = typeof av === "number"
        ? av - (bv as number)
        : String(av).localeCompare(String(bv));
      return sortDir === "asc" ? cmp : -cmp;
    });
    return list;
  }, [profiles, searchQuery, sortKey, sortDir]);

  const activeProfileId = config?.activeProfileId;

  const toggleSort = (key: SortKey) => {
    if (sortKey === key) {
      setSortDir((d) => (d === "asc" ? "desc" : "asc"));
    } else {
      setSortKey(key);
      setSortDir("asc");
    }
  };

  const allSelected =
    filtered.length > 0 && filtered.every((p) => selectedIds.has(p.id));

  const handleSelectAll = () => {
    if (allSelected) {
      setSelectedIds(new Set());
    } else {
      setSelectedIds(new Set(filtered.map((p) => p.id)));
    }
  };

  const handleEdit = (profile: model.ProfileItem) => {
    setEditingProfile(profile);
    setProfileDialogOpen(true);
  };

  const handleAdd = () => {
    setEditingProfile(null);
    setProfileDialogOpen(true);
  };

  const handleSaveProfile = async (profile: model.ProfileItem) => {
    const store = useProfileStore.getState();
    if (editingProfile) {
      await store.updateProfile(profile);
    } else {
      await store.addProfile(profile);
    }
  };

  const handleDelete = async (ids: string[]) => {
    if (!ids.length) return;
    await deleteProfiles(ids);
    toast.success(t("common.success"));
  };

  const handleSetActive = async (id: string) => {
    await setActiveProfile(id);
    const store = useAppStore.getState();
    await store.loadConfig();
  };

  const handleTestAll = async () => {
    setTesting(true);
    try {
      await testAllProfiles();
    } finally {
      setTesting(false);
    }
  };

  const handleTestSelected = async () => {
    const ids = Array.from(selectedIds);
    if (!ids.length) return;
    setTesting(true);
    try {
      await testProfiles(ids);
    } finally {
      setTesting(false);
    }
  };

  const handleCopyLink = async (id: string) => {
    const link = await exportShareLink(id);
    await ClipboardSetText(link);
    toast.success(t("common.success"));
  };

  const SortIndicator = ({ k }: { k: SortKey }) =>
    sortKey === k ? (
      <span className="ml-1">{sortDir === "asc" ? "↑" : "↓"}</span>
    ) : null;

  return (
    <div className="flex h-full flex-col">
      {/* Toolbar */}
      <div className="flex items-center gap-2 border-b border-border px-4 py-3">
        <div className="relative flex-1">
          <SearchIcon className="absolute left-3 top-1/2 size-4 -translate-y-1/2 text-muted-foreground" />
          <Input
            className="pl-9"
            placeholder={t("profiles.search")}
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
        </div>
        <Select value={filterSubId} onValueChange={setFilterSubId}>
          <SelectTrigger className="w-[160px]">
            <SelectValue placeholder={t("profiles.allGroups")} />
          </SelectTrigger>
          <SelectContent>
            <SelectItem value="all">{t("profiles.allGroups")}</SelectItem>
            <SelectItem value="manual">{t("profiles.manualGroup")}</SelectItem>
            {subscriptions.map((s) => (
              <SelectItem key={s.id} value={s.id}>
                {s.remarks}
              </SelectItem>
            ))}
          </SelectContent>
        </Select>
        <Button size="sm" onClick={handleAdd}>
          <PlusIcon className="mr-1 size-4" />
          {t("profiles.add")}
        </Button>
        <Button size="sm" variant="outline" onClick={() => setImportDialogOpen(true)}>
          <ImportIcon className="mr-1 size-4" />
          {t("profiles.import")}
        </Button>
        <Button
          size="sm"
          variant="outline"
          onClick={handleTestAll}
          disabled={testing || profiles.length === 0}
        >
          <ZapIcon className="mr-1 size-4" />
          {t("profiles.testAll")}
        </Button>
      </div>

      {/* Table */}
      <div className="flex-1 overflow-auto">
        {filtered.length === 0 ? (
          <div className="flex h-full flex-col items-center justify-center text-muted-foreground">
            <p className="text-lg">{t("profiles.noProfiles")}</p>
            <p className="text-sm">{t("profiles.noProfilesHint")}</p>
          </div>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead className="w-10">
                  <Checkbox
                    checked={allSelected}
                    onCheckedChange={handleSelectAll}
                  />
                </TableHead>
                <TableHead className="w-8" />
                <TableHead
                  className="cursor-pointer select-none"
                  onClick={() => toggleSort("remarks")}
                >
                  {t("profiles.columns.remarks")}
                  <SortIndicator k="remarks" />
                </TableHead>
                <TableHead
                  className="cursor-pointer select-none"
                  onClick={() => toggleSort("address")}
                >
                  {t("profiles.columns.address")}
                  <SortIndicator k="address" />
                </TableHead>
                <TableHead
                  className="w-20 cursor-pointer select-none"
                  onClick={() => toggleSort("port")}
                >
                  {t("profiles.columns.port")}
                  <SortIndicator k="port" />
                </TableHead>
                <TableHead
                  className="w-28 cursor-pointer select-none"
                  onClick={() => toggleSort("configType")}
                >
                  {t("profiles.columns.protocol")}
                  <SortIndicator k="configType" />
                </TableHead>
                <TableHead className="w-24">
                  {t("profiles.columns.transport")}
                </TableHead>
                <TableHead className="w-28">
                  {t("profiles.columns.subscription")}
                </TableHead>
                <TableHead className="w-20 text-right">
                  {t("profiles.columns.latency")}
                </TableHead>
                <TableHead className="w-20 text-right">
                  {t("profiles.columns.speed")}
                </TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {filtered.map((p) => {
                const isActive = p.id === activeProfileId;
                const speed = speedResults.get(p.id);
                return (
                  <DropdownMenu key={p.id}>
                    <DropdownMenuTrigger asChild>
                      <TableRow
                        className={cn(
                          "cursor-pointer",
                          isActive && "bg-primary/5",
                          selectedIds.has(p.id) && "bg-accent",
                        )}
                        onDoubleClick={() => handleSetActive(p.id)}
                      >
                        <TableCell onClick={(e) => e.stopPropagation()}>
                          <Checkbox
                            checked={selectedIds.has(p.id)}
                            onCheckedChange={() => toggleSelected(p.id)}
                          />
                        </TableCell>
                        <TableCell>
                          {isActive && (
                            <CheckCircleIcon className="size-4 text-success" />
                          )}
                        </TableCell>
                        <TableCell className="font-medium">{p.remarks}</TableCell>
                        <TableCell className="max-w-[200px] truncate text-muted-foreground">
                          {p.address}
                        </TableCell>
                        <TableCell>{p.port}</TableCell>
                        <TableCell>
                          <Badge variant="secondary" className="font-mono text-xs">
                            {protocolName(p.configType)}
                          </Badge>
                        </TableCell>
                        <TableCell className="text-muted-foreground">
                          {p.network}
                        </TableCell>
                        <TableCell className="text-muted-foreground">
                          {p.subId ? (subMap.get(p.subId) || "-") : "-"}
                        </TableCell>
                        <TableCell className="text-right font-mono">
                          {speed
                            ? speed.latency >= 0
                              ? `${speed.latency}ms`
                              : "timeout"
                            : "-"}
                        </TableCell>
                        <TableCell className="text-right font-mono text-xs">
                          {speed ? formatSpeed(speed.speed || 0) : "-"}
                        </TableCell>
                      </TableRow>
                    </DropdownMenuTrigger>
                    <DropdownMenuContent align="end">
                      <DropdownMenuItem onClick={() => handleSetActive(p.id)}>
                        <PlayIcon className="mr-2 size-4" />
                        {t("profiles.setActive")}
                      </DropdownMenuItem>
                      <DropdownMenuItem onClick={() => handleEdit(p)}>
                        <PencilIcon className="mr-2 size-4" />
                        {t("common.edit")}
                      </DropdownMenuItem>
                      <DropdownMenuItem onClick={() => handleCopyLink(p.id)}>
                        <CopyIcon className="mr-2 size-4" />
                        {t("profiles.exportShareLink")}
                      </DropdownMenuItem>
                      <DropdownMenuSeparator />
                      <DropdownMenuItem
                        className="text-destructive"
                        onClick={() => handleDelete([p.id])}
                      >
                        <TrashIcon className="mr-2 size-4" />
                        {t("common.delete")}
                      </DropdownMenuItem>
                    </DropdownMenuContent>
                  </DropdownMenu>
                );
              })}
            </TableBody>
          </Table>
        )}
      </div>

      {/* Bottom bar */}
      {selectedIds.size > 0 && (
        <div className="flex items-center justify-between border-t border-border px-4 py-2">
          <span className="text-sm text-muted-foreground">
            {selectedIds.size} selected
          </span>
          <div className="flex gap-2">
            <Button
              size="sm"
              variant="outline"
              onClick={handleTestSelected}
              disabled={testing}
            >
              <ZapIcon className="mr-1 size-4" />
              {t("profiles.testSelected")}
            </Button>
            <Button
              size="sm"
              variant="destructive"
              onClick={() => handleDelete(Array.from(selectedIds))}
            >
              <TrashIcon className="mr-1 size-4" />
              {t("profiles.deleteSelected")}
            </Button>
          </div>
        </div>
      )}

      {/* Dialogs */}
      <ProfileDialog
        open={profileDialogOpen}
        onOpenChange={setProfileDialogOpen}
        profile={editingProfile}
        onSave={handleSaveProfile}
      />
      <ImportDialog
        open={importDialogOpen}
        onOpenChange={setImportDialogOpen}
        onImport={importFromText}
      />
    </div>
  );
}
