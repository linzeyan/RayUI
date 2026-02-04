import { useEffect, useState } from "react";
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
import { Textarea } from "@/components/ui/textarea";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { GetDNSConfig, UpdateDNSConfig } from "@wailsjs/go/app/App";
import { model } from "@wailsjs/go/models";
import { toast } from "sonner";

export function DNSPage() {
  const { t } = useTranslation();
  const [dns, setDns] = useState<model.DNSItem | null>(null);
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    GetDNSConfig().then(setDns);
  }, []);

  const update = (partial: Partial<model.DNSItem>) =>
    setDns((prev) => (prev ? { ...prev, ...partial } as model.DNSItem : prev));

  const handleSave = async () => {
    if (!dns) return;
    setLoading(true);
    try {
      await UpdateDNSConfig(dns);
      toast.success(t("dns.saved"));
    } finally {
      setLoading(false);
    }
  };

  if (!dns) return null;

  return (
    <div className="flex h-full flex-col">
      <div className="flex items-center justify-between border-b border-border px-6 py-3">
        <h1 className="text-xl font-semibold">{t("dns.title")}</h1>
        <Button size="sm" onClick={handleSave} disabled={loading}>
          {t("common.save")}
        </Button>
      </div>

      <div className="flex-1 overflow-auto p-6">
        <Card className="max-w-2xl">
          <CardContent className="space-y-5 pt-6">
            <div>
              <Label>{t("dns.remoteDns")}</Label>
              <Input
                className="mt-1"
                value={dns.remoteDns}
                onChange={(e) => update({ remoteDns: e.target.value })}
                placeholder="https://dns.google/dns-query"
              />
            </div>

            <div>
              <Label>{t("dns.directDns")}</Label>
              <Input
                className="mt-1"
                value={dns.directDns}
                onChange={(e) => update({ directDns: e.target.value })}
                placeholder="https://dns.alidns.com/dns-query"
              />
            </div>

            <div>
              <Label>{t("dns.bootstrapDns")}</Label>
              <Input
                className="mt-1"
                value={dns.bootstrapDns}
                onChange={(e) => update({ bootstrapDns: e.target.value })}
                placeholder="1.1.1.1"
              />
            </div>

            <div>
              <Label>{t("dns.domainStrategy")}</Label>
              <Select
                value={dns.domainStrategy}
                onValueChange={(v) => update({ domainStrategy: v })}
              >
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

            <div className="flex items-center justify-between">
              <Label>{t("dns.fakeIP")}</Label>
              <Switch
                checked={dns.fakeIP}
                onCheckedChange={(v) => update({ fakeIP: v })}
              />
            </div>

            <div className="flex items-center justify-between">
              <Label>{t("dns.useSystemHosts")}</Label>
              <Switch
                checked={dns.useSystemHosts}
                onCheckedChange={(v) => update({ useSystemHosts: v })}
              />
            </div>

            <div>
              <Label>{t("dns.hosts")}</Label>
              <Textarea
                className="mt-1 font-mono text-xs"
                rows={4}
                value={dns.hosts}
                onChange={(e) => update({ hosts: e.target.value })}
                placeholder="example.com 1.2.3.4"
              />
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
