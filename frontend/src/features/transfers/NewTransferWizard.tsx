"use client";

import Link from "next/link";
import { FormEvent, useState } from "react";

import { Alert } from "@/components/ui/Alert";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Icon } from "@/components/ui/Icon";
import { TextField } from "@/components/ui/TextField";
import {
  useCreateTransferMutation,
  useUploadTransferFileMutation,
} from "@/features/transfers/transfersApi";
import { getApiErrorMessage } from "@/lib/apiError";
import { formatBytes } from "@/lib/format";

export function NewTransferWizard() {
  const [file, setFile] = useState<File | null>(null);
  const [title, setTitle] = useState("");
  const [message, setMessage] = useState("");
  const [password, setPassword] = useState("");
  const [recipients, setRecipients] = useState("");
  const [expiry, setExpiry] = useState("7");
  const [shareUrl, setShareUrl] = useState<string | null>(null);
  const [transferId, setTransferId] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);

  const [createTransfer, createState] = useCreateTransferMutation();
  const [uploadFile, uploadState] = useUploadTransferFileMutation();
  const busy = createState.isLoading || uploadState.isLoading;

  async function onSubmit(e: FormEvent) {
    e.preventDefault();
    setError(null);
    if (!file) {
      setError("Choose a file to upload.");
      return;
    }

    const days = Number(expiry);
    const recipientList = recipients
      .split(/[,;\s]+/)
      .map((v) => v.trim())
      .filter(Boolean);

    try {
      const transfer = await createTransfer({
        title: title.trim() || undefined,
        message: message.trim() || undefined,
        password: password.trim() || undefined,
        expires_in_hours: Math.max(1, days) * 24,
        recipients: recipientList.length ? recipientList : undefined,
      }).unwrap();

      await uploadFile({ transferId: transfer.id, file }).unwrap();
      setShareUrl(transfer.share_url);
      setTransferId(transfer.id);
    } catch (err) {
      setError(getApiErrorMessage(err, "Could not create transfer"));
    }
  }

  async function copyLink() {
    if (!shareUrl) return;
    try {
      await navigator.clipboard.writeText(shareUrl);
      setCopied(true);
      setTimeout(() => setCopied(false), 1500);
    } catch {
      /* ignore */
    }
  }

  if (shareUrl) {
    return (
      <div className="mx-auto max-w-3xl space-y-6">
        <div>
          <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">
            Transfer ready
          </h1>
          <p className="mt-1 text-sm text-on-surface-variant">
            File uploaded. Copy the share link or open transfer details.
          </p>
        </div>
        <Alert variant="success">
          <p className="font-semibold">Share link created</p>
          <div className="mt-3 flex flex-col gap-2 sm:flex-row">
            <div className="flex min-h-11 min-w-0 flex-1 items-center rounded-xl bg-surface-container-low px-3">
              <input
                readOnly
                value={shareUrl}
                className="w-full truncate bg-transparent text-sm outline-none"
              />
            </div>
            <Button type="button" onClick={copyLink}>
              <Icon name={copied ? "check" : "content_copy"} />
              Copy link
            </Button>
          </div>
          <div className="mt-4 flex flex-wrap gap-3">
            {transferId ? (
              <Link
                href={`/transfers/${transferId}`}
                className="inline-flex font-semibold text-primary underline-offset-2 hover:underline"
              >
                Open transfer details
              </Link>
            ) : null}
            <Link
              href="/transfers"
              className="inline-flex font-semibold text-primary underline-offset-2 hover:underline"
            >
              Back to dashboard
            </Link>
          </div>
        </Alert>
      </div>
    );
  }

  return (
    <div className="mx-auto max-w-3xl space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">
          New Transfer
        </h1>
        <p className="mt-1 text-sm text-on-surface-variant">
          Send a file securely with a simple shareable link or email
        </p>
      </div>

      {error ? <Alert variant="error">{error}</Alert> : null}

      <form className="space-y-6" onSubmit={onSubmit}>
        <Card className="space-y-4">
          <h2 className="text-lg font-bold">File</h2>
          {!file ? (
            <label className="flex min-h-40 cursor-pointer flex-col items-center justify-center gap-2 rounded-xl border border-dashed border-outline-variant bg-surface-container-low px-4 text-center">
              <Icon name="upload_file" className="text-3xl text-primary" />
              <span className="font-semibold">Drag and drop your file here</span>
              <span className="text-sm text-on-surface-variant">
                or browse from your device (single file, up to 512 MB)
              </span>
              <input
                type="file"
                className="sr-only"
                onChange={(e) => {
                  const f = e.target.files?.[0];
                  if (!f) return;
                  setFile(f);
                }}
              />
            </label>
          ) : (
            <div className="flex flex-col gap-3 rounded-xl bg-surface-container-low p-4 sm:flex-row sm:items-center sm:justify-between">
              <div className="flex min-w-0 items-center gap-3">
                <Icon name="description" className="text-primary" />
                <div className="min-w-0">
                  <p className="truncate font-semibold">{file.name}</p>
                  <p className="text-sm text-on-surface-variant">
                    {formatBytes(file.size)} • Ready to upload
                  </p>
                </div>
              </div>
              <Button
                type="button"
                variant="ghost"
                onClick={() => setFile(null)}
                disabled={busy}
              >
                <Icon name="close" />
                Remove
              </Button>
            </div>
          )}
        </Card>

        <Card className="space-y-4">
          <h2 className="text-lg font-bold">Transfer details</h2>
          <TextField
            label="Transfer Title"
            name="title"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            placeholder="Optional title"
          />
          <label className="flex flex-col gap-1.5 text-sm">
            <span className="font-semibold">Message</span>
            <textarea
              name="message"
              rows={3}
              value={message}
              onChange={(e) => setMessage(e.target.value)}
              placeholder="Optional note for recipients"
              className="rounded-lg bg-surface-container-low px-4 py-3 text-sm shadow-sm outline-none focus:ring-2 focus:ring-primary/30"
            />
          </label>
          <label className="flex flex-col gap-1.5 text-sm">
            <span className="font-semibold">Expires in</span>
            <select
              value={expiry}
              onChange={(e) => setExpiry(e.target.value)}
              className="min-h-[var(--touch-min)] rounded-lg bg-surface-container-low px-4 text-sm shadow-sm outline-none focus:ring-2 focus:ring-primary/30"
            >
              <option value="1">1 day</option>
              <option value="3">3 days</option>
              <option value="7">7 days</option>
              <option value="14">14 days</option>
            </select>
          </label>
          <TextField
            label="Optional password"
            name="password"
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            leadingIcon={<Icon name="lock" />}
            hint="Min. 4 characters if set"
          />
          <TextField
            label="Recipient emails"
            name="recipients"
            placeholder="a@example.com, b@example.com"
            value={recipients}
            onChange={(e) => setRecipients(e.target.value)}
            leadingIcon={<Icon name="mail" />}
            hint="Optional — mail sending comes later"
          />
        </Card>

        <div className="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end">
          <Link href="/transfers" className="block sm:inline-block">
            <Button
              type="button"
              variant="secondary"
              fullWidth
              className="sm:w-auto"
              disabled={busy}
            >
              Cancel
            </Button>
          </Link>
          <Button
            type="submit"
            fullWidth
            className="sm:w-auto"
            disabled={!file || busy}
          >
            {busy ? "Uploading…" : "Create transfer"}
            <Icon name="arrow_forward" />
          </Button>
        </div>
      </form>
    </div>
  );
}
