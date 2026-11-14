"use client";

import Link from "next/link";
import { useMemo, useState } from "react";

import { Button } from "@/components/ui/Button";
import { Icon } from "@/components/ui/Icon";
import { Spinner } from "@/components/ui/Spinner";
import { TextField } from "@/components/ui/TextField";
import { TransferCard } from "@/features/transfers/TransferCard";
import {
  type TransferFilter,
  useListTransfersQuery,
} from "@/features/transfers/transfersApi";
import { cn } from "@/lib/cn";

const filters: { id: TransferFilter; label: string }[] = [
  { id: "active", label: "Active" },
  { id: "expired", label: "Expired" },
  { id: "all", label: "All" },
];

export function TransfersDashboard() {
  const [filter, setFilter] = useState<TransferFilter>("active");
  const [query, setQuery] = useState("");
  const { data, isLoading, isError } = useListTransfersQuery(filter);

  const items = useMemo(() => {
    const list = data ?? [];
    const q = query.trim().toLowerCase();
    if (!q) return list;
    return list.filter((t) => (t.title ?? "").toLowerCase().includes(q));
  }, [data, query]);

  return (
    <div className="space-y-6">
      <div className="flex flex-col gap-4 sm:flex-row sm:items-end sm:justify-between">
        <div>
          <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">
            My Transfers
          </h1>
          <p className="mt-1 text-sm text-on-surface-variant">
            {data ? `${data.length} transfers in this filter` : "Loading…"}
          </p>
        </div>
        <Link href="/transfers/new" className="block w-full sm:w-auto">
          <Button fullWidth className="sm:w-auto">
            <Icon name="add" />
            New Transfer
          </Button>
        </Link>
      </div>

      <div className="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <div className="flex gap-1 rounded-xl bg-surface-container-low p-1">
          {filters.map((f) => (
            <button
              key={f.id}
              type="button"
              onClick={() => setFilter(f.id)}
              className={cn(
                "min-h-10 flex-1 rounded-lg px-3 text-sm font-semibold transition sm:flex-none",
                filter === f.id
                  ? "bg-primary-container text-on-primary"
                  : "text-on-surface-variant hover:text-on-surface",
              )}
            >
              {f.label}
            </button>
          ))}
        </div>
        <div className="w-full lg:max-w-sm">
          <TextField
            name="search"
            placeholder="Search by title"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            leadingIcon={<Icon name="search" />}
            trailing={
              query ? (
                <button
                  type="button"
                  className="inline-flex size-10 items-center justify-center"
                  aria-label="Clear search"
                  onClick={() => setQuery("")}
                >
                  <Icon name="cancel" />
                </button>
              ) : null
            }
          />
        </div>
      </div>

      {isLoading ? (
        <div className="flex justify-center py-16">
          <Spinner className="text-3xl" />
        </div>
      ) : isError ? (
        <p className="rounded-xl bg-error-container p-4 text-sm text-on-error-container">
          Could not load transfers.
        </p>
      ) : items.length === 0 ? (
        <div className="rounded-xl border border-dashed border-outline-variant bg-surface-container-lowest px-6 py-14 text-center">
          <Icon name="inbox" className="mx-auto text-4xl text-outline" />
          <h2 className="mt-3 text-lg font-bold">No transfers here</h2>
          <p className="mt-1 text-sm text-on-surface-variant">
            Create a new transfer to share a file.
          </p>
          <Link href="/transfers/new" className="mt-5 inline-block">
            <Button>
              <Icon name="add" />
              New Transfer
            </Button>
          </Link>
        </div>
      ) : (
        <div className="flex flex-col gap-3">
          {items.map((t) => (
            <TransferCard key={t.id} transfer={t} />
          ))}
        </div>
      )}
    </div>
  );
}
