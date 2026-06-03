import { useState, useEffect } from "react";
import { useGlobalFilter, type Domain } from "../../context/GlobalFilterContext";
import { Card } from "../../components/ui/card";
import { Button } from "../../components/ui/button";
import { Badge } from "../../components/ui/badge";
import { CodeViewer } from "../../components/ui/code-viewer";

import { api } from "../../services/api";
import {
	AlertTriangle,
	ShieldAlert,
	Trash2,
	ChevronRight,
	ChevronLeft,
	RefreshCw,
	Loader2,
	CheckCircle,
	X,
	ArrowLeft
} from "lucide-react";

interface RemovalWizardProps {
	domain: Domain;
	onComplete: () => void;
	onCancel: () => void;
}

interface CleanupScriptResponse {
	domain_slug: string;
	cleanup_script: string;
	files_to_remove: Array<{ filename: string; path: string }>;
}

interface ValidationCheck {
	name: string;
	passed: boolean;
	message: string;
}

interface RemovalValidationResponse {
	checks: ValidationCheck[];
	overall_valid: boolean;
}

export function RemovalWizard({ domain, onComplete, onCancel }: RemovalWizardProps) {
	const { refreshDomains } = useGlobalFilter();
	const [step, setStep] = useState(1);

	// Step 2: Nginx cleanup
	const [nginxChecked, setNginxChecked] = useState(false);
	const [validatingNginx, setValidatingNginx] = useState(false);
	const [nginxValidationResult, setNginxValidationResult] = useState<RemovalValidationResponse | null>(null);
	const [nginxError, setNginxError] = useState<string | null>(null);

	// Step 3: Fail2Ban cleanup script
	const [cleanupData, setCleanupData] = useState<CleanupScriptResponse | null>(null);
	const [loadingCleanup, setLoadingCleanup] = useState(false);
	const [cb1, setCb1] = useState(false);
	const [cb2, setCb2] = useState(false);
	const [cb3, setCb3] = useState(false);
	const [cb4, setCb4] = useState(false);

	// Final Step: Deleting
	const [deleting, setDeleting] = useState(false);
	const [deleteError, setDeleteError] = useState<string | null>(null);

	const slug = domain.domain_name.replace(/\./g, "-").replace(/_/g, "-");

	// Parse previous config if stored, else fallback
	let previousNginxSnippet = "";
	let previousNginxZoneLine = "";
	try {
		if (domain.generated_config) {
			const parsed = JSON.parse(domain.generated_config);
			previousNginxSnippet = parsed.nginx_snippet || "";
			previousNginxZoneLine = parsed.nginx_zone_line || "";
		}
	} catch (e) {
		console.error("Failed to parse generated_config", e);
	}

	// Fallback snippets if they aren't stored
	if (!previousNginxSnippet) {
		previousNginxSnippet = `# ── Fail2ban auto-blocked IPs check ──\nif ($${slug.replace(/-/g, "_")}_blocked) {\n    return 403;\n}\n\nlimit_req zone=${slug}_rate_limit burst=5 nodelay;\nlimit_req_status 429;`;
		previousNginxZoneLine = `limit_req_zone $binary_remote_addr zone=${slug}_rate_limit:10m rate=5r/s;\n\ngeo $${slug.replace(/-/g, "_")}_blocked {\n    default 0;\n    include /etc/nginx/${slug}_blocked.conf;\n}`;
	}

	// Fetch Cleanup Script details
	const fetchCleanupScript = async () => {
		setLoadingCleanup(true);
		try {
			const res = await api<CleanupScriptResponse>("/api/setup/generate-cleanup", {
				method: "POST",
				body: JSON.stringify({ domain_name: domain.domain_name })
			});
			setCleanupData(res);
		} catch (err) {
			console.error("Failed to fetch cleanup script", err);
		} finally {
			setLoadingCleanup(false);
		}
	};

	useEffect(() => {
		if (step === 3) {
			fetchCleanupScript();
		}
	}, [step]);

	// Validate removal from Nginx
	const validateNginxRemoval = async () => {
		setValidatingNginx(true);
		setNginxError(null);
		try {
			const res = await api<RemovalValidationResponse>("/api/setup/validate-removal", {
				method: "POST",
				body: JSON.stringify({ domain_name: domain.domain_name })
			});
			setNginxValidationResult(res);
			if (res.overall_valid) {
				// Wait 1.5s to show successful verification, then advance
				setTimeout(() => {
					setStep(3);
				}, 1500);
			}
		} catch (err) {
			setNginxError(err instanceof Error ? err.message : "Validation failed.");
		} finally {
			setValidatingNginx(false);
		}
	};

	// Final Deletion
	const handleDeleteProtection = async () => {
		setDeleting(true);
		setDeleteError(null);
		try {
			// Trigger DB deletion
			await api(`/api/domains/${domain.id}`, { method: "DELETE" });
			await refreshDomains();
			setStep(4); // Success screen!
		} catch (err) {
			setDeleteError(err instanceof Error ? err.message : "Failed to delete domain record.");
		} finally {
			setDeleting(false);
		}
	};

	const allCheckboxesChecked = cb1 && cb2 && cb3 && cb4;

	return (
		<div className="max-w-3xl mx-auto space-y-6 select-none animate-fadeInUp">
			{/* Header */}
			{step <= 3 && (
				<div className="flex items-center justify-between border-b pb-4">
					<div className="space-y-1">
						<div className="flex items-center gap-2 text-danger">
							<ShieldAlert className="h-5 w-5" />
							<span className="text-xs font-bold uppercase tracking-wider">
								Safe Domain Removal Workflow
							</span>
						</div>
						<h1 className="text-xl font-bold text-foreground">Remove Domain Protection: {domain.domain_name}</h1>
					</div>
					<Button variant="ghost" size="icon" onClick={onCancel}>
						<X className="h-4 w-4" />
					</Button>
				</div>
			)}

			<Card className="p-6 md:p-8 border border-border/80 bg-card/60 backdrop-blur-md shadow-2xl rounded-2xl">
				{/* Step 1: Warning Screen */}
				{step === 1 && (
					<div className="space-y-6">
						<div className="p-4 rounded-xl border border-red-500/20 bg-red-500/10 text-red-500 flex gap-3.5">
							<AlertTriangle className="h-6 w-6 flex-shrink-0 mt-0.5 animate-pulse" />
							<div className="text-sm">
								<h3 className="font-bold text-base">Warning: High Risk Operation</h3>
								<p className="text-xs text-red-500/80 mt-1 leading-normal">
									Removing protection requires manually cleaning Nginx settings before deleting. Immediate deletion will trigger Nginx configurations crashes or leave active jails orphaned.
								</p>
							</div>
						</div>

						<div className="space-y-3">
							<span className="text-xs font-bold uppercase tracking-wider text-muted-foreground">
								This workflow will guide you through:
							</span>
							<ul className="space-y-2.5 text-xs text-muted-foreground">
								<li className="flex items-start gap-2">
									<div className="h-4 w-4 rounded bg-red-500/20 text-red-500 flex items-center justify-center font-bold flex-shrink-0 mt-0.5">•</div>
									<span>Remove dashboard domain configuration and statistics logs.</span>
								</li>
								<li className="flex items-start gap-2">
									<div className="h-4 w-4 rounded bg-red-500/20 text-red-500 flex items-center justify-center font-bold flex-shrink-0 mt-0.5">•</div>
									<span>Clean up generated Fail2Ban jail definitions, action and filter rules.</span>
								</li>
								<li className="flex items-start gap-2">
									<div className="h-4 w-4 rounded bg-red-500/20 text-red-500 flex items-center justify-center font-bold flex-shrink-0 mt-0.5">•</div>
									<span>Remove Nginx geo-blocking file paths and restore config state.</span>
								</li>
							</ul>
						</div>

						<div className="p-4 rounded-xl border border-border/80 bg-muted/20 text-xs leading-relaxed text-muted-foreground">
							<strong className="font-semibold text-foreground">Important Note:</strong> Before continuing, you must remove the generated Nginx configuration directives from your virtual hosts. We will validate this in the next step.
						</div>

						<div className="flex justify-end gap-3 pt-4 border-t">
							<Button variant="ghost" onClick={onCancel}>
								Cancel
							</Button>
							<Button variant="primary" onClick={() => setStep(2)} className="bg-red-600 hover:bg-red-500 border-red-700 hover:border-red-600 text-white flex items-center gap-1">
								Continue <ChevronRight className="h-4 w-4" />
							</Button>
						</div>
					</div>
				)}

				{/* Step 2: Nginx Cleanup Instructions & Validation */}
				{step === 2 && (
					<div className="space-y-6">
						<div>
							<h2 className="text-lg font-bold text-foreground">Remove Nginx Configuration Snippets</h2>
							<p className="text-xs text-muted-foreground mt-1">
								Please search and delete the following snippets previously added during wizard setup.
							</p>
						</div>

						{nginxError && (
							<div className="p-4 rounded-xl border border-red-500/20 bg-red-500/10 text-xs text-red-500">
								{nginxError}
							</div>
						)}

						<div className="space-y-4">
							<div className="space-y-2">
								<span className="text-xs font-bold text-foreground block">
									1. Delete from virtual host configurations:
								</span>
								<CodeViewer showDownload={false} files={[{ name: `${domain.domain_name}.conf`, content: previousNginxSnippet }]} />
							</div>

							<div className="space-y-2">
								<span className="text-xs font-bold text-foreground block">
									2. Delete from global nginx.conf http block:
								</span>
								<CodeViewer showDownload={false} files={[{ name: "nginx.conf", content: previousNginxZoneLine }]} />
							</div>
						</div>

						<div className="p-4 rounded-xl border border-border/80 bg-muted/10">
							<label className="flex items-start gap-3 cursor-pointer select-none">
								<input
									type="checkbox"
									checked={nginxChecked}
									onChange={(e) => setNginxChecked(e.target.checked)}
									className="h-4 w-4 mt-0.5 rounded border-border text-primary focus:ring-primary focus:outline-none"
								/>
								<span className="text-xs text-foreground font-semibold leading-relaxed">
									I have manually removed all snippets and reloaded Nginx on the host server.
								</span>
							</label>
						</div>

						{nginxValidationResult && (
							<div className="space-y-2.5">
								{nginxValidationResult.checks.map((check, idx) => (
									<div
										key={idx}
										className={`p-3.5 rounded-xl border flex items-center justify-between text-xs ${
											check.passed
												? "border-emerald-500/25 bg-emerald-500/5 text-emerald-500"
												: "border-red-500/25 bg-red-500/5 text-red-500"
										}`}
									>
										<div className="flex gap-2.5">
											{check.passed ? <CheckCircle className="h-4 w-4 mt-0.5" /> : <AlertTriangle className="h-4 w-4 mt-0.5" />}
											<div>
												<span className="font-bold text-foreground block">{check.name}</span>
												<span className="text-muted-foreground block mt-0.5">{check.message}</span>
											</div>
										</div>
										<Badge tone={check.passed ? "success" : "danger"}>
											{check.passed ? "Clean" : "Reference Found"}
										</Badge>
									</div>
								))}
							</div>
						)}

						<div className="flex justify-between items-center pt-4 border-t">
							<Button variant="ghost" onClick={() => setStep(1)} className="flex items-center gap-1.5">
								<ChevronLeft className="h-4 w-4" /> Back
							</Button>
							<Button
								variant="primary"
								disabled={!nginxChecked || validatingNginx}
								onClick={validateNginxRemoval}
								className="flex items-center gap-1.5"
							>
								{validatingNginx ? (
									<>
										<Loader2 className="h-4 w-4 animate-spin" /> Validating Removal...
									</>
								) : (
									<>
										Validate Removal <ChevronRight className="h-4 w-4" />
									</>
								)}
							</Button>
						</div>
					</div>
				)}

				{/* Step 3: Fail2Ban Cleanup & Automation Script */}
				{step === 3 && (
					<div className="space-y-6">
						<div>
							<h2 className="text-lg font-bold text-foreground">Clean Up Fail2Ban Resources</h2>
							<p className="text-xs text-muted-foreground mt-1">
								Generate and run the cleanup script to remove all generated jail/block files on the host.
							</p>
						</div>

						{loadingCleanup ? (
							<div className="flex flex-col items-center justify-center py-12 space-y-3">
								<Loader2 className="h-8 w-8 text-primary animate-spin" />
								<p className="text-xs text-muted-foreground">Generating cleanup script...</p>
							</div>
						) : cleanupData ? (
							<div className="space-y-6">
								<div className="space-y-3">
									<span className="text-xs font-bold text-foreground block">
										The following files must be deleted:
									</span>
									<div className="grid grid-cols-1 gap-2">
										{cleanupData.files_to_remove.map((file, idx) => (
											<div key={idx} className="flex justify-between items-center p-2 px-3.5 rounded-lg border border-border/80 bg-muted/10 font-mono text-xs">
												<span className="text-muted-foreground font-semibold">{file.filename}</span>
												<span className="text-[10px] text-muted-foreground truncate max-w-sm">{file.path}</span>
											</div>
										))}
									</div>
								</div>

								{/* Tabbed script viewer */}
								<div className="space-y-2">
									<span className="text-xs font-bold uppercase tracking-wider text-muted-foreground">
										Cleanup Script
									</span>
									<CodeViewer files={[{ name: `cleanup-${cleanupData.domain_slug}.sh`, content: cleanupData.cleanup_script }]} />
								</div>

								{/* A checklist representing confirmations */}
								<div className="p-4 rounded-xl border border-border/80 bg-muted/10 space-y-3 text-xs">
									<span className="font-bold text-foreground block">Confirm execution tasks:</span>
									<div className="space-y-2.5">
										<label className="flex items-center gap-3 cursor-pointer">
											<input type="checkbox" checked={cb1} onChange={(e) => setCb1(e.target.checked)} className="h-4 w-4 rounded text-primary focus:ring-primary focus:outline-none" />
											<span>I have run the cleanup script on host to delete Fail2Ban jail files.</span>
										</label>
										<label className="flex items-center gap-3 cursor-pointer">
											<input type="checkbox" checked={cb2} onChange={(e) => setCb2(e.target.checked)} className="h-4 w-4 rounded text-primary focus:ring-primary focus:outline-none" />
											<span>I have deleted the geo IP blocked file for this domain slug.</span>
										</label>
										<label className="flex items-center gap-3 cursor-pointer">
											<input type="checkbox" checked={cb3} onChange={(e) => setCb3(e.target.checked)} className="h-4 w-4 rounded text-primary focus:ring-primary focus:outline-none" />
											<span>I restarted/reloaded Fail2Ban successfully on the host server.</span>
										</label>
										<label className="flex items-center gap-3 cursor-pointer">
											<input type="checkbox" checked={cb4} onChange={(e) => setCb4(e.target.checked)} className="h-4 w-4 rounded text-primary focus:ring-primary focus:outline-none" />
											<span>I validated Nginx configuration is running clean (no error exits).</span>
										</label>
									</div>
								</div>

								{deleteError && (
									<div className="p-4 rounded-xl border border-red-500/25 bg-red-500/10 text-xs text-red-500">
										{deleteError}
									</div>
								)}

								<div className="flex justify-between items-center pt-4 border-t">
									<Button variant="ghost" onClick={() => setStep(2)} className="flex items-center gap-1.5">
										<ChevronLeft className="h-4 w-4" /> Back
									</Button>
									<Button
										variant="primary"
										disabled={!allCheckboxesChecked || deleting}
										onClick={handleDeleteProtection}
										className="bg-red-600 hover:bg-red-500 border-red-700 hover:border-red-600 text-white flex items-center gap-1.5"
									>
										{deleting ? (
											<>
												<Loader2 className="h-4 w-4 animate-spin" /> Deleting Protection...
											</>
										) : (
											<>
												Delete Protection <Trash2 className="h-4 w-4" />
											</>
										)}
									</Button>
								</div>
							</div>
						) : null}
					</div>
				)}

				{/* Farewell Screen */}
				{step === 4 && (
					<div className="space-y-8 text-center py-8 select-none animate-fadeInUp">
						{/* Farewell Particles — fading embers drifting upward */}
						<FarewellParticles />

						<div className="flex flex-col items-center space-y-5">
							{/* Sad waving hand icon */}
							<div className="relative flex h-28 w-28 items-center justify-center rounded-full bg-gradient-to-br from-slate-600 via-indigo-800 to-slate-700 text-white shadow-2xl shadow-indigo-900/30">
								<span className="text-5xl" style={{ animation: "farewell-wave 2s ease-in-out infinite" }}>👋</span>
								<div className="absolute -inset-3 rounded-full bg-indigo-500 opacity-10 blur-xl -z-10" style={{ animation: "farewell-pulse 4s ease-in-out infinite" }} />
							</div>

							<div className="flex items-center gap-1.5 text-xs font-semibold tracking-widest text-indigo-400/80 uppercase mt-3">
								<span style={{ animation: "farewell-fade-letters 3s ease-in-out infinite" }}>Domain Released</span>
							</div>
						</div>

						<div className="space-y-3 max-w-md mx-auto">
							<h1 className="text-3xl font-extrabold tracking-tight bg-clip-text text-transparent bg-gradient-to-b from-slate-300 via-indigo-200 to-slate-500 leading-tight">
								Until We Meet Again
							</h1>
							<p className="text-sm font-bold text-indigo-400/70 max-w-xs mx-auto">
								Farewell, {domain.domain_name} 🌙
							</p>
							<p className="text-xs text-muted-foreground max-w-sm mx-auto leading-relaxed pt-1">
								Your protection has been gracefully removed. All configurations, jails, and logs have been cleaned up. The shield has been lowered — but it can always be raised again.
							</p>
						</div>

						{/* Decorative divider */}
						<div className="flex items-center justify-center gap-3 pt-2">
							<div className="h-px w-16 bg-gradient-to-r from-transparent to-indigo-500/30" />
							<span className="text-lg">✨</span>
							<div className="h-px w-16 bg-gradient-to-l from-transparent to-indigo-500/30" />
						</div>

						<div className="pt-4 flex justify-center">
							<Button
								variant="primary"
								size="md"
								onClick={() => {
									const event = new CustomEvent("navigate-page", { detail: "dashboard" });
									window.dispatchEvent(event);
									onComplete();
								}}
								className="h-12 px-10 text-sm font-semibold tracking-wide flex items-center gap-2 bg-indigo-600 hover:bg-indigo-500 border-indigo-700 hover:border-indigo-600 shadow-lg shadow-indigo-500/20 hover:shadow-indigo-500/30 hover:-translate-y-0.5 active:translate-y-0 transition-all duration-200"
							>
								<ArrowLeft className="h-4 w-4" /> Back to Dashboard
							</Button>
						</div>
					</div>
				)}
			</Card>
		</div>
	);
}

