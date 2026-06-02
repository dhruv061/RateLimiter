import { Activity, Database, Gauge, Server, Shield, ShieldAlert } from "lucide-react";
import { Area, AreaChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import { PageHeader } from "../../components/layout/PageHeader";
import { AnimatedNumber } from "../../components/ui/animated-number";
import { Badge } from "../../components/ui/badge";
import { Card, CardTitle } from "../../components/ui/card";
import { useApi } from "../../hooks/useApi";
import type { DashboardStats, TrafficStat } from "../../types/api";

const emptyStats: DashboardStats = {
  total_requests_today: 0,
  requests_last_hour: 0,
  current_rps: 0,
  peak_rps: 0,
  total_429_today: 0,
  count_429_last_hour: 0,
  current_429_rate: 0,
  total_bans_today: 0,
  active_bans: 0,
  bans_24h: 0,
  unbans_today: 0,
  nginx_status: "unknown",
  fail2ban_status: "unknown",
  database_status: "unknown",
  service_status: "unknown"
};

export function DashboardPage() {
  const { data: stats } = useApi<DashboardStats>("/api/dashboard/stats", emptyStats);
  const { data: trends } = useApi<TrafficStat[]>("/api/analytics/traffic-trends?period=hour&hours=24", []);

  return (
    <>
      <PageHeader title="Security Overview" subtitle="Nginx rate limiting, bans, and service health at a glance." />
      <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
        <Metric icon={Gauge} label="Requests Today" value={stats.total_requests_today} detail={`${stats.requests_last_hour.toLocaleString()} last hour`} />
        <Metric icon={ShieldAlert} label="429 Responses" value={stats.total_429_today} detail={`${stats.count_429_last_hour.toLocaleString()} last hour`} tone="warning" />
        <Metric icon={Shield} label="Active Bans" value={stats.active_bans} detail={`${stats.bans_24h.toLocaleString()} in 24h`} tone="danger" />
        <Metric icon={Activity} label="Unbans Today" value={stats.unbans_today} detail={`${stats.total_bans_today.toLocaleString()} bans today`} tone="success" />
      </div>

      <div className="mt-4 grid gap-4 xl:grid-cols-[1.8fr_1fr]">
        <Card className="min-h-[360px]">
          <div className="mb-6 flex items-center justify-between">
            <div>
              <CardTitle>Traffic Trend</CardTitle>
              <div className="mt-1 text-xl font-semibold">Last 24 hours</div>
            </div>
            <Badge tone="info">Hourly</Badge>
          </div>
          <div className="h-[270px]">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={trends}>
                <defs>
                  <linearGradient id="requests" x1="0" x2="0" y1="0" y2="1">
                    <stop offset="0%" stopColor="#3B82F6" stopOpacity={0.45} />
                    <stop offset="100%" stopColor="#3B82F6" stopOpacity={0.04} />
                  </linearGradient>
                </defs>
                <XAxis dataKey="timestamp" tickFormatter={(v) => new Date(v).getHours().toString()} stroke="hsl(var(--muted-foreground))" fontSize={12} />
                <YAxis stroke="hsl(var(--muted-foreground))" fontSize={12} />
                <Tooltip contentStyle={{ background: "hsl(var(--card))", border: "1px solid hsl(var(--border))", borderRadius: 12 }} />
                <Area dataKey="total_requests" stroke="#3B82F6" fill="url(#requests)" strokeWidth={2} />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </Card>
        <Card>
          <CardTitle>System Status</CardTitle>
          <div className="mt-5 space-y-4">
            <Status icon={Server} label="Nginx" value={stats.nginx_status} />
            <Status icon={Shield} label="Fail2Ban" value={stats.fail2ban_status} />
            <Status icon={Database} label="SQLite" value={stats.database_status} />
            <Status icon={Activity} label="Dashboard" value={stats.service_status} />
          </div>
        </Card>
      </div>
    </>
  );
}

function Metric({ icon: Icon, label, value, detail, tone = "info" }: { icon: typeof Gauge; label: string; value: number; detail: string; tone?: "info" | "success" | "warning" | "danger" }) {
  return (
    <Card>
      <div className="mb-5 flex items-center justify-between">
        <CardTitle>{label}</CardTitle>
        <Badge tone={tone}>
          <Icon className="h-3 w-3" />
        </Badge>
      </div>
      <div className="text-[32px] font-semibold leading-none">
        <AnimatedNumber value={value} />
      </div>
      <div className="mt-3 text-sm text-muted-foreground">{detail}</div>
    </Card>
  );
}

function Status({ icon: Icon, label, value }: { icon: typeof Server; label: string; value: string }) {
  const running = value === "running";
  return (
    <div className="flex items-center justify-between rounded-md border p-3">
      <div className="flex items-center gap-3">
        <Icon className="h-4 w-4 text-muted-foreground" />
        <span className="font-medium">{label}</span>
      </div>
      <Badge tone={running ? "success" : "warning"}>{value}</Badge>
    </div>
  );
}
