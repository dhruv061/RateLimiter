import { Download, Eye, Search, ShieldCheck, Unlock } from "lucide-react";
import { useMemo, useState } from "react";
import { PageHeader } from "../../components/layout/PageHeader";
import { Badge } from "../../components/ui/badge";
import { Button } from "../../components/ui/button";
import { Card } from "../../components/ui/card";
import { Drawer } from "../../components/ui/drawer";
import { Input } from "../../components/ui/input";
import { EmptyRow, Table, Td, Th } from "../../components/ui/table";
import { api } from "../../services/api";
import { useApi } from "../../hooks/useApi";
import type { Ban, PaginatedResponse } from "../../types/api";

const emptyPage: PaginatedResponse<Ban> = { data: [], total: 0, page: 1, per_page: 20, total_pages: 0 };

export function BansPage({ history = false }: { history?: boolean }) {
  const [search, setSearch] = useState("");
  const [selected, setSelected] = useState<Ban | null>(null);
  const endpoint = history ? "/api/bans/history?per_page=50" : "/api/bans/active?per_page=50";
  const { data, setData } = useApi<PaginatedResponse<Ban>>(endpoint, emptyPage);

  const rows = useMemo(() => data.data.filter((ban) => ban.ip_address.includes(search) || ban.country.toLowerCase().includes(search.toLowerCase())), [data.data, search]);

  async function unban(ban: Ban) {
    await api(`/api/bans/${ban.id}/unban`, { method: "POST", body: "{}" });
    setData({ ...data, data: data.data.filter((item) => item.id !== ban.id) });
  }

  return (
    <>
      <PageHeader
        title={history ? "Ban History" : "Active Bans"}
        subtitle={history ? "Historical ban and unban records." : "Currently banned IPs with quick remediation actions."}
        actions={
          <>
            <Button variant="secondary">
              <Download className="h-4 w-4" />
              Export
            </Button>
          </>
        }
      />
      <Card className="p-0">
        <div className="flex items-center gap-3 border-b p-4">
          <div className="relative w-full max-w-sm">
            <Search className="pointer-events-none absolute left-3 top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground" />
            <Input className="pl-9" value={search} onChange={(e) => setSearch(e.target.value)} placeholder="Search IP or country" />
          </div>
          <Badge tone="info">{data.total.toLocaleString()} records</Badge>
        </div>
        <div className="overflow-auto">
          <Table>
            <thead>
              <tr>
                <Th>IP Address</Th>
                <Th>Country</Th>
                <Th>{history ? "Ban Time" : "Remaining"}</Th>
                <Th>Reason</Th>
                <Th>Requests</Th>
                <Th>Actions</Th>
              </tr>
            </thead>
            <tbody>
              {rows.length === 0 ? <EmptyRow colSpan={6}>{history ? "No historical bans match the current filters." : "No active bans detected. Traffic appears normal."}</EmptyRow> : null}
              {rows.map((ban) => (
                <tr key={ban.id} className="transition-colors hover:bg-muted/60">
                  <Td className="font-mono text-sm">{ban.ip_address}</Td>
                  <Td>{ban.country || "Unknown"}</Td>
                  <Td>{history ? new Date(ban.ban_time).toLocaleString() : formatRemaining(ban)}</Td>
                  <Td className="max-w-[300px] truncate">{ban.reason}</Td>
                  <Td>{ban.request_count.toLocaleString()}</Td>
                  <Td>
                    <div className="flex gap-2">
                      <Button size="icon" variant="ghost" onClick={() => setSelected(ban)} aria-label="View ban details">
                        <Eye className="h-4 w-4" />
                      </Button>
                      {!history ? (
                        <Button size="icon" variant="danger" onClick={() => unban(ban)} aria-label="Unban IP">
                          <Unlock className="h-4 w-4" />
                        </Button>
                      ) : null}
                    </div>
                  </Td>
                </tr>
              ))}
            </tbody>
          </Table>
        </div>
      </Card>
      <Drawer open={Boolean(selected)} title={selected?.ip_address ?? "IP Details"} onClose={() => setSelected(null)}>
        {selected ? (
          <div className="space-y-4">
            <Detail label="Country" value={`${selected.country || "Unknown"} ${selected.country_code ? `(${selected.country_code})` : ""}`} />
            <Detail label="Region" value={selected.region || "Unknown"} />
            <Detail label="City" value={selected.city || "Unknown"} />
            <Detail label="ASN" value={selected.asn || "Unknown"} />
            <Detail label="ISP" value={selected.isp || "Unknown"} />
            <Detail label="Jail" value={selected.jail || "Unknown"} />
            <Detail label="Requests" value={selected.request_count.toLocaleString()} />
            <Detail label="429 Count" value={selected.violation_count.toLocaleString()} />
            <div className="flex gap-2 pt-2">
              <Button variant="primary">
                <ShieldCheck className="h-4 w-4" />
                Whitelist
              </Button>
              {!history ? (
                <Button variant="danger" onClick={() => unban(selected)}>
                  <Unlock className="h-4 w-4" />
                  Unban
                </Button>
              ) : null}
            </div>
          </div>
        ) : null}
      </Drawer>
    </>
  );
}

function Detail({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-md border p-3">
      <div className="text-xs text-muted-foreground">{label}</div>
      <div className="mt-1 font-medium">{value}</div>
    </div>
  );
}

function formatRemaining(ban: Ban) {
  const end = new Date(new Date(ban.ban_time).getTime() + ban.ban_duration * 1000);
  const seconds = Math.max(0, Math.floor((end.getTime() - Date.now()) / 1000));
  const hours = Math.floor(seconds / 3600);
  const minutes = Math.floor((seconds % 3600) / 60);
  return `${hours}h ${minutes}m`;
}
