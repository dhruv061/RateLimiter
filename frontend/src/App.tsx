import { useEffect, useState, type FormEvent } from "react";
import { AppShell, type PageKey } from "./components/layout/AppShell";
import { Button } from "./components/ui/button";
import { Card } from "./components/ui/card";
import { Input } from "./components/ui/input";
import { AnalyticsPage } from "./features/analytics/AnalyticsPage";
import { AuditPage } from "./features/audit/AuditPage";
import { BansPage } from "./features/bans/BansPage";
import { DashboardPage } from "./features/dashboard/DashboardPage";
import { LiveRequestsPage } from "./features/live/LiveRequestsPage";
import { ReportsPage } from "./features/reports/ReportsPage";
import { SettingsPage } from "./features/settings/SettingsPage";
import { WhitelistPage } from "./features/whitelist/WhitelistPage";
import { api, getToken, setToken } from "./services/api";
import type { User } from "./types/api";

export function App() {
  const [page, setPage] = useState<PageKey>("dashboard");
  const [theme, setTheme] = useState<"dark" | "light">(() => (localStorage.getItem("theme") as "dark" | "light") || "dark");
  const [authenticated, setAuthenticated] = useState(Boolean(getToken()));

  useEffect(() => {
    document.documentElement.classList.toggle("dark", theme === "dark");
    localStorage.setItem("theme", theme);
  }, [theme]);

  if (!authenticated) {
    return <Login onLogin={() => setAuthenticated(true)} />;
  }

  return (
    <AppShell page={page} setPage={setPage} theme={theme} onThemeToggle={() => setTheme(theme === "dark" ? "light" : "dark")}>
      {page === "dashboard" && <DashboardPage />}
      {page === "active-bans" && <BansPage />}
      {page === "history" && <BansPage history />}
      {page === "live" && <LiveRequestsPage />}
      {page === "analytics" && <AnalyticsPage />}
      {page === "whitelist" && <WhitelistPage />}
      {page === "audit" && <AuditPage />}
      {page === "settings" && <SettingsPage />}
      {page === "reports" && <ReportsPage />}
    </AppShell>
  );
}

function Login({ onLogin }: { onLogin: () => void }) {
  const [username, setUsername] = useState("admin");
  const [password, setPassword] = useState("admin");
  const [error, setError] = useState("");

  async function submit(event: FormEvent) {
    event.preventDefault();
    setError("");
    try {
      const result = await api<{ token: string; user: User }>("/api/auth/login", {
        method: "POST",
        body: JSON.stringify({ username, password })
      });
      setToken(result.token);
      onLogin();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Login failed");
    }
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-background p-4">
      <Card className="w-full max-w-md">
        <div className="mb-6">
          <h1 className="text-[30px] font-semibold leading-tight">Fail2Ban Dashboard</h1>
          <p className="mt-2 text-sm text-muted-foreground">Sign in to continue.</p>
        </div>
        <form className="space-y-4" onSubmit={submit}>
          <Input value={username} onChange={(e) => setUsername(e.target.value)} placeholder="Username" />
          <Input value={password} onChange={(e) => setPassword(e.target.value)} placeholder="Password" type="password" />
          {error ? <div className="rounded-md border border-danger/30 bg-danger/10 p-3 text-sm text-danger">{error}</div> : null}
          <Button className="w-full" variant="primary" type="submit">
            Sign In
          </Button>
        </form>
      </Card>
    </div>
  );
}
