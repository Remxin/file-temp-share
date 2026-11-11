"use client";

import Link from "next/link";
import { useState } from "react";

import { BrandLogo } from "@/components/ui/BrandLogo";
import { Button } from "@/components/ui/Button";
import { Icon } from "@/components/ui/Icon";
import { cn } from "@/lib/cn";

type MarketingShellProps = {
  children: React.ReactNode;
  className?: string;
};

export function MarketingShell({ children, className }: MarketingShellProps) {
  const [open, setOpen] = useState(false);

  return (
    <div className={cn("flex min-h-dvh flex-col bg-surface", className)}>
      <header className="sticky top-0 z-50 border-b border-outline-variant/30 bg-surface/90 backdrop-blur-md">
        <div className="mx-auto flex h-14 max-w-5xl items-center justify-between px-4 sm:h-16 sm:px-6">
          <BrandLogo size="sm" />

          <nav className="hidden items-center gap-2 sm:flex">
            <Link
              href="/login"
              className="inline-flex min-h-[var(--touch-min)] items-center px-3 text-sm font-semibold text-on-surface-variant hover:text-on-surface"
            >
              Sign In
            </Link>
            <Link href="/register">
              <Button>Create Account</Button>
            </Link>
          </nav>

          <button
            type="button"
            className="inline-flex size-11 items-center justify-center rounded-lg text-on-surface sm:hidden"
            aria-expanded={open}
            aria-label={open ? "Close menu" : "Open menu"}
            onClick={() => setOpen((v) => !v)}
          >
            <Icon name={open ? "close" : "menu"} />
          </button>
        </div>

        {open ? (
          <div className="border-t border-outline-variant/30 px-4 py-3 sm:hidden">
            <div className="flex flex-col gap-2">
              <Link
                href="/login"
                className="inline-flex min-h-[var(--touch-min)] items-center rounded-lg px-3 text-sm font-semibold"
                onClick={() => setOpen(false)}
              >
                Sign In
              </Link>
              <Link href="/register" onClick={() => setOpen(false)}>
                <Button fullWidth>Create Account</Button>
              </Link>
            </div>
          </div>
        ) : null}
      </header>

      <main className="flex flex-1 flex-col">{children}</main>

      <footer className="border-t border-outline-variant/20 px-4 py-6 sm:px-6">
        <div className="mx-auto flex max-w-5xl flex-col gap-3 text-xs text-outline sm:flex-row sm:items-center sm:justify-between">
          <p>© {new Date().getFullYear()} FlashShare • Simple student project</p>
          <div className="flex flex-wrap gap-4">
            <Link href="/login" className="hover:text-on-surface">
              Sign In
            </Link>
            <Link href="/register" className="hover:text-on-surface">
              Create Account
            </Link>
          </div>
        </div>
      </footer>
    </div>
  );
}