/* ── Farewell Particles ── */
interface FarewellParticle {
	id: number;
	x: number;
	size: number;
	delay: number;
	duration: number;
	opacity: number;
	drift: number;
}

function FarewellParticles() {
	const [particles, setParticles] = useState<FarewellParticle[]>([]);

	useEffect(() => {
		const temp: FarewellParticle[] = [];
		for (let i = 0; i < 60; i++) {
			temp.push({
				id: i,
				x: Math.random() * 100,
				size: Math.random() * 6 + 3,
				delay: Math.random() * 3,
				duration: Math.random() * 5 + 4,
				opacity: Math.random() * 0.5 + 0.2,
				drift: Math.random() * 40 - 20
			});
		}
		setParticles(temp);
	}, []);

	return (
		<div className="pointer-events-none fixed inset-0 z-50 overflow-hidden">
			<style>{`
				@keyframes farewell-rise {
					0% {
						transform: translateY(0) translateX(0) scale(1);
						opacity: var(--fw-opacity);
					}
					50% {
						transform: translateY(-40vh) translateX(var(--fw-drift)) scale(0.6);
						opacity: calc(var(--fw-opacity) * 0.6);
					}
					100% {
						transform: translateY(-85vh) translateX(calc(var(--fw-drift) * -0.5)) scale(0.1);
						opacity: 0;
					}
				}
				@keyframes farewell-wave {
					0%, 100% { transform: rotate(0deg); }
					15% { transform: rotate(14deg); }
					30% { transform: rotate(-8deg); }
					45% { transform: rotate(14deg); }
					60% { transform: rotate(-4deg); }
					75% { transform: rotate(10deg); }
				}
				@keyframes farewell-pulse {
					0%, 100% { opacity: 0.08; transform: scale(1); }
					50% { opacity: 0.2; transform: scale(1.15); }
				}
				@keyframes farewell-fade-letters {
					0%, 100% { opacity: 0.6; }
					50% { opacity: 1; }
				}
			`}</style>
			{particles.map((p) => {
				const colors = ["#6366f1", "#818cf8", "#a5b4fc", "#94a3b8", "#c7d2fe", "#e2e8f0"];
				const color = colors[p.id % colors.length];
				return (
					<div
						key={p.id}
						className="absolute rounded-full"
						style={{
							left: `${p.x}%`,
							bottom: `${Math.random() * 15}%`,
							width: `${p.size}px`,
							height: `${p.size}px`,
							backgroundColor: color,
							boxShadow: `0 0 ${p.size * 2}px ${color}50`,
							"--fw-opacity": p.opacity,
							"--fw-drift": `${p.drift}px`,
							animation: `farewell-rise ${p.duration}s ease-out ${p.delay}s forwards`
						} as React.CSSProperties}
					/>
				);
			})}
		</div>
	);
}
