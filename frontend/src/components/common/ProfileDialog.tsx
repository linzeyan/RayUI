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
import { Switch } from "@/components/ui/switch";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Textarea } from "@/components/ui/textarea";
import { model } from "@wailsjs/go/models";

interface Props {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  profile: model.ProfileItem | null;
  onSave: (profile: model.ProfileItem) => void;
}

const PROTOCOLS = [
  { value: "1", label: "VMess" },
  { value: "5", label: "VLESS" },
  { value: "6", label: "Trojan" },
  { value: "3", label: "Shadowsocks" },
  { value: "7", label: "Hysteria2" },
  { value: "8", label: "TUIC" },
  { value: "9", label: "WireGuard" },
];

const TRANSPORTS = [
  { value: "tcp", label: "TCP" },
  { value: "ws", label: "WebSocket" },
  { value: "h2", label: "HTTP/2" },
  { value: "grpc", label: "gRPC" },
  { value: "quic", label: "QUIC" },
  { value: "kcp", label: "mKCP" },
  { value: "httpupgrade", label: "HTTPUpgrade" },
];

const SECURITY_OPTIONS = [
  { value: "none", label: "None" },
  { value: "tls", label: "TLS" },
  { value: "reality", label: "Reality" },
];

const VMESS_SECURITY = [
  "auto", "aes-128-gcm", "chacha20-poly1305", "none", "zero",
];

const SS_METHODS = [
  "aes-128-gcm", "aes-256-gcm", "chacha20-ietf-poly1305",
  "2022-blake3-aes-128-gcm", "2022-blake3-aes-256-gcm",
  "2022-blake3-chacha20-poly1305",
];

function newProfile(): model.ProfileItem {
  return {
    id: "",
    configType: 5,
    remarks: "",
    subId: "",
    shareUri: "",
    sort: 0,
    address: "",
    port: 443,
    uuid: "",
    security: "auto",
    network: "tcp",
    streamSecurity: "tls",
    allowInsecure: false,
  } as model.ProfileItem;
}

