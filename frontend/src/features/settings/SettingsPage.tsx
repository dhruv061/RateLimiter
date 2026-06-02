import { CheckCircle2, RotateCcw, Save, Server, Shield } from "lucide-react";
import { useState } from "react";
import { PageHeader } from "../../components/layout/PageHeader";
import { Button } from "../../components/ui/button";
import { Card, CardTitle } from "../../components/ui/card";
import { Input } from "../../components/ui/input";
import { api } from "../../services/api";
import { useApi } from "../../hooks/useApi";

type SettingsMap = Record<string, string>;

export function SettingsPage() {
  const { data, setData } = useApi<SettingsMap>("/api/settings", {});
  const [saving, setSaving] = useState(false);

  function update(key: string, value: string) {
    setData({ ...data, [key]: value });
  }

  async function save() {
    setSaving(true);
    try {
      await api("/api/settings", { method: "PUT", body: JSON.stringify({ settings: data }) });
    } finally {
      setSaving(false);
    }
  }

  return (
    <>
      <PageHeader
        title="Settings"
        subtitle="Rate-limit, Fail2Ban, and dashboard configuration."
        actions={
          <Button variant="primary" onClick={save} disabled={saving}>
            <Save className="h-4 w-4" />
            Save
          </Button>
        }
      />
      <div className="grid gap-4 xl:grid-cols-2">
        <Card>
          <div className="mb-5 flex items-center gap-2">
            <Server className="h-4 w-4 text-primary" />
            <CardTitle>Nginx Rate Limit</CardTitle>
          </div>
          <Field label="Requests Per Second" value={data.nginx_rate_limit_rps ?? ""} onChange={(value) => update("nginx_rate_limit_rps", value)} />
          <Field label="Burst Size" value={data.nginx_rate_limit_burst ?? ""} onChange={(value) => update("nginx_rate_limit_burst", value)} />
          <Field label="Response Code" value={data.nginx_rate_limit_status ?? "429"} onChange={(value) => update("nginx_rate_limit_status", value)} />
        </Card>
        <Card>
          <div className="mb-5 flex items-center gap-2">
            <Shield className="h-4 w-4 text-primary" />
            <CardTitle>Fail2Ban</CardTitle>
          </div>
          <Field label="Ban Time" value={data.fail2ban_ban_time ?? ""} onChange={(value) => update("fail2ban_ban_time", value)} />
          <Field label="Find Time" value={data.fail2ban_find_time ?? ""} onChange={(value) => update("fail2ban_find_time", value)} />
          <Field label="Max Retry" value={data.fail2ban_max_retry ?? ""} onChange={(value) => update("fail2ban_max_retry", value)} />
        </Card>
      </div>
      <Card className="mt-4">
        <CardTitle>Quick Actions</CardTitle>
        <div className="mt-5 flex flex-wrap gap-2">
          <Action path="/api/system/nginx/validate" icon={CheckCircle2} label="Validate Nginx" />
          <Action path="/api/system/nginx/reload" icon={RotateCcw} label="Reload Nginx" />
          <Action path="/api/system/fail2ban/reload" icon={RotateCcw} label="Reload Fail2Ban" />
          <Action path="/api/system/fail2ban/sync-bans" icon={Shield} label="Sync Bans" />
        </div>
      </Card>
    </>
  );
}

function Field({ label, value, onChange }: { label: string; value: string; onChange: (value: string) => void }) {
  return (
    <label className="mb-4 block">
      <span className="mb-2 block text-sm text-muted-foreground">{label}</span>
      <Input value={value} onChange={(e) => onChange(e.target.value)} />
    </label>
  );
}

function Action({ path, icon: Icon, label }: { path: string; icon: typeof Save; label: string }) {
  return (
    <Button variant="secondary" onClick={() => api(path, { method: "POST", body: "{}" })}>
      <Icon className="h-4 w-4" />
      {label}
    </Button>
  );
}
