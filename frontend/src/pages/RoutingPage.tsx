import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import {
  PlusIcon,
  PencilIcon,
  TrashIcon,
  CheckCircleIcon,
  LockIcon,
  ShieldIcon,
} from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import { Badge } from "@/components/ui/badge";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { Separator } from "@/components/ui/separator";
import { cn } from "@/lib/utils";
import { useRoutingStore } from "@/stores/routingStore";
import { useAppStore } from "@/stores/appStore";
import { model } from "@wailsjs/go/models";
import { toast } from "sonner";

export function RoutingPage() {
  const { t } = useTranslation();
  const { routings, loadRoutings, addRouting, updateRouting, deleteRouting, setActiveRouting } =
    useRoutingStore();
  const config = useAppStore((s) => s.config);
  const loadConfig = useAppStore((s) => s.loadConfig);

  const [dialogOpen, setDialogOpen] = useState(false);
  const [editing, setEditing] = useState<model.RoutingItem | null>(null);
  const [formName, setFormName] = useState("");
  const [formStrategy, setFormStrategy] = useState("prefer_ipv4");

  useEffect(() => {
    loadRoutings();
  }, [loadRoutings]);

  const activeRoutingId = config?.activeRoutingId;

  const handleSetActive = async (id: string) => {
    await setActiveRouting(id);
    await loadConfig();
  };

  const handleAdd = () => {
    setEditing(null);
    setFormName("");
    setFormStrategy("prefer_ipv4");
    setDialogOpen(true);
  };

  const handleEdit = (routing: model.RoutingItem) => {
    setEditing(routing);
    setFormName(routing.remarks);
    setFormStrategy(routing.domainStrategy);
    setDialogOpen(true);
  };

  const handleSave = async () => {
    if (editing) {
      await updateRouting({
        ...editing,
        remarks: formName,
        domainStrategy: formStrategy,
      } as model.RoutingItem);
    } else {
      await addRouting({
        id: "",
        remarks: formName,
        domainStrategy: formStrategy,
        rules: [],
        enabled: true,
        locked: false,
        sort: 0,
      } as model.RoutingItem);
    }
    setDialogOpen(false);
  };

  const handleDelete = async (id: string) => {
    await deleteRouting(id);
    toast.success(t("common.success"));
  };

  return (
    <div className="flex h-full flex-col">
      <div className="flex items-center justify-between border-b border-border px-6 py-3">
        <h1 className="text-xl font-semibold">{t("routing.title")}</h1>
        <Button size="sm" onClick={handleAdd}>
          <PlusIcon className="mr-1 size-4" />
          {t("routing.add")}
        </Button>
      </div>

      <div className="flex-1 overflow-auto p-6">
        {routings.length === 0 ? (
          <div className="flex h-full flex-col items-center justify-center text-muted-foreground">
            <p className="text-lg">{t("routing.noRouting")}</p>
          </div>
        ) : (
          <div className="space-y-3">
            {routings.map((r) => {
              const isActive = r.id === activeRoutingId;
              return (
                <Card
                  key={r.id}
                  className={cn(
                    "cursor-pointer transition-colors hover:bg-accent/50",
                    isActive && "ring-2 ring-primary",
                  )}
                  onClick={() => handleSetActive(r.id)}
                >
                  <CardContent className="flex items-center justify-between p-4">
                    <div className="flex items-center gap-3">
                      {isActive ? (
                        <CheckCircleIcon className="size-5 text-primary" />
                      ) : (
                        <ShieldIcon className="size-5 text-muted-foreground" />
                      )}
                      <div>
                        <div className="flex items-center gap-2">
                          <span className="font-medium">{r.remarks}</span>
                          {r.locked && (
                            <Badge variant="secondary" className="text-xs">
                              <LockIcon className="mr-1 size-3" />
                              {t("routing.builtIn")}
                            </Badge>
                          )}
                        </div>
                        <p className="text-sm text-muted-foreground">
                          {t("routing.rules", { count: r.rules?.length || 0 })}
                          {" · "}
                          {r.domainStrategy}
                        </p>
                      </div>
                    </div>
                    {!r.locked && (
                      <div className="flex gap-2" onClick={(e) => e.stopPropagation()}>
                        <Button
                          size="sm"
                          variant="outline"
                          onClick={() => handleEdit(r)}
                        >
                          <PencilIcon className="size-4" />
                        </Button>
                        <Button
                          size="sm"
                          variant="outline"
                          className="text-destructive hover:text-destructive"
                          onClick={() => handleDelete(r.id)}
                        >
                          <TrashIcon className="size-4" />
                        </Button>
                      </div>
                    )}
                  </CardContent>
                </Card>
              );
            })}
          </div>
        )}
      </div>

      <Dialog open={dialogOpen} onOpenChange={setDialogOpen}>
        <DialogContent className="max-w-sm">
          <DialogHeader>
            <DialogTitle>
              {editing ? t("routing.edit") : t("routing.add")}
            </DialogTitle>
          </DialogHeader>
          <div className="space-y-3">
            <div>
              <Label>{t("routing.form.name")}</Label>
              <Input
                className="mt-1"
                value={formName}
                onChange={(e) => setFormName(e.target.value)}
                placeholder="My Routing"
              />
            </div>
            <div>
              <Label>{t("routing.form.domainStrategy")}</Label>
              <Select value={formStrategy} onValueChange={setFormStrategy}>
                <SelectTrigger className="mt-1">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="prefer_ipv4">prefer_ipv4</SelectItem>
                  <SelectItem value="prefer_ipv6">prefer_ipv6</SelectItem>
                  <SelectItem value="ipv4_only">ipv4_only</SelectItem>
                  <SelectItem value="ipv6_only">ipv6_only</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>
          <DialogFooter className="mt-4">
            <Button variant="outline" onClick={() => setDialogOpen(false)}>
              {t("common.cancel")}
            </Button>
            <Button onClick={handleSave} disabled={!formName}>
              {t("common.save")}
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </div>
  );
}
