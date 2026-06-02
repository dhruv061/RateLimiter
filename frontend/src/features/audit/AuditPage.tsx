import { PageHeader } from "../../components/layout/PageHeader";
import { Badge } from "../../components/ui/badge";
import { Card } from "../../components/ui/card";
import { EmptyRow, Table, Td, Th } from "../../components/ui/table";
import { useApi } from "../../hooks/useApi";
import type { AuditLog, PaginatedResponse } from "../../types/api";

const empty: PaginatedResponse<AuditLog> = { data: [], total: 0, page: 1, per_page: 20, total_pages: 0 };

export function AuditPage() {
  const { data } = useApi<PaginatedResponse<AuditLog>>("/api/audit-logs?per_page=100", empty);

  return (
    <>
      <PageHeader title="Audit Logs" subtitle="Administrative actions and configuration changes." />
      <Card className="p-0">
        <div className="overflow-auto">
          <Table>
            <thead>
              <tr>
                <Th>User</Th>
                <Th>Action</Th>
                <Th>Target</Th>
                <Th>IP</Th>
                <Th>Timestamp</Th>
              </tr>
            </thead>
            <tbody>
              {data.data.length === 0 ? <EmptyRow colSpan={5}>No dashboard actions have been recorded yet.</EmptyRow> : null}
              {data.data.map((log) => (
                <tr key={log.id} className="hover:bg-muted/60">
                  <Td>{log.username || "system"}</Td>
                  <Td>
                    <Badge tone="info">{log.action}</Badge>
                  </Td>
                  <Td>{log.target || "-"}</Td>
                  <Td className="font-mono">{log.ip_address}</Td>
                  <Td>{new Date(log.created_at).toLocaleString()}</Td>
                </tr>
              ))}
            </tbody>
          </Table>
        </div>
      </Card>
    </>
  );
}
