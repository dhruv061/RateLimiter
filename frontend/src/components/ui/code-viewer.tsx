import { useState } from "react";
import { Check, Copy, Download } from "lucide-react";
import { Button } from "./button";

export interface CodeFile {
	name: string;
	content: string;
	language?: string;
}

interface CodeViewerProps {
	files: CodeFile[];
	showDownload?: boolean;
}

export function CodeViewer({ files, showDownload = true }: CodeViewerProps) {
	const [activeTabIndex, setActiveTabIndex] = useState(0);
	const [copied, setCopied] = useState(false);

	if (files.length === 0) return null;
	const activeFile = files[activeTabIndex];

	const handleCopy = async () => {
		try {
			await navigator.clipboard.writeText(activeFile.content);
			setCopied(true);
			setTimeout(() => setCopied(false), 2000);
		} catch (err) {
			console.error("Failed to copy text: ", err);
		}
	};

	const handleDownload = () => {
		const blob = new Blob([activeFile.content], { type: "text/plain;charset=utf-8" });
		const url = URL.createObjectURL(blob);
		const link = document.createElement("a");
		link.href = url;
		link.download = activeFile.name;
		document.body.appendChild(link);
		link.click();
		document.body.removeChild(link);
		URL.revokeObjectURL(url);
	};

	const lines = activeFile.content.split("\n");
	// Remove trailing newline line count if empty
	if (lines.length > 1 && lines[lines.length - 1] === "") {
		lines.pop();
	}

	return (
		<div className="flex flex-col rounded-lg border border-border bg-[#0b0f19] text-[#e2e8f0] shadow-2xl overflow-hidden font-mono">
			{/* Tabs & Toolbar */}
			<div className="flex flex-wrap items-center justify-between border-b border-border bg-[#101726] px-4 py-2 gap-2 select-none">
				<div className="flex flex-wrap gap-1">
					{files.map((file, idx) => (
						<button
							key={idx}
							onClick={() => {
								setActiveTabIndex(idx);
								setCopied(false);
							}}
							className={`rounded-md px-3 py-1 text-xs font-semibold tracking-wide transition-all ${
								activeTabIndex === idx
									? "bg-[#1e293b] text-white border border-[#334155] shadow-inner"
									: "text-[#64748b] hover:text-[#94a3b8] hover:bg-[#1e293b]/50"
							}`}
						>
							{file.name}
						</button>
					))}
				</div>
				<div className="flex items-center gap-2">
					<Button
						size="sm"
						variant="ghost"
						onClick={handleCopy}
						className="h-8 text-[#94a3b8] hover:text-white hover:bg-[#1e293b] px-2 text-xs flex items-center gap-1.5"
					>
						{copied ? <Check className="h-3.5 w-3.5 text-emerald-500" /> : <Copy className="h-3.5 w-3.5" />}
						{copied ? "Copied" : "Copy"}
					</Button>
					{showDownload && (
						<Button
							size="sm"
							variant="ghost"
							onClick={handleDownload}
							className="h-8 text-[#94a3b8] hover:text-white hover:bg-[#1e293b] px-2 text-xs flex items-center gap-1.5"
						>
							<Download className="h-3.5 w-3.5" />
							Download
						</Button>
					)}
				</div>
			</div>

			{/* Code Content */}
			<div className="relative font-mono text-[11px] sm:text-xs overflow-auto max-h-[400px] p-4 flex leading-relaxed select-text bg-[#070b13]">
				<div className="text-[#334155] select-none text-right pr-4 border-r border-[#1e293b] flex-shrink-0 font-medium">
					{lines.map((_, i) => (
						<div key={i}>{i + 1}</div>
					))}
				</div>
				<pre className="pl-4 text-[#94a3b8] flex-1 whitespace-pre overflow-x-auto scrollbar-thin">
					<code className="text-[#cbd5e1]">{activeFile.content}</code>
				</pre>
			</div>
		</div>
	);
}
