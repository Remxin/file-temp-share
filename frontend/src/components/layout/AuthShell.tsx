import { BrandLogo } from "@/components/ui/BrandLogo";
import { Card } from "@/components/ui/Card";
import { cn } from "@/lib/cn";

type AuthShellProps = {
  children: React.ReactNode;
  className?: string;
  cardClassName?: string;
  showBrandAboveCard?: boolean;
};

export function AuthShell({
  children,
  className,
  cardClassName,
  showBrandAboveCard = false,
}: AuthShellProps) {
  return (
    <div
      className={cn(
        "relative flex min-h-dvh w-full flex-col items-center justify-center overflow-x-hidden px-4 py-8 sm:px-6 sm:py-12",
        className,
      )}
    >
      <div
        aria-hidden
        className="pointer-events-none absolute -left-24 -top-24 size-72 rounded-full bg-primary-fixed/25 blur-3xl sm:size-96"
      />
      <div
        aria-hidden
        className="pointer-events-none absolute -right-20 top-1/3 size-64 rounded-full bg-secondary-container/25 blur-3xl sm:size-80"
      />

      <div className="relative z-10 flex w-full max-w-md flex-col items-center gap-6">
        {showBrandAboveCard ? <BrandLogo size="md" /> : null}
        <Card className={cn("space-y-6", cardClassName)}>{children}</Card>
        <p className="text-center text-xs text-outline">
          © {new Date().getFullYear()} FlashShare • Simple student project
        </p>
      </div>
    </div>
  );
}
