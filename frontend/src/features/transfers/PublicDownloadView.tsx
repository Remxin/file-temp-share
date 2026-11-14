"use client";

import Link from "next/link";
import { FormEvent, useState } from "react";

import { Alert } from "@/components/ui/Alert";
import { BrandLogo } from "@/components/ui/BrandLogo";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Icon } from "@/components/ui/Icon";
import { Spinner } from "@/components/ui/Spinner";
import { TextField } from "@/components/ui/TextField";
import {
  useGetPublicTransferQuery,
  useUnlockPublicTransferMutation,
} from "@/features/transfers/transfersApi";
import { getApiErrorMessage } from "@/lib/apiError";
import { downloadFile } from "@/lib/downloadFile";
import {
  displayTitle,
  formatBytes,
  formatExpiresAt,
} from "@/lib/format";
import type { TransferFile } from "@/lib/types/transfer";

export function PublicDownloadView({ slug }: { slug: string }) {
  const { data, isLoading, isError } = useGetPublicTransferQuery(slug);
  const [unlock, unlockState] = useUnlockPublicTransferMutation();
  const [password, setPassword] = useState("");
  const [unlockedPassword, setUnlockedPassword] = useState<string | null>(null);
  const [formError, setFormError] = useState<string | null>(null);
  const [downloadError, setDownloadError] = useState<string | null>(null);
  const [downloadingId, setDownloadingId] = useState<string | null>(null);

  async function onUnlock(e: FormEvent) {
    e.preventDefault();
    setFormError(null);
    try {
      await unlock({ shortCode: slug, password }).unwrap();
      setUnlockedPassword(password);
    } catch (err) {
      setFormError(getApiErrorMessage(err, "Invalid password"));
    }
  }

  async function onDownload(file: TransferFile) {
    setDownloadError(null);
    setDownloadingId(file.id);
    try {
      await downloadFile({
        path: `/api/v1/public/transfers/${slug}/files/${file.id}`,
        filename: file.original_name,
        transferPassword: unlockedPassword ?? undefined,
      });
    } catch (err) {
      setDownloadError(
        err instanceof Error ? err.message : "Download failed",
      );
    } finally {
      setDownloadingId(null);
    }
  }

  async function onDownloadAll() {
    if (!data?.files.length) return;
    for (const file of data.files) {
      await onDownload(file);
    }
  }

  if (isLoading) {
    return (
      <div className="flex min-h-dvh items-center justify-center">
        <Spinner className="text-3xl" />
      </div>
    );
  }

  if (isError || !data) {
    return (
      <div className="flex min-h-dvh flex-col items-center justify-center gap-4 px-4">
        <BrandLogo />
        <Card className="max-w-md text-center">
          <h1 className="text-xl font-bold">Transfer not found</h1>
          <p className="mt-2 text-sm text-on-surface-variant">
            This link is invalid or was removed.
          </p>
          <Link href="/" className="mt-5 inline-block">
            <Button>Back to FlashShare</Button>
          </Link>
        </Card>
      </div>
    );
  }

  if (data.status === "expired") {
    return (
      <div className="flex min-h-dvh flex-col">
        <header className="border-b border-outline-variant/30 px-4 py-3 sm:px-6">
          <div className="mx-auto flex max-w-3xl items-center justify-between">
            <BrandLogo size="sm" />
            <Link href="/login" className="text-sm font-semibold text-primary">
              Sign In
            </Link>
          </div>
        </header>
        <main className="mx-auto flex w-full max-w-lg flex-1 items-center px-4 py-10">
          <Card className="w-full space-y-5 text-center">
            <span className="mx-auto inline-flex size-16 items-center justify-center rounded-full bg-error-container text-error">
              <Icon name="hourglass_disabled" className="text-3xl" />
            </span>
            <div>
              <h1 className="text-2xl font-bold">
                This transfer is no longer available
              </h1>
              <p className="mt-2 text-sm text-on-surface-variant">
                The file shared at this link has expired or was removed by the
                sender.
              </p>
            </div>
            <div className="rounded-xl bg-surface-container-low p-4 text-left">
              <p className="font-semibold">{displayTitle(data.title)}</p>
              <p className="text-sm text-on-surface-variant">
                {formatExpiresAt(data.expires_at)}
              </p>
            </div>
            <Link href="/" className="inline-block w-full">
              <Button fullWidth variant="secondary">
                Back to FlashShare
              </Button>
            </Link>
          </Card>
        </main>
      </div>
    );
  }

  const needsPassword = data.has_password && !data.unlocked;

  return (
    <div className="flex min-h-dvh flex-col">
      <header className="border-b border-outline-variant/30 px-4 py-3 sm:px-6">
        <div className="mx-auto flex max-w-3xl items-center justify-between gap-3">
          <BrandLogo size="sm" />
          <div className="flex items-center gap-2">
            <Link href="/transfers/new" className="hidden sm:block">
              <Button variant="secondary">Send a File</Button>
            </Link>
            <Link href="/login">
              <Button>Sign In</Button>
            </Link>
          </div>
        </div>
      </header>

      <main className="mx-auto flex w-full max-w-lg flex-1 items-center px-4 py-10">
        <Card className="w-full space-y-5">
          <div className="text-center sm:text-left">
            <p className="text-sm text-on-surface-variant">
              <span className="font-semibold text-on-surface">
                {data.sender_email || "Someone"}
              </span>{" "}
              shared a file with you
            </p>
            <h1 className="mt-2 text-2xl font-bold">
              {displayTitle(data.title)}
            </h1>
            <p className="mt-2 inline-flex items-center gap-1 text-sm text-on-surface-variant">
              <Icon name="schedule" className="text-[1rem]" />
              {formatExpiresAt(data.expires_at)}
            </p>
          </div>

          {needsPassword ? (
            <form className="space-y-3" onSubmit={onUnlock}>
              <Alert variant="info">
                This transfer is password protected. Enter the password to view
                and download files.
              </Alert>
              {formError ? <Alert variant="error">{formError}</Alert> : null}
              <TextField
                label="Password"
                name="password"
                type="password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                leadingIcon={<Icon name="lock" />}
              />
              <Button
                type="submit"
                fullWidth
                disabled={unlockState.isLoading || !password.trim()}
              >
                {unlockState.isLoading ? "Unlocking…" : "Unlock transfer"}
              </Button>
            </form>
          ) : (
            <>
              {data.message ? (
                <blockquote className="rounded-xl bg-surface-container-low p-4 text-sm text-on-surface-variant">
                  {data.message}
                </blockquote>
              ) : null}

              {downloadError ? (
                <Alert variant="error">{downloadError}</Alert>
              ) : null}

              <div className="space-y-2">
                {data.files.map((f) => (
                  <div
                    key={f.id}
                    className="flex flex-col gap-2 rounded-xl bg-surface-container-low p-3 sm:flex-row sm:items-center sm:justify-between"
                  >
                    <div className="flex min-w-0 items-center gap-3">
                      <Icon name="description" className="text-primary" />
                      <div className="min-w-0">
                        <p className="truncate font-semibold">
                          {f.original_name}
                        </p>
                        <p className="text-xs text-on-surface-variant">
                          {formatBytes(f.size_bytes)}
                        </p>
                      </div>
                    </div>
                    <Button
                      variant="secondary"
                      onClick={() => onDownload(f)}
                      disabled={downloadingId === f.id}
                    >
                      <Icon name="download" />
                      {downloadingId === f.id ? "…" : "Download"}
                    </Button>
                  </div>
                ))}
              </div>

              <Button
                fullWidth
                onClick={onDownloadAll}
                disabled={!data.files.length || downloadingId !== null}
              >
                <Icon name="download" />
                Download{" "}
                {data.files.length === 1
                  ? "File"
                  : "All"}{" "}
                ({formatBytes(data.total_size_bytes)})
              </Button>
            </>
          )}

          <p className="text-center text-xs text-outline">
            <Icon
              name="verified_user"
              className="mr-1 align-middle text-[1rem]"
            />
            No account needed. This link expires automatically.
          </p>

          <Link
            href="/"
            className="inline-flex min-h-[var(--touch-min)] w-full items-center justify-center gap-1 text-sm font-semibold text-primary"
          >
            <Icon name="arrow_back" />
            Back to FlashShare
          </Link>
        </Card>
      </main>
    </div>
  );
}
