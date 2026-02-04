import { useEffect, useState } from "react";
import { useTranslation } from "react-i18next";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import {
  Select,
  SelectContent,
  SelectItem,
  SelectTrigger,
  SelectValue,
} from "@/components/ui/select";
import { model } from "@wailsjs/go/models";

interface Props {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  subscription: model.SubItem | null;
  onSave: (sub: model.SubItem) => void;
}

const AUTO_UPDATE_OPTIONS = [
  { value: "0", label: "Disabled" },
  { value: "30", label: "30 min" },
  { value: "60", label: "1 hour" },
  { value: "360", label: "6 hours" },
  { value: "720", label: "12 hours" },
  { value: "1440", label: "24 hours" },
];

function newSub(): model.SubItem {
  return {
    id: "",
    remarks: "",
    url: "",
    enabled: true,
    sort: 0,
    autoUpdateInterval: 0,
    updateTime: 0,
  } as model.SubItem;
}

export function SubscriptionDialog({ open, onOpenChange, subscription, onSave }: Props) {
  const { t } = useTranslation();
  const [form, setForm] = useState<model.SubItem>(newSub());

  useEffect(() => {
    if (open) {
      setForm(subscription ? { ...subscription } : newSub());
    }
  }, [open, subscription]);

  const update = (partial: Partial<model.SubItem>) =>
    setForm((prev) => ({ ...prev, ...partial }) as model.SubItem);

  const handleSave = () => {
    onSave(form);
    onOpenChange(false);
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>
            {subscription ? t("subscriptions.edit") : t("subscriptions.add")}
          </DialogTitle>
        </DialogHeader>

        <div className="space-y-3">
          <div>
            <Label>{t("subscriptions.form.name")}</Label>
            <Input
              className="mt-1"
              value={form.remarks}
              onChange={(e) => update({ remarks: e.target.value })}
              placeholder="My Subscription"
            />
          </div>
          <div>
            <Label>{t("subscriptions.form.url")}</Label>
            <Input
              className="mt-1"
              value={form.url}
              onChange={(e) => update({ url: e.target.value })}
              placeholder="https://example.com/sub?token=..."
            />
          </div>
          <div>
            <Label>{t("subscriptions.form.autoUpdate")}</Label>
            <Select
              value={String(form.autoUpdateInterval)}
              onValueChange={(v) => update({ autoUpdateInterval: Number(v) })}
            >
              <SelectTrigger className="mt-1">
                <SelectValue />
              </SelectTrigger>
              <SelectContent>
                {AUTO_UPDATE_OPTIONS.map((o) => (
                  <SelectItem key={o.value} value={o.value}>
                    {o.label}
                  </SelectItem>
                ))}
              </SelectContent>
            </Select>
          </div>
          <div>
            <Label>{t("subscriptions.form.filter")}</Label>
            <Input
              className="mt-1"
              value={form.filter ?? ""}
              onChange={(e) => update({ filter: e.target.value })}
              placeholder="regex filter (optional)"
            />
          </div>
          <div>
            <Label>{t("subscriptions.form.userAgent")}</Label>
            <Input
              className="mt-1"
              value={form.userAgent ?? ""}
              onChange={(e) => update({ userAgent: e.target.value })}
              placeholder="(optional)"
            />
          </div>
        </div>

        <DialogFooter className="mt-4">
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            {t("common.cancel")}
          </Button>
          <Button onClick={handleSave} disabled={!form.remarks || !form.url}>
            {t("common.save")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
