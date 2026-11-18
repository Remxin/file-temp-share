"use client";

import { Alert } from "@/components/ui/Alert";
import { Badge } from "@/components/ui/Badge";
import { Card } from "@/components/ui/Card";
import { Spinner } from "@/components/ui/Spinner";
import { useMeQuery } from "@/features/auth/authApi";
import { formatDateTime } from "@/lib/format";

export function AccountOverview() {
  const { data: user, isLoading, isError, isFetching } = useMeQuery();

  if (isLoading) {
    return (
      <div className="flex justify-center py-16">
        <Spinner className="text-3xl" />
      </div>
    );
  }

  if (isError || !user) {
    return (
      <Alert variant="error">Could not load your account. Try again later.</Alert>
    );
  }

  return (
    <Card className="max-w-xl space-y-6">
      <div className="flex items-start justify-between gap-4">
        <div>
          <h1 className="text-2xl font-bold tracking-tight">Settings</h1>
          <p className="mt-1 text-sm text-on-surface-variant">
            Your account details
            {isFetching ? " · Refreshing…" : null}
          </p>
        </div>
        <Badge
          dot
          dotClassName={user.is_active ? "bg-tertiary" : "bg-error"}
          className={
            user.is_active
              ? "bg-tertiary-fixed/40 text-tertiary"
              : "bg-error-container text-on-error-container"
          }
        >
          {user.is_active ? "Active" : "Inactive"}
        </Badge>
      </div>

      <dl className="divide-y divide-outline-variant/30">
        <Row label="Email" value={user.email} />
        <Row
          label="Role"
          value={
            <Badge className="normal-case tracking-normal">
              {user.role === "ADMIN" ? "Admin" : "User"}
            </Badge>
          }
        />
        <Row label="Member since" value={formatDateTime(user.created_at)} />
      </dl>
    </Card>
  );
}

function Row({
  label,
  value,
}: {
  label: string;
  value: React.ReactNode;
}) {
  return (
    <div className="flex flex-col gap-1 py-4 first:pt-0 last:pb-0 sm:flex-row sm:items-center sm:justify-between sm:gap-6">
      <dt className="text-sm font-medium text-on-surface-variant">{label}</dt>
      <dd className="min-w-0 text-sm font-semibold text-on-surface sm:text-right">
        {value}
      </dd>
    </div>
  );
}
