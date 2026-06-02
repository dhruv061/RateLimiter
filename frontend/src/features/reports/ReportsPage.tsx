import { Download } from "lucide-react";
import { PageHeader } from "../../components/layout/PageHeader";
import { Button } from "../../components/ui/button";
import { Card, CardTitle } from "../../components/ui/card";
import { useApi } from "../../hooks/useApi";

type Report = {
  generated_at?: string;
  stats?: Record<string, number | string>;
};

export function ReportsPage() {
  const { data } = useApi<Report>("/api/reports/security", {});

  return (
    <>
      <PageHeader
        title="Reports"
        subtitle="Security, traffic, and compliance exports."
        actions={
          <Button variant="secondary">
            <Download className="h-4 w-4" />
            Export
          </Button>
        }
      />
      <div className="grid gap-4 lg:grid-cols-3">
        <ReportCard title="Security Report" value={String(data.stats?.active_bans ?? 0)} label="Active bans" />
        <ReportCard title="Traffic Report" value={String(data.stats?.total_requests_today ?? 0)} label="Requests today" />
        <ReportCard title="Compliance Report" value={data.generated_at ? new Date(data.generated_at).toLocaleDateString() : "-"} label="Generated" />
      </div>
    </>
  );
}

function ReportCard({ title, value, label }: { title: string; value: string; label: string }) {
  return (
    <Card>
      <CardTitle>{title}</CardTitle>
      <div className="mt-5 text-[32px] font-semibold leading-none">{value}</div>
      <div className="mt-3 text-sm text-muted-foreground">{label}</div>
    </Card>
  );
}
