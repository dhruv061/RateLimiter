import { useState } from "react";
import { useGlobalFilter, type Domain } from "../../context/GlobalFilterContext";
import { Card } from "../../components/ui/card";
import { Button } from "../../components/ui/button";
import { Badge } from "../../components/ui/badge";
import { Drawer } from "../../components/ui/drawer";
import { Input } from "../../components/ui/input";
import { api } from "../../services/api";
import { RemovalWizard } from "./RemovalWizard";
import { cn } from "../../utils/cn";
import {
	Plus,
	Globe,
	Trash2,
	Edit3,
	RefreshCcw,
	Activity,
	Shield,
	Clock,
	Settings,
	HelpCircle,
	AlertTriangle,
	CheckCircle,
	ArrowRight,
	Loader2
} from "lucide-react";

export function DomainsPage() {
	const { domains, loadingDomains, refreshDomains, setSelectedDomain, setSelectedRange } = useGlobalFilter();
	const [editingDomain, setEditingDomain] = useState<Domain | null>(null);
	const [deletingDomain, setDeletingDomain] = useState<Domain | null>(null);

	// Edit Drawer Form State
	const [serverName, setServerName] = useState("");
	const [description, setDescription] = useState("");
	const [accessLogPath, setAccessLogPath] = useState("");
	const [errorLogPath, setErrorLogPath] = useState("");
	const [blockedIPFilePath, setBlockedIPFilePath] = useState("");
	const [saving, setSaving] = useState(false);
	const [saveError, setSaveError] = useState<string | null>(null);

	const handleOpenEdit = (d: Domain) => {
		setEditingDomain(d);
		setServerName(d.server_name || "");
		setDescription(d.description || "");
		setAccessLogPath(d.access_log_path || "");
		setErrorLogPath(d.error_log_path || "");
		setBlockedIPFilePath(d.blocked_ip_file_path || "");
		setSaveError(null);
	};

	const handleSaveEdit = async () => {
		if (!editingDomain) return;
		setSaving(true);
		setSaveError(null);
		try {
			const payload = {
				server_name: serverName,
				description: description,
				access_log_path: accessLogPath,
				error_log_path: errorLogPath,
				blocked_ip_file_path: blockedIPFilePath
			};
			await api(`/api/domains/${editingDomain.id}`, {
				method: "PUT",
				body: JSON.stringify(payload)
			});
			await refreshDomains();
			setEditingDomain(null);
		} catch (err) {
			setSaveError(err instanceof Error ? err.message : "Failed to update domain");
		} finally {
			setSaving(false);
		}
	};

	// Navigate to setup page
	const triggerSetupWizard = () => {
		// Dispatch event or window location change to navigate to setup wizard.
		// App.tsx handles page routes, we will trigger page set to 'setup' using a custom event
		const event = new CustomEvent("navigate-page", { detail: "setup" });
		window.dispatchEvent(event);
	};

	const triggerReconfigure = (d: Domain) => {
		const event = new CustomEvent("navigate-page", {
			detail: { page: "setup", domainName: d.domain_name, rateLimit: d.rate_limit, burstSize: d.burst_size, banTime: d.ban_time }
		});
		window.dispatchEvent(event);
	};

	if (deletingDomain) {
		return (
			<RemovalWizard
				domain={deletingDomain}
				onComplete={() => {
					setDeletingDomain(null);
					refreshDomains();
				}}
				onCancel={() => setDeletingDomain(null)}
			/>
		);
	}

	return (
		<div className="space-y-6 select-none animate-fadeInUp">
			{/* Header */}
			<div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between border-b border-border pb-4">
				<div>
					<h1 className="text-2xl font-bold tracking-tight">Protected Domains</h1>
					<p className="text-sm text-muted-foreground flex items-center gap-1.5 mt-1 leading-normal">
						Configure multi-domain rate limits, manage access paths, and monitor jail activity.
						<span className="group relative cursor-pointer text-muted-foreground hover:text-foreground">
							<HelpCircle className="h-4 w-4 inline" />
							<span className="pointer-events-none absolute bottom-full left-1/2 z-50 mb-2 w-72 -translate-x-1/2 rounded-md border bg-card p-3.5 text-xs text-foreground shadow-lg opacity-0 transition-opacity group-hover:opacity-100 leading-normal font-normal">
								Domains isolate access logs, error logs, and blacklists. Statistics, whitelist entries, and audit logs are isolated on a per-domain basis.
							</span>
						</span>
					</p>
				</div>
				<Button variant="primary" onClick={triggerSetupWizard} className="flex items-center gap-2 h-11 px-5 shadow-lg shadow-primary/20">
					<Plus className="h-4 w-4" /> Configure Domain
				</Button>
			</div>

			{loadingDomains ? (
				<div className="flex flex-col items-center justify-center py-20 space-y-3 bg-card/20 rounded-2xl border border-border/80">
					<Loader2 className="h-8 w-8 text-primary animate-spin" />
					<p className="text-sm text-muted-foreground">Fetching domain configurations...</p>
				</div>
			) : (
				<div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
					{domains.map((d) => (
						<Card
							key={d.id}
							className="relative border border-border/85 bg-card/40 hover:bg-card/75 backdrop-blur-md transition-all duration-300 hover:shadow-2xl hover:shadow-primary/5 group rounded-2xl p-6 flex flex-col justify-between min-h-[220px]"
						>
							<div className="space-y-4">
								{/* Card Title & Badge */}
								<div className="flex items-start justify-between">
									<div className="flex items-center gap-2.5">
										<div className="h-9 w-9 rounded-lg bg-primary/10 text-primary flex items-center justify-center flex-shrink-0">
											<Globe className="h-4 w-4" />
										</div>
										<div>
											<h3 className="font-bold text-base text-foreground leading-tight">{d.domain_name}</h3>
											<span className="text-xs text-muted-foreground block mt-1">{d.server_name || "Nginx Server"}</span>
										</div>
									</div>
									<Badge
										tone={d.status === "active" ? "success" : d.status === "removing" ? "danger" : "warning"}
										className="text-[10px] px-2 py-0.5"
									>
										<span className="inline-block h-1.5 w-1.5 rounded-full mr-1.5 bg-current" />
										{d.status === "active" ? "Active" : d.status === "removing" ? "Removing" : "Pending"}
									</Badge>
								</div>

								{/* Details checklist */}
								<div className="grid grid-cols-2 gap-y-3 pt-2 text-xs border-t border-border/60">
									<div className="flex items-center gap-2 text-muted-foreground">
										<Shield className="h-3.5 w-3.5 flex-shrink-0 text-muted-foreground" />
										<span className="truncate">Jail: <span className="font-semibold text-foreground">{d.fail2ban_jail_name}</span></span>
									</div>
									<div className="flex items-center gap-2 text-muted-foreground">
										<Activity className="h-3.5 w-3.5 flex-shrink-0 text-muted-foreground" />
										<span>Limit: <span className="font-semibold text-foreground">{d.rate_limit || 5} r/s</span></span>
									</div>
									<div className="flex items-center gap-2 text-muted-foreground">
										<Clock className="h-3.5 w-3.5 flex-shrink-0 text-muted-foreground" />
										<span className="truncate">Last Ban: <span className="font-semibold text-foreground">No recent bans</span></span>
									</div>
									<div className="flex items-center gap-2 text-muted-foreground">
										<CheckCircle className="h-3.5 w-3.5 flex-shrink-0 text-emerald-500" />
										<span>Config: <span className={cn("font-semibold", d.is_valid ? "text-emerald-500" : "text-danger")}>{d.is_valid ? "Valid" : "Invalid"}</span></span>
									</div>
								</div>
							</div>

							{/* Actions Row */}
							<div className="flex gap-2.5 pt-6 border-t border-border/40 mt-6">
								<Button
									variant="ghost"
									size="sm"
									onClick={() => handleOpenEdit(d)}
									className="flex-1 h-9 text-xs flex items-center justify-center gap-1 hover:bg-muted border border-border"
								>
									<Edit3 className="h-3.5 w-3.5" /> Edit
								</Button>
								<Button
									variant="ghost"
									size="sm"
									onClick={() => triggerReconfigure(d)}
									className="flex-1 h-9 text-xs flex items-center justify-center gap-1 hover:bg-muted border border-border"
								>
									<Settings className="h-3.5 w-3.5" /> Reconfigure
								</Button>
								<Button
									variant="ghost"
									size="sm"
									onClick={() => setDeletingDomain(d)}
									className="h-9 px-3 text-xs flex items-center justify-center text-danger hover:bg-danger/10 hover:text-danger border border-danger/25"
									aria-label="Delete Protection"
								>
									<Trash2 className="h-3.5 w-3.5" />
								</Button>
							</div>
						</Card>
					))}

					{/* Add Domain Card CTA */}
					<button
						onClick={triggerSetupWizard}
						className="relative flex flex-col items-center justify-center p-8 rounded-2xl border-2 border-dashed border-border/80 bg-card/10 hover:bg-card/30 hover:border-primary/50 text-muted-foreground hover:text-foreground transition-all duration-300 min-h-[220px] gap-3"
					>
						<div className="h-12 w-12 rounded-full border border-dashed border-border flex items-center justify-center bg-card/25 group-hover:scale-105 transition-transform">
							<Plus className="h-6 w-6 text-muted-foreground" />
						</div>
						<div className="text-center space-y-1">
							<span className="font-bold text-sm text-foreground">Configure New Domain</span>
							<p className="text-xs text-muted-foreground max-w-[200px] leading-normal mt-1">
								Step-by-step onboarding guide to setup Fail2Ban on your server.
							</p>
						</div>
					</button>
				</div>
			)}

			{/* Edit Domain Drawer */}
			<Drawer open={editingDomain !== null} title="Edit Domain Information" onClose={() => setEditingDomain(null)}>
				{editingDomain && (
					<div className="flex flex-col h-full space-y-6">
						<div className="space-y-4">
							<h3 className="text-sm font-bold text-foreground uppercase tracking-wide">
								Domain: {editingDomain.domain_name}
							</h3>

							{saveError && (
								<div className="p-3 bg-danger/10 border border-danger/30 text-xs text-danger rounded-lg">
									{saveError}
								</div>
							)}

							<div className="space-y-2">
								<label className="text-xs font-semibold text-muted-foreground">Server Alias</label>
								<Input
									placeholder="e.g. Primary Nginx, API Gateway"
									value={serverName}
									onChange={(e) => setServerName(e.target.value)}
								/>
							</div>

							<div className="space-y-2">
								<label className="text-xs font-semibold text-muted-foreground">Description</label>
								<Input
									placeholder="e.g. Production site"
									value={description}
									onChange={(e) => setDescription(e.target.value)}
								/>
							</div>

							<div className="space-y-2">
								<label className="text-xs font-semibold text-muted-foreground">Access Log Path</label>
								<Input
									placeholder="/var/log/nginx/access.log"
									value={accessLogPath}
									onChange={(e) => setAccessLogPath(e.target.value)}
									className="font-mono text-xs"
								/>
							</div>

							<div className="space-y-2">
								<label className="text-xs font-semibold text-muted-foreground">Error Log Path</label>
								<Input
									placeholder="/var/log/nginx/error.log"
									value={errorLogPath}
									onChange={(e) => setErrorLogPath(e.target.value)}
									className="font-mono text-xs"
								/>
							</div>

							<div className="space-y-2">
								<label className="text-xs font-semibold text-muted-foreground">Blocked IP Conf File</label>
								<Input
									placeholder="/etc/nginx/blocked_ips.conf"
									value={blockedIPFilePath}
									onChange={(e) => setBlockedIPFilePath(e.target.value)}
									className="font-mono text-xs"
								/>
							</div>
						</div>

						<div className="flex justify-end gap-3 border-t border-border pt-4 mt-auto">
							<Button variant="secondary" onClick={() => setEditingDomain(null)}>
								Cancel
							</Button>
							<Button variant="primary" disabled={saving} onClick={handleSaveEdit} className="flex items-center gap-1">
								{saving ? <Loader2 className="h-4 w-4 animate-spin" /> : "Save Changes"}
							</Button>
						</div>
					</div>
				)}
			</Drawer>
		</div>
	);
}
