import Link from "next/link";

import { BrandLogo } from "@/components/ui/BrandLogo";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Icon } from "@/components/ui/Icon";

export default function NotFound() {
  return (
    <div className="relative flex min-h-dvh flex-col items-center justify-center overflow-x-hidden px-4 py-10">
      <div
        aria-hidden
        className="pointer-events-none absolute left-1/2 top-10 size-72 -translate-x-1/2 rounded-full bg-primary-fixed/30 blur-3xl sm:size-96"
      />
      <BrandLogo className="mb-6" />
      <Card className="relative z-10 max-w-lg space-y-6 text-center">
        <div className="relative mx-auto flex h-28 w-full max-w-xs items-center justify-center">
          <span className="absolute left-4 text-7xl font-black text-primary-fixed-dim/80 sm:left-8">
            4
          </span>
          <span className="relative z-10 inline-flex size-20 items-center justify-center rounded-full border-4 border-primary/30 bg-surface-container">
            <Icon name="cloud_off" className="text-3xl text-error" />
          </span>
          <span className="absolute right-4 text-7xl font-black text-primary-fixed-dim/80 sm:right-8">
            4
          </span>
        </div>
        <div>
          <h1 className="text-2xl font-bold sm:text-3xl">Page Not Found</h1>
          <p className="mt-2 text-sm text-on-surface-variant">
            The page you are looking for doesn&apos;t exist or has been moved.
          </p>
        </div>
        <div className="grid grid-cols-1 gap-2 sm:grid-cols-2">
          <Link href="/" className="block">
            <Button fullWidth>
              <Icon name="arrow_back" />
              Back to Home
            </Button>
          </Link>
          <Link href="/transfers" className="block">
            <Button variant="secondary" fullWidth>
              <Icon name="dashboard" />
              Dashboard
            </Button>
          </Link>
        </div>
      </Card>
    </div>
  );
}