export function ProfileDialog({ open, onOpenChange, profile, onSave }: Props) {
  const { t } = useTranslation();
  const [form, setForm] = useState<model.ProfileItem>(newProfile());

  useEffect(() => {
    if (open) {
      setForm(profile ? { ...profile } : newProfile());
    }
  }, [open, profile]);

  const update = (partial: Partial<model.ProfileItem>) =>
    setForm((prev) => ({ ...prev, ...partial }) as model.ProfileItem);

  const handleSave = () => {
    onSave(form);
    onOpenChange(false);
  };

  const isVMess = form.configType === 1;
  const isSS = form.configType === 3;
  const isReality = form.streamSecurity === "reality";

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-lg">
        <DialogHeader>
          <DialogTitle>
            {profile ? t("profiles.edit") : t("profiles.add")}
          </DialogTitle>
        </DialogHeader>

        <Tabs defaultValue="basic" className="mt-2">
          <TabsList className="w-full">
            <TabsTrigger value="basic" className="flex-1">
              {t("profiles.form.basic")}
            </TabsTrigger>
            <TabsTrigger value="transport" className="flex-1">
              {t("profiles.form.transport")}
            </TabsTrigger>
            <TabsTrigger value="tls" className="flex-1">
              {t("profiles.form.tls")}
            </TabsTrigger>
            <TabsTrigger value="advanced" className="flex-1">
              {t("profiles.form.advanced")}
            </TabsTrigger>
          </TabsList>

          <TabsContent value="basic" className="space-y-3 pt-2">
            <div className="grid grid-cols-2 gap-3">
              <div className="col-span-2">
                <Label>{t("profiles.form.protocol")}</Label>
                <Select
                  value={String(form.configType)}
                  onValueChange={(v) => update({ configType: Number(v) })}
                >
                  <SelectTrigger className="mt-1">
                    <SelectValue />
                  </SelectTrigger>
                  <SelectContent>
                    {PROTOCOLS.map((p) => (
                      <SelectItem key={p.value} value={p.value}>
                        {p.label}
                      </SelectItem>
                    ))}
                  </SelectContent>
                </Select>
              </div>
              <div className="col-span-2">
                <Label>{t("profiles.form.remarks")}</Label>
                <Input
                  className="mt-1"
                  value={form.remarks}
                  onChange={(e) => update({ remarks: e.target.value })}
                  placeholder="My Server"
                />
              </div>
              <div>
                <Label>{t("profiles.form.address")}</Label>
                <Input
                  className="mt-1"
                  value={form.address}
                  onChange={(e) => update({ address: e.target.value })}
                  placeholder="example.com"
                />
              </div>
              <div>
                <Label>{t("profiles.form.port")}</Label>
                <Input
                  className="mt-1"
                  type="number"
                  value={form.port}
                  onChange={(e) => update({ port: Number(e.target.value) })}
                />
              </div>
              <div className="col-span-2">
                <Label>{t("profiles.form.uuid")}</Label>
                <Input
                  className="mt-1"
                  value={form.uuid}
                  onChange={(e) => update({ uuid: e.target.value })}
                  placeholder={isSS ? "password" : "uuid"}
                />
              </div>
              {isVMess && (
                <>
                  <div>
                    <Label>{t("profiles.form.security")}</Label>
                    <Select
                      value={form.security}
                      onValueChange={(v) => update({ security: v })}
                    >
                      <SelectTrigger className="mt-1">
                        <SelectValue />
                      </SelectTrigger>
                      <SelectContent>
                        {VMESS_SECURITY.map((s) => (
                          <SelectItem key={s} value={s}>
                            {s}
                          </SelectItem>
                        ))}
                      </SelectContent>
                    </Select>
                  </div>
                  <div>
                    <Label>{t("profiles.form.alterId")}</Label>
                    <Input
                      className="mt-1"
                      type="number"
                      value={form.alterId ?? 0}
                      onChange={(e) =>
                        update({ alterId: Number(e.target.value) })
                      }
                    />
                  </div>
                </>
              )}
              {isSS && (
                <div className="col-span-2">
                  <Label>{t("profiles.form.security")}</Label>
                  <Select
                    value={form.security}
                    onValueChange={(v) => update({ security: v })}
                  >
                    <SelectTrigger className="mt-1">
                      <SelectValue />
                    </SelectTrigger>
                    <SelectContent>
                      {SS_METHODS.map((m) => (
                        <SelectItem key={m} value={m}>
                          {m}
                        </SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>
              )}
            </div>
          </TabsContent>

          <TabsContent value="transport" className="space-y-3 pt-2">
            <div>
              <Label>{t("profiles.form.network")}</Label>
              <Select
                value={form.network}
                onValueChange={(v) => update({ network: v })}
              >
                <SelectTrigger className="mt-1">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {TRANSPORTS.map((t) => (
                    <SelectItem key={t.value} value={t.value}>
                      {t.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            {(form.network === "ws" ||
              form.network === "h2" ||
              form.network === "httpupgrade") && (
              <>
                <div>
                  <Label>{t("profiles.form.host")}</Label>
                  <Input
                    className="mt-1"
                    value={form.host ?? ""}
                    onChange={(e) => update({ host: e.target.value })}
                  />
                </div>
                <div>
                  <Label>{t("profiles.form.path")}</Label>
                  <Input
                    className="mt-1"
                    value={form.path ?? ""}
                    onChange={(e) => update({ path: e.target.value })}
                    placeholder="/"
                  />
                </div>
              </>
            )}
            {form.network === "grpc" && (
              <div>
                <Label>Service Name</Label>
                <Input
                  className="mt-1"
                  value={form.path ?? ""}
                  onChange={(e) => update({ path: e.target.value })}
                />
              </div>
            )}
          </TabsContent>

          <TabsContent value="tls" className="space-y-3 pt-2">
            <div>
              <Label>{t("profiles.form.streamSecurity")}</Label>
              <Select
                value={form.streamSecurity}
                onValueChange={(v) => update({ streamSecurity: v })}
              >
                <SelectTrigger className="mt-1">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  {SECURITY_OPTIONS.map((s) => (
                    <SelectItem key={s.value} value={s.value}>
                      {s.label}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            {form.streamSecurity !== "none" && (
              <>
                <div>
                  <Label>{t("profiles.form.sni")}</Label>
                  <Input
                    className="mt-1"
                    value={form.sni ?? ""}
                    onChange={(e) => update({ sni: e.target.value })}
                  />
                </div>
                <div>
                  <Label>{t("profiles.form.alpn")}</Label>
                  <Input
                    className="mt-1"
                    value={form.alpn ?? ""}
                    onChange={(e) => update({ alpn: e.target.value })}
                    placeholder="h2,http/1.1"
                  />
                </div>
                <div>
                  <Label>{t("profiles.form.fingerprint")}</Label>
                  <Input
                    className="mt-1"
                    value={form.fingerprint ?? ""}
                    onChange={(e) => update({ fingerprint: e.target.value })}
                    placeholder="chrome"
                  />
                </div>
                <div className="flex items-center gap-2">
                  <Switch
                    checked={form.allowInsecure}
                    onCheckedChange={(v) => update({ allowInsecure: v })}
                  />
                  <Label>{t("profiles.form.allowInsecure")}</Label>
                </div>
              </>
            )}
            {isReality && (
              <>
                <div>
                  <Label>{t("profiles.form.publicKey")}</Label>
                  <Input
                    className="mt-1"
                    value={form.publicKey ?? ""}
                    onChange={(e) => update({ publicKey: e.target.value })}
                  />
                </div>
                <div>
                  <Label>{t("profiles.form.shortId")}</Label>
                  <Input
                    className="mt-1"
                    value={form.shortId ?? ""}
                    onChange={(e) => update({ shortId: e.target.value })}
                  />
                </div>
                <div>
                  <Label>{t("profiles.form.spiderX")}</Label>
                  <Input
                    className="mt-1"
                    value={form.spiderX ?? ""}
                    onChange={(e) => update({ spiderX: e.target.value })}
                  />
                </div>
              </>
            )}
          </TabsContent>

          <TabsContent value="advanced" className="space-y-3 pt-2">
            <div>
              <Label>{t("profiles.form.coreType")}</Label>
              <Select
                value={String(form.coreType ?? 0)}
                onValueChange={(v) => update({ coreType: Number(v) })}
              >
                <SelectTrigger className="mt-1">
                  <SelectValue />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="0">Auto</SelectItem>
                  <SelectItem value="1">xray</SelectItem>
                  <SelectItem value="2">sing-box</SelectItem>
                </SelectContent>
              </Select>
            </div>
            <div className="flex items-center gap-2">
              <Switch
                checked={form.muxEnabled ?? false}
                onCheckedChange={(v) => update({ muxEnabled: v })}
              />
              <Label>{t("profiles.form.mux")}</Label>
            </div>
            <div>
              <Label>Extra JSON</Label>
              <Textarea
                className="mt-1 font-mono text-xs"
                rows={4}
                value={form.extra ?? ""}
                onChange={(e) => update({ extra: e.target.value })}
                placeholder="{}"
              />
            </div>
          </TabsContent>
        </Tabs>

        <DialogFooter className="mt-4">
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            {t("common.cancel")}
          </Button>
          <Button onClick={handleSave}>{t("common.save")}</Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
