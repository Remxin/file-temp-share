"use client";

import { Alert } from "@/components/ui/Alert";
import { Badge } from "@/components/ui/Badge";
import { Button } from "@/components/ui/Button";
import { Card } from "@/components/ui/Card";
import { Icon } from "@/components/ui/Icon";
import { Spinner } from "@/components/ui/Spinner";
import { useMeQuery } from "@/features/auth/authApi";
import { useLogout } from "@/features/auth/useLogout";
import { formatDateTime } from "@/lib/format";

export function AccountOverview() {
  const { data: user, isLoading, isError, isFetching } = useMeQuery();
  const logout = useLogout();

  if (isLoading) {
    return (
      <div className="flex justify-center py-16">
        <Spinner className="text-3xl" />
      </div>
    );
  }

  if (isError || !user) {
    return (
      <Alert variant="error">
        Could not load your account. Try again later.
      </Alert>
    );
  }

  return (
    <div className="space-y-6">
      <header className="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
        <div className="min-w-0">
          <h1 className="text-2xl font-bold tracking-tight sm:text-3xl">
            Settings
          </h1>
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
              ? "w-fit bg-tertiary-fixed/40 text-tertiary"
              : "w-fit bg-error-container text-on-error-container"
          }
        >
          {user.is_active ? "Active" : "Inactive"}
        </Badge>
      </header>

      <Card className="space-y-1">
        <h2 className="text-base font-bold tracking-tight sm:text-lg">
          Account
        </h2>
        <dl className="mt-3 divide-y divide-outline-variant/30">
          <Row label="Email" value={user.email} breakAll />
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

      <Card className="space-y-3">
        <div>
          <h2 className="text-base font-bold tracking-tight sm:text-lg">
            Session
          </h2>
          <p className="mt-1 text-sm text-on-surface-variant">
            Sign out on this device. You can sign back in anytime.
          </p>
        </div>
        <Button
          variant="danger"
          fullWidth
          className="sm:w-auto"
          onClick={logout}
        >
          <Icon name="logout" />
          Sign out
        </Button>
      </Card>
    </div>
  );
}

function Row({
  label,
  value,
  breakAll,
}: {
  label: string;
  value: React.ReactNode;
  breakAll?: boolean;
}) {
  return (
    <div className="flex flex-col gap-1.5 py-3.5 first:pt-0 last:pb-0 sm:flex-row sm:items-center sm:justify-between sm:gap-6">
      <dt className="shrink-0 text-sm font-medium text-on-surface-variant">
        {label}
      </dt>
      <dd
        className={
          breakAll
            ? "min-w-0 break-all text-sm font-semibold text-on-surface sm:text-right"
            : "min-w-0 text-sm font-semibold text-on-surface sm:text-right"
        }
      >
        {value}
      </dd>
    </div>
  );
}
