import { cn } from "@/lib/cn";
import { Icon } from "@/components/ui/Icon";

type AlertProps = {
  children: React.ReactNode;
  variant?: "success" | "error" | "info";
  className?: string;
  icon?: string;
};

const styles = {
  success: "bg-tertiary-fixed/30 text-tertiary",
  error: "bg-error-container text-on-error-container",
  info: "bg-surface-container text-on-surface",
};

export function Alert({
  children,
  variant = "info",
  className,
  icon,
}: AlertProps) {
  const defaultIcon =
    variant === "success"
      ? "mark_email_read"
      : variant === "error"
        ? "error"
        : "info";

  return (
    <div
      className={cn(
        "flex items-start gap-3 rounded-lg p-4 text-sm",
        styles[variant],
        className,
      )}
      role="status"
    >
      <Icon name={icon ?? defaultIcon} className="mt-0.5 shrink-0" />
      <div className="min-w-0 flex-1">{children}</div>
    </div>
  );
}
