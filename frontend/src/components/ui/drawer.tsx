import { X } from "lucide-react";
import { motion } from "framer-motion";
import type { ReactNode } from "react";
import { Button } from "./button";

export function Drawer({ open, title, children, onClose }: { open: boolean; title: string; children: ReactNode; onClose: () => void }) {
  if (!open) {
    return null;
  }

  return (
    <div className="fixed inset-0 z-50 bg-background/70 backdrop-blur-sm" role="dialog" aria-modal="true">
      <motion.aside
        initial={{ x: 420, opacity: 0 }}
        animate={{ x: 0, opacity: 1 }}
        exit={{ x: 420, opacity: 0 }}
        transition={{ duration: 0.2 }}
        className="ml-auto flex h-full w-full max-w-[420px] flex-col border-l bg-card shadow-soft"
      >
        <div className="flex h-16 items-center justify-between border-b px-5">
          <h2 className="text-lg font-semibold">{title}</h2>
          <Button size="icon" variant="ghost" onClick={onClose} aria-label="Close drawer">
            <X className="h-4 w-4" />
          </Button>
        </div>
        <div className="flex-1 overflow-y-auto p-5">{children}</div>
      </motion.aside>
    </div>
  );
}
