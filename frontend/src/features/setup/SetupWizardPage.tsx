import { useState, useEffect } from "react";
import { useGlobalFilter } from "../../context/GlobalFilterContext";
import { Card } from "../../components/ui/card";
import { Button } from "../../components/ui/button";
import { Input } from "../../components/ui/input";
import { Badge } from "../../components/ui/badge";
import { CodeViewer } from "../../components/ui/code-viewer";
import { Confetti } from "../../components/ui/confetti";
import { api } from "../../services/api";
import {
	Shield,
	Cpu,
	Globe,
	Settings2,
	CheckCircle,
	AlertTriangle,
	Play,
	Copy,
	Download,
	RefreshCw,
	ChevronRight,
	ChevronLeft,
	Loader2,
	Lock,
	Check,
	HelpCircle
} from "lucide-react";

interface SetupWizardPageProps {
	onComplete: (domainId: number) => void;
	onCancel: () => void;
}

interface Fail2BanStatus {
	installed: boolean;
	running: boolean;
	version: string;
	active_jails: string[];
	jail_count: number;
}

interface DiscoveredDomain {
	server_name: string;
	config_file: string;
	has_ssl: boolean;
}

interface GeneratedFile {
	filename: string;
	path: string;
	content: string;
	type: string;
}

interface GeneratedConfigResponse {
	domain_slug: string;
	files: GeneratedFile[];
	setup_script: string;
	nginx_snippet: string;
	nginx_zone_line: string;
}

interface ValidationCheck {
	name: string;
	passed: boolean;
	message: string;
}

interface ValidationResponse {
	checks: ValidationCheck[];
	overall_valid: boolean;
}

