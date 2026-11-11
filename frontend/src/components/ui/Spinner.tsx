import { cn } from "@/lib/cn";
import { Icon } from "@/components/ui/Icon";

export function Spinner({ className }: { className?: string }) {
  return (
    <Icon
      name="progress_activity"
      className={cn("animate-spin text-primary", className)}
    />
  );
}
