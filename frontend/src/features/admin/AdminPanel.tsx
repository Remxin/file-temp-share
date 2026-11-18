"use client";

import { useMemo, useState } from "react";

import { Button } from "@/components/ui/Button";
import { Icon } from "@/components/ui/Icon";
import { Spinner } from "@/components/ui/Spinner";
import { TextField } from "@/components/ui/TextField";
import {
  useListAdminTransfersQuery,
  useListAdminUsersQuery,
  useSetTransferStatusMutation,
  useSetUserActiveMutation,
} from "@/features/admin/adminApi";
import { selectCurrentUser } from "@/features/auth/selectors";
import { cn } from "@/lib/cn";
import {
  displayTitle,
  formatBytes,
  formatDateTime,
  formatExpiresAt,
} from "@/lib/format";
import { useAppSelector } from "@/store/hooks";

type Tab = "users" | "transfers";

function statusBadgeClass(status: string): string {
  switch (status) {
    case "blocked":
      return "bg-error-container text-on-error-container";
    case "expired":
      return "bg-surface-container-high text-on-surface-variant";
    default:
      return "bg-tertiary-fixed/40 text-tertiary";
  }
}

export function AdminPanel() {
  const me = useAppSelector(selectCurrentUser);
  const [tab, setTab] = useState<Tab>("users");
  const [query, setQuery] = useState("");

  const usersQuery = useListAdminUsersQuery(undefined, { skip: tab !== "users" });
  const transfersQuery = useListAdminTransfersQuery(undefined, {
    skip: tab !== "transfers",
  });
  const [setUserActive, setUserActiveState] = useSetUserActiveMutation();
  const [setTransferStatus, setTransferStatusState] =
    useSetTransferStatusMutation();

  const actionError =
    (setUserActiveState.error && "Failed to update user") ||
    (setTransferStatusState.error && "Failed to update transfer") ||
    null;

  const users = useMemo(() => {
    const list = usersQuery.data ?? [];
    const q = query.trim().toLowerCase();
    if (!q) return list;
    return list.filter((u) => u.email.toLowerCase().includes(q));
  }, [usersQuery.data, query]);

  const transfers = useMemo(() => {
    const list = transfersQuery.data ?? [];
    const q = query.trim().toLowerCase();
    if (!q) return list;
    return list.filter(
      (t) =>
        (t.title ?? "").toLowerCase().includes(q) ||
        t.owner_email.toLowerCase().includes(q) ||
        t.short_code.toLowerCase().includes(q),
    );
  }, [transfersQuery.data, query]);

  const loading = tab === "users" ? usersQuery.isLoading : transfersQuery.isLoading;
  const errored = tab === "users" ? usersQuery.isError : transfersQuery.isError;

  return (
    <div className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">Admin</h1>
        <p className="mt-1 text-sm text-on-surface-variant">
          Block users and transfers across the platform.
        </p>
      </div>

      <div className="flex flex-col gap-3 lg:flex-row lg:items-center lg:justify-between">
        <div className="flex gap-1 rounded-xl bg-surface-container-low p-1">
          {(
            [
              { id: "users", label: "Users" },
              { id: "transfers", label: "Transfers" },
            ] as const
          ).map((item) => (
            <button
              key={item.id}
              type="button"
              onClick={() => {
                setTab(item.id);
                setQuery("");
              }}
              className={cn(
                "min-h-10 flex-1 rounded-lg px-3 text-sm font-semibold transition sm:flex-none",
                tab === item.id
                  ? "bg-primary-container text-on-primary"
                  : "text-on-surface-variant hover:text-on-surface",
              )}
            >
              {item.label}
            </button>
          ))}
        </div>
        <div className="w-full lg:max-w-sm">
          <TextField
            name="search"
            placeholder={
              tab === "users"
                ? "Search by email"
                : "Search by title, owner, or code"
            }
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

      {actionError ? (
        <p className="rounded-xl bg-error-container p-4 text-sm text-on-error-container">
          {actionError}
        </p>
      ) : null}

      {loading ? (
        <div className="flex justify-center py-16">
          <Spinner className="text-3xl" />
        </div>
      ) : errored ? (
        <p className="rounded-xl bg-error-container p-4 text-sm text-on-error-container">
          Could not load {tab}. Check that you are signed in as an admin.
        </p>
      ) : tab === "users" ? (
        <div className="flex flex-col gap-3">
          {users.length === 0 ? (
            <EmptyState message="No users match this search." />
          ) : (
            users.map((user) => {
              const isSelf = user.id === me?.id;
              const busy =
                setUserActiveState.isLoading &&
                setUserActiveState.originalArgs?.id === user.id;
              return (
                <article
                  key={user.id}
                  className="rounded-xl border border-outline-variant/30 bg-surface-container-lowest p-4 shadow-sm sm:p-5"
                >
                  <div className="flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
                    <div className="min-w-0 space-y-1">
                      <div className="flex flex-wrap items-center gap-2">
                        <h3 className="truncate text-base font-bold text-on-surface">
                          {user.email}
                        </h3>
                        <span
                          className={cn(
                            "rounded-full px-2 py-0.5 text-xs font-semibold",
                            user.is_active
                              ? "bg-tertiary-fixed/40 text-tertiary"
                              : "bg-error-container text-on-error-container",
                          )}
                        >
                          {user.is_active ? "Active" : "Inactive"}
                        </span>
                        <span className="rounded-full bg-surface-container px-2 py-0.5 text-xs font-semibold text-on-surface-variant">
                          {user.role}
                        </span>
                      </div>
                      <p className="text-sm text-on-surface-variant">
                        Joined {formatDateTime(user.created_at)}
                        {isSelf ? " · you" : ""}
                      </p>
                    </div>
                    <Button
                      variant={user.is_active ? "danger" : "secondary"}
                      loading={busy}
                      disabled={isSelf && user.is_active}
                      onClick={() =>
                        void setUserActive({
                          id: user.id,
                          body: { is_active: !user.is_active },
                        })
                      }
                    >
                      <Icon name={user.is_active ? "block" : "check_circle"} />
                      {user.is_active ? "Deactivate" : "Activate"}
                    </Button>
                  </div>
                </article>
              );
            })
          )}
        </div>
      ) : (
        <div className="flex flex-col gap-3">
          {transfers.length === 0 ? (
            <EmptyState message="No transfers match this search." />
          ) : (
            transfers.map((transfer) => {
              const blocked = transfer.status === "blocked";
              const busy =
                setTransferStatusState.isLoading &&
                setTransferStatusState.originalArgs?.id === transfer.id;
              return (
                <article
                  key={transfer.id}
                  className="rounded-xl border border-outline-variant/30 bg-surface-container-lowest p-4 shadow-sm sm:p-5"
                >
                  <div className="flex flex-col gap-3 sm:flex-row sm:items-start sm:justify-between">
                    <div className="min-w-0 space-y-1">
                      <div className="flex flex-wrap items-center gap-2">
                        <h3 className="truncate text-base font-bold text-on-surface">
                          {displayTitle(transfer.title)}
                        </h3>
                        <span
                          className={cn(
                            "rounded-full px-2 py-0.5 text-xs font-semibold capitalize",
                            statusBadgeClass(transfer.status),
                          )}
                        >
                          {transfer.status}
                        </span>
                      </div>
                      <p className="text-sm text-on-surface-variant">
                        {transfer.owner_email} · {transfer.short_code}
                      </p>
                      <p className="text-sm text-on-surface-variant">
                        {transfer.file_count}{" "}
                        {transfer.file_count === 1 ? "file" : "files"} ·{" "}
                        {formatBytes(transfer.total_size_bytes)} ·{" "}
                        {transfer.download_count} downloads
                      </p>
                      <p className="inline-flex items-center gap-1 text-sm text-on-surface-variant">
                        <Icon name="schedule" className="text-[1rem]" />
                        {formatExpiresAt(transfer.expires_at)}
                      </p>
                    </div>
                    <Button
                      variant={blocked ? "secondary" : "danger"}
                      loading={busy}
                      onClick={() =>
                        void setTransferStatus({
                          id: transfer.id,
                          body: { status: blocked ? "active" : "blocked" },
                        })
                      }
                    >
                      <Icon name={blocked ? "lock_open" : "block"} />
                      {blocked ? "Unblock" : "Block"}
                    </Button>
                  </div>
                </article>
              );
            })
          )}
        </div>
      )}
    </div>
  );
}

function EmptyState({ message }: { message: string }) {
  return (
    <div className="rounded-xl border border-dashed border-outline-variant bg-surface-container-lowest px-6 py-12 text-center">
      <p className="text-sm text-on-surface-variant">{message}</p>
    </div>
  );
}
