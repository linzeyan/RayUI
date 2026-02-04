import { useState } from "react";
import { useTranslation } from "react-i18next";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogFooter,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { ClipboardPasteIcon } from "lucide-react";
import { ClipboardGetText } from "@wailsjs/runtime/runtime";
import { toast } from "sonner";

interface Props {
  open: boolean;
  onOpenChange: (open: boolean) => void;
  onImport: (text: string) => Promise<number>;
}

export function ImportDialog({ open, onOpenChange, onImport }: Props) {
  const { t } = useTranslation();
  const [text, setText] = useState("");
  const [loading, setLoading] = useState(false);

  const handleImport = async (input: string) => {
    if (!input.trim()) return;
    setLoading(true);
    try {
      const count = await onImport(input);
      toast.success(t("profiles.imported", { count }));
      setText("");
      onOpenChange(false);
    } catch {
      toast.error(t("profiles.importFailed"));
    } finally {
      setLoading(false);
    }
  };

  const handleClipboard = async () => {
    const clip = await ClipboardGetText();
    if (clip) {
      await handleImport(clip);
    }
  };

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-w-md">
        <DialogHeader>
          <DialogTitle>{t("profiles.import")}</DialogTitle>
        </DialogHeader>

        <div className="space-y-3">
          <Button
            variant="outline"
            className="w-full"
            onClick={handleClipboard}
            disabled={loading}
          >
            <ClipboardPasteIcon className="mr-2 size-4" />
            {t("profiles.importFromClipboard")}
          </Button>

          <div className="relative">
            <div className="absolute inset-x-0 top-0 flex items-center justify-center">
              <span className="bg-background px-2 text-xs text-muted-foreground">
                or
              </span>
            </div>
            <div className="border-t border-border" />
          </div>

          <Textarea
            rows={6}
            value={text}
            onChange={(e) => setText(e.target.value)}
            placeholder="vmess://...\nvless://...\nss://..."
            className="font-mono text-xs"
          />
        </div>

        <DialogFooter>
          <Button variant="outline" onClick={() => onOpenChange(false)}>
            {t("common.cancel")}
          </Button>
          <Button
            onClick={() => handleImport(text)}
            disabled={!text.trim() || loading}
          >
            {t("profiles.import")}
          </Button>
        </DialogFooter>
      </DialogContent>
    </Dialog>
  );
}
