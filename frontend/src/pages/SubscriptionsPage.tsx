import { useEffect, useMemo, useState } from "react";
import { useTranslation } from "react-i18next";
import {
  PlusIcon,
  RefreshCwIcon,
  PencilIcon,
  TrashIcon,
  Loader2Icon,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { cn } from "@/lib/utils";
import { formatDate } from "@/lib/format";
import { useSubscriptionStore } from "@/stores/subscriptionStore";
import { useProfileStore } from "@/stores/profileStore";
import { SubscriptionDialog } from "@/components/common/SubscriptionDialog";
import { model } from "@wailsjs/go/models";
import { toast } from "sonner";

export function SubscriptionsPage() {
  const { t } = useTranslation();
  const {
    subscriptions,
    syncingIds,
    loadSubscriptions,
    addSubscription,
    updateSubscription,
    deleteSubscription,
    syncSubscription,
    syncAllSubscriptions,
  } = useSubscriptionStore();

  const profiles = useProfileStore((s) => s.profiles);
  const loadProfiles = useProfileStore((s) => s.loadProfiles);

  const [dialogOpen, setDialogOpen] = useState(false);
  const [editing, setEditing] = useState<model.SubItem | null>(null);
  const [syncingAll, setSyncingAll] = useState(false);

  useEffect(() => {
    loadSubscriptions();
    loadProfiles();
  }, [loadSubscriptions, loadProfiles]);

  const profileCounts = useMemo(() => {
    const counts = new Map<string, number>();
    profiles.forEach((p) => {
      if (p.subId) {
        counts.set(p.subId, (counts.get(p.subId) || 0) + 1);
      }
    });
    return counts;
  }, [profiles]);

  const handleAdd = () => {
    setEditing(null);
    setDialogOpen(true);
  };

  const handleEdit = (sub: model.SubItem) => {
    setEditing(sub);
    setDialogOpen(true);
  };

  const handleSave = async (sub: model.SubItem) => {
    if (editing) {
      await updateSubscription(sub);
    } else {
      await addSubscription(sub);
    }
  };

  const handleDelete = async (id: string) => {
    await deleteSubscription(id);
    await loadProfiles();
    toast.success(t("common.success"));
  };

  const handleSync = async (id: string, name: string) => {
    try {
      const count = await syncSubscription(id);
      await loadProfiles();
      toast.success(t("subscriptions.syncSuccess", { name, count }));
    } catch {
      toast.error(t("subscriptions.syncFailed", { name }));
    }
  };

  const handleSyncAll = async () => {
    setSyncingAll(true);
    try {
      await syncAllSubscriptions();
      await loadProfiles();
      toast.success(t("common.success"));
    } catch {
      toast.error(t("common.error"));
    } finally {
      setSyncingAll(false);
    }
  };

  return (
    <div className="flex h-full flex-col">
      <div className="flex items-center justify-between border-b border-border px-6 py-3">
        <h1 className="text-xl font-semibold">{t("subscriptions.title")}</h1>
        <div className="flex gap-2">
          <Button size="sm" variant="outline" onClick={handleSyncAll} disabled={syncingAll}>
            {syncingAll ? (
              <Loader2Icon className="mr-1 size-4 animate-spin" />
            ) : (
              <RefreshCwIcon className="mr-1 size-4" />
            )}
            {t("subscriptions.updateAll")}
          </Button>
          <Button size="sm" onClick={handleAdd}>
            <PlusIcon className="mr-1 size-4" />
            {t("subscriptions.add")}
          </Button>
        </div>
      </div>

      <div className="flex-1 overflow-auto p-6">
        {subscriptions.length === 0 ? (
          <div className="flex h-full flex-col items-center justify-center text-muted-foreground">
            <p className="text-lg">{t("subscriptions.noSubscriptions")}</p>
            <p className="text-sm">{t("subscriptions.noSubscriptionsHint")}</p>
          </div>
        ) : (
          <div className="space-y-3">
            {subscriptions.map((sub) => {
              const isSyncing = syncingIds.has(sub.id);
              const count = profileCounts.get(sub.id) || 0;
              return (
                <Card key={sub.id}>
                  <CardContent className="flex items-center justify-between p-4">
                    <div className="min-w-0 flex-1">
                      <h3 className="font-medium">{sub.remarks}</h3>
                      <p className="truncate text-sm text-muted-foreground">
                        {sub.url}
                      </p>
                      <p className="mt-1 text-xs text-muted-foreground">
                        {t("subscriptions.serverCount", { count })}
                        {sub.updateTime
                          ? ` | ${t("subscriptions.lastUpdated", { time: formatDate(sub.updateTime) })}`
                          : ""}
                      </p>
                    </div>
                    <div className="ml-4 flex shrink-0 gap-2">
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => handleSync(sub.id, sub.remarks)}
                        disabled={isSyncing}
                      >
                        {isSyncing ? (
                          <Loader2Icon className="size-4 animate-spin" />
                        ) : (
                          <RefreshCwIcon className="size-4" />
                        )}
                      </Button>
                      <Button
                        size="sm"
                        variant="outline"
                        onClick={() => handleEdit(sub)}
                      >
                        <PencilIcon className="size-4" />
                      </Button>
                      <Button
                        size="sm"
                        variant="outline"
                        className="text-destructive hover:text-destructive"
                        onClick={() => handleDelete(sub.id)}
                      >
                        <TrashIcon className="size-4" />
                      </Button>
                    </div>
                  </CardContent>
                </Card>
              );
            })}
          </div>
        )}
      </div>

      <SubscriptionDialog
        open={dialogOpen}
        onOpenChange={setDialogOpen}
        subscription={editing}
        onSave={handleSave}
      />
    </div>
  );
}