export function SetupWizardPage({ onComplete, onCancel }: SetupWizardPageProps) {
	const { refreshDomains, setSelectedDomain, setSelectedRange } = useGlobalFilter();
	const [step, setStep] = useState(1);

	// Step 1: Fail2Ban Status
	const [f2bStatus, setF2bStatus] = useState<Fail2BanStatus | null>(null);
	const [loadingF2b, setLoadingF2b] = useState(false);

	// Step 2: Domain Discovery
	const [discoveredDomains, setDiscoveredDomains] = useState<DiscoveredDomain[]>([]);
	const [loadingDiscovery, setLoadingDiscovery] = useState(false);
	const [selectedDomainName, setSelectedDomainName] = useState("");
	const [isManualDomain, setIsManualDomain] = useState(false);
	const [manualDomainName, setManualDomainName] = useState("");

	// Step 3: Settings & Snippet Generation
	const [rateLimit, setRateLimit] = useState(5);
	const [burstSize, setBurstSize] = useState(5);
	const [banTime, setBanTime] = useState(86400); // 24 hours in seconds
	const [generatedConfigs, setGeneratedConfigs] = useState<GeneratedConfigResponse | null>(null);
	const [generating, setGenerating] = useState(false);

	// Step 4: Nginx snippets applied confirmation
	const [nginxApplied, setNginxApplied] = useState(false);

	// Step 5: Validation & Activation
	const [validating, setValidating] = useState(false);
	const [validationChecks, setValidationChecks] = useState<ValidationCheck[]>([]);
	const [validationPassed, setValidationPassed] = useState(false);
	const [saving, setSaving] = useState(false);
	const [saveError, setSaveError] = useState<string | null>(null);

	const activeDomain = isManualDomain ? manualDomainName : selectedDomainName;

	// Load Fail2Ban Status
	const fetchF2bStatus = async () => {
		setLoadingF2b(true);
		try {
			const res = await api<Fail2BanStatus>("/api/setup/fail2ban-status");
			setF2bStatus(res);
		} catch (err) {
			console.error("Failed to load Fail2Ban status", err);
		} finally {
			setLoadingF2b(false);
		}
	};

	// Discover domains
	const discoverDomains = async () => {
		setLoadingDiscovery(true);
		try {
			const res = await api<{ domains: DiscoveredDomain[] }>("/api/setup/discover-domains");
			setDiscoveredDomains(res.domains || []);
			if (res.domains && res.domains.length > 0) {
				setSelectedDomainName(res.domains[0].server_name);
			} else {
				setIsManualDomain(true);
			}
		} catch (err) {
			console.error("Failed to discover domains", err);
			setIsManualDomain(true);
		} finally {
			setLoadingDiscovery(false);
		}
	};

	useEffect(() => {
		if (step === 1) {
			fetchF2bStatus();
		} else if (step === 2) {
			discoverDomains();
		}
	}, [step]);

	// Generate configs
	const handleGenerateConfigs = async () => {
		if (!activeDomain) return;
		setGenerating(true);
		try {
			const payload = {
				domain_name: activeDomain,
				rate_limit: rateLimit,
				burst_size: burstSize,
				ban_time: banTime
			};
			const res = await api<GeneratedConfigResponse>("/api/setup/generate-config", {
				method: "POST",
				body: JSON.stringify(payload)
			});
			setGeneratedConfigs(res);
			setStep(3);
		} catch (err) {
			alert(err instanceof Error ? err.message : "Failed to generate configuration files.");
		} finally {
			setGenerating(false);
		}
	};

	// Validate setup
	const runValidation = async () => {
		if (!activeDomain) return;
		setValidating(true);
		setSaveError(null);
		try {
			const payload = { domain_name: activeDomain };
			const res = await api<ValidationResponse>("/api/setup/validate", {
				method: "POST",
				body: JSON.stringify(payload)
			});
			setValidationChecks(res.checks || []);
			setValidationPassed(res.overall_valid);

			if (res.overall_valid) {
				// Auto save domain record
				await saveDomainRecord();
			}
		} catch (err) {
			setSaveError(err instanceof Error ? err.message : "Validation failed.");
		} finally {
			setValidating(false);
		}
	};

	// Save domain record to DB
	const saveDomainRecord = async () => {
		if (!activeDomain || !generatedConfigs) return;
		setSaving(true);
		try {
			const slug = generatedConfigs.domain_slug;
			const payload = {
				domain_name: activeDomain,
				access_log_path: `/var/log/nginx/${slug}_access.log`,
				error_log_path: `/var/log/nginx/${slug}_error.log`,
				blocked_ip_file_path: `/etc/nginx/${slug}_blocked.conf`,
				fail2ban_jail_name: `${slug}-429`,
				server_name: activeDomain,
				description: "Auto-configured protection via ShieldWatch Setup Wizard",
				rate_limit: rateLimit,
				burst_size: burstSize,
				ban_time: banTime,
				generated_config: JSON.stringify(generatedConfigs),
				status: "active"
			};

			const savedDomain = await api<{ id: number }>("/api/domains", {
				method: "POST",
				body: JSON.stringify(payload)
			});
			await refreshDomains();
			setSelectedDomain(savedDomain.id);
			setSelectedRange("last_24h");
			// Complete wizard! Success screen will trigger next step
			setStep(6);
			return savedDomain.id;
		} catch (err) {
			setSaveError(err instanceof Error ? err.message : "Failed to save domain database record.");
			setValidationPassed(false); // revert so they can retry
		} finally {
			setSaving(false);
		}
	};

	const handleSuccessComplete = () => {
		// Callback to App shell
		onComplete(0);
	};

	const progressPercentage = (Math.min(step, 5) / 5) * 100;

	return (
		<div className="max-w-4xl mx-auto space-y-8 py-4 select-none">
			{/* Breadcrumbs / Progress */}
			{step <= 5 && (
				<div className="space-y-4">
					<div className="flex items-center justify-between">
						<div className="space-y-1">
							<span className="text-xs font-semibold uppercase tracking-wider text-primary">
								Setup Onboarding
							</span>
							<h1 className="text-2xl font-bold tracking-tight">Configure New Domain Protection</h1>
						</div>
						<span className="text-sm font-medium text-muted-foreground">Step {step} of 5</span>
					</div>
					<div className="h-2 w-full rounded-full bg-muted overflow-hidden">
						<div
							className="h-full bg-gradient-to-r from-primary to-blue-500 rounded-full transition-all duration-300 ease-out"
							style={{ width: `${progressPercentage}%` }}
						/>
					</div>
				</div>
			)}

			<Card className="p-6 md:p-8 border border-border/80 bg-card/60 backdrop-blur-md shadow-xl rounded-2xl relative overflow-hidden">
				{/* Step 1: Fail2Ban Status */}
				{step === 1 && (
					<div className="space-y-6">
						<div className="flex items-center gap-3 border-b pb-4">
							<div className="h-10 w-10 rounded-lg bg-primary/10 text-primary flex items-center justify-center flex-shrink-0">
								<Cpu className="h-5 w-5" />
							</div>
							<div>
								<h2 className="text-lg font-bold text-foreground">Detect Fail2Ban Status</h2>
								<p className="text-sm text-muted-foreground">We need Fail2Ban installed and running on the host system.</p>
							</div>
						</div>

						{loadingF2b ? (
							<div className="flex flex-col items-center justify-center py-12 space-y-3">
								<Loader2 className="h-8 w-8 text-primary animate-spin" />
								<p className="text-sm text-muted-foreground">Checking Fail2Ban service status...</p>
							</div>
						) : f2bStatus ? (
							<div className="space-y-6">
								<div className="flex flex-col sm:flex-row items-start sm:items-center justify-between gap-4 p-4 rounded-xl border border-border bg-muted/30">
									<div className="flex items-center gap-3">
										<div
											className={`h-3 w-3 rounded-full animate-pulse ${
												f2bStatus.running ? "bg-emerald-500" : "bg-red-500"
											}`}
										/>
										<div>
											<span className="font-semibold text-sm">
												Fail2Ban Service: {f2bStatus.running ? "Running" : "Stopped"}
											</span>
											<p className="text-xs text-muted-foreground mt-0.5">
												{f2bStatus.installed
													? `Version ${f2bStatus.version} detected`
													: "Service not installed"}
											</p>
										</div>
									</div>
									<Button variant="ghost" size="sm" onClick={fetchF2bStatus} className="h-8 flex items-center gap-1 text-xs">
										<RefreshCw className="h-3 w-3" /> Refresh
									</Button>
								</div>

								{f2bStatus.running ? (
									<div className="space-y-4">
										<div className="rounded-lg bg-emerald-500/10 border border-emerald-500/20 p-4 flex gap-3 text-emerald-500">
											<CheckCircle className="h-5 w-5 flex-shrink-0 mt-0.5" />
											<div className="text-sm">
												<h4 className="font-semibold">Ready to Continue</h4>
												<p className="text-xs text-emerald-500/80 mt-1 leading-normal">
													Fail2Ban daemon is running on your server. Currently has{" "}
													<strong className="font-bold">{f2bStatus.jail_count}</strong> active jails:{" "}
													<code className="font-mono bg-emerald-500/20 px-1 py-0.5 rounded">
														{f2bStatus.active_jails.join(", ") || "None"}
													</code>
												</p>
											</div>
										</div>

										<div className="flex justify-end gap-3 pt-4 border-t">
											<Button variant="ghost" onClick={onCancel}>
												Cancel
											</Button>
											<Button variant="primary" onClick={() => setStep(2)} className="flex items-center gap-1.5">
												Continue <ChevronRight className="h-4 w-4" />
											</Button>
										</div>
									</div>
								) : (
									<div className="space-y-6">
										<div className="rounded-lg bg-amber-500/10 border border-amber-500/20 p-4 flex gap-3 text-amber-500">
											<AlertTriangle className="h-5 w-5 flex-shrink-0 mt-0.5" />
											<div className="text-sm">
												<h4 className="font-semibold">Action Required: Install/Start Fail2Ban</h4>
												<p className="text-xs text-amber-500/80 mt-1 leading-normal">
													We couldn't connect to Fail2Ban. If it's not installed, run the installation script below on the host. If running in containerized dev mode, you can bypass this step.
												</p>
											</div>
										</div>

										<div className="space-y-2">
											<label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
												Install Fail2Ban (Debian/Ubuntu)
											</label>
											<div className="flex items-center justify-between rounded-lg border border-border bg-[#0f141c] p-3 font-mono text-xs text-[#a5b4fc]">
												<code>sudo apt update && sudo apt install fail2ban -y</code>
												<Button
													size="icon"
													variant="ghost"
													className="h-8 w-8 hover:bg-[#1e293b] text-[#94a3b8] hover:text-white"
													onClick={() =>
														navigator.clipboard.writeText(
															"sudo apt update && sudo apt install fail2ban -y"
														)
													}
												>
													<Copy className="h-4 w-4" />
												</Button>
											</div>
										</div>

										<div className="flex items-center justify-between gap-3 pt-4 border-t">
											<Button variant="ghost" onClick={() => setStep(2)} className="text-xs text-muted-foreground hover:text-foreground">
												Bypass Check (Dev Mode)
											</Button>
											<div className="flex gap-3">
												<Button variant="ghost" onClick={onCancel}>
													Cancel
												</Button>
												<Button variant="primary" disabled className="flex items-center gap-1.5">
													Continue <ChevronRight className="h-4 w-4" />
												</Button>
											</div>
										</div>
									</div>
								)}
							</div>
						) : null}
					</div>
				)}

				{/* Step 2: Domain Discovery */}
				{step === 2 && (
					<div className="space-y-6 animate-fadeInUp">
						<div className="flex items-center gap-3 border-b pb-4">
							<div className="h-10 w-10 rounded-lg bg-primary/10 text-primary flex items-center justify-center flex-shrink-0">
								<Globe className="h-5 w-5" />
							</div>
							<div>
								<h2 className="text-lg font-bold text-foreground">Select Domain to Protect</h2>
								<p className="text-sm text-muted-foreground">
									Scan active Nginx site configurations for domains, or enter a domain name manually.
								</p>
							</div>
						</div>

						{loadingDiscovery ? (
							<div className="flex flex-col items-center justify-center py-12 space-y-3">
								<Loader2 className="h-8 w-8 text-primary animate-spin" />
								<p className="text-sm text-muted-foreground">Scanning Nginx configs for virtual hosts...</p>
							</div>
						) : (
							<div className="space-y-6">
								{/* Manual Entry Toggle */}
								<div className="flex items-center justify-between p-3 rounded-lg border border-border/80 bg-muted/10">
									<span className="text-sm font-semibold">Configure domain manually</span>
									<input
										type="checkbox"
										checked={isManualDomain}
										onChange={(e) => setIsManualDomain(e.target.checked)}
										className="h-4 w-4 rounded border-border text-primary focus:ring-primary focus:outline-none"
									/>
								</div>

								{isManualDomain ? (
									<div className="space-y-2 max-w-md">
										<label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
											Domain Name
										</label>
										<Input
											value={manualDomainName}
											onChange={(e) => setManualDomainName(e.target.value)}
											placeholder="e.g. example.com"
										/>
									</div>
								) : (
									<div className="space-y-4">
										<label className="text-xs font-semibold uppercase tracking-wider text-muted-foreground">
											Discovered Nginx Server Blocks
										</label>

										{discoveredDomains.length === 0 ? (
											<div className="p-8 text-center border border-dashed rounded-xl border-border bg-muted/10 text-muted-foreground">
												<AlertTriangle className="h-8 w-8 mx-auto text-amber-500/80 mb-2 animate-bounce" />
												<p className="text-sm font-medium">No active Nginx domains discovered.</p>
												<p className="text-xs text-muted-foreground mt-1">
													Ensure your Nginx server files are mounted correctly or configure domain manually.
												</p>
											</div>
										) : (
											<div className="grid grid-cols-1 sm:grid-cols-2 gap-4 max-h-[260px] overflow-y-auto pr-1">
												{discoveredDomains.map((d) => (
													<button
														key={d.server_name}
														type="button"
														onClick={() => setSelectedDomainName(d.server_name)}
														className={`flex flex-col items-start p-4 rounded-xl border text-left transition-all ${
															selectedDomainName === d.server_name
																? "border-primary bg-primary/5 ring-1 ring-primary/30"
																: "border-border hover:bg-muted/40"
														}`}
													>
														<span className="font-bold text-sm text-foreground">
															{d.server_name}
														</span>
														<span className="text-xs text-muted-foreground mt-1 truncate w-full">
															Config: {d.config_file}
														</span>
														<div className="flex gap-2 mt-3.5">
															{d.has_ssl && (
																<Badge tone="success" className="text-[10px] px-1.5 py-0.5">
																	SSL
																</Badge>
															)}
															<Badge tone="info" className="text-[10px] px-1.5 py-0.5">
																Discovered
															</Badge>
														</div>
													</button>
												))}
											</div>
										)}
									</div>
								)}

								<div className="flex justify-between items-center pt-4 border-t">
									<Button variant="ghost" onClick={() => setStep(1)} className="flex items-center gap-1.5">
										<ChevronLeft className="h-4 w-4" /> Back
									</Button>
									<Button
										variant="primary"
										onClick={handleGenerateConfigs}
										disabled={!activeDomain}
										className="flex items-center gap-1.5"
									>
										{generating ? (
											<Loader2 className="h-4 w-4 animate-spin" />
										) : (
											<>
												Configure Domain <ChevronRight className="h-4 w-4" />
											</>
										)}
									</Button>
								</div>
							</div>
						)}
					</div>
				)}

				{/* Step 3: Settings & Configuration Generation */}
				{step === 3 && (
					<div className="space-y-6 animate-fadeInUp">
						<div className="flex items-center gap-3 border-b pb-4">
							<div className="h-10 w-10 rounded-lg bg-primary/10 text-primary flex items-center justify-center flex-shrink-0">
								<Settings2 className="h-5 w-5" />
							</div>
							<div>
								<h2 className="text-lg font-bold text-foreground">Fail2Ban Protection Setup</h2>
								<p className="text-sm text-muted-foreground">
									Adjust rate limits and download your customized configuration files & scripts.
								</p>
							</div>
						</div>

						{/* Settings inputs with help tooltips */}
						<div className="grid grid-cols-1 md:grid-cols-3 gap-6 p-4 rounded-xl border border-border/80 bg-muted/10">
							<div className="space-y-2">
								<div className="flex items-center gap-1.5">
									<label className="text-xs font-bold uppercase text-muted-foreground">
										Rate Limit (r/s)
									</label>
									<span className="group relative cursor-help text-muted-foreground hover:text-foreground">
										<HelpCircle className="h-3.5 w-3.5" />
										<span className="pointer-events-none absolute bottom-full left-1/2 z-50 mb-2 w-64 -translate-x-1/2 rounded-md border bg-card p-3 text-[10px] text-foreground shadow-lg opacity-0 transition-opacity group-hover:opacity-100 leading-normal font-normal">
											The maximum number of requests allowed from a single client IP address per second. Excess requests beyond this will be blocked with a HTTP 429 status.
										</span>
									</span>
								</div>
								<Input
									type="number"
									min={1}
									max={1000}
									value={rateLimit}
									onChange={(e) => {
										setRateLimit(Number(e.target.value));
										setGeneratedConfigs(null); // invalidate generated
									}}
									className="bg-card border border-border rounded-md w-full px-3 py-2 text-sm focus:ring-1 focus:ring-primary focus:outline-none"
								/>
								<p className="text-[10px] text-muted-foreground">Requests allowed per client IP per second.</p>
							</div>

							<div className="space-y-2">
								<div className="flex items-center gap-1.5">
									<label className="text-xs font-bold uppercase text-muted-foreground">
										Burst Size
									</label>
									<span className="group relative cursor-help text-muted-foreground hover:text-foreground">
										<HelpCircle className="h-3.5 w-3.5" />
										<span className="pointer-events-none absolute bottom-full left-1/2 z-50 mb-2 w-64 -translate-x-1/2 rounded-md border bg-card p-3 text-[10px] text-foreground shadow-lg opacity-0 transition-opacity group-hover:opacity-100 leading-normal font-normal">
											The number of excess requests allowed to buffer before blocking occurs. This accommodates bursty legitimate browsing patterns.
										</span>
									</span>
								</div>
								<Input
									type="number"
									min={1}
									max={1000}
									value={burstSize}
									onChange={(e) => {
										setBurstSize(Number(e.target.value));
										setGeneratedConfigs(null); // invalidate generated
									}}
									className="bg-card border border-border rounded-md w-full px-3 py-2 text-sm focus:ring-1 focus:ring-primary focus:outline-none"
								/>
								<p className="text-[10px] text-muted-foreground">Allow excess requests momentarily up to this count.</p>
							</div>

							<div className="space-y-2">
								<div className="flex items-center gap-1.5">
									<label className="text-xs font-bold uppercase text-muted-foreground">
										Ban Time (seconds)
									</label>
									<span className="group relative cursor-help text-muted-foreground hover:text-foreground">
										<HelpCircle className="h-3.5 w-3.5" />
										<span className="pointer-events-none absolute bottom-full left-1/2 z-50 mb-2 w-64 -translate-x-1/2 rounded-md border bg-card p-3 text-[10px] text-foreground shadow-lg opacity-0 transition-opacity group-hover:opacity-100 leading-normal font-normal">
											The duration in seconds that client IP remains blocked in Nginx (e.g. 3600 for 1 hour, 86400 for 24 hours).
										</span>
									</span>
								</div>
								<Input
									type="number"
									min={10}
									max={31536000}
									value={banTime}
									onChange={(e) => {
										setBanTime(Number(e.target.value));
										setGeneratedConfigs(null); // invalidate generated
									}}
									className="bg-card border border-border rounded-md w-full px-3 py-2 text-sm focus:ring-1 focus:ring-primary focus:outline-none"
								/>
								<p className="text-[10px] text-muted-foreground">Duration that client IP remains blocked in Nginx.</p>
							</div>
						</div>

						{/* If settings changed, show regenerate option */}
						{!generatedConfigs ? (
							<div className="flex flex-col items-center justify-center p-8 border border-dashed border-border rounded-xl bg-card/25 gap-3">
								<p className="text-xs text-muted-foreground font-medium">Settings have changed. Please regenerate configurations to continue.</p>
								<Button variant="primary" onClick={handleGenerateConfigs} className="flex items-center gap-1.5 h-10 px-4">
									{generating ? <Loader2 className="h-4 w-4 animate-spin" /> : <RefreshCw className="h-4 w-4" />}
									Generate Configuration Snippets
								</Button>
							</div>
						) : (
							<div className="space-y-4">
								<div className="flex items-center justify-between">
									<span className="text-xs font-bold uppercase tracking-wider text-muted-foreground">
										Generated Configurations & Script
									</span>
									<span className="text-xs text-muted-foreground font-medium">
										Run setup script on the host to configure automatically
									</span>
								</div>

								{/* Tabbed Code Viewer */}
								<CodeViewer
									files={[
										{
											name: `setup-${generatedConfigs.domain_slug}.sh`,
											content: generatedConfigs.setup_script
										},
										...generatedConfigs.files.map((f) => ({
											name: f.filename,
											content: f.content
										}))
									]}
								/>

								<div className="p-4 rounded-xl border border-primary/20 bg-primary/5 text-sm flex gap-3 text-primary">
									<Lock className="h-5 w-5 flex-shrink-0 mt-0.5" />
									<div className="text-xs leading-relaxed">
										<p className="font-semibold">Manual Action Required on Host Server</p>
										<p className="mt-1">
											To automatically generate and write files under <code className="font-mono bg-primary/20 px-1 py-0.5 rounded">/etc/fail2ban/</code>, copy/download the setup script above, save it as <code className="font-mono bg-primary/20 px-1 py-0.5 rounded">setup.sh</code> on your server host, and execute: <code className="font-mono bg-primary/20 px-1.5 py-0.5 rounded select-all">sudo bash setup.sh</code>.
										</p>
									</div>
								</div>

								<div className="flex justify-between items-center pt-4 border-t">
									<Button variant="ghost" onClick={() => setStep(2)} className="flex items-center gap-1.5">
										<ChevronLeft className="h-4 w-4" /> Back
									</Button>
									<Button variant="primary" onClick={() => setStep(4)} className="flex items-center gap-1.5">
										I've Run the Script & Continue <ChevronRight className="h-4 w-4" />
									</Button>
								</div>
							</div>
						)}
					</div>
				)}

				{/* Step 4: Nginx Configurations */}
				{step === 4 && generatedConfigs && (
					<div className="space-y-6 animate-fadeInUp">
						<div className="flex items-center gap-3 border-b pb-4">
							<div className="h-10 w-10 rounded-lg bg-primary/10 text-primary flex items-center justify-center flex-shrink-0">
								<Globe className="h-5 w-5" />
							</div>
							<div>
								<h2 className="text-lg font-bold text-foreground">Update Nginx Configuration</h2>
								<p className="text-sm text-muted-foreground">
									Paste the generated rate limit rules inside Nginx server and http blocks.
								</p>
							</div>
						</div>

						<div className="space-y-6">
							{/* http block instructions */}
							<div className="space-y-3">
								<h3 className="text-sm font-semibold text-foreground">
									1. Add Zone Limit to Nginx Global Config
								</h3>
								<p className="text-xs text-muted-foreground leading-normal">
									Add this line inside the <code className="font-mono bg-muted p-0.5 px-1 rounded">http &#123; ... &#125;</code> block in <code className="font-mono bg-muted p-0.5 px-1 rounded">/etc/nginx/nginx.conf</code>:
								</p>
								<CodeViewer
									showDownload={false}
									files={[
										{
											name: "nginx.conf (http block)",
											content: generatedConfigs.nginx_zone_line
										}
									]}
								/>
							</div>

							{/* server block instructions */}
							<div className="space-y-3">
								<h3 className="text-sm font-semibold text-foreground">
									2. Add Protection Snippet to Virtual Host
								</h3>
								<p className="text-xs text-muted-foreground leading-normal">
									Add these directives inside the <code className="font-mono bg-muted p-0.5 px-1 rounded">server &#123; ... &#125;</code> block of your domain configuration file (e.g. <code className="font-mono bg-muted p-0.5 px-1 rounded">/etc/nginx/sites-enabled/{activeDomain}.conf</code>):
								</p>
								<CodeViewer
									showDownload={false}
									files={[
										{
											name: `${activeDomain}.conf (server block)`,
											content: generatedConfigs.nginx_snippet
										}
									]}
								/>
							</div>

							<div className="rounded-xl border border-border/80 bg-muted/10 p-4">
								<label className="flex items-start gap-3 cursor-pointer select-none">
									<input
										type="checkbox"
										checked={nginxApplied}
										onChange={(e) => setNginxApplied(e.target.checked)}
										className="h-4 w-4 mt-0.5 rounded border-border text-primary focus:ring-primary focus:outline-none"
									/>
									<span className="text-xs text-foreground font-semibold leading-relaxed">
										I have pasted both Nginx config snippets, and tested & reloaded Nginx (e.g. via <code className="font-mono bg-muted px-1 rounded">sudo nginx -t && sudo nginx -s reload</code>) on the host.
									</span>
								</label>
							</div>

							<div className="flex justify-between items-center pt-4 border-t">
								<Button variant="ghost" onClick={() => setStep(3)} className="flex items-center gap-1.5">
									<ChevronLeft className="h-4 w-4" /> Back
								</Button>
								<Button
									variant="primary"
									disabled={!nginxApplied}
									onClick={() => {
										setStep(5);
										runValidation();
									}}
									className="flex items-center gap-1.5"
								>
									I Have Updated Nginx <ChevronRight className="h-4 w-4" />
								</Button>
							</div>
						</div>
					</div>
				)}

				{/* Step 5: Validation & Activation */}
				{step === 5 && (
					<div className="space-y-6 animate-fadeInUp">
						<div className="flex items-center gap-3 border-b pb-4">
							<div className="h-10 w-10 rounded-lg bg-primary/10 text-primary flex items-center justify-center flex-shrink-0">
								<Shield className="h-5 w-5" />
							</div>
							<div>
								<h2 className="text-lg font-bold text-foreground">Verify & Activate Protection</h2>
								<p className="text-sm text-muted-foreground">
									Dashboard is validating that all Fail2Ban configs and Nginx variables are configured correctly.
								</p>
							</div>
						</div>

						{validating ? (
							<div className="flex flex-col items-center justify-center py-12 space-y-4">
								<Loader2 className="h-10 w-10 text-primary animate-spin" />
								<p className="text-sm text-muted-foreground">Running verification checks on host configuration...</p>
							</div>
						) : (
							<div className="space-y-6">
								{saveError && (
									<div className="p-4 rounded-xl border border-red-500/20 bg-red-500/10 text-sm text-red-500 flex gap-2">
										<AlertTriangle className="h-5 w-5 flex-shrink-0" />
										<div>
											<h4 className="font-semibold">Setup Error</h4>
											<p className="text-xs text-red-500/80 mt-1 leading-normal">{saveError}</p>
										</div>
									</div>
								)}

								<div className="space-y-3">
									<span className="text-xs font-bold uppercase tracking-wider text-muted-foreground">
										Verification Checklist
									</span>
									<div className="space-y-2">
										{validationChecks.map((check, idx) => (
											<div
												key={idx}
												className={`flex items-start justify-between p-3.5 rounded-xl border transition-all ${
													check.passed
														? "border-emerald-500/20 bg-emerald-500/5 text-emerald-500"
														: "border-red-500/20 bg-red-500/5 text-red-500"
												}`}
											>
												<div className="flex gap-3">
													{check.passed ? (
														<CheckCircle className="h-5 w-5 flex-shrink-0 mt-0.5" />
													) : (
														<AlertTriangle className="h-5 w-5 flex-shrink-0 mt-0.5" />
													)}
													<div className="text-xs">
														<span className="font-bold block text-foreground">
															{check.name}
														</span>
														<span className="text-muted-foreground mt-0.5 block">
															{check.message}
														</span>
													</div>
												</div>
												<Badge tone={check.passed ? "success" : "danger"} className="text-[10px]">
													{check.passed ? "Passed" : "Failed"}
												</Badge>
											</div>
										))}
									</div>
								</div>

								{validationPassed ? (
									<div className="rounded-xl border border-emerald-500/25 bg-emerald-500/10 p-4 flex gap-3 text-emerald-500 text-sm animate-pulse-glow">
										<CheckCircle className="h-5 w-5 flex-shrink-0 mt-0.5" />
										<div>
											<h4 className="font-semibold">All Verification Checks Passed!</h4>
											<p className="text-xs text-emerald-500/80 mt-1 leading-normal">
												Protection variables validated. Saving configurations and registering domain in dashboard...
											</p>
										</div>
									</div>
								) : (
									<div className="rounded-xl border border-amber-500/20 bg-amber-500/5 p-4 flex gap-3 text-amber-500 text-sm">
										<AlertTriangle className="h-5 w-5 flex-shrink-0 mt-0.5 animate-bounce" />
										<div>
											<h4 className="font-semibold">Validation Checklist Failed</h4>
											<p className="text-xs text-amber-500/80 mt-1 leading-normal">
												One or more configuration checks did not pass. Please verify that you've run the setup script and reloaded Nginx on the host correctly.
											</p>
										</div>
									</div>
								)}

								<div className="flex justify-between items-center pt-4 border-t">
									<Button variant="ghost" onClick={() => setStep(4)} className="flex items-center gap-1.5">
										<ChevronLeft className="h-4 w-4" /> Back
									</Button>
									<div className="flex gap-2">
										<Button variant="ghost" onClick={runValidation} className="flex items-center gap-1 text-xs">
											<RefreshCw className="h-3 w-3" /> Retry Check
										</Button>
										{!validationPassed && (
											<Button
												variant="ghost"
												onClick={() => {
													setSaveError(null);
													saveDomainRecord();
												}}
												className="text-xs text-muted-foreground hover:text-foreground border border-border"
											>
												Bypass & Force Save
											</Button>
										)}
									</div>
								</div>
							</div>
						)}
					</div>
				)}

				{/* Success Screen */}
				{step === 6 && (
					<div className="space-y-8 text-center py-6 animate-fadeInUp select-none">
						<Confetti />
						<div className="flex flex-col items-center space-y-4">
							<div className="relative flex h-24 w-24 items-center justify-center rounded-full bg-gradient-to-tr from-emerald-500 to-teal-400 text-white shadow-2xl shadow-emerald-500/20 scale-110">
								<Check className="h-12 w-12" />
								<div className="absolute -inset-2 rounded-full bg-emerald-500 opacity-20 blur-md -z-10 animate-ping" style={{ animationDuration: "3s" }} />
							</div>
							<div className="flex items-center gap-1 text-xs font-semibold tracking-wider text-emerald-500 uppercase mt-4">
								<Lock className="h-3.5 w-3.5" />
								Server Secured
							</div>
						</div>

						<div className="space-y-2">
							<h1 className="text-3xl font-extrabold tracking-tight bg-clip-text text-transparent bg-gradient-to-b from-white to-neutral-400 leading-tight">
								Protection Enabled
							</h1>
							<p className="text-sm text-emerald-500/80 font-bold max-w-sm mx-auto">
								{activeDomain} is now fully protected!
							</p>
							<p className="text-xs text-muted-foreground max-w-md mx-auto leading-relaxed pt-2">
								Rate limiting (5 requests/sec) and automatic 24-hour IP banning are active. Client violations and ban analytics will now feed into your console.
							</p>
						</div>

						<div className="pt-6 flex justify-center">
							<Button
								variant="primary"
								size="md"
								onClick={handleSuccessComplete}
								className="h-12 px-10 text-sm font-semibold tracking-wide shadow-lg shadow-emerald-500/25 hover:shadow-emerald-500/35 hover:-translate-y-0.5 active:translate-y-0 transition-all duration-200"
							>
								Go to Dashboard
							</Button>
						</div>
					</div>
				)}
			</Card>
		</div>
	);
}
