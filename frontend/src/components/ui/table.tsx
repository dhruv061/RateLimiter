import type { HTMLAttributes, TdHTMLAttributes, ThHTMLAttributes } from "react";
import { cn } from "../../utils/cn";

export function Table({ className, ...props }: HTMLAttributes<HTMLTableElement>) {
  return <table className={cn("w-full border-collapse text-sm", className)} {...props} />;
}

export function Th({ className, ...props }: ThHTMLAttributes<HTMLTableCellElement>) {
  return (
    <th
      className={cn(
        "sticky top-0 z-10 h-11 border-b bg-card px-4 text-left text-xs font-medium uppercase text-muted-foreground",
        className
      )}
      {...props}
    />
  );
}

export function Td({ className, ...props }: TdHTMLAttributes<HTMLTableCellElement>) {
  return <td className={cn("h-14 border-b px-4 align-middle", className)} {...props} />;
}

export function EmptyRow({ children, colSpan }: { children: string; colSpan: number }) {
  return (
    <tr>
      <Td colSpan={colSpan} className="h-28 text-center text-muted-foreground">
        {children}
      </Td>
    </tr>
  );
}
