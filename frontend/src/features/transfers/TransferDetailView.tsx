"use client";

import Link from "next/link";
import { useState } from "react";

import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Icon } from "@/components/ui/Icon";
import { Spinner } from "@/components/ui/Spinner";
import { useGetTransferQuery } from "@/features/transfers/transfersApi";
import { cn } from "@/lib/cn";

export function TransferDetailView({ id }: { id: string }) {
  const { data, isLoading, isError } = useGetTransferQuery(id);
  const [copied, setCopied] = useState(false);

  if (isLoading) {
    return (
      <div className="flex justify-center py-20">
        <Spinner className="text-3xl" />
      </div>
    );
  }

  if (isError || !data) {
    return (
      <Card className="text-center">
        <h1 className="text-xl font-bold">Transfer not found</h1>
        <Link href="/transfers" className="mt-4 inline-block">
          <Button variant="secondary">Back to transfers</Button>
        </Link>
      </Card>
    );
  }

  const expired = data.status === "expired";

  async function copyLink() {
    await navigator.clipboard.writeText(data!.shareUrl);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  }

  return (
    <div className="mx-auto max-w-3xl space-y-6">
      <Link
        href="/transfers"
        className="inline-flex min-h-[var(--touch-min)] items-center gap-1 text-sm font-semibold text-primary"
      >
        <Icon name="arrow_back" />
        Back to Transfers
      </Link>

      <Card className="space-y-5">
        <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
          <div className="min-w-0 space-y-2">
            <div className="flex flex-wrap items-center gap-2">
              <h1 className="truncate text-2xl font-bold">{data.title}</h1>
              <span
                className={cn(
                  "rounded-full px-2.5 py-0.5 text-xs font-semibold",
                  expired
                    ? "bg-error-container text-on-error-container"
                    : "bg-tertiary-fixed/40 text-tertiary",
                )}
              >
                {expired ? "Expired" : "Active"}
              </span>
            </div>
            <p className="text-sm text-on-surface-variant">
              {data.totalSizeLabel} • {data.uploadedAtLabel}
            </p>
            <p className="text-sm text-on-surface-variant">
              {data.expiresAtLabel} • {data.downloadCount} downloads
            </p>
          </div>
        </div>

        <div>
          <p className="mb-2 text-sm font-semibold">Share link</p>
          <div className="flex flex-col gap-2 sm:flex-row">
            <div className="flex min-h-11 min-w-0 flex-1 items-center rounded-xl bg-surface-container-low px-3">
              <input
                readOnly
                value={data.shareUrl}
                className="w-full truncate bg-transparent text-sm outline-none"
              />
            </div>
            <Button onClick={copyLink} disabled={expired}>
              <Icon name={copied ? "check" : "content_copy"} />
              Copy Link
            </Button>
          </div>
        </div>

        {data.message ? (
          <blockquote className="rounded-xl bg-surface-container-low p-4 text-sm text-on-surface-variant">
            <Icon name="format_quote" className="mb-1 text-primary" />
            <p>{data.message}</p>
          </blockquote>
        ) : null}

        <div className="space-y-2">
          <h2 className="text-sm font-semibold">Files</h2>
          {data.files.map((f) => (
            <div
              key={f.id}
              className="flex flex-col gap-2 rounded-xl bg-surface-container-low p-3 sm:flex-row sm:items-center sm:justify-between"
            >
              <div className="flex min-w-0 items-center gap-2">
                <Icon name="picture_as_pdf" />
                <div className="min-w-0">
                  <p className="truncate font-medium">{f.name}</p>
                  <p className="text-xs text-on-surface-variant">{f.sizeLabel}</p>
                </div>
              </div>
              <Button variant="secondary" disabled={expired}>
                <Icon name="download" />
                Download
              </Button>
            </div>
          ))}
        </div>

        <Button variant="danger" className="w-full sm:w-auto">
          <Icon name="delete" />
          Delete transfer
        </Button>
      </Card>
    </div>
  );
}
