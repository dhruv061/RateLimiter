import { ShieldCheck, ArrowRight, Server } from "lucide-react";
import { Button } from "../../components/ui/button";

interface WelcomeScreenProps {
	onConfigure: () => void;
}

export function WelcomeScreen({ onConfigure }: WelcomeScreenProps) {
	return (
		<div className="relative min-h-screen flex items-center justify-center bg-radial-glow overflow-hidden py-12 px-4 select-none">
			{/* Decorative glowing gradient bubbles */}
			<div className="absolute top-1/4 left-1/4 -translate-x-1/2 -translate-y-1/2 w-80 h-80 rounded-full bg-primary/10 blur-3xl animate-pulse-glow" style={{ animationDuration: "8s" }} />
			<div className="absolute bottom-1/4 right-1/4 translate-x-1/2 translate-y-1/2 w-96 h-96 rounded-full bg-emerald-500/5 blur-3xl animate-pulse-glow" style={{ animationDuration: "12s" }} />

			<div className="relative w-full max-w-xl text-center space-y-8 animate-fadeInUp">
				{/* Branding/Icon */}
				<div className="flex flex-col items-center space-y-4">
					<div className="relative flex h-24 w-24 items-center justify-center rounded-2xl bg-gradient-to-tr from-primary to-blue-600 text-primary-foreground shadow-2xl shadow-primary/20 animate-shield-float">
						<ShieldCheck className="h-12 w-12 text-white" />
						<div className="absolute -inset-1 rounded-2xl bg-gradient-to-tr from-primary to-blue-500 opacity-30 blur-sm -z-10" />
					</div>
					<div className="flex items-center gap-1.5 text-sm font-semibold tracking-wider text-primary uppercase mt-4">
						<Server className="h-4 w-4" />
						ShieldWatch Protection
					</div>
				</div>

				{/* Title and subtitle */}
				<div className="space-y-3">
					<h1 className="text-4xl sm:text-5xl font-bold tracking-tight bg-clip-text text-transparent bg-gradient-to-b from-foreground to-muted-foreground/80 leading-tight">
						Protect Your Servers
					</h1>
					<p className="text-base sm:text-lg text-muted-foreground max-w-md mx-auto leading-relaxed">
						Fail2Ban Dashboard has detected zero protected domains configured. Let's connect your first domain.
					</p>
				</div>

				{/* CTA Button */}
				<div className="pt-4 flex justify-center">
					<Button
						variant="primary"
						size="md"
						onClick={onConfigure}
						className="relative h-13 px-8 text-sm font-semibold tracking-wide shadow-lg shadow-primary/25 hover:shadow-primary/35 hover:-translate-y-0.5 active:translate-y-0 transition-all duration-200 group overflow-hidden"
					>
						<span className="relative z-10 flex items-center gap-2">
							Configure Domain
							<ArrowRight className="h-4 w-4 group-hover:translate-x-1 transition-transform" />
						</span>
						<div className="absolute inset-0 bg-gradient-to-r from-blue-600 to-primary opacity-0 hover:opacity-100 transition-opacity duration-300 -z-10" />
					</Button>
				</div>

				{/* Details card grid */}
				<div className="grid grid-cols-1 sm:grid-cols-2 gap-4 max-w-lg mx-auto pt-6 text-left">
					<div className="rounded-xl border border-border bg-card/40 p-4 backdrop-blur-md hover:bg-card/60 transition-colors">
						<h3 className="font-semibold text-sm text-foreground mb-1">Nginx Integration</h3>
						<p className="text-xs text-muted-foreground leading-normal">
							Scans Nginx configurations, discovers server block directives, and sets up high performance geo-blocking behind proxies.
						</p>
					</div>
					<div className="rounded-xl border border-border bg-card/40 p-4 backdrop-blur-md hover:bg-card/60 transition-colors">
						<h3 className="font-semibold text-sm text-foreground mb-1">Fail2Ban Auto-tuning</h3>
						<p className="text-xs text-muted-foreground leading-normal">
							Configures custom filters, actions, and local jail definitions to auto-ban offensive IP addresses for 24 hours.
						</p>
					</div>
				</div>
			</div>
		</div>
	);
}
