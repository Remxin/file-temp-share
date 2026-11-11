import Link from "next/link";

import { cn } from "@/lib/cn";
import { Icon } from "@/components/ui/Icon";

type BrandLogoProps = {
  /** Pass `null` to render without a link. */
  href?: string | null;
  size?: "sm" | "md" | "lg";
  showWordmark?: boolean;
  className?: string;
};

const sizes = {
  sm: { box: "size-8 rounded-lg", icon: "text-[1.1rem]", text: "text-base" },
  md: { box: "size-9 rounded-xl", icon: "text-[1.25rem]", text: "text-lg" },
  lg: { box: "size-12 rounded-xl", icon: "text-[1.5rem]", text: "text-xl" },
};

export function BrandLogo({
  href = "/",
  size = "md",
  showWordmark = true,
  className,
}: BrandLogoProps) {
  const s = sizes[size];
  const content = (
    <>
      <span
        className={cn(
          "inline-flex items-center justify-center bg-primary-container text-on-primary shadow-md shadow-primary/20",
          s.box,
        )}
      >
        <Icon name="bolt" className={s.icon} filled />
      </span>
      {showWordmark ? (
        <span className={cn("font-bold tracking-tight text-on-surface", s.text)}>
          Flash<span className="text-primary">Share</span>
        </span>
      ) : null}
    </>
  );

  if (href === null) {
    return (
      <span className={cn("inline-flex items-center gap-2", className)}>
        {content}
      </span>
    );
  }

  return (
    <Link
      href={href}
      className={cn(
        "inline-flex items-center gap-2 transition-opacity hover:opacity-90",
        className,
      )}
    >
      {content}
    </Link>
  );
}
