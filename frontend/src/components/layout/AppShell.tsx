"use client";

import Link from "next/link";
import { usePathname, useRouter } from "next/navigation";
import { useState } from "react";

import { BrandLogo } from "@/components/ui/BrandLogo";
import { Button } from "@/components/ui/Button";
import { Icon } from "@/components/ui/Icon";
import { cn } from "@/lib/cn";
import { useAppDispatch, useAppSelector } from "@/store/hooks";
import { clearCredentials } from "@/features/auth/authSlice";
import { selectCurrentUser } from "@/features/auth/selectors";

const nav = [
  { href: "/transfers", label: "Dashboard", icon: "dashboard" },
  { href: "/transfers/new", label: "New Transfer", icon: "cloud_upload" },
  { href: "/settings", label: "Settings", icon: "settings" },
] as const;

type AppShellProps = {
  children: React.ReactNode;
};

export function AppShell({ children }: AppShellProps) {
  const pathname = usePathname();
  const router = useRouter();
  const dispatch = useAppDispatch();
  const user = useAppSelector(selectCurrentUser);
  const [menuOpen, setMenuOpen] = useState(false);

  function logout() {
    dispatch(clearCredentials());
    router.replace("/login");
  }

  return (
    <div className="flex min-h-dvh flex-col bg-surface">
      <header className="sticky top-0 z-50 border-b border-outline-variant/20 bg-surface/85 shadow-[0_1px_8px_rgba(0,0,0,0.04)] backdrop-blur-xl">
        <div className="mx-auto flex h-14 max-w-7xl items-center justify-between gap-3 px-4 sm:h-16 sm:px-6 lg:px-8">
          <div className="flex min-w-0 items-center gap-3">
            <BrandLogo href="/transfers" size="sm" />
            <nav className="hidden items-center rounded-xl bg-surface-container-low p-1 md:flex">
              {nav.map((item) => {
                const active =
                  item.href === "/transfers"
                    ? pathname === "/transfers"
                    : pathname.startsWith(item.href);
                return (
                  <Link
                    key={item.href}
                    href={item.href}
                    className={cn(
                      "inline-flex min-h-10 items-center gap-1.5 rounded-lg px-3 text-sm font-semibold transition-all",
                      active
                        ? "bg-primary-container text-on-primary"
                        : "text-on-surface-variant hover:text-on-surface",
                    )}
                  >
                    <Icon name={item.icon} className="text-[1.05rem]" />
                    {item.label}
                  </Link>
                );
              })}
            </nav>
          </div>

          <div className="flex items-center gap-2">
            <div className="hidden items-center gap-2 sm:flex">
              <div className="text-right leading-tight">
                <p className="truncate text-sm font-semibold text-on-surface">
                  {user?.email?.split("@")[0] ?? "Account"}
                </p>
                <p className="max-w-[10rem] truncate text-xs text-on-surface-variant">
                  {user?.email}
                </p>
              </div>
              <span className="inline-flex size-10 items-center justify-center rounded-full bg-surface-container text-on-surface-variant">
                <Icon name="person" />
              </span>
            </div>
            <Button
              variant="ghost"
              className="hidden sm:inline-flex"
              onClick={logout}
              aria-label="Sign out"
            >
              <Icon name="logout" />
            </Button>
            <button
              type="button"
              className="inline-flex size-11 items-center justify-center rounded-lg md:hidden"
              aria-label={menuOpen ? "Close menu" : "Open menu"}
              onClick={() => setMenuOpen((v) => !v)}
            >
              <Icon name={menuOpen ? "close" : "menu"} />
            </button>
          </div>
        </div>

        {menuOpen ? (
          <div className="border-t border-outline-variant/30 px-4 py-3 md:hidden">
            <div className="flex flex-col gap-1">
              {nav.map((item) => (
                <Link
                  key={item.href}
                  href={item.href}
                  className="inline-flex min-h-[var(--touch-min)] items-center gap-2 rounded-lg px-3 text-sm font-semibold"
                  onClick={() => setMenuOpen(false)}
                >
                  <Icon name={item.icon} />
                  {item.label}
                </Link>
              ))}
              <button
                type="button"
                className="inline-flex min-h-[var(--touch-min)] items-center gap-2 rounded-lg px-3 text-left text-sm font-semibold text-error"
                onClick={logout}
              >
                <Icon name="logout" />
                Sign out
              </button>
            </div>
          </div>
        ) : null}
      </header>

      <main className="mx-auto w-full max-w-7xl flex-1 px-4 py-6 sm:px-6 sm:py-8 lg:px-8">
        {children}
      </main>
    </div>
  );
}
