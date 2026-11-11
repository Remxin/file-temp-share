"use client";

import Link from "next/link";
import { FormEvent, useState } from "react";

import { Alert } from "@/components/ui/Alert";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Icon } from "@/components/ui/Icon";
import { TextField } from "@/components/ui/TextField";

type StagedFile = {
  name: string;
  sizeLabel: string;
};

export function NewTransferWizard() {
  const [file, setFile] = useState<StagedFile | null>({
    name: "branding_guide_v2.pdf",
    sizeLabel: "24 MB",
  });
  const [title, setTitle] = useState("Brand guidelines");
  const [message, setMessage] = useState("");
  const [password, setPassword] = useState("");
  const [recipients, setRecipients] = useState("");
  const [expiry, setExpiry] = useState("7");
  const [submitted, setSubmitted] = useState(false);

  function onSubmit(e: FormEvent) {
    e.preventDefault();
    setSubmitted(true);
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

      {submitted ? (
        <Alert variant="success">
          <p className="font-semibold">Transfer ready (UI stub)</p>
          <p className="mt-1">
            Backend upload API is not wired yet. Share link will appear here
            after CRUD is implemented.
          </p>
          <Link
            href="/transfers"
            className="mt-3 inline-flex font-semibold text-primary underline-offset-2 hover:underline"
          >
            Back to dashboard
          </Link>
        </Alert>
      ) : null}

      <form className="space-y-6" onSubmit={onSubmit}>
        <Card className="space-y-4">
          <h2 className="text-lg font-bold">File</h2>
          {!file ? (
            <label className="flex min-h-40 cursor-pointer flex-col items-center justify-center gap-2 rounded-xl border border-dashed border-outline-variant bg-surface-container-low px-4 text-center">
              <Icon name="upload_file" className="text-3xl text-primary" />
              <span className="font-semibold">Drag and drop your file here</span>
              <span className="text-sm text-on-surface-variant">
                or browse from your device (single file, up to 2 GB)
              </span>
              <input
                type="file"
                className="sr-only"
                onChange={(e) => {
                  const f = e.target.files?.[0];
                  if (!f) return;
                  setFile({
                    name: f.name,
                    sizeLabel: `${Math.max(1, Math.round(f.size / (1024 * 1024)))} MB`,
                  });
                }}
              />
            </label>
          ) : (
            <div className="flex flex-col gap-3 rounded-xl bg-surface-container-low p-4 sm:flex-row sm:items-center sm:justify-between">
              <div className="flex min-w-0 items-center gap-3">
                <Icon name="picture_as_pdf" className="text-primary" />
                <div className="min-w-0">
                  <p className="truncate font-semibold">{file.name}</p>
                  <p className="text-sm text-on-surface-variant">
                    {file.sizeLabel} • Ready to beam
                  </p>
                </div>
              </div>
              <Button
                type="button"
                variant="ghost"
                onClick={() => setFile(null)}
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
            required
            value={title}
            onChange={(e) => setTitle(e.target.value)}
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
            <Button type="button" variant="secondary" fullWidth className="sm:w-auto">
              Cancel
            </Button>
          </Link>
          <Button type="submit" fullWidth className="sm:w-auto" disabled={!file}>
            Create transfer
            <Icon name="arrow_forward" />
          </Button>
        </div>
      </form>
    </div>
  );
}
