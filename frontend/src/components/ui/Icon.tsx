import { cn } from "@/lib/cn";

type IconProps = {
  name: string;
  className?: string;
  filled?: boolean;
  "aria-hidden"?: boolean;
};

export function Icon({
  name,
  className,
  filled = false,
  "aria-hidden": ariaHidden = true,
}: IconProps) {
  return (
    <span
      className={cn("material-symbols-outlined text-[1.25rem] leading-none", className)}
      style={
        filled
          ? { fontVariationSettings: "'FILL' 1, 'wght' 400, 'GRAD' 0, 'opsz' 24" }
          : undefined
      }
      aria-hidden={ariaHidden}
    >
      {name}
    </span>
  );
}
