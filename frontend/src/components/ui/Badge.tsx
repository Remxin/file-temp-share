import { cn } from "@/lib/cn";

type BadgeProps = {
  children: React.ReactNode;
  className?: string;
  dot?: boolean;
  dotClassName?: string;
};

export function Badge({ children, className, dot, dotClassName }: BadgeProps) {
  return (
    <span
      className={cn(
        "inline-flex items-center gap-2 rounded-full bg-surface-container-high px-3 py-1.5 text-xs font-semibold uppercase tracking-wide text-on-surface-variant",
        className,
      )}
    >
      {dot ? (
        <span
          className={cn("inline-block size-2 rounded-full bg-secondary", dotClassName)}
        />
      ) : null}
      {children}
    </span>
  );
}
