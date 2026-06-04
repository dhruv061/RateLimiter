import { Search } from "lucide-react";
import { useMemo, useState } from "react";
import { PageHeader } from "../../components/layout/PageHeader";
import { Badge } from "../../components/ui/badge";
import { Card } from "../../components/ui/card";
import { Input } from "../../components/ui/input";
import { EmptyRow, Table, Td, Th } from "../../components/ui/table";
import { useApi } from "../../hooks/useApi";
import type { LiveRequest } from "../../types/api";

export function LiveRequestsPage() {
  const [search, setSearch] = useState("");
  const { data } = useApi<LiveRequest[]>("/api/live-requests?limit=100", []);
  const rows = useMemo(
    () =>
      data.filter((item) => {
        const ip = item.ip_address || "";
        const url = item.url || "";
        return ip.includes(search) || url.toLowerCase().includes(search.toLowerCase());
      }),
    [data, search]
  );

  return (
    <>
      <PageHeader title="Live Requests" subtitle="Recent access log activity with rate-limit and ban signals." />
      <Card className="p-0">
        <div className="flex items-center gap-3 border-b p-4">
          <div className="relative w-full max-w-sm">
            <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input className="pl-9" value={search} onChange={(e) => setSearch(e.target.value)} placeholder="Search IP or URL" />
          </div>
          <Badge tone="success">Streaming</Badge>
        </div>
        <div className="overflow-auto">
          <Table>
            <thead>
              <tr>
                <Th>Timestamp</Th>
                <Th>IP</Th>
                <Th>Method</Th>
                <Th>URL</Th>
                <Th>Status</Th>
                <Th>Response</Th>
                <Th>User Agent</Th>
              </tr>
            </thead>
            <tbody>
              {rows.length === 0 ? <EmptyRow colSpan={7}>No recent requests found in the configured access log.</EmptyRow> : null}
              {rows.map((request, index) => (
                <tr key={`${request.timestamp}-${index}`} className="hover:bg-muted/60">
                  <Td>{request.timestamp ? new Date(request.timestamp).toLocaleTimeString() : "-"}</Td>
                  <Td className="font-mono">{request.ip_address || "-"}</Td>
                  <Td>{request.method || "-"}</Td>
                  <Td className="max-w-[260px] truncate">{request.url || "-"}</Td>
                  <Td>
                    <Badge tone={request.status_code === 429 || request.status_code === 403 ? "danger" : "success"}>{request.status_code}</Badge>
                  </Td>
                  <Td>{Number(request.response_time || 0).toFixed(0)}ms</Td>
                  <Td className="max-w-[260px] truncate text-muted-foreground">{request.user_agent || "-"}</Td>
                </tr>
              ))}
            </tbody>
          </Table>
        </div>
      </Card>
    </>
  );
}
