import { Activity, Database, Gauge, Server, Shield, ShieldAlert, HelpCircle } from "lucide-react";
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
        <Metric
          icon={Gauge}
          label="Requests Today"
          value={stats.total_requests_today}
          detail={`${stats.requests_last_hour.toLocaleString()} last hour`}
          help="Total Nginx HTTP requests logged for the selected domain and date/time range."
        />
        <Metric
          icon={ShieldAlert}
          label="429 Responses"
          value={stats.total_429_today}
          detail={`${stats.count_429_last_hour.toLocaleString()} last hour`}
          tone="warning"
          help="Total rate limit hits (HTTP 429 status code) logged for the selected domain and date/time range."
        />
        <Metric
          icon={Shield}
          label="Active Bans"
          value={stats.active_bans}
          detail={`${stats.bans_24h.toLocaleString()} in 24h`}
          tone="danger"
          help="The count of IP addresses currently blocked by Fail2Ban for the selected domain."
        />
        <Metric
          icon={Activity}
          label="Unbans Today"
          value={stats.unbans_today}
          detail={`${stats.total_bans_today.toLocaleString()} bans today`}
          tone="success"
          help="Total IPs unbanned or whose bans expired today for the selected domain."
        />
      </div>

      <div className="mt-4 grid gap-4 xl:grid-cols-[1.8fr_1fr]">
        <Card className="min-h-[360px]">
          <div className="mb-6 flex items-center justify-between">
            <div>
              <div className="flex items-center gap-1.5">
                <CardTitle>Traffic Trend</CardTitle>
                <span className="group relative cursor-pointer text-muted-foreground hover:text-foreground">
                  <HelpCircle className="h-3.5 w-3.5" />
                  <span className="pointer-events-none absolute bottom-full left-1/2 z-50 mb-2 w-48 -translate-x-1/2 rounded border bg-card p-2 text-xs text-foreground shadow-lg opacity-0 transition-opacity group-hover:opacity-100 leading-normal font-normal">
                    Aggregated request counts over time for the selected domain.
                  </span>
                </span>
              </div>
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
          <div className="mb-5 flex items-center gap-1.5">
            <CardTitle>System Status</CardTitle>
            <span className="group relative cursor-pointer text-muted-foreground hover:text-foreground">
              <HelpCircle className="h-3.5 w-3.5" />
              <span className="pointer-events-none absolute bottom-full left-1/2 z-50 mb-2 w-48 -translate-x-1/2 rounded border bg-card p-2 text-xs text-foreground shadow-lg opacity-0 transition-opacity group-hover:opacity-100 leading-normal font-normal">
                Real-time check on status of local running services.
              </span>
            </span>
          </div>
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

function Metric({
  icon: Icon,
  label,
  value,
  detail,
  tone = "info",
  help
}: {
  icon: typeof Gauge;
  label: string;
  value: number;
  detail: string;
  tone?: "info" | "success" | "warning" | "danger";
  help: string;
}) {
  return (
    <Card>
      <div className="mb-5 flex items-center justify-between">
        <div className="flex items-center gap-1.5">
          <CardTitle>{label}</CardTitle>
          <span className="group relative cursor-pointer text-muted-foreground hover:text-foreground">
            <HelpCircle className="h-3.5 w-3.5" />
            <span className="pointer-events-none absolute bottom-full left-1/2 z-50 mb-2 w-48 -translate-x-1/2 rounded border bg-card p-2 text-xs text-foreground shadow-lg opacity-0 transition-opacity group-hover:opacity-100 leading-normal font-normal">
              {help}
            </span>
          </span>
        </div>
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
