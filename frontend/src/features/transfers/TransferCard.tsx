"use client";

import Link from "next/link";
import { useState } from "react";

import { Button } from "@/components/ui/Button";
import { Icon } from "@/components/ui/Icon";
import { cn } from "@/lib/cn";
import {
  displayTitle,
  formatBytes,
  formatExpiresAt,
} from "@/lib/format";
import type { TransferListItem } from "@/lib/types/transfer";

type TransferCardProps = {
  transfer: TransferListItem;
};

export function TransferCard({ transfer }: TransferCardProps) {
  const [copied, setCopied] = useState(false);
  const expired = transfer.status === "expired";
  const fileName = transfer.first_file_name?.trim() || "No files yet";
  const shareMode = transfer.recipient_count > 0 ? "recipients" : "public";

  async function copyLink() {
    try {
      await navigator.clipboard.writeText(transfer.share_url);
      setCopied(true);
      setTimeout(() => setCopied(false), 1500);
    } catch {
      /* ignore */
    }
  }

  return (
    <article
      className={cn(
        "rounded-xl border border-outline-variant/30 bg-surface-container-lowest p-4 shadow-sm sm:p-5",
        expired && "opacity-80",
      )}
    >
      <div className="flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between">
        <div className="flex min-w-0 flex-1 gap-3">
          <span className="inline-flex size-11 shrink-0 items-center justify-center rounded-xl bg-surface-container text-primary">
            <Icon
              name={
                fileName.endsWith(".mp4")
                  ? "video_file"
                  : fileName.endsWith(".pdf")
                    ? "picture_as_pdf"
                    : "description"
              }
            />
          </span>
          <div className="min-w-0 space-y-1">
            <div className="flex flex-wrap items-center gap-2">
              <h3 className="truncate text-base font-bold text-on-surface">
                {displayTitle(transfer.title)}
              </h3>
              <span
                className={cn(
                  "rounded-full px-2 py-0.5 text-xs font-semibold",
                  expired
                    ? "bg-error-container text-on-error-container"
                    : "bg-tertiary-fixed/40 text-tertiary",
                )}
              >
                {expired ? "Expired" : "Active"}
              </span>
            </div>
            <p className="text-sm text-on-surface-variant">
              {transfer.file_count}{" "}
              {transfer.file_count === 1 ? "file" : "files"} •{" "}
              {formatBytes(transfer.total_size_bytes)}
            </p>
            <p className="inline-flex items-center gap-1 text-sm text-on-surface-variant">
              <Icon name="schedule" className="text-[1rem]" />
              {formatExpiresAt(transfer.expires_at)}
            </p>
            <p className="inline-flex items-center gap-1 text-sm text-on-surface-variant">
              <Icon name="download" className="text-[1rem]" />
              {transfer.download_count} downloads
            </p>
            <p className="inline-flex items-center gap-1 text-xs font-medium uppercase tracking-wide text-outline">
              <Icon
                name={shareMode === "public" ? "public" : "mail"}
                className="text-[1rem]"
              />
              {shareMode === "public" ? "Public Link" : "Recipients"}
            </p>
          </div>
        </div>

        <div className="flex w-full flex-col gap-2 sm:w-auto sm:min-w-[11rem]">
          <div className="flex min-h-11 items-center gap-1 rounded-xl bg-surface-container-low px-3">
            <input
              readOnly
              value={transfer.share_url}
              className="min-w-0 flex-1 truncate bg-transparent text-xs text-on-surface-variant outline-none"
            />
            <button
              type="button"
              className="inline-flex size-9 shrink-0 items-center justify-center rounded-md text-primary"
              onClick={copyLink}
              aria-label="Copy link"
            >
              <Icon name={copied ? "check" : "content_copy"} />
            </button>
          </div>
          <div className="grid grid-cols-2 gap-2">
            <Link href={`/transfers/${transfer.id}`} className="block">
              <Button variant="secondary" fullWidth className="text-xs sm:text-sm">
                <Icon name="info" />
                Details
              </Button>
            </Link>
            <Button
              variant="ghost"
              fullWidth
              className="text-error text-xs sm:text-sm"
              disabled
              title="Delete comes next"
            >
              <Icon name="delete" />
              Delete
            </Button>
          </div>
        </div>
      </div>
    </article>
  );
}
