import {
	Activity,
	BarChart3,
	Bell,
	FileClock,
	Gauge,
	ListTree,
	LockKeyhole,
	Search,
	Settings,
	Shield,
	ShieldCheck,
	SunMoon,
	UserCircle,
	Wifi,
	Globe,
	RefreshCw,
	Calendar
} from "lucide-react";
import type { ReactNode } from "react";
import { Button } from "../ui/button";
import { Input } from "../ui/input";
import { Badge } from "../ui/badge";
import { cn } from "../../utils/cn";
import { useWebSocket } from "../../hooks/useWebSocket";
import { useGlobalFilter } from "../../context/GlobalFilterContext";

export type PageKey =
	| "dashboard"
	| "active-bans"
	| "history"
	| "live"
	| "analytics"
	| "whitelist"
	| "audit"
	| "settings"
	| "reports"
	| "domains";

const navItems: Array<{ key: PageKey; label: string; icon: typeof Gauge }> = [
	{ key: "dashboard", label: "Dashboard", icon: Gauge },
	{ key: "active-bans", label: "Active Bans", icon: Shield },
	{ key: "history", label: "Ban History", icon: FileClock },
	{ key: "live", label: "Live Requests", icon: Activity },
	{ key: "analytics", label: "Analytics", icon: BarChart3 },
	{ key: "whitelist", label: "Whitelist", icon: ShieldCheck },
	{ key: "domains", label: "Domains", icon: Globe },
	{ key: "audit", label: "Audit Logs", icon: ListTree },
	{ key: "settings", label: "Settings", icon: Settings },
	{ key: "reports", label: "Reports", icon: LockKeyhole }
];

function toLocalISOString(date: Date): string {
	const offset = date.getTimezoneOffset();
	const localDate = new Date(date.getTime() - offset * 60 * 1000);
	return localDate.toISOString().substring(0, 16);
}

