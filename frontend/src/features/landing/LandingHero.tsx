"use client";

import Link from "next/link";
import { useRouter } from "next/navigation";
import { FormEvent, useState } from "react";

import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Icon } from "@/components/ui/Icon";
import { TextField } from "@/components/ui/TextField";

export function LandingHero() {
  const router = useRouter();
  const [code, setCode] = useState("");

  function onGo(e: FormEvent) {
    e.preventDefault();
    const raw = code.trim();
    if (!raw) return;
    const slug = raw.includes("/")
      ? raw.split("/").filter(Boolean).pop()!
      : raw.replace(/^flash-/, "");
    router.push(`/t/${encodeURIComponent(slug)}`);
  }

  return (
    <div className="mx-auto flex w-full max-w-5xl flex-1 items-center px-4 py-10 sm:px-6 sm:py-16 lg:py-24">
      <div className="grid w-full grid-cols-1 items-center gap-10 lg:grid-cols-12 lg:gap-12">
        <div className="flex flex-col gap-5 text-left lg:col-span-6">
          <Badge className="self-start normal-case tracking-normal" dot>
            Single-file transfer
          </Badge>
          <h1 className="text-4xl font-extrabold leading-tight tracking-tight text-on-surface sm:text-5xl">
            Send big files{" "}
            <span className="text-primary">in a flash.</span>
          </h1>
          <p className="max-w-md text-base text-on-surface-variant sm:text-lg">
            A simple, temporary single-file sharing tool. Free for recipients
            with zero registration required.
          </p>
          <ul className="flex flex-col gap-2 sm:flex-row sm:flex-wrap sm:gap-x-5 sm:gap-y-2">
            {[
              "Single file transfers",
              "Auto-expiring links",
              "No account to download",
            ].map((item) => (
              <li
                key={item}
                className="inline-flex items-center gap-2 text-sm font-medium text-on-surface"
              >
                <Icon
                  name="check_circle"
                  className="text-tertiary"
                  filled
                />
                {item}
              </li>
            ))}
          </ul>
        </div>

        <div className="lg:col-span-6">
          <Card className="space-y-6 border border-outline-variant/30 shadow-lg">
            <div className="space-y-3">
              <div className="flex flex-wrap items-center gap-2">
                <h2 className="text-lg font-bold">Send a file</h2>
                <span className="rounded-full bg-surface-container px-2.5 py-1 text-xs font-semibold text-on-surface-variant">
                  Account required
                </span>
              </div>
              <p className="text-sm text-on-surface-variant">
                Sign in to upload and configure your temporary share link.
              </p>
              <div className="grid grid-cols-1 gap-2 sm:grid-cols-2">
                <Link href="/login" className="block">
                  <Button fullWidth>
                    <Icon name="login" />
                    Sign In
                  </Button>
                </Link>
                <Link href="/register" className="block">
                  <Button variant="secondary" fullWidth>
                    <Icon name="person_add" />
                    Register
                  </Button>
                </Link>
              </div>
            </div>

            <div className="relative flex items-center gap-3">
              <div className="h-px flex-1 bg-outline-variant/50" />
              <span className="text-xs font-semibold uppercase tracking-wider text-outline">
                or
              </span>
              <div className="h-px flex-1 bg-outline-variant/50" />
            </div>

            <form className="space-y-3" onSubmit={onGo}>
              <h2 className="text-lg font-bold">Have a download link?</h2>
              <p className="text-sm text-on-surface-variant">
                Enter the code or paste the link to download instantly.
              </p>
              <div className="flex flex-col gap-2 sm:flex-row sm:items-end">
                <div className="min-w-0 flex-1">
                  <TextField
                    name="code"
                    placeholder="e.g. flash-924x"
                    value={code}
                    onChange={(e) => setCode(e.target.value)}
                    leadingIcon={<Icon name="link" />}
                    aria-label="Download code or link"
                  />
                </div>
                <Button type="submit" className="sm:min-w-[5.5rem]">
                  Go
                  <Icon name="arrow_forward" />
                </Button>
              </div>
              <p className="text-xs text-outline">
                Recipients never need to sign in or pay.
              </p>
            </form>
          </Card>
        </div>
      </div>
    </div>
  );
}
