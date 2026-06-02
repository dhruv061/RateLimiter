import { useState } from "react";
import { useGlobalFilter, type Domain } from "../../context/GlobalFilterContext";
import { Card } from "../../components/ui/card";
import { Button } from "../../components/ui/button";
import { Table, Th, Td, EmptyRow } from "../../components/ui/table";
import { Badge } from "../../components/ui/badge";
import { Drawer } from "../../components/ui/drawer";
import { Input } from "../../components/ui/input";
import { api } from "../../services/api";
import { cn } from "../../utils/cn";
import {
	Plus,
	Trash2,
	CheckCircle2,
	AlertTriangle,
	ShieldCheck,
	Globe,
	Server,
	RefreshCw,
	HelpCircle
} from "lucide-react";

export function DomainsPage() {
	const { domains, loadingDomains, refreshDomains } = useGlobalFilter();
	const [drawerOpen, setDrawerOpen] = useState(false);
	const [step, setStep] = useState(1);
	const [deletingId, setDeletingId] = useState<number | null>(null);

	// Form State
	const [domainName, setDomainName] = useState("");
	const [serverName, setServerName] = useState("");
	const [description, setDescription] = useState("");
	const [accessLogPath, setAccessLogPath] = useState("");
	const [errorLogPath, setErrorLogPath] = useState("");
	const [blockedIPFilePath, setBlockedIPFilePath] = useState("");
	const [fail2banJailName, setFail2banJailName] = useState("nginx-429");

	// Validation & Saving State
	const [validating, setValidating] = useState(false);
	const [validationResult, setValidationResult] = useState<{
		access_log_exists: boolean;
		access_log_msg: string;
		error_log_exists: boolean;
		error_log_msg: string;
		block_file_exists: boolean;
		block_file_msg: string;
		fail2ban_jail_ok: boolean;
		fail2ban_jail_msg: string;
		overall_valid: boolean;
	} | null>(null);
	const [errorMsg, setErrorMsg] = useState<string | null>(null);

	const resetForm = () => {
		setDomainName("");
		setServerName("");
		setDescription("");
		setAccessLogPath("");
		setErrorLogPath("");
		setBlockedIPFilePath("");
		setFail2banJailName("nginx-429");
		setValidationResult(null);
		setErrorMsg(null);
		setStep(1);
	};

	const handleOpenDrawer = () => {
		resetForm();
		setDrawerOpen(true);
	};

	const handleDelete = async (id: number) => {
		if (!confirm("Are you sure you want to delete this domain? This will remove its scoped statistics and configuration.")) {
			return;
		}
		try {
			setDeletingId(id);
			await api(`/api/domains/${id}`, { method: "DELETE" });
			await refreshDomains();
		} catch (err) {
			alert(err instanceof Error ? err.message : "Failed to delete domain");
		} finally {
			setDeletingId(null);
		}
	};

	const runValidation = async () => {
		setValidating(true);
		setValidationError(null);
		try {
			const payload = {
				domain_name: domainName,
				access_log_path: accessLogPath,
				error_log_path: errorLogPath,
				blocked_ip_file_path: blockedIPFilePath,
				fail2ban_jail_name: fail2banJailName,
				server_name: serverName,
				description
			};
			const res = await api<typeof validationResult>("/api/domains/validate", {
				method: "POST",
				body: JSON.stringify(payload)
			});
			setValidationResult(res);
		} catch (err) {
			setValidationError(err instanceof Error ? err.message : "Validation request failed");
		} finally {
			setValidating(false);
		}
	};

	const [validationError, setValidationError] = useState<string | null>(null);

	const handleSave = async () => {
		setErrorMsg(null);
		try {
			const payload = {
				domain_name: domainName,
				access_log_path: accessLogPath,
				error_log_path: errorLogPath,
				blocked_ip_file_path: blockedIPFilePath,
				fail2ban_jail_name: fail2banJailName,
				server_name: serverName,
				description
			};
			await api("/api/domains", {
				method: "POST",
				body: JSON.stringify(payload)
			});
			await refreshDomains();
			setDrawerOpen(false);
		} catch (err) {
			setErrorMsg(err instanceof Error ? err.message : "Failed to save domain configuration");
		}
	};

	return (
		<div className="space-y-6">
			{/* Header */}
			<div className="flex flex-col gap-4 sm:flex-row sm:items-center sm:justify-between">
				<div>
					<h1 className="text-2xl font-bold tracking-tight">Domain Configurations</h1>
					<p className="text-sm text-muted-foreground flex items-center gap-1 mt-1">
						Configure multi-domain hosting, scope metrics, and configure custom log paths.
						<span className="group relative cursor-pointer text-muted-foreground hover:text-foreground">
							<HelpCircle className="h-4 w-4 inline" />
							<span className="pointer-events-none absolute bottom-full left-1/2 z-50 mb-2 w-72 -translate-x-1/2 rounded-md border bg-card p-3.5 text-xs text-foreground shadow-lg opacity-0 transition-opacity group-hover:opacity-100 leading-normal">
								Domains isolate access logs, error logs, and blacklists. Statistics, whitelist entries, and audit logs are isolated on a per-domain basis.
							</span>
						</span>
					</p>
				</div>
				<Button variant="primary" onClick={handleOpenDrawer} className="flex items-center gap-2">
					<Plus className="h-4 w-4" /> Add Domain
				</Button>
			</div>

			{/* Listing Table */}
			<Card className="overflow-hidden">
				<div className="overflow-x-auto">
					<Table>
						<thead>
							<tr>
								<Th>Domain</Th>
								<Th>Server / Details</Th>
								<Th>Access Log Path</Th>
								<Th>Jail Name</Th>
								<Th>Validation Status</Th>
								<Th className="text-right">Actions</Th>
							</tr>
						</thead>
						<tbody>
							{loadingDomains ? (
								<tr>
									<Td colSpan={6} className="h-28 text-center text-muted-foreground">
										<RefreshCw className="h-5 w-5 animate-spin mx-auto mb-2" />
										Loading domains...
									</Td>
								</tr>
							) : domains.length === 0 ? (
								<EmptyRow colSpan={6}>No domains configured. Create a domain to start tracking traffic isolation.</EmptyRow>
							) : (
								domains.map((d) => (
									<tr key={d.id}>
										<Td className="font-semibold text-foreground">
											<div className="flex items-center gap-2">
												<Globe className="h-4 w-4 text-muted-foreground" />
												{d.domain_name}
											</div>
										</Td>
										<Td>
											<div className="text-sm font-medium">{d.server_name || "N/A"}</div>
											<div className="text-xs text-muted-foreground truncate max-w-xs">{d.description || "No description"}</div>
										</Td>
										<Td className="text-xs font-mono text-muted-foreground select-all">{d.access_log_path}</Td>
										<Td>
											<Badge tone="info" className="font-mono">
												{d.fail2ban_jail_name}
											</Badge>
										</Td>
										<Td>
											{d.is_valid ? (
												<Badge tone="success" className="gap-1">
													<CheckCircle2 className="h-3 w-3" /> Valid
												</Badge>
											) : (
												<Badge tone="danger" className="gap-1">
													<AlertTriangle className="h-3 w-3" /> Invalid Paths
												</Badge>
											)}
										</Td>
										<Td className="text-right">
											<Button
												size="icon"
												variant="ghost"
												onClick={() => handleDelete(d.id)}
												disabled={deletingId === d.id}
												aria-label="Delete domain"
												className="text-danger hover:bg-danger/10"
											>
												<Trash2 className="h-4 w-4" />
											</Button>
										</Td>
									</tr>
								))
							)}
						</tbody>
					</Table>
				</div>
			</Card>

			{/* Add Domain Wizard Drawer */}
			<Drawer open={drawerOpen} title="Configure New Domain" onClose={() => setDrawerOpen(false)}>
				<div className="flex flex-col h-full space-y-6">
					{/* Progress Indicator */}
					<div className="mb-4 flex justify-between items-center text-xs text-muted-foreground border-b border-border pb-4">
						<div className={cn("flex flex-col items-center", step >= 1 && "text-primary font-semibold")}>
							<span className={cn("flex h-6 w-6 items-center justify-center rounded-full border", step >= 1 ? "border-primary bg-primary/10 text-primary" : "border-muted-foreground")}>1</span>
							<span className="mt-1">Info</span>
						</div>
						<div className={cn("h-[1px] flex-1 mx-2", step >= 2 ? "bg-primary" : "bg-border")} />
						<div className={cn("flex flex-col items-center", step >= 2 && "text-primary font-semibold")}>
							<span className={cn("flex h-6 w-6 items-center justify-center rounded-full border", step >= 2 ? "border-primary bg-primary/10 text-primary" : "border-muted-foreground")}>2</span>
							<span className="mt-1">Logs</span>
						</div>
						<div className={cn("h-[1px] flex-1 mx-2", step >= 3 ? "bg-primary" : "bg-border")} />
						<div className={cn("flex flex-col items-center", step >= 3 && "text-primary font-semibold")}>
							<span className={cn("flex h-6 w-6 items-center justify-center rounded-full border", step >= 3 ? "border-primary bg-primary/10 text-primary" : "border-muted-foreground")}>3</span>
							<span className="mt-1">Jail</span>
						</div>
						<div className={cn("h-[1px] flex-1 mx-2", step >= 4 ? "bg-primary" : "bg-border")} />
						<div className={cn("flex flex-col items-center", step >= 4 && "text-primary font-semibold")}>
							<span className={cn("flex h-6 w-6 items-center justify-center rounded-full border", step >= 4 ? "border-primary bg-primary/10 text-primary" : "border-muted-foreground")}>4</span>
							<span className="mt-1">Verify</span>
						</div>
					</div>

					{/* Wizard Step Content */}
					<div className="flex-1 space-y-4">
						{step === 1 && (
							<div className="space-y-4 animate-fadeIn">
								<h3 className="text-sm font-semibold text-foreground uppercase tracking-wide">Domain Information</h3>
								<div className="space-y-2">
									<label className="text-xs font-semibold text-muted-foreground">Domain Name *</label>
									<Input
										placeholder="e.g. company.com"
										value={domainName}
										onChange={(e) => setDomainName(e.target.value)}
									/>
								</div>
								<div className="space-y-2">
									<label className="text-xs font-semibold text-muted-foreground">Server Name</label>
									<Input
										placeholder="e.g. Primary Nginx, API Gateway"
										value={serverName}
										onChange={(e) => setServerName(e.target.value)}
									/>
								</div>
								<div className="space-y-2">
									<label className="text-xs font-semibold text-muted-foreground">Description</label>
									<Input
										placeholder="e.g. Staging server domain configuration"
										value={description}
										onChange={(e) => setDescription(e.target.value)}
									/>
								</div>
							</div>
						)}

						{step === 2 && (
							<div className="space-y-4 animate-fadeIn">
								<h3 className="text-sm font-semibold text-foreground uppercase tracking-wide">Log Configuration Paths</h3>
								<div className="space-y-2">
									<label className="text-xs font-semibold text-muted-foreground">Access Log Path *</label>
									<Input
										placeholder="e.g. /var/log/nginx/access.log"
										value={accessLogPath}
										onChange={(e) => setAccessLogPath(e.target.value)}
										className="font-mono text-xs"
									/>
								</div>
								<div className="space-y-2">
									<label className="text-xs font-semibold text-muted-foreground">Error Log Path *</label>
									<Input
										placeholder="e.g. /var/log/nginx/error.log"
										value={errorLogPath}
										onChange={(e) => setErrorLogPath(e.target.value)}
										className="font-mono text-xs"
									/>
								</div>
								<div className="space-y-2">
									<label className="text-xs font-semibold text-muted-foreground">Blocked IP Config File *</label>
									<Input
										placeholder="e.g. /etc/nginx/blocked_ips.conf"
										value={blockedIPFilePath}
										onChange={(e) => setBlockedIPFilePath(e.target.value)}
										className="font-mono text-xs"
									/>
								</div>
							</div>
						)}

						{step === 3 && (
							<div className="space-y-4 animate-fadeIn">
								<h3 className="text-sm font-semibold text-foreground uppercase tracking-wide">Security Jail Mapping</h3>
								<div className="space-y-2">
									<label className="text-xs font-semibold text-muted-foreground">Fail2Ban Jail Name *</label>
									<Input
										placeholder="e.g. nginx-429"
										value={fail2banJailName}
										onChange={(e) => setFail2banJailName(e.target.value)}
										className="font-mono text-xs"
									/>
								</div>
								<div className="p-3 bg-muted rounded-lg border border-border text-xs text-muted-foreground leading-relaxed flex gap-2">
									<ShieldCheck className="h-4 w-4 text-primary shrink-0" />
									Specify the Fail2Ban jail assigned to rate-limiting this domain. Ensure the jail exists in your Fail2Ban configuration.
								</div>
							</div>
						)}

						{step === 4 && (
							<div className="space-y-4 animate-fadeIn">
								<h3 className="text-sm font-semibold text-foreground uppercase tracking-wide">Verification</h3>

								{!validationResult && !validating && (
									<div className="py-6 text-center space-y-4">
										<p className="text-sm text-muted-foreground">Ready to validate domain paths and configuration.</p>
										<Button variant="primary" onClick={runValidation}>
											Run Validation Checks
										</Button>
									</div>
								)}

								{validating && (
									<div className="py-8 text-center space-y-3">
										<RefreshCw className="h-6 w-6 animate-spin text-primary mx-auto" />
										<p className="text-sm text-muted-foreground">Connecting to host client to validate configuration...</p>
									</div>
								)}

								{validationError && (
									<div className="p-3.5 bg-danger/10 border border-danger/30 text-sm text-danger rounded-lg">
										{validationError}
									</div>
								)}

								{validationResult && (
									<div className="space-y-3 animate-fadeIn">
										<div className="flex items-center justify-between border-b pb-2 text-xs font-semibold text-muted-foreground">
											<span>Path Check</span>
											<span>Status</span>
										</div>

										<div className="flex justify-between items-center text-sm py-1">
											<span className="font-medium">Access Log Existence</span>
											{validationResult.access_log_exists ? (
												<span className="text-success flex items-center gap-1"><CheckCircle2 className="h-4 w-4" /> Exist</span>
											) : (
												<span className="text-danger flex items-center gap-1"><AlertTriangle className="h-4 w-4" /> Missing</span>
											)}
										</div>
										<div className="text-xs text-muted-foreground font-mono bg-muted p-1.5 rounded truncate">{accessLogPath}</div>

										<div className="flex justify-between items-center text-sm py-1 border-t pt-2">
											<span className="font-medium">Error Log Existence</span>
											{validationResult.error_log_exists ? (
												<span className="text-success flex items-center gap-1"><CheckCircle2 className="h-4 w-4" /> Exist</span>
											) : (
												<span className="text-danger flex items-center gap-1"><AlertTriangle className="h-4 w-4" /> Missing</span>
											)}
										</div>
										<div className="text-xs text-muted-foreground font-mono bg-muted p-1.5 rounded truncate">{errorLogPath}</div>

										<div className="flex justify-between items-center text-sm py-1 border-t pt-2">
											<span className="font-medium">Block File Writable</span>
											{validationResult.block_file_exists ? (
												<span className="text-success flex items-center gap-1"><CheckCircle2 className="h-4 w-4" /> Exist</span>
											) : (
												<span className="text-danger flex items-center gap-1"><AlertTriangle className="h-4 w-4" /> Missing</span>
											)}
										</div>
										<div className="text-xs text-muted-foreground font-mono bg-muted p-1.5 rounded truncate">{blockedIPFilePath}</div>

										<div className="flex justify-between items-center text-sm py-1 border-t pt-2">
											<span className="font-medium">Fail2Ban Jail Active</span>
											{validationResult.fail2ban_jail_ok ? (
												<span className="text-success flex items-center gap-1"><CheckCircle2 className="h-4 w-4" /> OK</span>
											) : (
												<span className="text-danger flex items-center gap-1"><AlertTriangle className="h-4 w-4" /> Inactive</span>
											)}
										</div>

										{validationResult.overall_valid ? (
											<div className="mt-4 p-3 bg-success/10 border border-success/30 rounded-lg text-sm text-success flex items-center gap-2">
												<CheckCircle2 className="h-5 w-5 shrink-0" />
												All validation checks passed successfully.
											</div>
										) : (
											<div className="mt-4 p-3 bg-warning/10 border border-warning/30 rounded-lg text-xs text-warning flex items-start gap-2">
												<AlertTriangle className="h-4 w-4 shrink-0 mt-0.5" />
												<div>
													<p className="font-semibold text-sm">Path validation failed</p>
													<p className="mt-1 leading-normal">
														One or more log file paths do not exist. In demo mode or development, you can still save the configuration to populate demo stats.
													</p>
												</div>
											</div>
										)}
									</div>
								)}

								{errorMsg && (
									<div className="p-3 bg-danger/10 border border-danger/30 text-sm text-danger rounded-lg">
										{errorMsg}
									</div>
								)}
							</div>
						)}
					</div>

					{/* Navigation Buttons */}
					<div className="flex justify-between items-center border-t border-border pt-4 mt-auto">
						<Button
							variant="secondary"
							disabled={step === 1}
							onClick={() => setStep((s) => s - 1)}
						>
							Back
						</Button>

						{step < 4 ? (
							<Button
								variant="primary"
								disabled={step === 1 && !domainName}
								onClick={() => setStep((s) => s + 1)}
							>
								Next
							</Button>
						) : (
							<Button
								variant="primary"
								disabled={!validationResult}
								onClick={handleSave}
							>
								Save Configuration
							</Button>
						)}
					</div>
				</div>
			</Drawer>
		</div>
	);
}