export function AppShell({
	page,
	setPage,
	theme,
	onThemeToggle,
	children
}: {
	page: PageKey;
	setPage: (page: PageKey) => void;
	theme: "dark" | "light";
	onThemeToggle: () => void;
	children: ReactNode;
}) {
	const { connected } = useWebSocket();
	const {
		selectedDomain,
		setSelectedDomain,
		selectedRange,
		setSelectedRange,
		customRange,
		setCustomRange,
		domains,
		triggerRefresh
	} = useGlobalFilter();

	return (
		<div className="min-h-screen bg-background text-foreground">
			<aside className="fixed inset-y-0 left-0 z-30 hidden w-[240px] border-r bg-card md:block">
				<div className="flex h-16 items-center gap-3 border-b px-5">
					<div className="flex h-9 w-9 items-center justify-center rounded-md bg-primary text-primary-foreground">
						<Shield className="h-5 w-5" />
					</div>
					<div>
						<div className="font-semibold">Fail2Ban</div>
						<div className="text-xs text-muted-foreground">Rate Limit Console</div>
					</div>
				</div>
				<nav className="space-y-1 p-3">
					{navItems.map((item) => {
						const Icon = item.icon;
						return (
							<button
								key={item.key}
								className={cn(
									"flex h-10 w-full items-center gap-3 rounded-md px-3 text-left text-sm text-muted-foreground transition-colors hover:bg-muted hover:text-foreground",
									page === item.key && "bg-muted text-foreground"
								)}
								onClick={() => setPage(item.key)}
							>
								<Icon className="h-4 w-4" />
								{item.label}
							</button>
						);
					})}
				</nav>
			</aside>

			<div className="md:pl-[240px]">
				<header className="sticky top-0 z-20 flex h-16 items-center gap-3 border-b bg-background/90 px-4 backdrop-blur md:px-6">
					{/* Global Filter Bar */}
					<div className="flex flex-wrap items-center gap-3">
						{/* Domain Selector */}
						<div className="flex items-center gap-2">
							<Globe className="h-4 w-4 text-muted-foreground" />
							<select
								value={selectedDomain}
								onChange={(e) => setSelectedDomain(Number(e.target.value))}
								className="bg-card border border-border rounded-md px-2.5 py-1 text-sm font-medium focus:ring-1 focus:ring-primary focus:outline-none"
							>
								<option value={0}>All Domains</option>
								{domains.map((d) => (
									<option key={d.id} value={d.id}>
										{d.domain_name}
									</option>
								))}
							</select>
						</div>

						{/* Date Picker Presets */}
						<div className="flex items-center gap-2">
							<Calendar className="h-4 w-4 text-muted-foreground" />
							<select
								value={selectedRange}
								onChange={(e) => setSelectedRange(e.target.value)}
								className="bg-card border border-border rounded-md px-2.5 py-1 text-sm font-medium focus:ring-1 focus:ring-primary focus:outline-none"
							>
								<option value="last_1h">Last 1 Hour</option>
								<option value="last_24h">Last 24 Hours</option>
								<option value="last_7d">Last 7 Days</option>
								<option value="last_30d">Last 30 Days</option>
								<option value="custom">Custom Range</option>
							</select>
						</div>

						{/* Custom Range Date Pickers */}
						{selectedRange === "custom" && (
							<div className="flex items-center gap-2 text-xs sm:text-sm">
								<input
									type="datetime-local"
									value={customRange.start ? toLocalISOString(customRange.start) : ""}
									onChange={(e) =>
										setCustomRange({
											...customRange,
											start: e.target.value ? new Date(e.target.value) : null
										})
									}
									className="bg-card border border-border rounded-md px-2 py-0.5 text-xs text-foreground focus:ring-1 focus:ring-primary focus:outline-none"
								/>
								<span className="text-muted-foreground">to</span>
								<input
									type="datetime-local"
									value={customRange.end ? toLocalISOString(customRange.end) : ""}
									onChange={(e) =>
										setCustomRange({
											...customRange,
											end: e.target.value ? new Date(e.target.value) : null
										})
									}
									className="bg-card border border-border rounded-md px-2 py-0.5 text-xs text-foreground focus:ring-1 focus:ring-primary focus:outline-none"
								/>
							</div>
						)}

						{/* Manual Force Refresh */}
						<Button size="icon" variant="ghost" onClick={triggerRefresh} aria-label="Refresh data">
							<RefreshCw className="h-4 w-4" />
						</Button>
					</div>

					<div className="ml-auto flex items-center gap-2">
						<Badge tone={connected ? "success" : "warning"}>
							<Wifi className="mr-1 h-3 w-3" />
							{connected ? "Live" : "Polling"}
						</Badge>
						<Button size="icon" variant="ghost" onClick={onThemeToggle} aria-label="Toggle theme">
							<SunMoon className="h-4 w-4" />
						</Button>
						<Button size="icon" variant="ghost" aria-label="Notifications">
							<Bell className="h-4 w-4" />
						</Button>
						<Button variant="secondary" className="hidden sm:inline-flex">
							<UserCircle className="h-4 w-4" />
							Admin
						</Button>
					</div>
				</header>
				<main className="mx-auto max-w-[1600px] px-4 py-6 md:px-6">{children}</main>
			</div>

			<nav className="fixed inset-x-0 bottom-0 z-30 grid grid-cols-5 border-t bg-card md:hidden">
				{navItems.slice(0, 5).map((item) => {
					const Icon = item.icon;
					return (
						<button key={item.key} className="flex h-14 flex-col items-center justify-center gap-1 text-xs" onClick={() => setPage(item.key)}>
							<Icon className={cn("h-4 w-4", page === item.key && "text-primary")} />
							<span className={cn("text-muted-foreground", page === item.key && "text-foreground")}>{item.label.split(" ")[0]}</span>
						</button>
					);
				})}
			</nav>
		</div>
	);
}
