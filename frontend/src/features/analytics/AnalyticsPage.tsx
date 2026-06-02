import { Bar, BarChart, CartesianGrid, Line, LineChart, ResponsiveContainer, Tooltip, XAxis, YAxis } from "recharts";
import { PageHeader } from "../../components/layout/PageHeader";
import { Card, CardTitle } from "../../components/ui/card";
import { EmptyRow, Table, Td, Th } from "../../components/ui/table";
import { useApi } from "../../hooks/useApi";
import type { CountryStats, TopOffender, TrafficStat } from "../../types/api";

export function AnalyticsPage() {
  const { data: trends } = useApi<TrafficStat[]>("/api/analytics/traffic-trends?period=hour&hours=72", []);
  const { data: countries } = useApi<CountryStats[]>("/api/analytics/countries", []);
  const { data: offenders } = useApi<TopOffender[]>("/api/analytics/top-offenders?limit=10", []);

  return (
    <>
      <PageHeader title="Analytics" subtitle="Traffic trends, geographic distribution, and abusive clients." />
      <div className="grid gap-4 xl:grid-cols-2">
        <Card>
          <CardTitle>Requests and 429s</CardTitle>
          <div className="mt-5 h-[300px]">
            <ResponsiveContainer width="100%" height="100%">
              <LineChart data={trends}>
                <CartesianGrid strokeDasharray="3 3" stroke="hsl(var(--border))" />
                <XAxis dataKey="timestamp" tickFormatter={(v) => new Date(v).getHours().toString()} stroke="hsl(var(--muted-foreground))" fontSize={12} />
                <YAxis stroke="hsl(var(--muted-foreground))" fontSize={12} />
                <Tooltip contentStyle={{ background: "hsl(var(--card))", border: "1px solid hsl(var(--border))", borderRadius: 12 }} />
                <Line type="monotone" dataKey="total_requests" stroke="#3B82F6" strokeWidth={2} dot={false} />
                <Line type="monotone" dataKey="status_429" stroke="#EF4444" strokeWidth={2} dot={false} />
              </LineChart>
            </ResponsiveContainer>
          </div>
        </Card>
        <Card>
          <CardTitle>Country Bans</CardTitle>
          <div className="mt-5 h-[300px]">
            <ResponsiveContainer width="100%" height="100%">
              <BarChart data={countries.slice(0, 8)}>
                <XAxis dataKey="country_code" stroke="hsl(var(--muted-foreground))" fontSize={12} />
                <YAxis stroke="hsl(var(--muted-foreground))" fontSize={12} />
                <Tooltip contentStyle={{ background: "hsl(var(--card))", border: "1px solid hsl(var(--border))", borderRadius: 12 }} />
                <Bar dataKey="bans" fill="#3B82F6" radius={[6, 6, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </Card>
      </div>
      <Card className="mt-4 p-0">
        <div className="border-b p-4 font-medium">Top Offenders</div>
        <div className="overflow-auto">
          <Table>
            <thead>
              <tr>
                <Th>IP</Th>
                <Th>Country</Th>
                <Th>Requests</Th>
                <Th>429 Count</Th>
                <Th>Ban Count</Th>
              </tr>
            </thead>
            <tbody>
              {offenders.length === 0 ? <EmptyRow colSpan={5}>No offenders detected in the current data set.</EmptyRow> : null}
              {offenders.map((offender) => (
                <tr key={offender.ip_address} className="hover:bg-muted/60">
                  <Td className="font-mono">{offender.ip_address}</Td>
                  <Td>{offender.country || "Unknown"}</Td>
                  <Td>{offender.total_requests.toLocaleString()}</Td>
                  <Td>{offender.violation_count.toLocaleString()}</Td>
                  <Td>{offender.ban_count.toLocaleString()}</Td>
                </tr>
              ))}
            </tbody>
          </Table>
        </div>
      </Card>
    </>
  );
}
