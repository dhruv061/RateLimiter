import { Plus, Trash2 } from "lucide-react";
import { useState } from "react";
import { PageHeader } from "../../components/layout/PageHeader";
import { Button } from "../../components/ui/button";
import { Card } from "../../components/ui/card";
import { Input } from "../../components/ui/input";
import { EmptyRow, Table, Td, Th } from "../../components/ui/table";
import { api } from "../../services/api";
import { useApi } from "../../hooks/useApi";
import type { PaginatedResponse, WhitelistEntry } from "../../types/api";

const empty: PaginatedResponse<WhitelistEntry> = { data: [], total: 0, page: 1, per_page: 50, total_pages: 0 };

export function WhitelistPage() {
  const { data, setData } = useApi<PaginatedResponse<WhitelistEntry>>("/api/whitelist?per_page=100", empty);
  const [ip, setIp] = useState("");
  const [description, setDescription] = useState("");

  async function add() {
    if (!ip) {
      return;
    }
    const entry = await api<WhitelistEntry>("/api/whitelist", { method: "POST", body: JSON.stringify({ ip_address: ip, description }) });
    setData({ ...data, data: [entry, ...data.data], total: data.total + 1 });
    setIp("");
    setDescription("");
  }

  async function remove(entry: WhitelistEntry) {
    await api(`/api/whitelist/${entry.id}`, { method: "DELETE", body: "{}" });
    setData({ ...data, data: data.data.filter((item) => item.id !== entry.id), total: data.total - 1 });
  }

  return (
    <>
      <PageHeader title="Whitelist" subtitle="Trusted IPs that should not be banned." />
      <Card className="mb-4">
        <div className="grid gap-3 md:grid-cols-[1fr_2fr_auto]">
          <Input value={ip} onChange={(e) => setIp(e.target.value)} placeholder="IP address" />
          <Input value={description} onChange={(e) => setDescription(e.target.value)} placeholder="Description" />
          <Button variant="primary" onClick={add}>
            <Plus className="h-4 w-4" />
            Add
          </Button>
        </div>
      </Card>
      <Card className="p-0">
        <div className="overflow-auto">
          <Table>
            <thead>
              <tr>
                <Th>IP Address</Th>
                <Th>Description</Th>
                <Th>Added By</Th>
                <Th>Date</Th>
                <Th>Actions</Th>
              </tr>
            </thead>
            <tbody>
              {data.data.length === 0 ? <EmptyRow colSpan={5}>No trusted IPs have been whitelisted.</EmptyRow> : null}
              {data.data.map((entry) => (
                <tr key={entry.id} className="hover:bg-muted/60">
                  <Td className="font-mono">{entry.ip_address}</Td>
                  <Td>{entry.description}</Td>
                  <Td>{entry.added_by}</Td>
                  <Td>{new Date(entry.created_at).toLocaleDateString()}</Td>
                  <Td>
                    <Button size="icon" variant="ghost" onClick={() => remove(entry)} aria-label="Remove whitelist entry">
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  </Td>
                </tr>
              ))}
            </tbody>
          </Table>
        </div>
      </Card>
    </>
  );
}
