import { useEffect, useState } from "react";

interface Particle {
	id: number;
	x: number;
	y: number;
	color: string;
	size: number;
	delay: number;
	duration: number;
	tilt: number;
}

export function Confetti() {
	const [particles, setParticles] = useState<Particle[]>([]);

	useEffect(() => {
		const colors = ["#3b82f6", "#10b981", "#f59e0b", "#ef4444", "#8b5cf6", "#ec4899", "#06b6d4"];
		const temp: Particle[] = [];
		for (let i = 0; i < 120; i++) {
			temp.push({
				id: i,
				x: Math.random() * 100, // percentage width
				y: -20 - Math.random() * 30, // start above screen
				color: colors[Math.floor(Math.random() * colors.length)],
				size: Math.random() * 10 + 6, // 6px to 16px
				delay: Math.random() * 1.5, // 0s to 1.5s
				duration: Math.random() * 3.5 + 2.5, // 2.5s to 6s
				tilt: Math.random() * 20 - 10
			});
		}
		setParticles(temp);
	}, []);

	return (
		<div className="pointer-events-none fixed inset-0 z-50 overflow-hidden">
			<style>{`
				@keyframes confetti-fall {
					0% {
						transform: translateY(0) rotate(0deg) translateX(0);
						opacity: 1;
					}
					50% {
						transform: translateY(50vh) rotate(180deg) translateX(15px);
						opacity: 0.9;
					}
					100% {
						transform: translateY(105vh) rotate(360deg) translateX(-15px);
						opacity: 0;
					}
				}
			`}</style>
			{particles.map((p) => (
				<div
					key={p.id}
					className="absolute rounded-sm"
					style={{
						left: `${p.x}%`,
						top: `${p.y}px`,
						width: `${p.size}px`,
						height: `${p.size * 0.6}px`, // slightly rectangular for flutter effect
						backgroundColor: p.color,
						opacity: 0.9,
						transform: `rotate(${p.tilt}deg)`,
						animation: `confetti-fall ${p.duration}s linear ${p.delay}s forwards`
					}}
				/>
			))}
		</div>
	);
}
